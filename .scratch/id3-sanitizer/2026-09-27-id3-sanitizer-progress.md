# SDD ledger — plan: .scratch/id3-sanitizer/2026-09-27-id3-sanitizer.md
Ruling: User requests one-shot implementation; execute inline in existing checkout on new lb branch, with no approval pauses or automatic commits.
Pre-flight: CLI consumes sanitizer.File(path string) (bool, error); consistent.
Baseline: make test passed with network-enabled dependency download.
Ruling: Preserve raw frames and remove encoded ASCII suffix directly, retaining title encoding and all other bytes. ID3 library used only as independent test reader; avoids lossy whole-tag reserialization.
Task 1: complete. Stub tests failed on missing state changes, implemented sanitizer and CLI, focused tests passed.
Final review: independent reviewer found no significant issues. Added coverage for big-endian BOM, duplicate titles, save failure cleanup and malformed UTF-16.
Task 2: complete. Root make format, make lint, make test, make build and final make verify passed. Built binary updated a sample copy, preserved audio bytes, and was unchanged on repeat. Branch retained uncommitted for user review.
