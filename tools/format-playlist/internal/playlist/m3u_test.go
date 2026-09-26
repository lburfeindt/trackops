package playlist

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTitles_FromFixture(t *testing.T) {
	// Given a playlist fixture
	file, err := os.Open("../../testdata/playlist.m3u8")
	require.NoError(t, err)
	defer file.Close()

	// When the playlist titles are extracted
	got, err := Titles(file)
	require.NoError(t, err)

	want := []string{
		"Example Artist - Example Track",
		"Artist One, Artist Two - Café Song, Club Mix",
	}
	// Then the extracted titles match the fixture
	assert.Equal(t, want, got)
}

func TestTitles_PreservesCommasAndUnicode(t *testing.T) {
	// Given an entry containing commas and Unicode characters
	input := strings.NewReader("#EXTINF:240,Artist One, Artist Two - Café Song, Club Mix\n")

	// When the playlist title is extracted
	got, err := Titles(input)
	require.NoError(t, err)
	want := []string{"Artist One, Artist Two - Café Song, Club Mix"}
	// Then commas and Unicode characters are preserved
	assert.Equal(t, want, got)
}

func TestTitles_IgnoresMalformedAndNonEntryLines(t *testing.T) {
	// Given malformed and non-entry playlist lines
	input := strings.NewReader("#EXTM3U\n/music/track.mp3\n#EXTINF:12\n")

	// When the playlist titles are extracted
	got, err := Titles(input)
	require.NoError(t, err)
	// Then no titles are returned
	assert.Empty(t, got)
}

func TestTitles_HandlesCRLF(t *testing.T) {
	// Given a playlist with CRLF line endings
	input := strings.NewReader("#EXTM3U\r\n#EXTINF:180,Example Artist - Example Track\r\n/music/example-track.mp3\r\n")

	// When the playlist titles are extracted
	got, err := Titles(input)
	require.NoError(t, err)
	want := []string{"Example Artist - Example Track"}
	// Then the title is extracted correctly
	assert.Equal(t, want, got)
}

func TestTitles_ReturnsReaderError(t *testing.T) {
	// Given a reader that returns an error
	wantErr := errors.New("reader failure")
	input := errorReader{err: wantErr}

	// When the playlist titles are extracted
	got, err := Titles(input)

	// Then the reader error is returned and no titles are produced
	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, got)
}

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}
