package main

import (
	"context"
	"path/filepath"

	"github.com/lburfeindt/trackops/internal/playlist"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	context context.Context
}

type PlaylistResult struct {
	Name      string   `json:"name"`
	Titles    []string `json:"titles"`
	Error     string   `json:"error"`
	Cancelled bool     `json:"cancelled"`
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
		return PlaylistResult{Error: err.Error()}
	}
	if path == "" {
		return PlaylistResult{Cancelled: true}
	}

	titles, err := playlist.Open(path)
	if err != nil {
		return PlaylistResult{Name: filepath.Base(path), Error: err.Error()}
	}
	return PlaylistResult{Name: filepath.Base(path), Titles: titles}
}
