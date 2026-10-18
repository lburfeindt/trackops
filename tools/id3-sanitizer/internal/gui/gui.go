// Package gui provides the sanitizer's native dialog workflow.
package gui

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lburfeindt/trackops/tools/id3-sanitizer/internal/batch"
	"github.com/ncruces/zenity"
)

// Progress is the portion of a native progress dialog used by the workflow.
type Progress interface {
	Text(string) error
	Done() <-chan struct{}
	Close() error
}

// Dialogs supplies the user interactions independently of batch processing.
type Dialogs interface {
	Choose() ([]string, error)
	Progress() (Progress, error)
	Result(batch.Result) error
}

// Run selects inputs, processes them, and reports the outcome. Canceling a
// chooser is a no-op; canceling progress stops after the current file.
func Run(dialogs Dialogs) error {
	paths, err := dialogs.Choose()
	if errors.Is(err, zenity.ErrCanceled) {
		return nil
	}
	if err != nil || len(paths) == 0 {
		return err
	}
	progress, err := dialogs.Progress()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-progress.Done():
			cancel()
		case <-ctx.Done():
		}
	}()
	var updateErr error
	update := func(text string) {
		select {
		case <-progress.Done():
			cancel()
			return
		default:
		}
		if err := progress.Text(text); err != nil {
			updateErr = ignoreCancel(err)
			cancel()
		}
	}
	update("Finding MP3 files…")
	result := batch.Run(ctx, paths, func(path string) { update("Sanitizing: " + path) })
	closeErr := ignoreCancel(progress.Close())
	resultErr := ignoreCancel(dialogs.Result(result))
	return errors.Join(updateErr, closeErr, resultErr)
}

func ignoreCancel(err error) error {
	if errors.Is(err, zenity.ErrCanceled) {
		return nil
	}
	return err
}

func summary(result batch.Result) string {
	var text strings.Builder
	if result.Canceled {
		text.WriteString("Canceled. Completed changes have been saved.\n\n")
	} else if result.Updated+result.Unchanged+len(result.Failures) == 0 {
		text.WriteString("No MP3 files found.\n\n")
	}
	fmt.Fprintf(&text, "Updated: %d\nUnchanged: %d\nFailed: %d", result.Updated, result.Unchanged, len(result.Failures))
	for i, failure := range result.Failures {
		if i == 10 {
			fmt.Fprintf(&text, "\n…and %d more errors.", len(result.Failures)-i)
			break
		}
		fmt.Fprintf(&text, "\n\n%s: %v", failure.Path, failure.Err)
	}
	return text.String()
}
