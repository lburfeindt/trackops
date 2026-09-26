package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_PrintsTitlesFromPathWithSpaces(t *testing.T) {
	// Given
	playlistFixture, err := os.ReadFile("testdata/playlist.m3u8")
	require.NoError(t, err)
	playlistPath := filepath.Join(t.TempDir(), "frequency sessions.m3u8")
	require.NoError(t, os.WriteFile(playlistPath, playlistFixture, 0o600))

	// When
	var stdout, stderr bytes.Buffer
	status := run([]string{playlistPath}, &stdout, &stderr)

	// Then
	want := "Example Artist - Example Track\n" +
		"Artist One, Artist Two - Café Song, Club Mix\n"
	assert.Zero(t, status)
	assert.Equal(t, want, stdout.String())
	assert.Empty(t, stderr.String())
}

func TestRun_RejectsMissingOrExtraArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "missing"},
		{name: "extra", args: []string{"playlist.m3u8", "extra"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given

			// When
			var stdout, stderr bytes.Buffer
			status := run(test.args, &stdout, &stderr)

			// Then
			assert.NotZero(t, status)
			assert.Contains(t, stderr.String(), "usage: format-playlist <playlist>")
			assert.Empty(t, stdout.String())
		})
	}
}

func TestRun_ReportsUnreadablePlaylist(t *testing.T) {
	// Given
	invalidPlaylistPath := filepath.Join(t.TempDir(), "missing playlist.m3u8")

	// When
	var stdout, stderr bytes.Buffer
	status := run([]string{invalidPlaylistPath}, &stdout, &stderr)

	// Then
	assert.NotZero(t, status)
	assert.Contains(t, stderr.String(), "format-playlist: open playlist:")
	assert.Empty(t, stdout.String())
}
