package playlist

import (
	"bufio"
	"io"
	"strings"
)

// Titles returns the titles from M3U extended-information entries.
func Titles(r io.Reader) ([]string, error) {
	var titles []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if !strings.HasPrefix(line, "#EXTINF:") {
			continue
		}

		metadata := strings.TrimPrefix(line, "#EXTINF:")
		separator := strings.IndexByte(metadata, ',')
		if separator < 0 {
			continue
		}
		titles = append(titles, metadata[separator+1:])
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return titles, nil
}
