# Issue tracker: GitHub

Issues and work specs live in this repository's GitHub Issues. Run `gh` in this checkout to select its remote.

Run networked `gh` calls as standalone commands. In this Codex environment, the same read succeeds directly but shell redirection or a subprocess can run without DNS access and fail immediately with `lookup api.github.com: no such host`. Capture `gh` output from the tool result, then write local files in a separate step. If a compound command is necessary, request network-enabled execution for that call.

## Operations

- Create: `gh issue create --title "..." --body-file <file>`.
- Read: `gh issue view <number> --comments`.
- List: `gh issue list --state open`, adding label or state filters as needed.
- Comment: `gh issue comment <number> --body-file <file>`.
- Label: `gh issue edit <number> --add-label "..."` or `--remove-label "..."`.
- Close: `gh issue close <number>`.

When a skill says to publish to the issue tracker, create a GitHub issue. When it says to fetch a ticket, read the issue and its comments.

## Pull requests

**PRs as a request surface: no.** Triage issues, not external pull requests.

## Wayfinding

Use one issue labelled `wayfinder:map` and child issues labelled `wayfinder:<type>` (`research`, `prototype`, `grilling`, or `task`). Link each child as a GitHub sub-issue when available; otherwise add it to a task list in the map and put `Part of #<map>` in the child body. Use GitHub issue dependencies for blockers; otherwise put `Blocked by: #<number>` in the child body. A ticket is available when open, unassigned, and without open blockers; first in map order wins. Claim with `gh issue edit <number> --add-assignee @me`. To resolve, comment with the answer, close the issue, and add a short pointer and link to the map's Decisions-so-far section.
