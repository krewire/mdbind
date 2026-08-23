# mdbind — The Krewire Book Builder

**mdbind** assembles folders of Markdown into book-shaped static websites — `github.com/krewire/mdbind`. It implements the `book` project kind in the unified Krewire framework and is the engine behind the [Krewire documentation site](https://github.com/krewire/docs).

In the unified workload matrix ([`KWF-M8K2Q`](../framework/docs/specs/KWF-ARCH-M8K2Q-unified-framework-vision.md)), `mdbind` covers `book` while `framework/web/ssg` covers `site`. Both are built by `krewire build` and share routing, theming (`framework/ui`), and the single `krewire.yaml`.

## Features

- **Book-shaped output** — ordered chapters, table of contents, prev/next navigation.
- **Dir-based routing** — URLs mirror the manuscript filesystem, no `/chapters/` segment, trailing-slash normalization for static hosts.
- **Markdown native** — GFM via a pure-Go processor.
- **Static export** — complete website from one folder via `framework/web`.
- **Library or CLI** — `book.Build` powers both `krewire build` (book mode) and standalone `mdbind`.
- **Dogfooded** — built on `framework` and `libs`, following unified conventions.

## Workspace Layout

| Path | Description |
|------|-------------|
| `book/` | Site builder: loading, rendering, export. |
| `cmd/mdbind/` | Standalone `mdbind` CLI (superseded by `krewire build` for Krewire projects). |
| `internal/commands/` | CLI sub-commands (`build`, `init`, `serve`). |
| `docs/` | Specifications (`KWM-*`). |

## Getting Started

### Prerequisites

- Go 1.22+ — https://go.dev/dl/

### Building

```sh
go build ./...
go test ./...
gofmt -l . && go vet ./...
```

### Using via `krewire` (recommended for Krewire projects)

```sh
krewire new mybook
krewire init --book mybook
# manuscript/ already populated
krewire build        # book mode: manuscript/ → site/
krewire serve        # preview at :8080
```

### Using standalone `mdbind`

```sh
mdbind init
mdbind build --title "My Book" --author "Me"
open site/index.html

mdbind serve --addr :8080
```

Settings honor Krewire precedence: flags > `MDBIND_*` env > defaults (`MDBIND_INPUT`, `MDBIND_OUTPUT`, `MDBIND_TITLE`, `MDBIND_AUTHOR`, `MDBIND_BASE`, `MDBIND_ADDR`). Pass `--base /guide/` for subdirectory deploys.

### Using the library

```go
created, err := book.Build(book.Config{
    Input:  "manuscript",
    Output: "site",
    Title:  "My Book",
    Author: "Me",
})
```

## Specifications

- [`KWM-BUILDER-FX9H2`](./docs/specs/KWM-BUILDER-FX9H2-mdbind-site-builder.md) — Site Builder
- [`KWM-CLI-4TCPA`](./docs/specs/KWM-CLI-4TCPA-cli-workflows.md) — CLI & Workflows

## Related Repositories

- [framework](https://github.com/krewire/framework) — unified framework (`web`/`ssg` for `site`, `runtime` for frontend, `ui` for theming)
- [libs](https://github.com/krewire/libs) — shared primitives (`core`, `term`, `config`, `validate`)
- [docs](https://github.com/krewire/docs) — documentation site, built with `mdbind` (`book` kind)

## License

MIT — see [LICENSE](./LICENSE).
