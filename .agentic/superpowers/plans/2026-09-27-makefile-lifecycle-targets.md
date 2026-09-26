# Makefile Lifecycle Targets Implementation Plan

**Goal:** Add root and tool-level `format`, `lint`, and `verify` targets, replacing `check` with `verify`.

**Architecture:** The root Makefile dispatches lifecycle targets to registered tools. The Go tool owns the commands that format, lint, and verify its sources.

**Tech Stack:** GNU Make, Go 1.27.1, gofmt, go vet, go test.

**Spec:** `.agentic/superpowers/specs/2026-09-27-makefile-lifecycle-targets.md`

## Global Constraints

- Preserve the registered-tool dispatch pattern in the root Makefile.
- `verify` checks formatting, lint, and tests without modifying source files.
- `format` applies gofmt; `lint` runs go vet.

## Review Focus

- Root targets dispatch the requested lifecycle command to each registered tool.
- `verify` fails if Go files are unformatted and runs both lint and tests.
- Documentation contains no stale `make check` lifecycle command.

---

### Task 1: Add lifecycle targets and update references

**Files:**
- Modify: `Makefile`
- Modify: `tools/format-playlist/Makefile`
- Modify: `README.md`
- Modify: `tools/format-playlist/README.md`
- Modify: `AGENTS.md`

**Interfaces:**
- Consumes: Existing root `TOOLS` list and tool-local Go commands.
- Produces: Root and tool-local `format`, `lint`, and `verify` Make targets.

- [x] Add tool targets: `format` runs gofmt, `lint` runs go vet, and `verify` checks formatting then runs lint and tests.
- [x] Add matching root dispatch targets and remove `check`.
- [x] Replace lifecycle references to `make check` with `make verify`.
- [x] Run `make format`, `make lint`, `make verify`, and `make test`.
