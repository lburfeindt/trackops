# Instructions for Coding Agents

## General

- Create a Git branch for each task
- Keep changes focused and follow the conventions and lifecycle documented by the affected tool.
- This repository contains independent tools. Check the tool's README before changing it; avoid assuming one language or workflow applies to every tool.
- From the repository root, `make format`, `make lint`, `make test`, `make verify`, and `make build` run the registered tools' corresponding targets. Tool-specific instructions may define additional requirements.
- Keep commit messages and pull request descriptions short (one to three lines).
- Read [CONTEXT.md](CONTEXT.md) and relevant ADRs in `docs/adr/` before changing domain behavior. Add terms to the shared glossary when they are resolved; record lasting architectural decisions as concise ADRs.
- Use [docs/README.md](docs/README.md) to find lasting documentation.

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

For agent planning, use one tracked `.scratch/<feature-or-epic>/` directory per feature or epic. Keep specs, plans, and progress notes there, including Superpowers artifacts. Update paths in these files when moving them. The user clears completed scratch folders manually. Move lasting decisions into `docs/adr/`; keep a full design in `docs/design/` only when it remains useful.

**IMPORTANT**: Never consider a code change complete until all tests pass successfully and the code has been formatted and linted properly. Always verify your changes with the appropriate tests before finishing.
