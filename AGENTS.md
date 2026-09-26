# Instructions for Coding Agents

## General

- Create a Git branch for each task
- Keep changes focused and follow the conventions and lifecycle documented by the affected tool.
- This repository contains independent tools. Check the tool's README before changing it; avoid assuming one language or workflow applies to every tool.
- From the repository root, `make format`, `make lint`, `make test`, `make verify`, and `make build` run the registered tools' corresponding targets. Tool-specific instructions may define additional requirements.
- Keep commit messages and pull request descriptions short (one to three lines).

## Agent workflow

When making code changes or adding new functionality, always follow this workflow:
1. Understand the requirements and context first
2. Plan your changes in accordance with the existing architecture
3. Implement the changes following the established code style and patterns
4. Write or update tests for any new or modified functionality
5. Before considering the task complete, always:
- Format the code: `make format`
- Run the linter: `make lint`
- Run the tests to ensure nothing is broken:
    - Run specific tests related to your changes: `go test ./package/path -run TestSpecificFunction -v`
    - Run all tests to ensure full compatibility: `make test`

When using Superpowers:

1. Brainstorm/design → `.agentic/superpowers/specs/`
2. After design approval → `.agentic/superpowers/plans/`
3. Execute the plan using the Superpowers workflow.
4. Keep implementation artifacts in the locations required by the active
   Superpowers skill.
5. Never put planning artifacts in the repository root.

**IMPORTANT**: Never consider a code change complete until all tests pass successfully and the code has been formatted and linted properly. Always verify your changes with the appropriate tests before finishing.
