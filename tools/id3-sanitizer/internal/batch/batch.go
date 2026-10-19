// Package batch sanitizes selected MP3 files and directory trees.
package batch

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/lburfeindt/trackops/tools/id3-sanitizer/internal/sanitizer"
)

// Failure identifies an input that could not be read or sanitized.
type Failure struct {
	Path string
	Err  error
}

// Result describes a batch, including any work completed before cancellation.
type Result struct {
	Updated   int
	Unchanged int
	Failures  []Failure
	Canceled  bool
}

// Run discovers MP3s recursively without following symlinks, then sanitizes each
// unique path once. Errors do not stop other files. Cancellation is checked
// between files; an in-progress file replacement is allowed to finish.
// BeforeFile, when supplied, runs synchronously before each file is modified.
func Run(ctx context.Context, paths []string, beforeFile func(string)) Result {
	var result Result
	var files []string
	seen := make(map[string]bool)
	fail := func(path string, err error) {
		result.Failures = append(result.Failures, Failure{Path: path, Err: err})
	}
	for _, path := range paths {
		if ctx.Err() != nil {
			break
		}
		absolute, err := filepath.Abs(path)
		if err != nil {
			fail(path, err)
			continue
		}
		info, err := os.Lstat(absolute)
		if err != nil {
			fail(path, err)
			continue
		}
		if !info.IsDir() && !info.Mode().IsRegular() {
			fail(path, fmt.Errorf("selection must be a regular file or directory"))
			continue
		}
		if info.Mode().IsRegular() && !strings.EqualFold(filepath.Ext(absolute), ".mp3") {
			fail(path, fmt.Errorf("selection must be an MP3 file"))
			continue
		}
		err = filepath.WalkDir(absolute, func(current string, entry fs.DirEntry, walkErr error) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if walkErr != nil {
				fail(current, walkErr)
				return nil
			}
			if entry.Type().IsRegular() && strings.EqualFold(filepath.Ext(current), ".mp3") && !seen[current] {
				seen[current] = true
				files = append(files, current)
			}
			return nil
		})
		if err != nil && ctx.Err() == nil {
			fail(path, err)
		}
	}
	for _, path := range files {
		if ctx.Err() != nil {
			break
		}
		if beforeFile != nil {
			beforeFile(path)
		}
		if ctx.Err() != nil {
			break
		}
		changed, err := sanitizer.File(path)
		switch {
		case err != nil:
			fail(path, err)
		case changed:
			result.Updated++
		default:
			result.Unchanged++
		}
	}
	result.Canceled = ctx.Err() != nil
	return result
}
