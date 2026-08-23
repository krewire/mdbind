# Architecture — mdbind — Book Builder

## Module Structure

```
mdbind/
├── book/                 # Core: load manuscript, render Markdown (GFM), scoped CSS, dir-based routing, export to site/
│   ├── assets/           # mdbind.css
│   └── templates/        # chrome.tmpl, chapter.tmpl, index.tmpl
├── cmd/mdbind/           # Standalone CLI (for non-Krewire use; Krewire projects use `krewire build`)
├── internal/commands/    # build/serve subcommands
└── docs/
```

**Design decisions:**

- **Book = one workload of the unified matrix.** `site` (`framework/web/ssg`) and `book` (`mdbind`) are siblings, both driven by `krewire build`; `docs` (this documentation site) is a `book` showcase.
- **Dir-based routing.** URLs mirror `manuscript/` filesystem, no `/chapters/` segment, trailing-slash normalization for static hosts.
- **Library + CLI.** `book.Build(Config{Input, Output, Title, Author})` powers both `krewire build` (book mode) and standalone `mdbind`.


## Conventions

- Documentation in English, Markdown, spec-driven (`docs/specs/`).
- Quality gates: `gofmt -l .`, `go vet ./...`, `go test ./...` in each Go repo; per-kind `krewire build` / `krewire build --plan` spot-checks.
- Cross-repo testing via temporary `replace` in `go.mod`; never `go.work`.
