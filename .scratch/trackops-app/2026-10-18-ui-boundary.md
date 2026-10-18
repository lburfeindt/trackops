# TrackOps UI and integration boundary

Approved planning note for [issue #20](https://github.com/lburfeindt/trackops/issues/20).

## Established context

- The first app is macOS only: a small window with ID3 Sanitizer and Format Playlist buttons. ID3 accepts files or a folder; playlist titles appear in the app with Copy and are written to a sibling text file. See [map #18](https://github.com/lburfeindt/trackops/issues/18) and the [approved window flow](2026-10-18-window-interactions.md).
- [UI research #19](https://github.com/lburfeindt/trackops/issues/19) found SwiftUI, Wails v2, and Fyne v2 viable. SwiftUI and Wails offer native pickers directly; Fyne needs a bridge to preserve one.
- Both tools are separate Go modules. Reusable code lives in module-local `internal/` packages. The ID3 GUI currently combines Zenity dialogs with batch processing; the playlist CLI prints results from `playlist.Titles`.

## Decisions made

- Use Wails v2 for the unified app. Wails supports macOS and Windows 10/11, with native file and folder dialogs and Windows builds. Windows is not in the first release scope; verify its runtime separately when added. [Platform support](https://v2.wails.io/docs/gettingstarted/installation/), [dialogs](https://v2.wails.io/docs/reference/runtime/dialog/), [build CLI](https://v2.wails.io/docs/reference/cli/).
- Consolidate the two existing Go modules into one TrackOps module; remove old CLI and standalone GUI entry points when the unified app replaces them.
- Wails owns pickers, progress, and result display. Headless Go processing accepts inputs and returns structured results, so a future CLI can call the same operations.
- Separate the ID3 and playlist features internally, without cross-dependencies or extra DDD layers until needed.
- The first release targets macOS only. Windows and Linux are later targets; keep the core and UI portable without taking on their packaging now.

## Proposed code boundary

```text
TrackOps Go module
├── Wails adapter: native pickers, calls, UI-facing progress and results
├── internal/id3: file sanitization and cancellable batch workflow
└── internal/playlist: playlist title extraction
```

The Wails adapter passes selected paths to the feature packages and receives structured results. ID3 processing retains its current file and batch behavior; playlist processing retains its current title extraction behavior. The app writes the extracted titles to a sibling text file only when the user chooses **Save titles**, as decided in [issue #28](https://github.com/lburfeindt/trackops/issues/28). A later CLI can call the same packages through another entry point. Existing CLI and Zenity entry points and their standalone packaging leave the build when the unified app replaces them. The root lifecycle targets will run the one app module.

Window interaction details belong to [issue #21](https://github.com/lburfeindt/trackops/issues/21); macOS packaging details belong to [issue #22](https://github.com/lburfeindt/trackops/issues/22). This ticket does not change processing behavior.

Lasting decision: [ADR 0001](../../docs/adr/0001-one-wails-module.md).

## Windows follow-up

- Wails uses WebView2 on Windows; installation and frontend behavior need Windows runtime verification. [Windows guide](https://v2.wails.io/docs/guides/windows/).
- The sanitizer replaces files with `os.Rename`. Go documents that replacement is not atomic on non-Unix platforms, so validate this path on Windows when that release is planned. [Go Rename](https://pkg.go.dev/os#Rename).

## Linux follow-up

- Wails v2 supports Linux AMD64/ARM64 and native file/folder dialogs. [Installation](https://v2.wails.io/docs/gettingstarted/installation/), [dialogs](https://v2.wails.io/docs/reference/runtime/dialog/).
- Linux distribution needs a Linux build and package metadata for GTK3 and WebKit2GTK runtime dependencies. Package names and WebKit ABI versions differ among distributions. [Linux distro guide](https://v2.wails.io/docs/guides/linux-distro-support/).
- Linux remains outside the first release, as in [map #18](https://github.com/lburfeindt/trackops/issues/18).
