# TrackOps

TrackOps is one macOS app for two music chores:

- **ID3 Sanitizer** removes the exact `(Original Mix)` suffix from MP3 titles. Choose MP3 files or a folder; folder selection includes subfolders. The app shows progress, supports Stop between files, and reports updated, unchanged, and failed counts.
- **Format Playlist** opens an M3U or M3U8 playlist and shows its `#EXTINF` titles in order. Copy them or use **Save titles** to write `<playlist-stem>-formatted.txt` beside the playlist. Existing output requires Replace or Cancel.

The app edits selected MP3s in place. Saving playlist titles happens only when you choose **Save titles**.

## Build and launch

Requires Apple Silicon macOS 13 or later, Go 1.27.1 or later, Make, and the Xcode command line tools. The first build downloads the pinned Wails v2 CLI and Go dependencies.

```sh
make package-macos
open build/bin/TrackOps.app
```

`package-macos` makes a clean Apple Silicon build at `build/bin/TrackOps.app` and signs it ad hoc for local use. To install it, copy the app to `/Applications`, then launch it from Finder or Spotlight. The app is for local use; it is not notarized for distribution.

## Development

The root Go module contains both headless features and the Wails app. Run the lifecycle from the repository root:

```sh
make format
make lint
make test
make verify
make build
```

`make build` makes a development app. `make test-ui` packages the app and runs the macOS UI suite, including native pickers and file assertions. Give the terminal running it Accessibility access. `make clean` removes generated files in `build/`.

Run direct Wails CLI commands from `cmd/trackops`.

## Code layout

- `cmd/trackops` holds the Wails entry point and project configuration. `frontend/assets.go` embeds the static files. `internal/desktop` connects the window and native pickers to the headless features.
- `internal/id3sanitizer` discovers MP3s, runs the batch, and edits each file.
- `internal/playlist` reads M3U titles and saves formatted output.
- `frontend/dist` contains the source HTML, CSS, and JavaScript embedded in the app. `test/ui` contains the macOS UI checks.

See [documentation](docs/README.md) for project guidance and decisions.
