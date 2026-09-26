package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/bogem/id3v2/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_SanitizesFileAtSamePath(t *testing.T) {
	// Given
	fixture, err := os.ReadFile("testdata/Silence.mp3")
	require.NoError(t, err)
	trackPath := filepath.Join(t.TempDir(), "test track.mp3")
	require.NoError(t, os.WriteFile(trackPath, fixture, 0o600))

	// When
	var stdout, stderr bytes.Buffer
	status := run([]string{trackPath}, &stdout, &stderr)

	// Then
	assert.Zero(t, status)
	assert.Contains(t, stdout.String(), "updated")
	assert.Empty(t, stderr.String())
	updated, err := os.ReadFile(trackPath)
	require.NoError(t, err)
	tag, err := id3v2.ParseReader(bytes.NewReader(updated), id3v2.Options{Parse: true})
	require.NoError(t, err)
	assert.Equal(t, "What a tune", tag.Title())
	stdout.Reset()
	status = run([]string{trackPath}, &stdout, &stderr)
	assert.Zero(t, status)
	assert.Contains(t, stdout.String(), "unchanged")
	assert.Empty(t, stderr.String())
}

func TestRun_ReportsUsageAndFileErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		message string
	}{
		{"missing argument", nil, "usage: id3-sanitizer <mp3>"},
		{"extra argument", []string{"a.mp3", "b.mp3"}, "usage: id3-sanitizer <mp3>"},
		{"missing file", []string{filepath.Join(t.TempDir(), "missing.mp3")}, "id3-sanitizer:"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			var stdout, stderr bytes.Buffer

			// When
			status := run(test.args, &stdout, &stderr)

			// Then
			assert.NotZero(t, status)
			assert.Contains(t, stderr.String(), test.message)
			assert.Empty(t, stdout.String())
		})
	}
}
