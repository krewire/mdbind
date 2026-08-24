# Specification — mdbind Site Builder

| Field       | Value                                       |
| ----------- | ------------------------------------------- |
| SpecID      | KWM-FX9H2                                   |
| Title       | mdbind Site Builder — Initial Specification |
| Status      | Draft                                       |
| Date        | 2026-08-18                                  |
| Author      | Krewire Contributors                         |
| Domain      | Site Builders                              |

## 1. Context

**mdbind** is the Krewire site builder: it assembles many Markdown files into a
single, book-shaped static website with a table of contents, subchapter
sections, breadcrumbs, and prev/next navigation. It is the first consumer of
the Krewire Web and UI Frameworks (`krewire/framework/web`, `krewire/framework/ui`)
and the reference dogfooding project in the ecosystem.

The builder is also a library: downstream projects such as the Krewire website
(`krewire/docs`) import it to generate their sites, making the dependency chain
`krewire/docs -> mdbind -> framework/web + framework/ui -> libs` explicit.

### 1.1 Current State

- Module root `github.com/krewire/mdbind`; public `book` package for site generation.
- Built on `github.com/krewire/framework` (web + cli + ui) and `github.com/krewire/libs`.
- Markdown rendering through a pure-Go markdown processor (goldmark).
- Shared module root: Go 1.22, MIT license.

## 2. Problem Statement

Publishing a set of Markdown documents as a cohesive reading experience — with
ordering, navigation, and styling — currently means hand-rolling HTML
generation or adopting a full static-site generator with its own model and
plugin surface. Every project re-invents chapter ordering, page templates, and
asset handling, and none of it reuses the Krewire web foundation.

Long works are never flat: real manuals have parts, sections, and subsections.
mdbind removes the friction by letting a folder of Markdown files — plus
sub-folders for subchapters — become a book-shaped website in one command,
rendered through the Krewire Web Framework and reusable as a library.

## 3. Goals

- G1 — Turn a folder of Markdown files and sub-folders into an ordered, book-shaped website with subchapter sections.
- G2 — Generate pages and assets through `framework/web` (routes, templates, export).
- G3 — Expose the builder as a public library usable by other Krewire projects.
- G4 — Keep output deterministic and the reading experience first-class, with breadcrumbs and section-aware navigation.

## 4. Non-Goals

The following are explicitly out of scope for the initial phase:

- NG1 — A full static-site generator (taxonomies, feeds, drafts, plugins).
- NG2 — A content-management or WYSIWYG editing surface.
- NG3 — Internationalization or multi-language book support.

## 5. Requirements

### 5.1 Manuscript Model

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-SB-001 | Load every `*.md` file in the manuscript directory and in its direct sub-directories. | Must     |
| MDB-SB-002 | Order chapters by numeric filename prefix, then filename; top-level entries and subchapter files each order independently. | Must |
| MDB-SB-003 | Derive the chapter title from its first `#` heading, falling back to its slug. | Must |
| MDB-SB-004 | Render Markdown bodies to HTML before inclusion in pages.         | Must     |
| MDB-SB-019 | Treat a manuscript sub-directory as a chapter with subchapters; number subchapters "N.M" within their chapter. | Must     |
| MDB-SB-020 | Serve chapters and subchapters at extensionless file URLs mirroring the manuscript tree (`/{slug}`, `/{chapter}/{sub}`), each emitted as a sibling `.html` file, under the configured base path. Internal links carry the same extensionless form. | Must |
| MDB-SB-021 | An `index.md` or `_index.md` in a sub-directory becomes its chapter page body; without one, the chapter page auto-lists its subchapters. | Must |

### 5.2 Site Generation

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-SB-005 | Build a page tree through `web.Router` and `web.Page` primitives. | Must     |
| MDB-SB-006 | Export the page tree and assets with `web.Export`.                | Must     |
| MDB-SB-007 | Provide a root page listing the table of contents.                | Must     |
| MDB-SB-008 | Provide one page per chapter and subchapter with prev/next navigation across the full reading order. | Must |
| MDB-SB-009 | Ship a default reader stylesheet as an asset.                     | Must     |
| MDB-SB-010 | Produce deterministic output for identical input.                 | Must     |

### 5.3 Library API

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-SB-011 | Expose a public `Build` entry point returning created paths.      | Must     |
| MDB-SB-012 | Keep the CLI and the library API backed by the same builder.      | Must     |
| MDB-SB-013 | Support a base path (`/guide/` etc.) that prefixes all generated links and asset references, leaving exported file paths root-relative. | Must |

### 5.4 Page Chrome

| ID         | Requirement                                                       | Priority |
| ---------- | ----------------------------------------------------------------- | -------- |
| MDB-SB-014 | Every page renders a navbar: brand link to the site root plus optional extra links. | Must |
| MDB-SB-015 | Every page renders a sidebar listing all chapters; the active chapter's subchapters are shown indented beneath it. | Must |
| MDB-SB-016 | Every page renders a footer with footer text or the author copyright, plus a credit line. | Must |
| MDB-SB-017 | Chrome settings are configurable per build (`NavLinks`, `FooterText`). | Must |
| MDB-SB-018 | When enabled, pages render a light/dark theme switcher built on the UI Framework theming system (`ui.Theme`), with all colors driven by configurable light/dark palettes. | Must |
| MDB-SB-022 | Chapter and subchapter pages render a breadcrumb trail: Home, the containing chapter, then the current page. | Must |

## 6. Non-Functional Requirements

- NFR1 — **Memory safety.** The `unsafe` package must not be used.
- NFR2 — **Performance.** Building a book-sized manuscript must complete quickly.
- NFR3 — **Portability.** Linux, macOS, and Windows.
- NFR4 — **Quality gates.** `gofmt`, `go vet ./...`, and `go test ./...` must pass.

## 7. Success Criteria

- S1 — A sample manuscript builds into `index.html`, chapter and subchapter pages, and assets.
- S2 — Chapter pages link to their neighbors, the table of contents, and a breadcrumb trail.
- S3 — The `krewire/docs` repository builds its website through `book.Build` without importing mdbind internals.
- S4 — Identical manuscripts produce byte-identical output.
- S5 — A directory chapter renders its subchapters in the sidebar and auto-lists them when no index file is present.

## 8. Related Specifications

| SpecID    | Title                                           |
| --------- | ----------------------------------------------- |
| [KWM-4TCPA](./KWM-CLI-4TCPA-cli-workflows.md)        | mdbind CLI & Workflows              |
| [KWF-M07QS](https://github.com/krewire/framework/blob/main/docs/specs/KWF-WEB-M07QS-krewire-web-framework.md) | Krewire Web Framework |
| [KWF-5XJFC](https://github.com/krewire/framework/blob/main/docs/specs/KWF-CLI-5XJFC-cli-application-model.md) | CLI Application Model |

## 9. References

- [Framework — KWF-CMBZJ](https://github.com/krewire/framework/blob/main/docs/specs/KWF-META-CMBZJ-krewire-meta-framework.md) — Krewire Framework initial specification.
- [KWL-M1ZKS](https://github.com/krewire/libs/blob/main/docs/specs/KWL-CORE-M1ZKS-krewire-libraries.md) — Krewire Libraries initial specification.
- [KWL-R934Y](https://github.com/krewire/libs/blob/main/docs/specs/KWL-TERM-R934Y-terminal-io-rendering.md) — Terminal I/O & Rendering.