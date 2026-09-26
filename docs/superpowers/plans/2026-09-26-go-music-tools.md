# Go Music Tools Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [x]`) syntax.

**Goal:** Replace the Bash playlist formatter with a tested Go CLI and establish a repeatable local and CI workflow.

**Architecture:** `cmd/format-playlist` handles arguments and files, while `internal/playlist` parses M3U metadata. Go's standard library supplies tests and build tooling; Make targets and GitHub Actions expose the lifecycle.

**Tech Stack:** Go 1.24+, Go `testing`, Make, GitHub Actions.

**Spec:** `docs/superpowers/specs/2026-09-26-go-music-tools-design.md`

## Global Constraints

- Use Go's standard library for the initial implementation and tests.
- Print one title per `#EXTINF` entry; preserve commas and Unicode by splitting at the first metadata comma.
- Ignore non-`#EXTINF` lines and malformed entries without a title separator.
- Return a clear error and nonzero status for missing or extra arguments and unreadable files.
- Do not add third-party command, test, or runtime dependencies.
- Provide `make test`, `make check`, and `make build`; CI runs checks and builds on pushes and pull requests.

## Review Focus

- A title with commas and Unicode remains intact: parser fixture test asserts exact titles.
- CRLF input does not leave carriage returns in output: parser test asserts exact title.
- Malformed `#EXTINF` entries without a comma are ignored: parser test asserts no title is emitted.
- A playlist path containing spaces works: CLI test invokes `run` with a temporary path containing spaces.
- Missing/extra arguments and unreadable paths fail with a useful stderr message and nonzero result: CLI tests cover each case.

---

### Task 1: Parse M3U entries

**Files:**
- Create: `go.mod`
- Create: `internal/playlist/m3u.go`
- Create: `internal/playlist/m3u_test.go`
- Create: `testdata/frequency_sessions.m3u8`

**Interfaces:**
- Produces: `playlist.Titles(r io.Reader) ([]string, error)`; returns titles from `#EXTINF` lines, splitting at the first comma after the prefix.

- [x] **Step 1: Initialize the module and add a compact fixture.** Set `go.mod` to `module github.com/lburfeindt/trackops` and `go 1.24`. Add `testdata/frequency_sessions.m3u8` using representative entries from the supplied playlist, including comma-separated artists, Unicode, punctuation, and titles with commas. Replace all absolute media paths with harmless placeholders.
- [x] **Step 2: Write failing parser tests** named `TestTitlesFromFrequencyFixture`, `TestTitlesPreservesCommasAndUnicode`, `TestTitlesIgnoresMalformedAndNonEntryLines`, and `TestTitlesHandlesCRLF`. Assert exact output slices, including an empty result for a playlist with no valid entries.
- [x] **Step 3: Run `go test ./internal/playlist`** and confirm it fails because `playlist.Titles` is not implemented.
- [x] **Step 4: Implement `Titles(r io.Reader) ([]string, error)`** in `internal/playlist/m3u.go`. Scan line by line, recognize the exact `#EXTINF:` prefix, find the first comma, append the remainder as the title, and return scanner errors.
- [x] **Step 5: Run `go test ./internal/playlist`** and confirm all parser tests pass.
- [x] **Step 6: Commit** the parser and fixture as `feat: add M3U playlist parser`.

### Task 2: Add the CLI command

**Files:**
- Create: `cmd/format-playlist/main.go`
- Create: `cmd/format-playlist/main_test.go`

**Interfaces:**
- Consumes: `playlist.Titles(r io.Reader) ([]string, error)` from Task 1.
- Produces: `run(args []string, stdout, stderr io.Writer) int`; returns zero on success and one after writing a useful error to stderr on failure. `main` passes process arguments and standard streams to `run` and exits with its result.

- [x] **Step 1: Write failing command tests** named `TestRunPrintsTitlesFromPathWithSpaces`, `TestRunRejectsMissingOrExtraArguments`, and `TestRunReportsUnreadablePlaylist`. The success case uses the fixture copied to a temporary playlist path containing spaces and asserts exact stdout and empty stderr; failure cases assert nonzero status and a useful stderr message.
- [x] **Step 2: Run `go test ./cmd/format-playlist`** and confirm it fails because the command's `run` function is not implemented.
- [x] **Step 3: Implement the command** in `cmd/format-playlist/main.go`. Require exactly one argument; open and close the file; pass it to `playlist.Titles`; print each title on its own line; report usage or file/parser errors to stderr.
- [x] **Step 4: Run `go test ./cmd/format-playlist` and `go test ./...`** and confirm both pass.
- [x] **Step 5: Commit** the CLI and command tests as `feat: add format-playlist Go command`.

### Task 3: Establish the project workflow

**Files:**
- Create: `Makefile`
- Create: `.github/workflows/ci.yml`
- Modify: `.gitignore`
- Modify: `README.md`
- Modify: `tools/format-playlist/README.md`
- Delete: `tools/format-playlist/format-playlist`
- Delete: `tools/format-playlist/test.sh`

**Interfaces:**
- Consumes: the Go command and module from Tasks 1 and 2.
- Produces: `make test`, `make check`, and `make build`; documented source and binary invocations; CI that runs checks and build.

- [x] **Step 1: Define the Make targets.** `test` runs `go test ./...`; `check` fails if `gofmt -l` reports any Go files and then runs `go vet ./...` and `go test ./...`; `build` creates `bin/` and builds `bin/format-playlist` with `go build -o bin/format-playlist ./cmd/format-playlist`.
- [x] **Step 2: Configure CI** in `.github/workflows/ci.yml` for pushes and pull requests. Set up the Go version from `go.mod`, then run `make check` and `make build`.
- [x] **Step 3: Update documentation and ignore rules.** Document the Go prerequisite, `make test`, `make check`, `make build`, `go run ./cmd/format-playlist <playlist>`, and the built binary invocation. Add `/bin/` to `.gitignore`.
- [x] **Step 4: Remove the superseded Bash executable and temporary shell test.**
- [x] **Step 5: Run `make check`, `make build`, and `git diff --check`** and confirm all succeed, the binary exists at `bin/format-playlist`, and no generated binary is tracked.
- [x] **Step 6: Commit** the project workflow and migration as `chore: establish Go project workflow`.
