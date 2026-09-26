# format-playlist

Prints track titles from the `#EXTINF` entries in an M3U playlist.

## Usage

Run from the tool directory with Go:

```sh
cd tools/format-playlist
go run . path/to/playlist.m3u8
```

From the repository root, build and run a standalone binary:

```sh
make build
./bin/format-playlist path/to/playlist.m3u8
```

This tool requires Go 1.24 or later. Its local targets are `make -C tools/format-playlist test`, `check`, and `build`. The root [project README](../../README.md) describes the language-neutral lifecycle commands. The command prints one title per `#EXTINF` entry, preserving commas and Unicode in titles.
