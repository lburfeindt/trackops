package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPrintsTitlesFromPathWithSpaces(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/playlist.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	playlistPath := filepath.Join(t.TempDir(), "frequency sessions.m3u8")
	if err := os.WriteFile(playlistPath, fixture, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	status := run([]string{playlistPath}, &stdout, &stderr)
	if status != 0 {
		t.Fatalf("run() status = %d, stderr = %q; want 0", status, stderr.String())
	}
	want := "Example Artist - Example Track\n" +
		"Artist One, Artist Two - Café Song, Club Mix\n"
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunRejectsMissingOrExtraArguments(t *testing.T) {
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
			status := run(test.args, &stdout, &stderr)
			if status == 0 {
				t.Fatal("run() status = 0, want failure")
			}
			if !strings.Contains(stderr.String(), "usage: format-playlist <playlist>") {
				t.Fatalf("stderr = %q, want usage message", stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
		})
	}
}

func TestRunReportsUnreadablePlaylist(t *testing.T) {
	playlistPath := filepath.Join(t.TempDir(), "missing playlist.m3u8")
	var stdout, stderr bytes.Buffer
	status := run([]string{playlistPath}, &stdout, &stderr)

	if status == 0 {
		t.Fatal("run() status = 0, want failure")
	}
	if !strings.Contains(stderr.String(), "format-playlist: open playlist:") {
		t.Fatalf("stderr = %q, want playlist open error", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}
