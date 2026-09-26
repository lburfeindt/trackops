# Instructions for Codex

## General

- Create a Git branch for each task, prefixed with `lb/`.
- Keep changes focused and follow the conventions and lifecycle documented by the affected tool.
- This repository contains independent tools. Check the tool's README before changing it; avoid assuming one language or workflow applies to every tool.
- From the repository root, `make test`, `make check`, and `make build` run the registered tools' corresponding targets. Tool-specific instructions may define additional requirements.
- Keep commit messages and pull request descriptions short (one to three lines).

## Agent workflow

When using Superpowers:

1. Brainstorm/design → `.agentic/superpowers/specs/`
2. After design approval → `.agentic/superpowers/plans/`
3. Execute the plan using the Superpowers workflow.
4. Keep implementation artifacts in the locations required by the active
   Superpowers skill.
5. Never put planning artifacts in the repository root.
