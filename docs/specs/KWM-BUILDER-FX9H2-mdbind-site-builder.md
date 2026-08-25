# Specification — mdbind Site Builder

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | KWM-FX9H2                                   |
| Title       | mdbind Site Builder — Content-First Book Builder |
| Status      | Draft                                       |
| Date        | 2026-08-25                                  |
| Author      | Krewire Contributors                         |
| Domain      | Site Builders                              |

## 1. Context

**mdbind** is the Krewire book builder: it assembles many Markdown files into a
single, book-shaped static website with a table of contents, subchapter
sections, breadcrumbs, and prev/next navigation. It is the engine behind
`kiw init --book` and is intentionally **content-first and lightweight** — a
developer writes `content/**/*.md` and gets a complete docs/book site without
learning `framework/web/ssg` components, layouts, or scoped CSS.

This contrasts with `framework/web/ssg` (`site` kind, `KWF-PT8OD`), which is
**design-first and powerful**: `.kiw` DSL, components, layouts, collections,
`ssg:` key. Both are driven by `kiw build`, both default to `.krewire/build`,
and both share the Markdown renderer via `libs/markdown` (`KWL-Q3N8P`), but
**neither depends on the other**. A user project can `require
github.com/krewire/mdbind` and `require github.com/krewire/framework` together;
a docs site can start as `book` and progressively enhance to `site` (add
`pages/*.kiw` or `ssg:`) without rewrite — `kiw build` then emits both into the
same output.

The builder is also a library: downstream projects import `book.Build` to
generate their sites. The dependency chain is now `krewire/docs -> mdbind ->
libs/markdown + stdlib` and separately `krewire/app -> framework -> libs/markdown`,
with `krewire/docs+app -> mdbind + framework` converging on `libs/markdown`.

### 1.1 Current State

- Module root `github.com/krewire/mdbind`; public `book` package for site generation.
- Built on `github.com/krewire/libs` (`libs/markdown` for Goldmark GFM) and stdlib (`net/http`, `html/template`, `os`), **no `framework/web` or `framework/ui`** in public API.
- Markdown rendering through `libs/markdown` (Goldmark GFM + AutoHeadingID + base-path `PrefixLinks`), shared with `framework/web/ssg` and `framework/dsl`.
- Public `book.Theme`/`book.Palette` mirrors `framework/ui.Theme` shape but lives locally, so `framework` is not pulled when only `mdbind` is needed.
- Default output `.krewire/build` (aligned with `kiw` `config.DefaultOutput`).

## 2. Problem Statement

Publishing a set of Markdown documents as a cohesive reading experience — with
ordering, navigation, and styling — currently means hand-rolling HTML
generation or adopting a full static-site generator with its own model and
plugin surface. Every project re-invents chapter ordering, page templates, and
asset handling.

`framework/ssg` solves this for **custom sites** but forces every docs site to
adopt its component model and pull its surface, even when only content matters.
Long works are never flat: real manuals have parts, sections, and subsections.
`mdbind` removes the friction by letting a folder of Markdown files — plus
sub-folders for subchapters — become a book-shaped website in one command,
with batteries-included chrome, and remains co-installable with `framework` for
progressive enhancement.

## 3. Goals

- G1 — Turn a folder of Markdown files and sub-folders into an ordered, book-shaped website with subchapter sections.
- G2 — Generate pages and assets via **stdlib + `libs/markdown` only** (no `framework/web`/`framework/ui` in public API), so `mdbind` stays lightweight and co-exists with `framework` in the same `go.mod`.
- G3 — Expose the builder as a public library usable by other Krewire projects without pulling `framework`.
- G4 — Keep output deterministic and the reading experience first-class, with breadcrumbs and section-aware navigation.
- G5 — Support progressive enhancement: a project with both `content/` and `pages/*.kiw`/`ssg:` builds both into the same `.krewire/build` without conflict.

## 4. Non-Goals

The following are explicitly out of scope for the initial phase:

- NG1 — A full static-site generator (taxonomies, feeds, drafts, plugins).
- NG2 — A content-management or WYSIWYG editing surface.
- NG3 — Internationalization or multi-language book support.
- NG4 — Generic component/layout authoring or scoped CSS — use `framework/web/ssg` for that.
- NG5 — Reusing `framework/ui.Theme` type directly — `book.Theme` is local (shape-mirrored) to avoid the dependency.

## 5. Requirements

### 5.1 Manuscript Model

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-SB-001 | Load every `*.md` file in the content directory and in its direct sub-directories. | Must     |
| MDB-SB-002 | Order chapters by numeric filename prefix, then filename; top-level entries and subchapter files each order independently. | Must |
| MDB-SB-003 | Derive the chapter title from its first `#` heading, falling back to its slug. | Must |
| MDB-SB-004 | Render Markdown bodies to HTML via `libs/markdown.RenderWithBase` before inclusion in pages. | Must     |
| MDB-SB-019 | Treat a content sub-directory as a chapter with subchapters; number subchapters "N.M" within their chapter. | Must     |
| MDB-SB-020 | Serve chapters and subchapters at extensionless file URLs mirroring the content tree (`/{slug}`, `/{chapter}/{sub}`), each emitted as a sibling `.html` file, under the configured base path. Internal links carry the same extensionless form. | Must |
| MDB-SB-021 | An `index.md` or `_index.md` in a sub-directory becomes its chapter page body; without one, the chapter page auto-lists its subchapters. | Must |
| MDB-SB-025 | Content filtering via include/exclude glob rules (`**` crosses segments): nil include = `**/*.md`; nil exclude = `**/README.md`, `**/readme.md` (developer-internal notes); an empty non-nil list disables that side; a directory whose files are all filtered out is dropped. Configurable from `kiw build` via `build.include/exclude` in krewire.yaml or `--include/--exclude`. | Must |
| MDB-SB-026 | A leading YAML frontmatter block (`--- … ---`, optional BOM) is stripped before title/render so content files stay shared with framework/web/ssg collections. | Must |
| MDB-SB-027 | Default input is the project `content/` directory (shared with the ssg layout); legacy `manuscript/` remains accepted by `kiw build` dispatch. | Must |

### 5.2 Site Generation

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-SB-005 | Build a page tree via local `book.Page{Path, Render}` (stdlib `html/template`), not `web.Page`/`web.Router`. | Must     |
| MDB-SB-006 | Export the page tree and assets via stdlib `os.WriteFile` + `exportAssets` (sorted for determinism), to `Config.Output` defaulting to `.krewire/build`. | Must     |
| MDB-SB-007 | Provide a root page listing the table of contents.                | Must     |
| MDB-SB-008 | Provide one page per chapter and subchapter with prev/next navigation across the full reading order. | Must |
| MDB-SB-009 | Ship a default reader stylesheet as an asset (`assets/mdbind.css` includes local `ThemeModeVarsCSS`+`ThemeToggleCSS`). | Must     |
| MDB-SB-010 | Produce deterministic output for identical input.                 | Must     |

### 5.3 Library API

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-SB-011 | Expose a public `Build(Config) ([]string, error)` returning created paths, sorted. | Must     |
| MDB-SB-012 | Keep the library and the `kiw build` (book mode) backed by the same builder; `book.Config.Output` defaults to `.krewire/build` when empty. | Must     |
| MDB-SB-013 | Support a base path (`/guide/` etc.) that prefixes all generated links and asset references, leaving exported file paths root-relative. | Must |
| MDB-SB-023 | `book.Config` and `Book` must not import `framework/*`; public `Theme` is `book.Theme`/`Palette` (local, shape-mirrored from `framework/ui`), and `Handler() (http.Handler, error)` uses stdlib `ServeMux`. | Must |
| MDB-SB-024 | A `go.mod` can `require github.com/krewire/framework` and `require github.com/krewire/mdbind` together without import cycle; `kiw build` with both `content/` and `ssg` present builds both into the same output (progressive). | Must |

### 5.4 Page Chrome

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-SB-014 | Every page renders a navbar: brand link to the site root plus optional extra links. | Must |
| MDB-SB-015 | Every page renders a sidebar listing all chapters; the active chapter's subchapters are shown indented beneath it. | Must |
| MDB-SB-016 | Every page renders a footer with footer text or the author copyright, plus a credit line. | Must |
| MDB-SB-017 | Chrome settings are configurable per build (`NavLinks`, `FooterText`). | Must |
| MDB-SB-018 | When enabled, pages render a light/dark theme switcher using local `book.Theme` (light/dark palettes, `Script()`/`Button()`), with all colors driven by configurable palettes. | Must |
| MDB-SB-022 | Chapter and subchapter pages render a breadcrumb trail: Home, the containing chapter, then the current page. | Must |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Performance.** Building a book-sized manuscript must complete quickly.
- NFR3 — **Portability.** Linux, macOS, and Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.
- NFR5 — **No framework leakage.** `go list -f '{{.Deps}}' github.com/krewire/mdbind/book` must not contain `github.com/krewire/framework`.

## 7. Success Criteria

- S1 — A sample content tree builds into `.krewire/build/index.html`, chapter and subchapter pages, and `assets/mdbind.css`.
- S2 — Chapter pages link to their neighbors, the table of contents, and a breadcrumb trail.
- S3 — A user project with `require mdbind` + `require framework` builds both via `kiw build` (merged `.krewire/build`) and `go test ./...` passes in both modules.
- S4 — Identical content trees produce byte-identical output.
- S5 — A directory chapter renders its subchapters in the sidebar and auto-lists them when no index file is present.
- S6 — `book.Handler()` serves book and assets via `net/http` without `framework/web` and returns 404 for missing paths.

## 8. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [KWL-Q3N8P](./KWL-MARKDOWN-Q3N8P-shared-markdown-renderer.md) | Shared Markdown Renderer (libs) |
| [KWF-PT8OD](https://github.com/krewire/framework/blob/main/docs/specs/KWF-SSG-PT8OD-static-site-generator.md) | Krewire Static Site Generator (sibling, not parent) |
| [KWN-1QGI2](https://github.com/krewire/kiw/blob/main/docs/specs/KWN-BUILD-1QGI2-project-building.md) | Project Building (progressive book+site) |
| [KWM-4TCPA](./KWM-CLI-4TCPA-cli-workflows.md)        | mdbind CLI & Workflows (superseded by kiw) |

## 9. References

- [Framework — KWF-CMBZJ](https://github.com/krewire/framework/blob/main/docs/specs/KWF-META-CMBZJ-krewire-meta-framework.md) — Krewire Framework initial specification.
- [KWL-M1ZKS](https://github.com/krewire/libs/blob/main/docs/specs/KWL-CORE-M1ZKS-krewire-libraries.md) — Krewire Libraries initial specification.
- [KWL-Q3N8P](https://github.com/krewire/libs/blob/main/docs/specs/KWL-MARKDOWN-Q3N8P-shared-markdown-renderer.md) — Shared Markdown Renderer.
