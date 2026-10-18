# ID3 Sanitizer Implementation Plan

> **For agentic workers:** Use superpowers:executing-plans inline. Track steps below.

**Goal:** Remove the exact Original Mix title suffix safely at the original MP3 path.
**Architecture:** Thin CLI and internal sanitizer, independent Go module following format-playlist.
**Tech Stack:** Go 1.27.1, ID3 decoding library, testify.
**Spec:** .scratch/id3-sanitizer/2026-09-27-id3-sanitizer-design.md

## Global Constraints
- Preserve audio and unrelated metadata; no write on no-op or invalid input.
- Follow existing Go style and Given/When/Then tests.
- Root format, lint, test and build must pass.

## Review Focus
- Unicode titles, trailing null terminators and both ID3 versions.
- Malformed tag/frame lengths must leave source intact.
- Unsupported tag or title flags must fail without loss.
- Audio, artwork, unknown frames and permissions must survive.
- Repeated runs and files with absent titles must be no-ops.

### Task 1: sanitizer and CLI
**Files:** tools/id3-sanitizer/{main.go,main_test.go,go.mod,go.sum,Makefile,README.md,AGENTS.md,testdata/Silence.mp3}, internal/sanitizer/{sanitizer.go,sanitizer_test.go}
**Interface:** sanitizer.File(path string) (bool, error); main.run(args []string, stdout, stderr io.Writer) int.
- [x] Write tests for fixture sanitizing, preservation, no-op, Unicode, corrupt input, file errors and CLI usage.
- [x] Run tests and observe missing implementation failure.
- [x] Implement sanitizer and CLI.
- [x] Run focused tests and fix failures.

### Task 2: repository integration and verification
**Files:** Makefile, README.md, .github/workflows/ci.yml, .github/dependabot.yml.
- [x] Register tool and document usage and supported formats.
- [x] Run make format, make lint, make test and make build.
- [x] Review changes and run built CLI against a disposable copy of Silence.mp3.
