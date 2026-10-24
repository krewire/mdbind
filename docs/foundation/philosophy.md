# Philosophy — mdbind — Book Builder

## Philosophy

**Markdown-native, book-shaped.** `content/**/*.md` is the source of truth; `book.Build` turns it into a complete, navigable site with TOC and prev/next. No custom CMS, no Node toolchain, no `packages/web` needed — focus on content.

**Principles:**

- **Lightweight & co-installable.** Built on `packages/markdown` (Goldmark, shared across the ecosystem) and standard library; `pkg/*` and `mdbind` can be required together for progressive enhancement (`book` → `site` without rewrite).
- **Static export.** Output is `.krewire/build` by default (each page as `<path>.html` + assets) — host anywhere, configurable via `krewire.yaml` `output` or `--output/-o`.


## Contribution

- Read [`project-vision.md`](https://github.com/krewire/internal/blob/main/docs/project-vision.md) before changing behavior.
- Add/update tests matching project patterns; keep suite green.
- Update `README.md` / `docs/` when public behavior changes.
