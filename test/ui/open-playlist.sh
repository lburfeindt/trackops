#!/bin/sh
set -eu

project_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
test_directory=$(mktemp -d)
trap 'rm -rf "$test_directory"' EXIT
playlist_fixture="$project_root/internal/playlist/testdata/playlist.m3u8"
cp "$playlist_fixture" "$test_directory/Set.m3u8"

open -a "$project_root/build/bin/TrackOps.app"
osascript "$project_root/test/ui/open-playlist.applescript" "$test_directory/Set.m3u8"

if ! cmp -s "$playlist_fixture" "$test_directory/Set.m3u8"; then
	printf '%s\n' 'Opening a playlist changed the source file' >&2
	exit 1
fi

if [ "$(find "$test_directory" -type f | wc -l | tr -d ' ')" != 1 ]; then
	printf '%s\n' 'Opening a playlist created another file' >&2
	exit 1
fi
