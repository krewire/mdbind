# Philosophy — mdbind — Book Builder

## Philosophy

**Markdown-native, book-shaped.** `manuscript/` is the source of truth; `book.Build` turns it into a complete, navigable site with TOC and prev/next. No custom CMS, no Node toolchain.

**Principles:**

- **Dogfooded.** Built on `framework/web` and `libs`; follows unified `krewire.yaml` precedence (flags > env > defaults).
- **Static export.** Output is `site/` (each page as `<path>/index.html` + assets) — host anywhere.
- **Spec-driven.** `KWM-*` specs in `krewire/internal` (`docs/specs/mdbind/`).


## Contribution

- Read [`project-vision.md`](https://github.com/krewire/internal/blob/main/docs/project-vision.md) and `docs/specs/index.md` before changing behavior.
- Add/update tests matching project patterns; keep suite green.
- Update `README.md` / `docs/` and specs when public behavior changes; follow ecosystem spec conventions.
