package gui

import (
	"errors"

	"github.com/lburfeindt/trackops/tools/id3-sanitizer/internal/batch"
	"github.com/ncruces/zenity"
)

const title = "ID3 Sanitizer"

// NativeDialogs uses the operating system's dialogs through Zenity.
type NativeDialogs struct{}

func (NativeDialogs) Choose() ([]string, error) {
	err := zenity.Question(
		"Remove (Original Mix) from MP3 title endings.\n\nSelected files are modified in place. Folders include all subfolders.",
		zenity.Title(title), zenity.OKLabel("Files…"), zenity.ExtraButton("Folder…"), zenity.CancelLabel("Cancel"),
	)
	switch {
	case errors.Is(err, zenity.ErrExtraButton):
		path, err := zenity.SelectFile(zenity.Title("Choose a folder to sanitize"), zenity.Directory())
		if err != nil {
			return nil, err
		}
		return []string{path}, nil
	case err != nil:
		return nil, err
	default:
		return zenity.SelectFileMultiple(
			zenity.Title("Choose MP3 files to sanitize"),
			zenity.FileFilter{Name: "MP3 files", Patterns: []string{"*.mp3"}, CaseFold: true},
		)
	}
}

func (NativeDialogs) Progress() (Progress, error) {
	return zenity.Progress(zenity.Title(title), zenity.Pulsate())
}

func (NativeDialogs) Result(result batch.Result) error {
	if len(result.Failures) > 0 {
		return zenity.Warning(summary(result), zenity.Title(title))
	}
	return zenity.Info(summary(result), zenity.Title(title))
}
