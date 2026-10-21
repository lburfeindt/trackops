package main

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lburfeindt/trackops/internal/id3sanitizer"
	"github.com/lburfeindt/trackops/internal/playlist"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	context      context.Context
	playlistPath string
	titles       []string
	sanitizeMu   sync.Mutex
	cancelBatch  context.CancelFunc
}

type PlaylistResult struct {
	Name      string   `json:"name"`
	Titles    []string `json:"titles"`
	Error     string   `json:"error"`
	Cancelled bool     `json:"cancelled"`
}

type SaveResult struct {
	Exists bool   `json:"exists"`
	Error  string `json:"error"`
}

type SanitizeFailure struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

type SanitizeResult struct {
	Updated   int               `json:"updated"`
	Unchanged int               `json:"unchanged"`
	Failures  []SanitizeFailure `json:"failures"`
	Error     string            `json:"error"`
	Cancelled bool              `json:"cancelled"`
	Stopped   bool              `json:"stopped"`
}

func (app *App) startup(context context.Context) {
	app.context = context
}

func (app *App) OpenPlaylist() PlaylistResult {
	path, err := runtime.OpenFileDialog(app.context, runtime.OpenDialogOptions{
		Title: "Choose an M3U playlist",
		Filters: []runtime.FileFilter{{
			DisplayName: "M3U playlists (*.m3u;*.m3u8)",
			Pattern:     "*.m3u;*.m3u8",
		}},
	})
	if err != nil {
		app.playlistPath = ""
		app.titles = nil
		return PlaylistResult{Error: err.Error()}
	}
	if path == "" {
		return PlaylistResult{Cancelled: true}
	}

	titles, err := playlist.Open(path)
	if err != nil {
		app.playlistPath = ""
		app.titles = nil
		return PlaylistResult{Name: filepath.Base(path), Error: err.Error()}
	}
	app.playlistPath = path
	app.titles = titles
	return PlaylistResult{Name: filepath.Base(path), Titles: titles}
}

func (app *App) SaveTitles(replace bool) SaveResult {
	if app.playlistPath == "" {
		return SaveResult{Error: "No playlist selected"}
	}
	err := playlist.SaveTitles(app.playlistPath, app.titles, replace)
	if errors.Is(err, playlist.ErrOutputExists) {
		return SaveResult{Exists: true}
	}
	if err != nil {
		return SaveResult{Error: err.Error()}
	}
	return SaveResult{}
}

func (app *App) CopyTitles(titles []string) error {
	return runtime.ClipboardSetText(app.context, strings.Join(titles, "\n"))
}

func (app *App) SanitizeMP3() SanitizeResult {
	paths, err := selectMP3Inputs(app.context)
	if err != nil {
		return SanitizeResult{Error: err.Error()}
	}
	if len(paths) == 0 {
		return SanitizeResult{Cancelled: true}
	}
	return app.sanitizePaths(paths)
}

func (app *App) StopSanitize() {
	app.sanitizeMu.Lock()
	defer app.sanitizeMu.Unlock()
	if app.cancelBatch != nil {
		app.cancelBatch()
	}
}

func (app *App) sanitizePaths(paths []string) SanitizeResult {
	app.sanitizeMu.Lock()
	if app.cancelBatch != nil {
		app.sanitizeMu.Unlock()
		return SanitizeResult{Error: "Sanitizer is already running"}
	}
	ctx, cancel := context.WithCancel(app.context)
	app.cancelBatch = cancel
	app.sanitizeMu.Unlock()
	defer func() {
		app.sanitizeMu.Lock()
		app.cancelBatch = nil
		app.sanitizeMu.Unlock()
		cancel()
	}()

	runtime.EventsEmit(app.context, "sanitize-progress", "Preparing selected files…")
	result := id3sanitizer.Run(ctx, paths, func(path string) {
		runtime.EventsEmit(app.context, "sanitize-progress", "Sanitizing: "+filepath.Base(path))
	})
	response := SanitizeResult{Updated: result.Updated, Unchanged: result.Unchanged, Failures: []SanitizeFailure{}, Stopped: result.Canceled}
	for _, failure := range result.Failures {
		response.Failures = append(response.Failures, SanitizeFailure{
			Name: filepath.Base(failure.Path), Error: failure.Err.Error(),
		})
	}
	return response
}
