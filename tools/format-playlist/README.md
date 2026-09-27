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
./tools/format-playlist/bin/format-playlist path/to/playlist.m3u8
```

This tool requires Go 1.27.1 or later. Its local targets are `format`, `lint`, `test`, `verify`, `build`, and `clear`, run as `make -C tools/format-playlist <target>`. The root [project README](../../README.md) describes the language-neutral lifecycle commands. The command prints one title per `#EXTINF` entry, preserving commas and Unicode in titles.

`make clear` removes the tool's `bin/` directory and Go's shared downloaded-module, build, and test-result caches. The next `make test` downloads dependencies again and rebuilds before running tests. Run `make build` to recreate the binary.
