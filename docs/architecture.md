# Architecture — mdbind — Book Builder

## Module Structure

```
mdbind/
├── book/                 # Core: load content/, render Markdown via libs/markdown (GFM), file-based routing, export to .krewire/build
│   ├── assets/           # mdbind.css (+ local ThemeModeVarsCSS/ThemeToggleCSS)
│   ├── templates/        # chrome.tmpl, chapter.tmpl, index.tmpl
│   └── theme.go          # local Theme/Palette (mirrors framework/ui shape, no framework import)
├── cmd/mdbind/           # Standalone CLI (for non-Krewire use; Krewire projects use `kiw build`)
├── internal/commands/    # build/serve subcommands
└── docs/
```

**Design decisions:**

- **Book = one workload of the unified matrix.** `site` (`framework/web/ssg`, design-first) and `book` (`mdbind`, content-first) are siblings, both driven by `kiw build` to `.krewire/build`; `docs` is a `book` showcase that can progressively add `site`.
- **File-based routing.** URLs mirror `content/` filesystem, no `/chapters/` segment; every route is extensionless and maps one-to-one onto a sibling `.html` file.
- **Library + CLI.** `book.Build(Config{Input, Output, Title, Author})` powers both `kiw build` (book mode) and standalone `mdbind`; local `book.Theme` avoids `framework/ui` leakage so both modules co-exist.
- **Shared Markdown.** Both `mdbind` and `framework` use `libs/markdown` (Goldmark GFM + `PrefixLinks`); no duplicate parsers, no `gomarkdown` divergence.
- **Progressive.** A project may `require mdbind` + `require framework` and have both `content/` and `pages/*.kiw`/`ssg:` — `kiw build` merges both into the same output.


## Conventions

- Documentation in English, Markdown, spec-driven (`docs/specs/`).
- Quality gates: `gofmt -l .`, `go vet ./...`, `go test ./...` in each Go repo; per-kind `kiw build` / `kiw build --plan` spot-checks.
- Cross-repo testing via the hub `go.work` workspace; temporary `replace` directives only for single-repo clones outside the workspace.
