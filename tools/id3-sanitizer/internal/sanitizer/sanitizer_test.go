package sanitizer

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"unicode/utf16"

	"github.com/bogem/id3v2/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFile_SanitizesFixtureAndPreservesAudioAndMetadata(t *testing.T) {
	// Given
	original, err := os.ReadFile("../../testdata/Silence.mp3")
	require.NoError(t, err)
	trackPath := writeTrack(t, original)
	require.NoError(t, os.Chmod(trackPath, 0o640))

	// When
	changed, err := File(trackPath)

	// Then
	require.NoError(t, err)
	assert.True(t, changed)
	updated, err := os.ReadFile(trackPath)
	require.NoError(t, err)
	tag, err := id3v2.ParseReader(bytes.NewReader(updated), id3v2.Options{Parse: true})
	require.NoError(t, err)
	assert.Equal(t, "What a tune", tag.Title())
	assert.Equal(t, "The Artist", tag.Artist())
	assert.Equal(t, "Some Album", tag.Album())
	// The fixture's tag ends at byte 149; every audio byte must survive.
	assert.Equal(t, original[149:], updated[149:])
	assert.Equal(t, original[10:31], updated[10:31])
	assert.True(t, bytes.Contains(updated[:149], original[68:149]))
	info, err := os.Stat(trackPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o640), info.Mode().Perm())
	changed, err = File(trackPath)
	require.NoError(t, err)
	assert.False(t, changed)
	repeated, err := os.ReadFile(trackPath)
	require.NoError(t, err)
	assert.Equal(t, updated, repeated)
	entries, err := os.ReadDir(filepath.Dir(trackPath))
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestFile_TitleEncodingsAndMatching(t *testing.T) {
	tests := []struct {
		name     string
		version  byte
		encoding byte
		title    string
		want     string
		changed  bool
	}{
		{"latin1", 3, 0, "Café (Original Mix)", "Café", true},
		{"utf16", 3, 1, "夜 🎵 (Original Mix)", "夜 🎵", true},
		{"utf16be", 4, 2, "夜 (Original Mix)", "夜", true},
		{"utf8", 4, 3, "Café 🎵 (Original Mix)", "Café 🎵", true},
		{"spaces", 4, 3, "Song  \t(Original Mix)", "Song", true},
		{"no space", 3, 0, "Song(Original Mix)", "Song", true},
		{"only suffix", 4, 3, "(Original Mix)", "", true},
		{"middle", 4, 3, "Song (Original Mix) Live", "Song (Original Mix) Live", false},
		{"case sensitive", 4, 3, "Song (original mix)", "Song (original mix)", false},
		{"other mix", 3, 0, "Song (Extended Mix)", "Song (Extended Mix)", false},
		{"empty", 4, 3, "", "", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			title := encodeTitle(test.title, test.encoding)
			original := makeTrack(test.version, makeFrame(test.version, "TIT2", title))
			trackPath := writeTrack(t, original)
			before, err := os.Stat(trackPath)
			require.NoError(t, err)

			// When
			changed, err := File(trackPath)

			// Then
			require.NoError(t, err)
			assert.Equal(t, test.changed, changed)
			updated, err := os.ReadFile(trackPath)
			require.NoError(t, err)
			tag, err := id3v2.ParseReader(bytes.NewReader(updated), id3v2.Options{Parse: true})
			require.NoError(t, err)
			assert.Equal(t, test.want, tag.Title())
			assert.Equal(t, test.encoding, updated[20])
			if !test.changed {
				assert.Equal(t, original, updated)
				after, err := os.Stat(trackPath)
				require.NoError(t, err)
				assert.True(t, os.SameFile(before, after))
				assert.Equal(t, before.ModTime(), after.ModTime())
			}
		})
	}
}

func TestFile_PreservesOpaqueFrames(t *testing.T) {
	for _, version := range []byte{3, 4} {
		t.Run(string(rune('0'+version)), func(t *testing.T) {
			// Given
			opaqueFrame := makeFrame(version, "PRIV", bytes.Repeat([]byte{0xff, 0x01}, 100))
			opaqueFrame[9] = 0x80
			artwork := makeFrame(version, "APIC", []byte{0, 'i', 'm', 'g', 0, 3, 0, 0xff, 0xd8})
			frames := append(makeFrame(version, "TIT2", encodeTitle("Song (Original Mix)", 0)), opaqueFrame...)
			frames = append(frames, artwork...)
			original := makeTrack(version, frames)
			trackPath := writeTrack(t, original)

			// When
			changed, err := File(trackPath)

			// Then
			require.NoError(t, err)
			assert.True(t, changed)
			updated, err := os.ReadFile(trackPath)
			require.NoError(t, err)
			assert.True(t, bytes.Contains(updated, opaqueFrame))
			assert.True(t, bytes.Contains(updated, artwork))
			assert.Equal(t, original[len(original)-5:], updated[len(updated)-5:])
		})
	}
}

func TestFile_RejectsMalformedTagsWithoutWriting(t *testing.T) {
	tests := []struct {
		name   string
		mutate func([]byte) []byte
	}{
		{"truncated header", func(b []byte) []byte { return b[:8] }},
		{"truncated body", func(b []byte) []byte { return b[:25] }},
		{"unsupported version", func(b []byte) []byte { b[3] = 2; return b }},
		{"tag flags", func(b []byte) []byte { b[5] = 0x80; return b }},
		{"invalid tag size", func(b []byte) []byte { b[6] = 0x80; return b }},
		{"frame overflow", func(b []byte) []byte { b[17] = 127; return b }},
		{"invalid frame size", func(b []byte) []byte { b[14] = 0x80; return b }},
		{"title flags", func(b []byte) []byte { b[19] = 0x08; return b }},
		{"unknown encoding", func(b []byte) []byte { b[20] = 5; return b }},
		{"invalid frame id", func(b []byte) []byte { b[10] = '!'; return b }},
		{"nonzero padding", func(b []byte) []byte { b[10] = 0; return b }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			original := test.mutate(makeTrack(4, makeFrame(4, "TIT2", encodeTitle("Song (Original Mix)", 3))))
			trackPath := writeTrack(t, original)

			// When
			changed, err := File(trackPath)

			// Then
			require.Error(t, err)
			assert.False(t, changed)
			updated, err := os.ReadFile(trackPath)
			require.NoError(t, err)
			assert.Equal(t, original, updated)
		})
	}
}

func TestFile_MissingTitleIsUnchanged(t *testing.T) {
	for _, original := range [][]byte{[]byte("untagged audio"), makeTrack(3, makeFrame(3, "TPE1", encodeTitle("Artist", 0)))} {
		// Given
		trackPath := writeTrack(t, original)

		// When
		changed, err := File(trackPath)

		// Then
		require.NoError(t, err)
		assert.False(t, changed)
		updated, err := os.ReadFile(trackPath)
		require.NoError(t, err)
		assert.Equal(t, original, updated)
	}
}

func TestFile_RejectsMissingFilesDirectoriesAndSymlinks(t *testing.T) {
	// Given
	directory := t.TempDir()
	target := writeTrack(t, []byte("audio"))
	link := filepath.Join(directory, "link.mp3")
	require.NoError(t, os.Symlink(target, link))

	// When
	for _, trackPath := range []string{filepath.Join(directory, "missing.mp3"), directory, link} {
		changed, err := File(trackPath)

		// Then
		assert.Error(t, err)
		assert.False(t, changed)
	}
}

func writeTrack(t *testing.T, data []byte) string {
	t.Helper()
	trackPath := filepath.Join(t.TempDir(), "track with spaces.mp3")
	require.NoError(t, os.WriteFile(trackPath, data, 0o600))
	return trackPath
}

func makeTrack(version byte, frames []byte) []byte {
	header := []byte{'I', 'D', '3', version, 0, 0, 0, 0, 0, 0}
	size := len(frames) + 12
	for i := 9; i >= 6; i-- {
		header[i] = byte(size & 0x7f)
		size >>= 7
	}
	track := append(header, frames...)
	track = append(track, make([]byte, 12)...)
	return append(track, 0xff, 0xfb, 0x01, 0x02, 0x03)
}

func makeFrame(version byte, id string, body []byte) []byte {
	header := make([]byte, 10)
	copy(header, id)
	binary.BigEndian.PutUint32(header[4:8], uint32(len(body)))
	if version == 4 {
		size := len(body)
		for i := 7; i >= 4; i-- {
			header[i] = byte(size & 0x7f)
			size >>= 7
		}
	}
	return append(header, body...)
}

func encodeTitle(title string, encoding byte) []byte {
	body := []byte{encoding}
	switch encoding {
	case 0:
		for _, r := range title {
			body = append(body, byte(r))
		}
	case 3:
		body = append(body, title...)
	default:
		var order binary.ByteOrder = binary.BigEndian
		if encoding == 1 {
			body = append(body, 0xff, 0xfe)
			order = binary.LittleEndian
		}
		for _, unit := range utf16.Encode([]rune(title)) {
			pair := make([]byte, 2)
			order.PutUint16(pair, unit)
			body = append(body, pair...)
		}
		body = append(body, 0)
	}
	return append(body, 0)
}

func TestReplaceFile_CopyFailureLeavesOriginalAndRemovesTemporaryFile(t *testing.T) {
	// Given
	original := []byte("original audio")
	trackPath := writeTrack(t, original)
	source, err := os.Open(trackPath)
	require.NoError(t, err)
	info, err := source.Stat()
	require.NoError(t, err)
	require.NoError(t, source.Close())

	// When
	err = replaceFile(trackPath, source, info, []byte("new tag"))

	// Then
	require.Error(t, err)
	unchanged, err := os.ReadFile(trackPath)
	require.NoError(t, err)
	assert.Equal(t, original, unchanged)
	entries, err := os.ReadDir(filepath.Dir(trackPath))
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestFile_RejectsDuplicateTitlesWithoutWriting(t *testing.T) {
	// Given
	frames := makeFrame(4, "TIT2", encodeTitle("First (Original Mix)", 3))
	frames = append(frames, makeFrame(4, "TIT2", encodeTitle("Second (Original Mix)", 3))...)
	original := makeTrack(4, frames)
	trackPath := writeTrack(t, original)

	// When
	changed, err := File(trackPath)

	// Then
	require.Error(t, err)
	assert.False(t, changed)
	unchanged, err := os.ReadFile(trackPath)
	require.NoError(t, err)
	assert.Equal(t, original, unchanged)
}

func TestFile_PreservesBigEndianBOMAndUnterminatedTitles(t *testing.T) {
	for _, title := range [][]byte{
		append([]byte{1, 0xfe, 0xff}, encodeTitle("夜 (Original Mix)", 2)[1:]...),
		append([]byte{3}, []byte("Café (Original Mix)")...),
	} {
		// Given
		original := makeTrack(4, makeFrame(4, "TIT2", title))
		trackPath := writeTrack(t, original)

		// When
		changed, err := File(trackPath)

		// Then
		require.NoError(t, err)
		assert.True(t, changed)
		updated, err := os.ReadFile(trackPath)
		require.NoError(t, err)
		tag, err := id3v2.ParseReader(bytes.NewReader(updated), id3v2.Options{Parse: true})
		require.NoError(t, err)
		if title[0] == 1 {
			assert.Equal(t, "夜", tag.Title())
			assert.Equal(t, []byte{1, 0xfe, 0xff}, updated[20:23])
		} else {
			assert.Equal(t, "Café", tag.Title())
		}
	}
}

func TestFile_RejectsMalformedUTF16WithoutWriting(t *testing.T) {
	for _, title := range [][]byte{{1}, {1, 0, 0}, {2, 0}, {1, 0xff, 0xfe, 0}} {
		// Given
		original := makeTrack(4, makeFrame(4, "TIT2", title))
		trackPath := writeTrack(t, original)

		// When
		changed, err := File(trackPath)

		// Then
		require.Error(t, err)
		assert.False(t, changed)
		unchanged, err := os.ReadFile(trackPath)
		require.NoError(t, err)
		assert.Equal(t, original, unchanged)
	}
}
