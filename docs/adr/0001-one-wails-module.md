# Build TrackOps as one Wails module

TrackOps will use Wails v2 and consolidate the existing Go tools into one module. Headless ID3 and playlist feature packages will contain the processing logic; Wails will handle the window and native selection, and the old CLIs and standalone ID3 app will be removed. This keeps a future CLI possible and leaves Windows and Linux UI support available after the macOS-only first release, without a helper-process protocol.

The macOS MP3 picker uses a small AppKit `NSOpenPanel` bridge so one dialog can select files and folders. Wails v2 exposes those as separate dialogs; other pickers remain on Wails.
