# Instructions for Coding Agents

## General Code Style

- Follow standard Go style and idiomatic patterns
- Use camelCase for variable/function names
- Use PascalCase for exported types, functions, and methods
- Keep functions focused on a single responsibility
- Use descriptive variable and function names that indicate purpose
- Use comments for package documentation and complex logic
- Package names should be lowercase, concise, and reflect functionality
- Use interfaces for dependency injection and testability

### Code Organization

- Organize code by domain concepts
- Keep implementation details private within packages
- Export only what's necessary from packages
- Use interfaces to define contracts between components
- Place common helpers and utilities in dedicated packages
- Use `internal` directory for code that shouldn't be imported outside the project

### Test Style and Naming

- Test files are placed next to the code they test with `_test.go` suffix
- Test functions use the naming pattern `Test<FunctionName>` or `Test<FunctionName>_<Scenario>`
- Use table-driven tests for similar test cases with different inputs
- Mark each test phase with the bare comments `// Given`, `// When`, and `// Then`, in that order
- Put a blank line between each phase marker and the preceding phase; keep setup, action, and assertions under their matching markers
- Keep phase markers as bare labels without explanatory text; a table-driven subtest may have an empty Given phase when the case table supplies its setup
- Use scenario-specific names for test data, such as `playlistFixture` or `invalidPlaylistPath`, rather than generic names
- Use subtests (t.Run) for organizing related test cases
- Use descriptive test names that explain the scenario being tested
- Use mock implementations for interface dependencies
- Test both success and error paths
- Use assertions from the testify package for cleaner test code (assert, require)
- Mock external dependencies to isolate unit tests
