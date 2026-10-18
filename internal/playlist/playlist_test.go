package playlist

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenShowsTitlesWithoutWritingOutput(t *testing.T) {
	// Given
	playlistPath := filepath.Join(t.TempDir(), "Set.m3u8")
	contents := "#EXTM3U\n#EXTINF:123,Artist One, Artist Two - Café Song, Club Mix\ntrack.mp3\n#EXTINF:99,Next - Title\nnext.mp3\n"
	if err := os.WriteFile(playlistPath, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}

	// When
	titles, err := Open(playlistPath)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Artist One, Artist Two - Café Song, Club Mix", "Next - Title"}
	if len(titles) != len(want) {
		t.Fatalf("titles = %q, want %q", titles, want)
	}
	for index := range want {
		if titles[index] != want[index] {
			t.Fatalf("titles = %q, want %q", titles, want)
		}
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(playlistPath), "Set-formatted.txt")); !os.IsNotExist(err) {
		t.Fatalf("opening playlist created output: %v", err)
	}
}
