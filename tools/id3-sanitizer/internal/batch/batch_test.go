package batch

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/bogem/id3v2/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_RecursesFiltersAndDeduplicates(t *testing.T) {
	// Given
	directory := t.TempDir()
	nested := filepath.Join(directory, "album")
	require.NoError(t, os.Mkdir(nested, 0o700))
	firstTrack := copyTrack(t, directory, "first.mp3")
	secondTrack := copyTrack(t, nested, "second.MP3")
	ignoredTrack := copyTrack(t, directory, "ignored.wav")
	original, err := os.ReadFile(ignoredTrack)
	require.NoError(t, err)

	// When
	result := Run(context.Background(), []string{directory, firstTrack, nested}, nil)

	// Then
	assert.Equal(t, 2, result.Updated)
	assert.Zero(t, result.Unchanged)
	assert.Empty(t, result.Failures)
	assert.False(t, result.Canceled)
	for _, path := range []string{firstTrack, secondTrack} {
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		tag, err := id3v2.ParseReader(bytes.NewReader(data), id3v2.Options{Parse: true})
		require.NoError(t, err)
		assert.Equal(t, "What a tune", tag.Title())
	}
	ignored, err := os.ReadFile(ignoredTrack)
	require.NoError(t, err)
	assert.Equal(t, original, ignored)
	repeated := Run(context.Background(), []string{directory}, nil)
	assert.Equal(t, 2, repeated.Unchanged)
	assert.Zero(t, repeated.Updated)
}

func TestRun_ContinuesAfterFileErrors(t *testing.T) {
	// Given
	directory := t.TempDir()
	broken := filepath.Join(directory, "broken.mp3")
	require.NoError(t, os.WriteFile(broken, []byte("ID3"), 0o600))
	missing := filepath.Join(directory, "missing.mp3")
	valid := copyTrack(t, directory, "valid.mp3")

	// When
	result := Run(context.Background(), []string{broken, missing, valid}, nil)

	// Then
	assert.Equal(t, 1, result.Updated)
	require.Len(t, result.Failures, 2)
	assert.Equal(t, missing, result.Failures[0].Path)
	assert.Equal(t, broken, result.Failures[1].Path)
	for _, failure := range result.Failures {
		assert.Error(t, failure.Err)
	}
	contents, err := os.ReadFile(broken)
	require.NoError(t, err)
	assert.Equal(t, []byte("ID3"), contents)
}

func TestRun_CancellationLeavesRemainingFilesUntouched(t *testing.T) {
	// Given
	directory := t.TempDir()
	first := copyTrack(t, directory, "first.mp3")
	second := copyTrack(t, directory, "second.mp3")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// When
	result := Run(ctx, []string{first, second}, func(path string) {
		if path == second {
			cancel()
		}
	})

	// Then
	assert.True(t, result.Canceled)
	assert.Equal(t, 1, result.Updated)
	remaining := Run(context.Background(), []string{first, second}, nil)
	assert.Equal(t, 1, remaining.Updated)
	assert.Equal(t, 1, remaining.Unchanged)
}

func TestRun_DoesNotFollowSymlinks(t *testing.T) {
	// Given
	outside := t.TempDir()
	target := copyTrack(t, outside, "target.mp3")
	directory := t.TempDir()
	link := filepath.Join(directory, "link.mp3")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	require.NoError(t, os.Symlink(outside, filepath.Join(directory, "linked-album")))

	// When
	result := Run(context.Background(), []string{directory, link}, nil)

	// Then
	assert.Zero(t, result.Updated)
	assert.Zero(t, result.Unchanged)
	require.Len(t, result.Failures, 1)
	assert.Equal(t, link, result.Failures[0].Path)
	untouched := Run(context.Background(), []string{target}, nil)
	assert.Equal(t, 1, untouched.Updated)
}

func TestRun_EmptyOrCanceledSelection(t *testing.T) {
	// Given
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	track := copyTrack(t, t.TempDir(), "track.mp3")

	// When
	empty := Run(context.Background(), []string{t.TempDir()}, nil)
	stopped := Run(canceled, []string{track}, nil)

	// Then
	assert.Equal(t, Result{}, empty)
	assert.True(t, stopped.Canceled)
	assert.Zero(t, stopped.Updated)
	assert.Empty(t, stopped.Failures)
}

func copyTrack(t *testing.T, directory, name string) string {
	t.Helper()
	fixture, err := os.ReadFile("../../testdata/Silence.mp3")
	require.NoError(t, err)
	path := filepath.Join(directory, name)
	require.NoError(t, os.WriteFile(path, fixture, 0o600))
	return path
}
