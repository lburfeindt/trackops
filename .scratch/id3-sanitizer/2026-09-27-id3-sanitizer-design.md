# ID3 sanitizer design

Add an independent Go CLI following format-playlist. Accept exactly one MP3 path, strip the case-sensitive trailing `(Original Mix)` and preceding ASCII spaces, tabs or line breaks from its ID3v2 title, and save at the same path. No matching title means no write. Keep audio and unrelated metadata. Support ID3v2.3 and v2.4 and common text encodings; reject malformed or unsupported metadata before writing. Untagged files remain unchanged.

Match the encoded ASCII suffix directly, preserving title encoding and raw unrelated frames. Use an ID3 library only as an independent test reader. Write to a temporary sibling, sync and close it, then rename over the source, preserving permissions. Reject symlinks and nonregular inputs to avoid replacing a link unexpectedly. Tests copy the supplied MP3 fixture and exercise Unicode, no-op, corruption and preservation cases.

Register in root Makefile, README, CI module cache and Dependabot. Go version and test style follow format-playlist. User requests complete implementation in one shot: proceed with routine design decisions, implement inline on lb/id3-sanitizer in the supplied checkout, and leave changes for review.
