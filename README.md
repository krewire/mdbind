# mdbind — Fast, Standalone Markdown Book & Docs Builder in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/krewire/mdbind.svg)](https://pkg.go.dev/github.com/krewire/mdbind)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)
[![Zero Dependency](https://img.shields.io/badge/Dependencies-Zero%20External-green.svg)](#)

**mdbind** is an ultra-fast, zero-dependency static site builder written in Go that compiles folders of Markdown files into beautiful, book-shaped documentation websites.

Inspired by tools like Rust's `mdBook` and `GitBook`, **mdbind** gives developers and technical authors a single, lightweight binary with **zero runtime dependencies** — no Node.js, no npm, no Python, and no complex bundlers. It compiles an entire multi-chapter book in **less than 15 milliseconds**.

`mdbind` is designed to be used **standalone** by any developer, writer, or team, while also serving as the official `book` workload engine across the Krewire ecosystem.

---

## Highlights

- ⚡ **Blazingly Fast** — Compiles dozens of chapters and nested pages in milliseconds using pure Go standard library.
- 📦 **Zero External Runtime Dependencies** — No `node_modules`, no npm scripts, no Python runtime. A single compiled binary is all you need.
- 📖 **Book-Shaped Hierarchy** — Automatically generates numbered chapters, collapsible sidebar navigation, breadcrumbs, and Previous/Next chapter pagers from your directory layout.
- 🔄 **Instant Live Authoring** — Built-in local HTTP server (`mdbind serve`) that re-renders Markdown files dynamically on page refresh.
- 🌓 **Automatic Dark / Light Theming** — Built-in theme switcher that respects user and operating system preferences with zero CSS setup.
- 🛠️ **Standalone CLI + Go Library** — Run standalone via `mdbind build` / `mdbind serve`, or import `github.com/krewire/mdbind/book` into any Go application.
- 📐 **Declarative `book.yaml`** — Configure metadata, navbar links, custom base URLs, and theme modes with a simple, human-friendly YAML file.

---

## Quickstart (Under 30 Seconds)

### 1. Installation

```bash
# Install via Go (Go 1.26+)
go install github.com/krewire/mdbind/cmd/mdbind@latest
```

Or build from source:
```bash
git clone https://github.com/krewire/mdbind.git
cd mdbind
go build -o /usr/local/bin/mdbind ./cmd/mdbind
```

### 2. Scaffold a New Book

```bash
mdbind init my-book
cd my-book
```

This creates a clean, ready-to-write book project:
```text
my-book/
├── book.yaml
└── content/
    ├── 01-introduction.md
    ├── 02-getting-started.md
    └── 03-deep-dive/
        └── 01-architecture.md
```

### 3. Live Preview & Writing

Start the local live development server:

```bash
mdbind serve
# → 📖 mdbind local preview server at http://localhost:8080/
```

Open `http://localhost:8080/` in your browser. Edit any `.md` file in `content/` and refresh to see instant updates!

### 4. Build for Production

Compile your book into static HTML, CSS, and assets ready to host on GitHub Pages, Netlify, Cloudflare Pages, or any web server:

```bash
mdbind build
# ✓ Built 6 files in 11ms → dist/
```

---

## Project Structure & Conventions

`mdbind` uses intuitive file-based ordering. No manual table-of-contents files needed.

```text
my-book/
├── book.yaml              # Optional project configuration
└── content/               # Source manuscript directory
    ├── 01-introduction.md # Chapter 1 (URL: /introduction)
    ├── 02-installation.md # Chapter 2 (URL: /installation)
    ├── 03-guides/         # Chapter 3 (Nested group)
    │   ├── 01-basic.md    # Subchapter 3.1 (URL: /guides/basic)
    │   └── 02-advanced.md # Subchapter 3.2 (URL: /guides/advanced)
    └── 04-appendix.md     # Chapter 4 (URL: /appendix)
```

- **Ordering**: Prefix filenames or folders with numbers (`01-`, `02-`, etc.) to control reading order. Numbers are stripped from URLs.
- **Headings & Title**: The first `# Heading` in each file automatically becomes the chapter title in the sidebar and navbar.
- **Frontmatter**: YAML frontmatter (`--- title: Custom ---`) is supported and stripped cleanly from output.

---

## Configuration (`book.yaml`)

Create a `book.yaml` (or `mdbind.yaml`) in your project root to customize behavior:

```yaml
title: "The Go Systems Blueprint"
author: "Alex Morgan"
input: "content"          # Input directory (default: content or .)
output: "dist"            # Output directory (default: dist)
base: "/"                 # Base URL (e.g. /guide/ for subpaths)
theme: "auto"             # auto | light | dark | off

nav:
  - text: "GitHub"
    url: "https://github.com/my-org/my-book"
  - text: "Release Notes"
    url: "https://my-org.github.io/releases"

footer: "© 2026 Alex Morgan. Built with mdbind."
```

---

## CLI Command Reference

### `mdbind build`
Compiles Markdown manuscripts into static web pages.

```bash
mdbind build [flags]

Flags:
  -i, --input    Input directory containing Markdown manuscript (default: content)
  -o, --output   Destination output directory (default: dist)
  -t, --title    Override book title
  -a, --author   Override author name
  -b, --base     Base URL path prefix (e.g. /docs/)
  --theme        Theme mode: auto, light, dark, or off
  -c, --config   Path to custom config file (default: book.yaml)
```

### `mdbind serve`
Starts an HTTP preview server with instant live re-rendering on browser refresh.

```bash
mdbind serve [flags]

Flags:
  --addr         Listen address (default: :8080)
  -i, --input    Input directory (default: content)
  -c, --config   Path to custom config file (default: book.yaml)
```

### `mdbind init`
Scaffolds a new book project with sample chapters and configuration.

```bash
mdbind init [path]
```

---

## Using as a Go Library

You can also use `mdbind` programmatically inside your own Go services or CLI tools:

```go
package main

import (
    "log"
    "github.com/krewire/mdbind/book"
)

func main() {
    created, err := book.Build(book.Config{
        Input:    "manuscript",
        Output:   "public",
        Title:    "Cloud Architecture Handbook",
        Author:   "Dev Team",
        BasePath: "/handbook/",
        Theme:    book.DefaultTheme(),
    })
    if err != nil {
        log.Fatalf("Build failed: %v", err)
    }
    log.Printf("Successfully created %d pages", len(created))
}
```

---

## Ecosystem Integration (`kiw build`)

For projects using the unified Krewire Framework:
- `mdbind` powers the `book` project kind (`kind: book` in `krewire.yaml`).
- `kiw build` automatically delegates to `mdbind` when a `content/` or `manuscript/` folder is present.
- A project can start as an `mdbind` standalone book and progressively add web apps or microservices without restructuring files.

---

## Specifications

- [`KWM-BUILDER-FX9H2`](./docs/specs/KWM-BUILDER-FX9H2-mdbind-site-builder.md) — Site Builder Specification
- [`KWM-CLI-4TCPA`](./docs/specs/KWM-CLI-4TCPA-cli-workflows.md) — Standalone CLI Specification

---

## License

MIT License — see [LICENSE](./LICENSE).
