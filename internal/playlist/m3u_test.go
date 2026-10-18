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
	// Given
	file, err := os.Open("testdata/playlist.m3u8")
	require.NoError(t, err)
	defer file.Close()

	// When
	got, err := Titles(file)
	require.NoError(t, err)

	want := []string{
		"Example Artist - Example Track",
		"Artist One, Artist Two - Café Song, Club Mix",
	}
	// Then
	assert.Equal(t, want, got)
}

func TestTitles_PreservesCommasAndUnicode(t *testing.T) {
	// Given
	input := strings.NewReader("#EXTINF:240,Artist One, Artist Two - Café Song, Club Mix\n")

	// When
	got, err := Titles(input)
	require.NoError(t, err)
	want := []string{"Artist One, Artist Two - Café Song, Club Mix"}
	// Then
	assert.Equal(t, want, got)
}

func TestTitles_IgnoresMalformedAndNonEntryLines(t *testing.T) {
	// Given
	input := strings.NewReader("#EXTM3U\n/music/track.mp3\n#EXTINF:12\n")

	// When
	got, err := Titles(input)
	require.NoError(t, err)
	// Then
	assert.Empty(t, got)
}

func TestTitles_HandlesCRLF(t *testing.T) {
	// Given
	input := strings.NewReader("#EXTM3U\r\n#EXTINF:180,Example Artist - Example Track\r\n/music/example-track.mp3\r\n")

	// When
	got, err := Titles(input)
	require.NoError(t, err)
	want := []string{"Example Artist - Example Track"}
	// Then
	assert.Equal(t, want, got)
}

func TestTitles_ReturnsReaderError(t *testing.T) {
	// Given
	wantErr := errors.New("reader failure")
	input := errorReader{err: wantErr}

	// When
	got, err := Titles(input)

	// Then
	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, got)
}

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}
