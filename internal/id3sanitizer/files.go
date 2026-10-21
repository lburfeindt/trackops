// Package id3sanitizer exposes batch file processing to the TrackOps app.
package id3sanitizer

import (
	"context"

	"github.com/lburfeindt/trackops/internal/id3sanitizer/batch"
)

type Result = batch.Result

// Run sanitizes selected files and reports each file before it is processed.
func Run(ctx context.Context, paths []string, beforeFile func(string)) Result {
	return batch.Run(ctx, paths, beforeFile)
}
