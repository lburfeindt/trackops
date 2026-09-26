# Independent Music Tools Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Keep each music CLI self-contained in its own language while providing common root commands for testing, checking, building, and CI.

**Architecture:** Each tool owns its module or runtime configuration, sources, tests, fixture, README, and `test`, `check`, and `build` targets. The root Makefile delegates those targets to the registered tool directories; root CI sets up current tool prerequisites and runs the root targets.

**Tech Stack:** Go 1.24+, Go `testing`, Make, GitHub Actions.

**Spec:** `../specs/2026-09-26-go-music-tools-design.md`

## Global Constraints

- Keep format-playlist as an independent Go module at `../../../tools/format-playlist`.
- Use Go's standard library for the formatter and its tests.
- Preserve commas and Unicode by splitting `#EXTINF` at the first metadata comma.
- Ignore non-`#EXTINF` lines and malformed entries without a title separator.
- Return a clear error and nonzero status for missing or extra arguments and unreadable files.
- Do not add a root language-specific module, third-party command/test frameworks, or shared runtime.
- Provide root `make test`, `make check`, and `make build` dispatching to tool-owned Makefiles.

## Review Focus

- The repository root has no Go module, and the Go module is under the formatter tool: verify root and tool `go.mod` placement.
- The tool tests and fixture work from the nested module: run `make -C tools/format-playlist test`.
- Root targets delegate successfully and build the binary in `../../../tools/format-playlist/bin`: run root `make test`, `make check`, and `make build`.
- CI selects the nested Go module version: inspect `go-version-file` against the relocated `go.mod`.
- The formatter behavior remains intact: run the parser and CLI tests for commas, Unicode, CRLF, malformed entries, path spaces, argument errors, and unreadable paths.

---

### Task 1: Scope the Go module to the playlist tool

**Files:**
- Move: `go.mod` to `../../../tools/format-playlist/go.mod`
- Move: `cmd/format-playlist/main.go` and `main_test.go` to `../../../tools/format-playlist`
- Move: `internal/playlist/` to `../../../tools/format-playlist/internal/playlist`
- Move: `testdata/playlist.m3u8` to `../../../tools/format-playlist/testdata/playlist.m3u8`
- Create: `../../../tools/format-playlist/Makefile`
- Delete: root `go.mod`, `cmd/`, `internal/`, and `testdata/`

**Interfaces:**
- Consumes: existing CLI `run(args []string, stdout, stderr io.Writer) int` and parser `playlist.Titles(r io.Reader) ([]string, error)`.
- Produces: self-contained module `github.com/lburfeindt/trackops/tools/format-playlist` with tool-local test, check, and build targets.

- [x] **Step 1: Verify the new module-location assertion fails.** Run `test -f tools/format-playlist/go.mod`; expect failure because the module is currently at the repository root.
- [x] **Step 2: Relocate the module and sources.** Move the current Go files and fixture under `../../../tools/format-playlist`; set the module path to `github.com/lburfeindt/trackops/tools/format-playlist`; update the parser import in `main.go` and the CLI test's fixture path to `testdata/playlist.m3u8`.
- [x] **Step 3: Add tool-local Make targets.** `test` runs `go test ./...`; `check` verifies gofmt, runs `go vet ./...`, then tests; `build` creates `bin/` and builds `bin/format-playlist` inside the tool directory.
- [x] **Step 4: Run `make -C tools/format-playlist test`** and confirm both packages pass from the nested module.
- [x] **Step 5: Verify `test ! -f go.mod` at the repository root** and `test -f tools/format-playlist/go.mod`.
- [x] **Step 6: Commit** as `refactor: scope Go module to playlist tool`.

### Task 2: Add root lifecycle dispatch

**Files:**
- Modify: `Makefile`
- Modify: `../../../.github/workflows/ci.yml`
- Modify: `../../../README.md`
- Modify: `../../../tools/format-playlist/README.md`

**Interfaces:**
- Consumes: `test`, `check`, and `build` targets from `../../../tools/format-playlist/Makefile`.
- Produces: root lifecycle commands that delegate to registered tools; CI configured from the Go module inside the tool.

- [x] **Step 1: Update root Makefile** to register `../../../tools/format-playlist` and delegate root `test`, `check`, and `build` targets to every registered tool's matching target.
- [x] **Step 2: Update CI** to read Go from `../../../tools/format-playlist/go.mod`, then run root `make check` and `make build` on pushes and pull requests.
- [x] **Step 3: Update documentation.** Explain that tools own their language prerequisites and lifecycle. Document root `make test/check/build`; document `go run . <playlist>` from `../../../tools/format-playlist` and root `make build` followed by `./tools/format-playlist/bin/format-playlist <playlist>`.
- [x] **Step 4: Run `make test`, `make check`, `make build`, and `git diff --check` from the repository root.** Confirm tests and build pass and `../../../tools/format-playlist/bin/format-playlist` is ignored by Git.
- [x] **Step 5: Commit** as `chore: delegate root workflow to tools`.
