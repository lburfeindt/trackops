# First TrackOps window

Approved interaction flow for [issue #21](https://github.com/lburfeindt/trackops/issues/21). The [clickable prototype](window-prototype.html) is kept with this planning note for review; it is simulated and does not ship with the app.

The small window has two buttons: **ID3 Sanitizer** and **Format Playlist**.

## ID3 Sanitizer

Choose MP3 files or a folder. Selection starts processing immediately. Show progress and a Stop action. Stop between files; edits already completed remain saved. Show updated, unchanged, and failed counts when processing ends, including after a partial stop or failure. Canceling the chooser returns to the ready state without processing.

## Format Playlist

Choose an M3U playlist. Selection immediately extracts its `#EXTINF` titles, writes them one per line to `<playlist-name>-titles.txt` beside the playlist, and shows the same titles in the window with **Copy titles**. Copy puts the displayed titles on the clipboard, one per line. Canceling a later chooser preserves the currently displayed result.

An empty playlist shows “No titles found” and writes an empty output file. A read error shows an error and a way to choose another playlist; it does not write an output file. If the file cannot be saved, keep the extracted titles visible for Copy and show the save error.

Whether an existing output file is overwritten needs a decision before implementation.
