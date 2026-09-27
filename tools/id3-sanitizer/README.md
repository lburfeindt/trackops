# id3-sanitizer

Removes the exact, case-sensitive `(Original Mix)` suffix and preceding spaces, tabs or line breaks from an MP3's ID3 title, then saves to the same path.

## Native app

The GUI uses [Zenity's Go library](https://github.com/ncruces/zenity) for native dialogs on macOS and Windows. No Go installation is needed to run the built app.

Choose **Files…** to select one or several MP3s, or **Folder…** to process a directory and its subfolders. File selection filters for `.mp3`; batch processing accepts the extension case-insensitively. Choosing inputs starts sanitization immediately and saves to the original paths, without backups.

A progress dialog stays open while processing. Canceling it stops between files; changes already completed remain saved. The result shows updated, unchanged and failed counts, with details for up to ten errors. A failed file does not stop the remaining files. Symlinks inside folders are skipped; directly selected symlinks are rejected. Overlapping selections process each identical absolute path once.

Build a double-clickable Apple Silicon app for macOS 13 or later:

```sh
make -C tools/id3-sanitizer package-macos
open "tools/id3-sanitizer/bin/ID3 Sanitizer.app"
```

Packaging requires the macOS command-line developer tools for ad-hoc signing. The bundle is intended for local use; signing with a Developer ID and notarization for distribution are not included. Native dialogs use macOS's built-in scripting tools.

For development, run `go run ./cmd/id3-sanitizer-gui` from the tool directory, or run `bin/id3-sanitizer-gui` after `make build`. The existing CLI is also built. Windows-compatible dialog code is included, but a Windows build, packaging and runtime verification are deferred. Linux GUI use requires `zenity`, `matedialog` or `qarma`; automated tests do not open dialogs.

## Usage

Run from the tool directory with Go:

```sh
cd tools/id3-sanitizer
go run . path/to/track.mp3
```

From the repository root, build and run a standalone binary:

```sh
make build
./tools/id3-sanitizer/bin/id3-sanitizer "path/to/track.mp3"
```

For example, `What a tune (Original Mix)` becomes `What a tune`. The command prints `updated: <path>` or `unchanged: <path>`. Usage and file errors go to stderr with a nonzero exit code.

## File handling

- Supports ID3v2.3 and ID3v2.4 titles encoded as ISO-8859-1, UTF-8, UTF-16 with BOM, or UTF-16BE.
- Matching text elsewhere in the title, titles with different capitalization, and files without an ID3v2 title are unchanged. ID3v1 tags are left untouched.
- Preserves audio, unrelated tag frames (including artwork), and file permission bits. Removed title bytes become tag padding, retaining audio offsets.
- Writes a temporary file in the same directory, syncs it, and replaces the original. The containing directory must be writable. No backup is kept; files with no matching suffix are never rewritten.
- Rejects malformed tags, unsupported ID3 versions, tag flags (including unsynchronisation, extended headers and footers), flagged or duplicate title frames, symlinks and nonregular files. These errors leave the original file unchanged.
- Replacement creates a new file at the same path; timestamps, ownership, extended attributes and hard-link relationships are not preserved. Avoid editing the file concurrently.

## Development

Requires Go 1.27.1 or later. Local targets are `format`, `lint`, `test`, `verify`, `build`, and `clean`, run as `make -C tools/id3-sanitizer <target>`. See the root [project README](../../README.md) for shared lifecycle commands.

`make clean` removes the tool's `bin/` directory and Go's shared downloaded-module, build, and test-result caches. The next `make test` downloads dependencies again and rebuilds before running tests. Run `make build` to recreate the binary.

`testdata/Silence.mp3` is the supplied sample. Tests sanitize temporary copies and check the resulting title, metadata, audio and repeat-run behavior. The ID3 library is used only by tests as an independent tag reader; the command preserves raw tag data.
