package playlist

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveTitlesCreatesSiblingOutput(t *testing.T) {
	// Given
	directory := t.TempDir()
	playlistPath := filepath.Join(directory, "Set.m3u8")
	titles := []string{"Artist - Café", "Second Artist, Club Mix"}

	// When
	err := SaveTitles(playlistPath, titles, false)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(filepath.Join(directory, "Set-formatted.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "Artist - Café\nSecond Artist, Club Mix\n" {
		t.Fatalf("output = %q", contents)
	}
}

func TestSaveTitlesCreatesEmptyOutput(t *testing.T) {
	// Given
	directory := t.TempDir()
	playlistPath := filepath.Join(directory, "Empty.m3u")

	// When
	err := SaveTitles(playlistPath, nil, false)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(filepath.Join(directory, "Empty-formatted.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(contents) != 0 {
		t.Fatalf("output = %q, want empty", contents)
	}
}

func TestSaveTitlesRequiresReplacementChoice(t *testing.T) {
	// Given
	directory := t.TempDir()
	outputPath := filepath.Join(directory, "Set-formatted.txt")
	if err := os.WriteFile(outputPath, []byte("prior\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// When
	err := SaveTitles(filepath.Join(directory, "Set.m3u8"), []string{"new"}, false)

	// Then
	if !errors.Is(err, ErrOutputExists) {
		t.Fatalf("error = %v, want ErrOutputExists", err)
	}
	contents, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "prior\n" {
		t.Fatalf("output changed: %q", contents)
	}
}

func TestSaveTitlesReplacesExistingOutput(t *testing.T) {
	// Given
	directory := t.TempDir()
	outputPath := filepath.Join(directory, "Set-formatted.txt")
	if err := os.WriteFile(outputPath, []byte("prior\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// When
	err := SaveTitles(filepath.Join(directory, "Set.m3u8"), []string{"new"}, true)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "new\n" {
		t.Fatalf("output = %q", contents)
	}
}

func TestSaveTitlesFailureKeepsExistingOutput(t *testing.T) {
	// Given
	directory := t.TempDir()
	outputPath := filepath.Join(directory, "Set-formatted.txt")
	if err := os.WriteFile(outputPath, []byte("prior\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(directory, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(directory, 0700) })

	// When
	err := SaveTitles(filepath.Join(directory, "Set.m3u8"), []string{"new"}, true)

	// Then
	if err == nil {
		t.Fatal("save succeeded in a read-only directory")
	}
	contents, readErr := os.ReadFile(outputPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(contents) != "prior\n" {
		t.Fatalf("output changed: %q", contents)
	}
}

func TestSaveTitlesOffersReplacementBeforeWriting(t *testing.T) {
	// Given
	directory := t.TempDir()
	outputPath := filepath.Join(directory, "Set-formatted.txt")
	if err := os.WriteFile(outputPath, []byte("prior\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(directory, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(directory, 0700) })

	// When
	err := SaveTitles(filepath.Join(directory, "Set.m3u8"), []string{"new"}, false)

	// Then
	if !errors.Is(err, ErrOutputExists) {
		t.Fatalf("error = %v, want ErrOutputExists", err)
	}
	contents, readErr := os.ReadFile(outputPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(contents) != "prior\n" {
		t.Fatalf("output changed: %q", contents)
	}
}
