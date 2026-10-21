// Package sanitizer removes the Original Mix suffix from ID3v2 track titles.
package sanitizer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// File sanitizes an ID3v2.3 or ID3v2.4 title and replaces the file at path.
// It returns whether the file changed. Nonmatching titles cause no write.
func File(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, fmt.Errorf("open mp3: %w", err)
	}
	if !info.Mode().IsRegular() {
		return false, fmt.Errorf("mp3 must be a regular file: %s", path)
	}
	source, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open mp3: %w", err)
	}
	defer source.Close()
	opened, err := source.Stat()
	if err != nil {
		return false, err
	}
	if !os.SameFile(info, opened) {
		return false, fmt.Errorf("mp3 changed while opening: %s", path)
	}
	tag, err := readTag(source, info.Size())
	if err != nil {
		return false, fmt.Errorf("read ID3 tag: %w", err)
	}
	if tag == nil {
		return false, nil
	}
	changed, err := sanitizeTag(tag)
	if err != nil {
		return false, fmt.Errorf("sanitize ID3 tag: %w", err)
	}
	if !changed {
		return false, nil
	}
	if err := replaceFile(path, source, info, tag); err != nil {
		return false, fmt.Errorf("save mp3: %w", err)
	}
	return true, nil
}

func readTag(source io.Reader, fileSize int64) ([]byte, error) {
	header := make([]byte, 10)
	n, err := io.ReadFull(source, header)
	if n < 3 || string(header[:3]) != "ID3" {
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return nil, err
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if (header[3] != 3 && header[3] != 4) || header[4] != 0 {
		return nil, fmt.Errorf("unsupported ID3 version: 2.%d.%d", header[3], header[4])
	}
	// Extended headers, unsynchronisation and footers require different layouts.
	if header[5] != 0 {
		return nil, fmt.Errorf("unsupported ID3 tag flags: %#x", header[5])
	}
	size, err := decodeSize(header[6:10], true)
	if err != nil {
		return nil, err
	}
	if int64(size)+10 > fileSize {
		return nil, io.ErrUnexpectedEOF
	}
	tag := make([]byte, 10+size)
	copy(tag, header)
	_, err = io.ReadFull(source, tag[10:])
	return tag, err
}

func sanitizeTag(tag []byte) (bool, error) {
	titleStart, titleEnd := 0, 0
	for offset := 10; offset < len(tag); {
		if tag[offset] == 0 {
			for _, b := range tag[offset:] {
				if b != 0 {
					return false, fmt.Errorf("nonzero ID3 padding")
				}
			}
			break
		}
		if len(tag)-offset < 10 {
			return false, io.ErrUnexpectedEOF
		}
		header := tag[offset : offset+10]
		for _, b := range header[:4] {
			if !(b >= 'A' && b <= 'Z') && !(b >= '0' && b <= '9') {
				return false, fmt.Errorf("invalid ID3 frame ID")
			}
		}
		size, err := decodeSize(header[4:8], tag[3] == 4)
		if err != nil {
			return false, err
		}
		if size == 0 || size > len(tag)-offset-10 {
			return false, fmt.Errorf("invalid ID3 frame size")
		}
		end := offset + 10 + size
		if string(header[:4]) == "TIT2" {
			if titleStart != 0 {
				return false, fmt.Errorf("multiple ID3 title frames")
			}
			if header[8] != 0 || header[9] != 0 {
				return false, fmt.Errorf("unsupported ID3 title flags")
			}
			titleStart, titleEnd = offset, end
		}
		offset = end
	}
	if titleStart == 0 {
		return false, nil
	}
	title, changed, err := sanitizeTitle(tag[titleStart+10 : titleEnd])
	if err != nil || !changed {
		return false, err
	}
	oldSize := titleEnd - titleStart - 10
	copy(tag[titleStart+10:], title)
	copy(tag[titleStart+10+len(title):], tag[titleEnd:])
	clear(tag[len(tag)-(oldSize-len(title)):])
	encodeSize(tag[titleStart+4:titleStart+8], len(title), tag[3] == 4)
	return true, nil
}

// sanitizeTitle matches encoded ASCII so Unicode text and its encoding remain
// byte-for-byte intact. Removed bytes become padding at the end of the tag.
func sanitizeTitle(body []byte) ([]byte, bool, error) {
	unit, start := 1, 1
	var order binary.ByteOrder = binary.BigEndian
	switch body[0] {
	case 0, 3: // ISO-8859-1 and UTF-8 share ASCII bytes.
	case 1, 2:
		unit = 2
		if body[0] == 1 {
			if len(body) < 3 {
				return nil, false, fmt.Errorf("missing UTF-16 byte order mark")
			}
			switch {
			case bytes.Equal(body[1:3], []byte{0xff, 0xfe}):
				order = binary.LittleEndian
			case bytes.Equal(body[1:3], []byte{0xfe, 0xff}):
			default:
				return nil, false, fmt.Errorf("invalid UTF-16 byte order mark")
			}
			start = 3
		}
		if (len(body)-start)%2 != 0 {
			return nil, false, fmt.Errorf("odd UTF-16 title length")
		}
	default:
		return nil, false, fmt.Errorf("unsupported title encoding: %d", body[0])
	}
	character := func(offset int) uint16 {
		if unit == 1 {
			return uint16(body[offset])
		}
		return order.Uint16(body[offset : offset+2])
	}
	end := len(body)
	for end > start && character(end-unit) == 0 {
		end -= unit
	}
	suffix := []byte("(Original Mix)")
	if unit == 2 {
		suffix = make([]byte, 2*len("(Original Mix)"))
		for i, c := range "(Original Mix)" {
			order.PutUint16(suffix[2*i:], uint16(c))
		}
	}
	if !bytes.HasSuffix(body[start:end], suffix) {
		return nil, false, nil
	}
	cut := end - len(suffix)
	for cut > start {
		c := character(cut - unit)
		if c != ' ' && c != '\t' && c != '\r' && c != '\n' {
			break
		}
		cut -= unit
	}
	title := append([]byte(nil), body[:cut]...)
	title = append(title, body[end:]...)
	return title, true, nil
}

func decodeSize(data []byte, synchsafe bool) (int, error) {
	if !synchsafe {
		return int(binary.BigEndian.Uint32(data)), nil
	}
	size := 0
	for _, b := range data {
		if b&0x80 != 0 {
			return 0, fmt.Errorf("invalid ID3 synchsafe size")
		}
		size = size<<7 | int(b)
	}
	return size, nil
}

func encodeSize(data []byte, size int, synchsafe bool) {
	if !synchsafe {
		binary.BigEndian.PutUint32(data, uint32(size))
		return
	}
	for i := 3; i >= 0; i-- {
		data[i] = byte(size & 0x7f)
		size >>= 7
	}
}

func replaceFile(path string, source *os.File, info os.FileInfo, tag []byte) error {
	temp, err := os.CreateTemp(filepath.Dir(path), ".id3-sanitizer-*")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	defer temp.Close()
	if _, err := temp.Write(tag); err != nil {
		return err
	}
	if _, err := io.Copy(temp, source); err != nil {
		return err
	}
	if err := temp.Chmod(info.Mode().Perm()); err != nil {
		return err
	}
	if err := temp.Sync(); err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	current, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !os.SameFile(info, current) || info.Size() != current.Size() || !info.ModTime().Equal(current.ModTime()) {
		return fmt.Errorf("mp3 changed while sanitizing: %s", path)
	}
	if err := source.Close(); err != nil {
		return err
	}
	return os.Rename(temp.Name(), path)
}
