# Instructions for Coding Agents

## General

- Create a Git branch for each task
- Keep changes focused and follow the conventions and lifecycle documented by the affected tool.
- TrackOps is one Wails app with headless ID3 and playlist features in the root Go module. Check the root README before changing it.
- From the repository root, `make format`, `make lint`, `make test`, `make verify`, and `make build` cover the app module. `make package-macos` builds the local release; `make test-ui` runs its macOS UI suite.
- Keep commit messages and pull request descriptions short (one to three lines).
- Use [docs/README.md](docs/README.md) to find lasting documentation.

## Agent skills

### Issue tracker

Work tickets live in GitHub Issues. See [issue tracker guidance](docs/agents/issue-tracker.md).

### Triage labels

Use the five default triage labels. See [triage labels](docs/agents/triage-labels.md).

### Domain docs

Use one root `CONTEXT.md` and `docs/adr/`. See [domain guidance](docs/agents/domain.md).

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

For each feature or epic, keep temporary descriptions, design drafts, specs, plans, task instructions, progress notes, and Superpowers artifacts in one tracked `.scratch/<feature-or-epic>/` directory. Update links when moving files; the user clears completed scratch folders manually. Reserve `docs/` for lasting project guidance and validated outcomes: record lasting decisions in `docs/adr/`, and keep a full design in `docs/design/` only when it remains useful after the work is done.

**IMPORTANT**: Never consider a code change complete until all tests pass successfully and the code has been formatted and linted properly. Always verify your changes with the appropriate tests before finishing.
