# Music Tools Repository Layout Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Organize the existing playlist formatter as the first standalone tool in a lightweight music CLI repository.

**Architecture:** Put each tool in its own directory under `../../../tools`, with its executable and a tool README together. Keep a project overview and lightweight metadata at the repository root; add no runtime, package manager, or build framework.

**Tech Stack:** Bash and standard command-line utilities already used by the script; Markdown; EditorConfig.

**Spec:** `../specs/2026-09-26-music-tools-layout-design.md`

## Global Constraints

- Preserve the existing `format-playlist` behavior and executable entry point while relocating it.
- Add a root project README and a tool README with invocation and requirements.
- Add only the agreed lightweight metadata files.
- Do not introduce Node.js, a shared runtime, build tooling, package management, or a monorepo manager.
- Do not select or create a license without the owner's preference.

## Review Focus

- The move preserves the script's executable bit and source contents.
- The tool README's invocation matches the final repository path and current positional file argument.
- The root README indexes the tool and describes this repository as a small collection of independent music CLIs.
- `../../../.gitignore` contains only relevant generated or local files and does not hide tool sources.
- `../../../.editorconfig` specifies consistent text encoding, line endings, final newlines, and trailing whitespace handling.

---

### Task 1: Move and document the playlist formatter

**Files:**
- Move: `format-playlist` to `tools/format-playlist/format-playlist`
- Create: `../../../tools/format-playlist/README.md`

**Interfaces:**
- Consumes: the current executable Bash script and its single positional playlist-file argument.
- Produces: the same executable at `tools/format-playlist/format-playlist`.

- [x] Create `../../../tools/format-playlist` and move the script there without editing its contents; preserve its executable permission.
- [x] Write the tool README with a short description, the command `./tools/format-playlist/format-playlist path/to/playlist.m3u`, and its Bash plus `cat`, `grep`, and `sed` requirements.
- [x] Review the moved script's contents and permission, and check that the documented invocation points to the moved file.

### Task 2: Add project-level README and metadata

**Files:**
- Create: `../../../README.md`
- Create: `../../../.gitignore`
- Create: `../../../.editorconfig`

**Interfaces:**
- Consumes: the tool path and description from Task 1.
- Produces: project overview and metadata for the repository root.

- [x] Write a root README that explains the project purpose and links to `../../../tools/format-playlist/README.md`.
- [x] Add a minimal `../../../.gitignore` for macOS Finder metadata (`.DS_Store`).
- [x] Add `../../../.editorconfig` with `root = true`, UTF-8 encoding, LF line endings, a final newline, and trailing-whitespace trimming.
- [x] Review the final file tree and README paths against the approved structure.
