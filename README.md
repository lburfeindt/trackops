# Track Ops

A small collection of independent command-line tools for chore tasks when working with music software like Rekordbox.

## Development

Requires Make and the language runtimes listed in each tool's README. The current Go tools require Go 1.27.1 or later.

Run the registered tools' lifecycle from the repository root:

```sh
make format
make lint
make test
make verify
make build
```

## Tools

- [format-playlist](tools/format-playlist/README.md) — prints track names from an M3U playlist's `#EXTINF` entries.
- [id3-sanitizer](tools/id3-sanitizer/README.md) — removes `(Original Mix)` from MP3 title suffixes and saves to the same file.
