package playlist

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrOutputExists = errors.New("output already exists")

// SaveTitles writes titles beside the playlist, replacing output only when requested.
func SaveTitles(playlistPath string, titles []string, replace bool) error {
	extension := filepath.Ext(playlistPath)
	outputPath := strings.TrimSuffix(playlistPath, extension) + "-formatted.txt"
	if !replace {
		if _, err := os.Lstat(outputPath); err == nil {
			return ErrOutputExists
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	temporary, err := os.CreateTemp(filepath.Dir(outputPath), ".trackops-titles-*")
	if err != nil {
		return err
	}
	defer os.Remove(temporary.Name())

	contents := ""
	if len(titles) > 0 {
		contents = strings.Join(titles, "\n") + "\n"
	}
	if _, err := temporary.WriteString(contents); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}

	if replace {
		return os.Rename(temporary.Name(), outputPath)
	}
	if err := os.Link(temporary.Name(), outputPath); err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrOutputExists
		}
		return err
	}
	return nil
}
