package playlist

import "os"

// ReadTitles reads titles from a playlist without changing files.
func ReadTitles(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseTitles(file)
}
