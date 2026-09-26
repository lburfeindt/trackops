package playlist

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestTitlesFromFixture(t *testing.T) {
	file, err := os.Open("../../testdata/playlist.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	got, err := Titles(file)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"Example Artist - Example Track",
		"Artist One, Artist Two - Café Song, Club Mix",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Titles() = %#v, want %#v", got, want)
	}
}

func TestTitlesPreservesCommasAndUnicode(t *testing.T) {
	input := strings.NewReader("#EXTINF:240,Artist One, Artist Two - Café Song, Club Mix\n")

	got, err := Titles(input)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Artist One, Artist Two - Café Song, Club Mix"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Titles() = %#v, want %#v", got, want)
	}
}

func TestTitlesIgnoresMalformedAndNonEntryLines(t *testing.T) {
	input := strings.NewReader("#EXTM3U\n/music/track.mp3\n#EXTINF:12\n")

	got, err := Titles(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("Titles() = %#v, want no titles", got)
	}
}

func TestTitlesHandlesCRLF(t *testing.T) {
	input := strings.NewReader("#EXTM3U\r\n#EXTINF:180,Example Artist - Example Track\r\n/music/example-track.mp3\r\n")

	got, err := Titles(input)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Example Artist - Example Track"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Titles() = %#v, want %#v", got, want)
	}
}
