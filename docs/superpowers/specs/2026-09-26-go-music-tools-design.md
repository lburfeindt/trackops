# Go Music Tools Project Design

## Purpose

Evolve this repository into a collection of independent music-related command-line tools. Each tool owns its implementation language, dependencies, tests, and local development commands. The root project provides a small common entry point for running each tool's lifecycle and CI.

## Goals

- Make local development and automated verification consistent and easy to run across tools.
- Keep each tool isolated so future tools can use different languages and runtimes.
- Use Go's standard library for the initial implementation and tests.
- Preserve the playlist formatter's useful behavior with real-world M3U data.

## Project structure

```text
Makefile
.github/workflows/ci.yml
tools/
  format-playlist/
    Makefile
    go.mod
    main.go
    main_test.go
    internal/
      playlist/
        m3u.go
        m3u_test.go
    testdata/
      playlist.m3u8
    README.md
```

The `format-playlist` directory is a self-contained Go module. Its command owns argument handling, file access, output, and process exit behavior. Its internal playlist package parses `#EXTINF` entries and returns titles, keeping parsing testable without launching a process. Tests use a small, hand-authored fixture with synthetic titles and paths, independent of the supplied real-world playlist file.

Each tool directory exposes `test`, `check`, and `build` Make targets for that tool's language. The root Makefile delegates those targets to the registered tool directories, without owning a language-specific module or shared runtime.

## Command behavior

- Read one playlist path argument and print one title per `#EXTINF` entry to standard output.
- Keep commas and Unicode in titles; split an `#EXTINF` entry at its first metadata comma.
- Ignore non-`#EXTINF` lines and malformed `#EXTINF` entries that lack a title separator.
- Return a clear error and nonzero status for missing or extra arguments and unreadable files.
- Retain a concise, direct CLI; no third-party command framework is needed for this interface.

From the tool directory, the command will be runnable with `go run . <playlist>` and as a built binary. The tool README will document both invocations.

## Testing and development lifecycle

Use Go's built-in `testing` package. Parser tests will cover normal entries, commas and Unicode in titles, ignored non-entry lines, and malformed entries. Command tests will cover argument and file errors plus output for the fixture.

The `tools/format-playlist/Makefile` will provide:

- `make test` — run `go test ./...`.
- `make check` — verify gofmt, run `go vet ./...`, and run all tests.
- `make build` — build the `format-playlist` executable into the repository's ignored `bin/` directory.

The root `Makefile` will delegate `test`, `check`, and `build` to registered tool Makefiles. Root CI will provision the Go version from `tools/format-playlist/go.mod` and run the root `make check` and `make build` targets on pushes and pull requests. The root README will state prerequisites for the tools currently in the repository and document the root commands. New tools in other languages can add their own Makefile targets, module or runtime metadata, tests, and CI runtime setup without changing the format-playlist module. No external test framework or runtime dependencies are introduced for format-playlist.

## Migration and scope

Move the Go module, command, parser, tests, and fixture under `tools/format-playlist/`; remove the Bash formatter and its temporary shell-only regression test. Add a tool-owned Makefile and keep the root Makefile as a language-neutral dispatcher. Update the tool and root READMEs to the new invocation and lifecycle. Defer packaging, release automation, versioning, third-party libraries, and additional tools until needed.

The previous music-tools layout design remains the record of the initial repository organization; this design supersedes its Bash implementation constraint while restoring its per-tool isolation goal.
