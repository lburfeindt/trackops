# Makefile Lifecycle Targets

## Goal

Expose consistent `format`, `lint`, and `verify` commands from the repository root and the Go tool directory.

## Approved Design

- `format` runs `gofmt -w` on the tool's Go files.
- `lint` runs `go vet ./...`.
- `verify` replaces `check` and checks formatting, runs lint, and runs tests.
- Root targets dispatch to each registered tool.
- README and AGENTS.md lifecycle references use `verify` instead of `check`.

## Scope

Update the root and `tools/format-playlist` Makefiles and the lifecycle references in `README.md`, `tools/format-playlist/README.md`, and `AGENTS.md`.
