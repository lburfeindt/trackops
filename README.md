# Music Tools

A small collection of independent command-line tools for music tasks. Each tool owns its language and development lifecycle.

## Development

Requires Make and the language runtimes listed in each tool's README. The current Go tool requires Go 1.27.1 or later.

Run the registered tools' lifecycle from the repository root:

```sh
make format
make lint
make test
make verify
make build
```

## Tools

- [format-playlist](tools/format-playlist/README.md) — prints track names from an M3U playlist's `#EXTINF` entries.
