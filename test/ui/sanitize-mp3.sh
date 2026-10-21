#!/bin/sh
set -eu

# Given
project_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
test_directory=$(mktemp -d)
trap 'rm -rf "$test_directory"' EXIT
selection_directory="$test_directory/selection"
mkdir "$selection_directory"
fixture="$project_root/internal/id3sanitizer/testdata/Silence.mp3"
cp "$fixture" "$selection_directory/updated.mp3"
cp "$fixture" "$selection_directory/unchanged.mp3"
cp "$project_root/internal/id3sanitizer/testdata/Silence-sanitized.mp3" "$test_directory/expected.mp3"
cp "$test_directory/expected.mp3" "$selection_directory/unchanged.mp3"
printf 'ID3' > "$selection_directory/broken.mp3"

# When
pkill -x TrackOps || true
sleep 1
open -a "$project_root/build/bin/TrackOps.app"
osascript "$project_root/test/ui/sanitize-mp3.applescript" "$selection_directory"

# Then
cmp "$test_directory/expected.mp3" "$selection_directory/updated.mp3"
cmp "$test_directory/expected.mp3" "$selection_directory/unchanged.mp3"
if [ "$(cat "$selection_directory/broken.mp3")" != ID3 ]; then
	printf '%s\n' 'Invalid MP3 was modified' >&2
	exit 1
fi
