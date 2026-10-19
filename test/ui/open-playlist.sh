#!/bin/sh
set -eu

# Given
project_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
test_directory=$(mktemp -d)
trap 'rm -rf "$test_directory"' EXIT
playlist_fixture="$project_root/internal/playlist/testdata/playlist.m3u8"
playlist_name="Set-$$.m3u8"
playlist_path="$test_directory/$playlist_name"
cp "$playlist_fixture" "$playlist_path"

# When
pkill -x TrackOps || true
sleep 1
open -a "$project_root/build/bin/TrackOps.app"
osascript "$project_root/test/ui/open-playlist.applescript" "$playlist_path" "$playlist_name"

# Then
if ! cmp -s "$playlist_fixture" "$playlist_path"; then
	printf '%s\n' 'Opening a playlist changed the source file' >&2
	exit 1
fi

if [ "$(find "$test_directory" -type f | wc -l | tr -d ' ')" != 1 ]; then
	printf '%s\n' 'Opening a playlist created another file' >&2
	exit 1
fi
