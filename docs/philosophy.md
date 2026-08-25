# Philosophy — mdbind — Book Builder

## Philosophy

**Markdown-native, book-shaped.** `content/**/*.md` is the source of truth; `book.Build` turns it into a complete, navigable site with TOC and prev/next. No custom CMS, no Node toolchain, no `framework/web` needed — focus on content.

**Principles:**

- **Lightweight & co-installable.** Built on `libs` + `libs/markdown` (Goldmark, shared with `framework`) and stdlib; `framework` and `mdbind` can be required together for progressive enhancement (`book` → `site` without rewrite).
- **Static export.** Output is `.krewire/build` by default (each page as `<path>.html` + assets) — host anywhere, configurable via `krewire.yaml` `output` or `--output/-o`.
- **Spec-driven.** `KWM-*` specs live in this repo (`docs/specs/`).


## Contribution

- Read [`project-vision.md`](https://github.com/krewire/internal/blob/main/docs/project-vision.md) and `docs/specs/index.md` before changing behavior.
- Add/update tests matching project patterns; keep suite green.
- Update `README.md` / `docs/` and specs when public behavior changes; follow ecosystem spec conventions.
