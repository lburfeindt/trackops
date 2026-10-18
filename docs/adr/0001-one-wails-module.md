# Build TrackOps as one Wails module

TrackOps will use Wails v2 and consolidate the existing Go tools into one module. Headless ID3 and playlist feature packages will contain the processing logic; Wails will handle native selection and the window, and the old CLIs and standalone ID3 app will be removed. This keeps a future CLI possible and leaves Windows and Linux UI support available after the macOS-only first release, without a helper-process protocol or a separate native-picker bridge.
