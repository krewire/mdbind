# mdbind — The Krewire Book Builder

**mdbind** assembles folders of Markdown into book-shaped static websites — `github.com/krewire/mdbind`. It implements the `book` project kind in the unified Krewire framework and is the engine behind the [Krewire documentation site](https://github.com/krewire/docs).

In the unified workload matrix ([`KWF-M8K2Q`](../framework/docs/specs/KWF-ARCH-M8K2Q-unified-framework-vision.md)), `mdbind` covers `book` (content-first, lightweight) while `framework/web/ssg` covers `site` (design-first, powerful). Both are built by `kiw build` to `.krewire/build`, share the Goldmark renderer via `libs/markdown`, and use the single `krewire.yaml`. A docs site can start as `book` and progressively enhance to `site` without rewrite — both modules are co-installable in the same `go.mod`.

## Features

- **Book-shaped output** — ordered chapters, table of contents, prev/next navigation.
- **File-based routing** — URLs mirror the content filesystem, no `/chapters/` segment, each chapter served as an extensionless sibling `.html` file (`/{slug}`, `/{chapter}/{sub}`).
- **Markdown native** — GFM via `libs/markdown` (Goldmark, shared with `framework`), deterministic.
- **Static export** — complete website from one folder via stdlib + `libs/markdown` (no `framework/web` needed).
- **Library or CLI** — `book.Build` powers both `kiw build` (book mode, canonical) and standalone `mdbind` (superseded by `kiw build` for Krewire projects).
- **Lightweight** — `mdbind` depends only on `libs` + `libs/markdown` + stdlib, not on `framework`; both can be required together for progressive enhancement.

## Workspace Layout

| Path | Description |
|------|-------------|
| `book/` | Site builder: loading, rendering, export. |
| `cmd/mdbind/` | Standalone `mdbind` CLI (superseded by `kiw build` for Krewire projects). |
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

### Using via `kiw` (recommended for Krewire projects)

```sh
kiw new mybook
kiw init --book mybook
# content/docs/ already populated
kiw build            # book mode: content/ → .krewire/build (README notes skipped)
kiw serve            # preview at :8080
```

### Using standalone `mdbind`

```sh
mdbind init
mdbind build --title "My Book" --author "Me"
open .krewire/build/index.html

mdbind serve --addr :8080
```

Settings honor Krewire precedence: flags > `MDBIND_*` env > defaults (`MDBIND_INPUT`, `MDBIND_OUTPUT`, `MDBIND_TITLE`, `MDBIND_AUTHOR`, `MDBIND_BASE`, `MDBIND_ADDR`). Pass `--base /guide/` for subdirectory deploys.

### Using the library

```go
created, err := book.Build(book.Config{
    Input:  "content",
    Output: ".krewire/build",
    Title:  "My Book",
    Author: "Me",
})
// Add framework progressively later:
// go get github.com/krewire/framework
// kiw init --site  // adds pages/*.kiw + ssg: without removing content/
// kiw build        // merges both into .krewire/build
```

## Specifications

- [`KWM-BUILDER-FX9H2`](./docs/specs/KWM-BUILDER-FX9H2-mdbind-site-builder.md) — Site Builder
- [`KWM-CLI-4TCPA`](./docs/specs/KWM-CLI-4TCPA-cli-workflows.md) — CLI & Workflows

## Related Repositories

- [framework](https://github.com/krewire/framework) — unified framework (`web`/`ssg` for `site`, `runtime` for frontend, `ui` for theming) — shares `libs/markdown` with `mdbind` for co-existence
- [libs](https://github.com/krewire/libs) — shared primitives (`core`, `term`, `config`, `validate`, `markdown`)
- [docs](https://github.com/krewire/docs) — documentation site, built with `mdbind` (`book` kind)

## License

MIT — see [LICENSE](./LICENSE).
