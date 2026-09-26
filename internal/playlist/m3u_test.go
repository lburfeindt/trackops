package playlist

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestTitlesFromFrequencyFixture(t *testing.T) {
	file, err := os.Open("../../testdata/frequency_sessions.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	got, err := Titles(file)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"Forbidden Society - Distanced (Original Mix)",
		"Bcee, S.P.Y - Is Anybody out There? (S.P.Y. VIP)",
		"Brian Brainstorm, Fú, Bomsh - Get You Down (Original Mix)",
		"Mason, Princess Superstar - Perfect (Exceeder) (1991 Remix)",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Titles() = %#v, want %#v", got, want)
	}
}

func TestTitlesPreservesCommasAndUnicode(t *testing.T) {
	input := strings.NewReader("#EXTINF:268,Brian Brainstorm, Fú, Bomsh - Get You Down (Original Mix)\n")

	got, err := Titles(input)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Brian Brainstorm, Fú, Bomsh - Get You Down (Original Mix)"}
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
	input := strings.NewReader("#EXTM3U\r\n#EXTINF:270,Forbidden Society - Distanced\r\n/music/track.mp3\r\n")

	got, err := Titles(input)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Forbidden Society - Distanced"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Titles() = %#v, want %#v", got, want)
	}
}
