# format-playlist

Prints the track names from the `#EXTINF` entries in an M3U playlist.

## Usage

Run it from the repository root and pass the playlist file as the first argument:

```sh
./tools/format-playlist/format-playlist path/to/playlist.m3u
```

Requires Bash and the standard `cat`, `grep`, and `sed` command-line utilities.
