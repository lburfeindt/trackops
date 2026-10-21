//go:build !darwin

package desktop

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func selectMP3Inputs(ctx context.Context) ([]string, error) {
	return runtime.OpenMultipleFilesDialog(ctx, runtime.OpenDialogOptions{
		Title: "Choose MP3 files to sanitize",
		Filters: []runtime.FileFilter{{
			DisplayName: "MP3 files (*.mp3)",
			Pattern:     "*.mp3;*.MP3",
		}},
	})
}
