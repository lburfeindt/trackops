# Music Tools Repository Layout

## Purpose

Create a small home for standalone command-line utilities related to music. The repository starts with the existing Bash script `format-playlist`, which extracts playlist track names from M3U `#EXTINF` lines. Additional tools may be added later.

## Design

Each tool lives in its own directory under `tools/`. A tool directory contains its executable and a README with its purpose, usage, and requirements. This keeps tools independently understandable and leaves room for a tool to grow without imposing a common runtime or framework.

The repository root contains a README that introduces the project and indexes the tools, plus lightweight project metadata files such as `.gitignore` and `.editorconfig`. A `LICENSE` file is deferred until the owner selects a license. The project does not add a package manifest, build system, or monorepo manager at this stage.

## Initial structure

```text
README.md
.gitignore
.editorconfig
tools/
  format-playlist/
    README.md
    format-playlist
```

The existing script remains an executable Bash CLI and is moved to `tools/format-playlist/format-playlist` as part of implementation. Future tools get their own directories under `tools/`.

## Scope and acceptance

- Preserve the existing `format-playlist` behavior and executable entry point while relocating it.
- Add a root project README and a tool README with invocation and requirements.
- Add only the agreed lightweight metadata files.
- Do not introduce Node.js, a shared runtime, build tooling, package management, or a monorepo manager.
- Do not select or create a license without the owner's preference.

## Repository state

The working directory was not a Git repository when this design was prepared, so this design note could not be committed. Git initialization and the structure implementation remain for the next stage after the written design is reviewed.
