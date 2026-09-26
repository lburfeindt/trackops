# format-playlist

Prints track titles from the `#EXTINF` entries in an M3U playlist.

## Usage

Run from the repository root with Go:

```sh
go run ./cmd/format-playlist path/to/playlist.m3u8
```

To build and run a standalone binary:

```sh
make build
./bin/format-playlist path/to/playlist.m3u8
```

The command prints one title per `#EXTINF` entry. Commas and Unicode in titles are preserved. See the [project README](../../README.md) for test and check commands.
