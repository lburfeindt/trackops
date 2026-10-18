# Playlist saving

Approved decision for [Decide playlist output file replacement](https://github.com/lburfeindt/trackops/issues/28). This revises the playlist auto-save behavior from the [first window decision](2026-10-18-window-interactions.md).

- Choosing an M3U playlist extracts and displays titles without writing a file. **Copy titles** works immediately.
- **Save titles** writes one title per line beside the playlist. Remove the playlist extension, then append `-formatted.txt`: `Set.m3u8` becomes `Set-formatted.txt`.
- If that output exists, ask **Replace** or **Cancel** before writing. **Replace** writes the new result. **Cancel** leaves the existing file untouched; the new titles stay visible with “Not saved”. A failed save also leaves the titles visible and reports the error. An existing output must remain intact if replacement fails.
- Empty results show “No titles found” and can be saved as an empty file only after **Save titles**; an existing output gets the same replacement prompt.
- A read error offers no Save and writes nothing. Choosing another playlist replaces an unsaved display without warning. Canceling the file chooser keeps the current display.

The [clickable window prototype](window-prototype.html) predates this decision and is retained as a historical sketch, not the current saving contract.
