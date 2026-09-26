# Go Music Tools Project Design

## Purpose

Evolve this repository from a home for standalone shell scripts into a small Go project for music-related command-line tools. The first Go command will replace the Bash playlist formatter and provide a repeatable test, check, build, and CI workflow.

## Goals

- Make local development and automated verification consistent and easy to run.
- Keep tools small, independently understandable, and suitable for adding future commands.
- Use Go's standard library for the initial implementation and tests.
- Preserve the playlist formatter's useful behavior with real-world M3U data.

## Project structure

```text
cmd/
  format-playlist/
    main.go
    main_test.go
internal/
  playlist/
    m3u.go
    m3u_test.go
testdata/
  frequency_sessions.m3u8
Makefile
go.mod
```

The command package owns argument handling, file access, output, and process exit behavior. The internal playlist package parses `#EXTINF` entries and returns titles, keeping parsing testable without launching a process. A compact, sanitized fixture based on the supplied real-world playlist covers representative metadata, punctuation, commas in titles, Unicode, and filenames with spaces without storing the user's absolute media paths.

## Command behavior

- Read one playlist path argument and print one title per `#EXTINF` entry to standard output.
- Keep commas and Unicode in titles; split an `#EXTINF` entry at its first metadata comma.
- Ignore non-`#EXTINF` lines and malformed `#EXTINF` entries that lack a title separator.
- Return a clear error and nonzero status for missing or extra arguments and unreadable files.
- Retain a concise, direct CLI; no third-party command framework is needed for this interface.

The command will be runnable from source with `go run ./cmd/format-playlist <playlist>` and as a built binary. The README will document both.

## Testing and development lifecycle

Use Go's built-in `testing` package. Parser tests will cover normal entries, commas and Unicode in titles, ignored non-entry lines, and malformed entries. Command tests will cover argument and file errors plus output for the fixture.

The Makefile will provide:

- `make test` — run `go test ./...`.
- `make check` — verify gofmt, run `go vet ./...`, and run all tests.
- `make build` — build the `format-playlist` executable into an ignored `bin/` directory.

GitHub Actions CI will use the repository's documented Go version and run `make check` and `make build` for pushes and pull requests. The README will state the Go prerequisite and the local commands. No external test framework or runtime dependencies are introduced.

## Migration and scope

Replace the Bash formatter with the Go command and remove its temporary shell-only regression test. Update the tool and root READMEs to the new invocation and lifecycle. Keep the existing repository organization, and defer packaging, release automation, versioning, third-party libraries, and additional tools until needed.

The previous music-tools layout design remains as the record of the initial repository organization; this design supersedes its Bash implementation constraint for the playlist formatter.
