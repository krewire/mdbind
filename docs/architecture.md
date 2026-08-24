# Architecture — mdbind — Book Builder

## Module Structure

```
mdbind/
├── book/                 # Core: load manuscript, render Markdown (GFM), scoped CSS, file-based routing, export to site/
│   ├── assets/           # mdbind.css
│   └── templates/        # chrome.tmpl, chapter.tmpl, index.tmpl
├── cmd/mdbind/           # Standalone CLI (for non-Krewire use; Krewire projects use `kiw build`)
├── internal/commands/    # build/serve subcommands
└── docs/
```

**Design decisions:**

- **Book = one workload of the unified matrix.** `site` (`framework/web/ssg`) and `book` (`mdbind`) are siblings, both driven by `kiw build`; `docs` (this documentation site) is a `book` showcase.
- **File-based routing.** URLs mirror `manuscript/` filesystem, no `/chapters/` segment; every route is extensionless and maps one-to-one onto a sibling `.html` file.
- **Library + CLI.** `book.Build(Config{Input, Output, Title, Author})` powers both `kiw build` (book mode) and standalone `mdbind`.


## Conventions

- Documentation in English, Markdown, spec-driven (`docs/specs/`).
- Quality gates: `gofmt -l .`, `go vet ./...`, `go test ./...` in each Go repo; per-kind `kiw build` / `kiw build --plan` spot-checks.
- Cross-repo testing via the hub `go.work` workspace; temporary `replace` directives only for single-repo clones outside the workspace.
