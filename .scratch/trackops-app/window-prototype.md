# First window prototype (#21)

Open [window-prototype.html](window-prototype.html) directly in a browser. No build or server needed.

Run the small interaction test with `node --test .scratch/trackops-app/window-prototype.test.js` from the repository root.

Question: Does one small window make the two tool flows clear, including cancel, empty playlists, and errors?

The playlist auto-save behavior below was superseded by the [issue #28 saving decision](2026-10-18-playlist-saving.md). The prototype remains a record of the earlier window review.

This is throwaway, simulated UI for human review. No file picker, processing, or file write is connected. The two top buttons open the tool flows. The lower controls simulate chooser outcomes, and walkthroughs cover the main paths. Copy writes the displayed titles, one per line, when browser clipboard access permits.

Assumptions to review:

- Cancel returns to the ready state for the first selection. Canceling a later playlist choice restores the prior result.
- ID3 processing begins when inputs are chosen. Its result shows counts; a partial failure states that earlier edits remain saved.
- Stopping ID3 processing happens between files; completed edits remain saved.
- A playlist with no `#EXTINF` titles shows an empty state and a retry button.
- A playlist read error shows a retry button. A valid playlist keeps titles visible until another selection.
- A valid playlist immediately writes the shown titles, one per line, to `<playlist-name>-titles.txt` beside the playlist. An empty playlist writes an empty text file. A write failure leaves the titles visible for Copy.
