#!/bin/sh
set -eu

# Given
project_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
test_directory=$(mktemp -d)
trap 'rm -rf "$test_directory"' EXIT
selection_directory="$test_directory/selection"
mkdir "$selection_directory"
fixture="$project_root/tools/id3-sanitizer/testdata/Silence.mp3"
cp "$fixture" "$selection_directory/updated.mp3"
cp "$fixture" "$selection_directory/unchanged.mp3"
cp "$fixture" "$test_directory/expected.mp3"
"$project_root/tools/id3-sanitizer/bin/id3-sanitizer" "$test_directory/expected.mp3" >/dev/null
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
