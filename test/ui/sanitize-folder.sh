#!/bin/sh
set -eu

# Given
project_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
test_directory=$(mktemp -d)
trap 'rm -rf "$test_directory"' EXIT
fixture="$project_root/tools/id3-sanitizer/testdata/Silence.mp3"
folder="$test_directory/folder"
empty_folder="$test_directory/empty"
stop_folder="$test_directory/stop"
mkdir -p "$folder/album" "$empty_folder" "$stop_folder" "$test_directory/outside"
cp "$fixture" "$folder/root.mp3"
cp "$fixture" "$folder/album/nested.MP3"
cp "$fixture" "$folder/ignored.wav"
cp "$fixture" "$test_directory/outside/outside.mp3"
ln -s "$test_directory/outside" "$folder/linked-album"
printf 'ID3' > "$folder/album/broken.mp3"
cp "$fixture" "$empty_folder/ignored.wav"
cp "$fixture" "$test_directory/expected.mp3"
"$project_root/tools/id3-sanitizer/bin/id3-sanitizer" "$test_directory/expected.mp3" >/dev/null
cp "$fixture" "$stop_folder/0001.mp3"
cp "$fixture" "$stop_folder/0002.mp3"
cp "$fixture" "$stop_folder/0003.mp3"
dd if=/dev/zero of="$stop_folder/0002.mp3" bs=1048576 seek=1 count=1024 conv=notrunc 2>/dev/null

# When
pkill -x TrackOps || true
sleep 1
open -a "$project_root/build/bin/TrackOps.app"
osascript "$project_root/test/ui/sanitize-folder.applescript" "$folder" "$empty_folder" "$stop_folder" "$test_directory/stop-status.txt"

# Then
cmp "$test_directory/expected.mp3" "$folder/root.mp3"
cmp "$test_directory/expected.mp3" "$folder/album/nested.MP3"
cmp "$fixture" "$folder/ignored.wav"
cmp "$fixture" "$test_directory/outside/outside.mp3"
test "$(cat "$folder/album/broken.mp3")" = ID3
cmp "$fixture" "$empty_folder/ignored.wav"
cmp "$test_directory/expected.mp3" "$stop_folder/0001.mp3"
cmp "$fixture" "$stop_folder/0003.mp3"
fixture_size=$(wc -c < "$fixture" | tr -d ' ')
head -c "$fixture_size" "$stop_folder/0002.mp3" > "$test_directory/second-prefix.mp3"
if cmp -s "$test_directory/expected.mp3" "$test_directory/second-prefix.mp3"; then
	updated=2
elif cmp -s "$fixture" "$test_directory/second-prefix.mp3"; then
	updated=1
else
	printf '%s\n' 'In-progress MP3 had unexpected content after Stop' >&2
	exit 1
fi
test "$(cat "$test_directory/stop-status.txt")" = "Stopped · Updated: $updated · Unchanged: 0 · Failed: 0"
