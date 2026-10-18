# First macOS release contract

Approved planning decision for [Decide the first macOS release contract](https://github.com/lburfeindt/trackops/issues/22).

## Bundle and build

- Build one `TrackOps.app` for Apple Silicon (`darwin/arm64`) and macOS 13 or later. Bundle ID: `com.lburfeindt.trackops`. Initial version: `0.1.0`.
- Root `make package-macos` is the build entry point. It runs a clean Wails v2 macOS ARM64 build and places the app at `build/bin/TrackOps.app`.
- Ad hoc sign the finished bundle for local use. The build needs the macOS developer tools and the app's Go and frontend build dependencies.

## Local use

Build once, drag `TrackOps.app` to Applications, then launch it from Finder or Spotlight. `open build/bin/TrackOps.app` is also useful when testing the build in place.

The first release is for local use by the developer. No DMG, installer, Developer ID signature, notarization, or web download flow is included. Those steps would need a separate decision before distributing the app to others.

## Implementation references

- [Wails v2 build CLI](https://wails.io/docs/reference/cli/) documents `wails build`, `-platform darwin/arm64`, `-clean`, and the app bundle output.
- [Wails v2 platform requirements](https://wails.io/docs/gettingstarted/installation/) list supported macOS architectures and versions.
- [Wails signing guide](https://wails.io/docs/guides/signing/) and [Apple Developer ID](https://developer.apple.com/developer-id/) cover distribution signing and notarization, which are outside this local release.

This ticket settles packaging scope. App implementation remains for the implementation phase of the [TrackOps app map](https://github.com/lburfeindt/trackops/issues/18).
