# First TrackOps window

Approved interaction flow for [issue #21](https://github.com/lburfeindt/trackops/issues/21), with playlist saving revised by [issue #28](https://github.com/lburfeindt/trackops/issues/28). The [clickable prototype](window-prototype.html) records the earlier interaction sketch; it is simulated and does not ship with the app.

The small window has two buttons: **ID3 Sanitizer** and **Format Playlist**.

## ID3 Sanitizer

Choose MP3 files or a folder. Selection starts processing immediately. Show progress and a Stop action. Stop between files; edits already completed remain saved. Show updated, unchanged, and failed counts when processing ends, including after a partial stop or failure. Canceling the chooser returns to the ready state without processing.

## Format Playlist

Choose an M3U playlist. Selection extracts its `#EXTINF` titles and shows them in the window without writing a file. **Copy titles** puts the displayed titles on the clipboard, one per line. **Save titles** writes them one per line to `<playlist-name>-formatted.txt` beside the playlist; for example, `Set.m3u8` saves to `Set-formatted.txt`. If the output already exists, ask before replacing it. Canceling replacement leaves the existing file untouched and the new titles visible with a “Not saved” status.

An empty playlist shows “No titles found”; **Save titles** can create an empty file, with the same replacement prompt if needed. A read error shows an error and a way to choose another playlist; it does not offer Save. If the file cannot be saved, keep the extracted titles visible for Copy and show the save error without changing an existing output file.

Choosing another playlist replaces the displayed result without an unsaved warning. Canceling the chooser preserves the currently displayed result. See the [saving decision](2026-10-18-playlist-saving.md) for the complete contract.
