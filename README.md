# Track Ops

A small collection of independent vibe coded tools for chore tasks when working with music software like Rekordbox.

See [documentation](docs/README.md) for shared terms and lasting project decisions.

## Development

Requires Make, Go 1.27.1 or later, and [Wails v2](https://v2.wails.io/docs/gettingstarted/installation/) for the TrackOps app. The existing tool commands remain available while their features move into the app.

Run the registered tools' lifecycle from the repository root:

```sh
make format
make lint
make test
make verify
make build
```

`make build` creates `build/bin/TrackOps.app` on macOS. `make test-ui` builds it and runs the native playlist picker test. UI automation requires Accessibility access for the terminal running the test.

Run `make clean` to remove build outputs and clear Go's downloaded-module, build, and test-result caches. These caches are shared with other Go projects. Then run `make test` to download pinned dependencies again and rerun tests. Run `make build` to recreate the app and remaining CLI.

## Tools

- [format-playlist](tools/format-playlist/README.md) — prints track names from an M3U playlist's `#EXTINF` entries; the TrackOps app also shows them.
- [id3-sanitizer](tools/id3-sanitizer/README.md) — removes `(Original Mix)` from MP3 title suffixes, with a CLI and native file/folder chooser.
