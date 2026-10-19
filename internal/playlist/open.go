package playlist

import "os"

// Open reads titles from a playlist without changing files.
func Open(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Titles(file)
}
