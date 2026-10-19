#!/bin/sh
set -eu

# Given
project_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
test_directory=$(mktemp -d)
trap 'chmod 700 "$test_directory/failure"; rm -rf "$test_directory"' EXIT
mkdir "$test_directory/failure"
new_playlist="$test_directory/New.m3u8"
empty_playlist="$test_directory/Empty.m3u8"
replace_playlist="$test_directory/Replace.m3u8"
cancel_playlist="$test_directory/Cancel.m3u8"
failure_playlist="$test_directory/failure/Failure.m3u8"
printf '#EXTM3U\n#EXTINF:1,Artist - Café\na.mp3\n#EXTINF:2,Second Artist, Club Mix\nb.mp3\n' > "$new_playlist"
printf '#EXTM3U\n' > "$empty_playlist"
printf '#EXTM3U\n#EXTINF:1,Replacement Title\na.mp3\n' > "$replace_playlist"
printf '#EXTM3U\n#EXTINF:1,Unsaved Title\na.mp3\n' > "$cancel_playlist"
printf '#EXTM3U\n#EXTINF:1,Still Visible\na.mp3\n' > "$failure_playlist"
printf 'prior\n' > "$test_directory/Replace-formatted.txt"
printf 'prior\n' > "$test_directory/Cancel-formatted.txt"
printf 'prior\n' > "$test_directory/failure/Failure-formatted.txt"
chmod 500 "$test_directory/failure"

# When
pkill -x TrackOps || true
sleep 1
open -a "$project_root/build/bin/TrackOps.app"
osascript "$project_root/test/ui/save-playlist.applescript" "$new_playlist" "$empty_playlist" "$replace_playlist" "$cancel_playlist" "$failure_playlist" "$test_directory/New-formatted.txt" "$test_directory/Empty-formatted.txt"

# Then
printf 'Artist - Café\nSecond Artist, Club Mix\n' > "$test_directory/expected-new.txt"
cmp "$test_directory/expected-new.txt" "$test_directory/New-formatted.txt"
test -f "$test_directory/Empty-formatted.txt"
test ! -s "$test_directory/Empty-formatted.txt"
test "$(cat "$test_directory/Replace-formatted.txt")" = 'Replacement Title'
test "$(cat "$test_directory/Cancel-formatted.txt")" = 'prior'
test "$(cat "$test_directory/failure/Failure-formatted.txt")" = 'prior'
test "$(find "$test_directory" -name '.trackops-titles-*' | wc -l | tr -d ' ')" = 0
