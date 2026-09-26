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
	// Given a playlist fixture and a path containing spaces
	fixture, err := os.ReadFile("testdata/playlist.m3u8")
	require.NoError(t, err)
	playlistPath := filepath.Join(t.TempDir(), "frequency sessions.m3u8")
	require.NoError(t, os.WriteFile(playlistPath, fixture, 0o600))

	var stdout, stderr bytes.Buffer
	// When the playlist is run
	status := run([]string{playlistPath}, &stdout, &stderr)
	// Then the command prints all titles and no errors
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
			var stdout, stderr bytes.Buffer
			// Given missing or extra command-line arguments
			// When the command is run
			status := run(test.args, &stdout, &stderr)
			// Then it fails with a usage message and no standard output
			assert.NotZero(t, status)
			assert.Contains(t, stderr.String(), "usage: format-playlist <playlist>")
			assert.Empty(t, stdout.String())
		})
	}
}

func TestRun_ReportsUnreadablePlaylist(t *testing.T) {
	// Given a path to a missing playlist
	playlistPath := filepath.Join(t.TempDir(), "missing playlist.m3u8")
	var stdout, stderr bytes.Buffer
	// When the command is run
	status := run([]string{playlistPath}, &stdout, &stderr)

	// Then it fails with a playlist open error and no standard output
	assert.NotZero(t, status)
	assert.Contains(t, stderr.String(), "format-playlist: open playlist:")
	assert.Empty(t, stdout.String())
}
