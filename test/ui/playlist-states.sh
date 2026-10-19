#!/bin/sh
set -eu

# Given
project_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
test_directory=$(mktemp -d)
trap 'rm -rf "$test_directory"' EXIT
first_playlist="$test_directory/First.m3u8"
second_playlist="$test_directory/Second.m3u8"
empty_playlist="$test_directory/Empty.m3u8"
invalid_playlist="$test_directory/Invalid.m3u8"
printf '#EXTM3U\n#EXTINF:1,First Artist - Café\nfirst.mp3\n#EXTINF:2,Second Artist, Club Mix\nsecond.mp3\n' > "$first_playlist"
printf '#EXTM3U\n#EXTINF:1,Replacement Title\nreplacement.mp3\n' > "$second_playlist"
printf '#EXTM3U\n' > "$empty_playlist"
awk 'BEGIN { printf "#EXTINF:1,"; for (i = 0; i < 70000; i++) printf "x"; printf "\n" }' > "$invalid_playlist"

# When
pkill -x TrackOps || true
sleep 1
open -a "$project_root/build/bin/TrackOps.app"
osascript "$project_root/test/ui/playlist-states.applescript" "$first_playlist" "$second_playlist" "$empty_playlist" "$invalid_playlist"

# Then
expected_clipboard=$(printf 'First Artist - Café\nSecond Artist, Club Mix')
if [ "$(pbpaste)" != "$expected_clipboard" ]; then
	printf '%s\n' 'Clipboard titles were not in display order' >&2
	exit 1
fi

if [ "$(find "$test_directory" -type f | wc -l | tr -d ' ')" != 4 ]; then
	printf '%s\n' 'Opening or copying playlist titles created a file' >&2
	exit 1
fi
