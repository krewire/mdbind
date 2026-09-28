package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/krewire/mdbind"
	"github.com/krewire/mdbind/book"
)

const usageHelp = `mdbind — Fast, Zero-Dependency Markdown Book & Docs Builder

Usage:
  mdbind <command> [flags] [arguments]

Commands:
  build       Build Markdown manuscript into a static website
  serve       Start local development server with instant live reload
  init        Scaffold a new Markdown book project
  version     Display mdbind version information

Global Flags:
  -h, --help  Show help for mdbind or a specific command

Examples:
  mdbind init my-book
  cd my-book && mdbind serve
  mdbind build --output dist
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usageHelp)
		os.Exit(2)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "build":
		os.Exit(runBuild(args))
	case "serve":
		os.Exit(runServe(args))
	case "init":
		os.Exit(runInit(args))
	case "version", "--version", "-v":
		fmt.Printf("mdbind %s (%s/%s)\n", mdbind.Version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	case "help", "--help", "-h":
		fmt.Print(usageHelp)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "mdbind: unknown command %q\nRun 'mdbind --help' for usage.\n", cmd)
		os.Exit(2)
	}
}

func runBuild(args []string) int {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	inputFlag := fs.String("input", "", "input directory containing Markdown manuscript")
	fs.StringVar(inputFlag, "i", "", "shorthand for --input")

	outputFlag := fs.String("output", "", "output directory for static website")
	fs.StringVar(outputFlag, "o", "", "shorthand for --output")

	titleFlag := fs.String("title", "", "book title")
	fs.StringVar(titleFlag, "t", "", "shorthand for --title")

	authorFlag := fs.String("author", "", "book author")
	fs.StringVar(authorFlag, "a", "", "shorthand for --author")

	baseFlag := fs.String("base", "", "base URL path (e.g. /guide/)")
	fs.StringVar(baseFlag, "b", "", "shorthand for --base")

	mountFlag := fs.String("mount", "", "subpath mount under output (e.g. docs)")
	themeFlag := fs.String("theme", "", "theme mode: auto, light, dark, or off")
	configFlag := fs.String("config", "", "path to book.yaml configuration file")
	fs.StringVar(configFlag, "c", "", "shorthand for --config")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	// 1. Locate and load configuration
	cfgFile := *configFlag
	if cfgFile == "" {
		cfgFile = book.FindConfigFile(".")
	}

	var cfg *book.Config
	if cfgFile != "" {
		loaded, err := book.LoadConfigFile(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mdbind: failed to load config %q: %v\n", cfgFile, err)
			return 1
		}
		cfg = loaded
	} else {
		cfg = &book.Config{
			Theme: book.DefaultTheme(),
		}
	}

	// 2. Flags override configuration
	if *inputFlag != "" {
		cfg.Input = *inputFlag
	}
	if *outputFlag != "" {
		cfg.Output = *outputFlag
	}
	if *titleFlag != "" {
		cfg.Title = *titleFlag
	}
	if *authorFlag != "" {
		cfg.Author = *authorFlag
	}
	if *baseFlag != "" {
		cfg.BasePath = *baseFlag
	}
	if *mountFlag != "" {
		cfg.MountPath = *mountFlag
	}
	if *themeFlag != "" {
		switch strings.ToLower(strings.TrimSpace(*themeFlag)) {
		case "auto", "default":
			cfg.Theme = book.DefaultTheme()
		case "light":
			th := book.DefaultTheme()
			th.Default = "light"
			cfg.Theme = th
		case "dark":
			th := book.DefaultTheme()
			th.Default = "dark"
			cfg.Theme = th
		case "off", "none", "false":
			cfg.Theme = nil
		}
	}

	// 3. Fallbacks
	if cfg.Input == "" {
		// If content/ exists, default to content/, otherwise .
		if st, err := os.Stat("content"); err == nil && st.IsDir() {
			cfg.Input = "content"
		} else {
			cfg.Input = "."
		}
	}
	if cfg.Output == "" {
		cfg.Output = "dist"
	}
	if cfg.Title == "" {
		cfg.Title = filepath.Base(filepath.Clean(cfg.Input))
		if cfg.Title == "." {
			wd, _ := os.Getwd()
			cfg.Title = filepath.Base(wd)
		}
	}

	start := time.Now()
	created, err := book.Build(*cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mdbind: build error: %v\n", err)
		return 1
	}

	elapsed := time.Since(start)
	for _, p := range created {
		fmt.Printf("created %s\n", p)
	}
	fmt.Printf("\n✓ Built %d files in %v → %s\n", len(created), elapsed.Round(time.Millisecond), cfg.Output)
	return 0
}

func runServe(args []string) int {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	addrFlag := fs.String("addr", ":8080", "HTTP listen address")
	inputFlag := fs.String("input", "", "input directory containing Markdown manuscript")
	fs.StringVar(inputFlag, "i", "", "shorthand for --input")

	configFlag := fs.String("config", "", "path to book.yaml configuration file")
	fs.StringVar(configFlag, "c", "", "shorthand for --config")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	cfgFile := *configFlag
	if cfgFile == "" {
		cfgFile = book.FindConfigFile(".")
	}

	var cfg *book.Config
	if cfgFile != "" {
		loaded, err := book.LoadConfigFile(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mdbind: failed to load config %q: %v\n", cfgFile, err)
			return 1
		}
		cfg = loaded
	} else {
		cfg = &book.Config{
			Theme: book.DefaultTheme(),
		}
	}

	if *inputFlag != "" {
		cfg.Input = *inputFlag
	}
	if cfg.Input == "" {
		if st, err := os.Stat("content"); err == nil && st.IsDir() {
			cfg.Input = "content"
		} else {
			cfg.Input = "."
		}
	}
	if cfg.Title == "" {
		cfg.Title = filepath.Base(filepath.Clean(cfg.Input))
		if cfg.Title == "." {
			wd, _ := os.Getwd()
			cfg.Title = filepath.Base(wd)
		}
	}

	addr := *addrFlag
	url := addr
	if strings.HasPrefix(url, ":") {
		url = "http://localhost" + url
	} else if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	fmt.Println()
	fmt.Println("  📖 mdbind local preview server")
	fmt.Printf("  • URL:      %s/\n", url)
	fmt.Printf("  • Source:   %s/\n", cfg.Input)
	fmt.Println("  • Live:     Auto-reload on browser refresh")
	fmt.Println("  • Stop:     Press Ctrl+C")
	fmt.Println()

	if err := book.ServeConfig(*cfg, addr); err != nil {
		fmt.Fprintf(os.Stderr, "mdbind: serve error: %v\n", err)
		return 1
	}
	return 0
}

func runInit(args []string) int {
	targetDir := "."
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		targetDir = args[0]
	}

	if targetDir != "." {
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "mdbind: failed to create directory %q: %v\n", targetDir, err)
			return 1
		}
	}

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mdbind: failed to read directory %q: %v\n", targetDir, err)
		return 1
	}
	if len(entries) > 0 {
		fmt.Fprintf(os.Stderr, "mdbind: directory %q is not empty. Please init in an empty folder.\n", targetDir)
		return 2
	}

	// Create book.yaml
	bookTitle := filepath.Base(filepath.Clean(targetDir))
	if bookTitle == "." {
		wd, _ := os.Getwd()
		bookTitle = filepath.Base(wd)
	}

	bookYaml := fmt.Sprintf(`title: "%s"
author: "Author Name"
input: "content"
output: "dist"
base: "/"
theme: "auto"
nav:
  - text: "GitHub"
    url: "https://github.com"
footer: "© 2026 %s. Built with mdbind."
`, bookTitle, bookTitle)

	if err := os.WriteFile(filepath.Join(targetDir, "book.yaml"), []byte(bookYaml), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "mdbind: error writing book.yaml: %v\n", err)
		return 1
	}

	contentDir := filepath.Join(targetDir, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mdbind: error creating content dir: %v\n", err)
		return 1
	}

	chapters := map[string]string{
		"01-introduction.md": `# Introduction

Welcome to your new book powered by **mdbind** — the fast, zero-dependency Markdown book builder.

## What is mdbind?

mdbind compiles directories of Markdown files into beautiful, responsive static documentation and book websites.

### Key Features

- **Blazingly Fast**: Compiles in milliseconds with Go standard library.
- **Ordered Chapters**: Name files with numeric prefixes like '01-introduction.md' for automatic chapter ordering.
- **Nested Subchapters**: Subdirectories automatically form collapsible sidebar menus.
- **Dark/Light Theme**: Built-in automatic theme switcher that respects system preferences.
`,
		"02-getting-started.md": `# Getting Started

Writing content with mdbind is as simple as creating Markdown files.

## Adding Chapters

Place markdown files inside the content/ directory:

` + "```text" + `
content/
├── 01-introduction.md
├── 02-getting-started.md
└── 03-deep-dive/
    ├── 01-architecture.md
    └── 02-internals.md
` + "```" + `

## Previewing Your Book

Run the live preview server:

` + "```bash" + `
mdbind serve
` + "```" + `

Edit any file and refresh your browser to see immediate changes!
`,
		"03-deep-dive/01-architecture.md": `# Architecture Overview

This is an example of a nested subchapter. Subdirectories in mdbind automatically turn into expandable groups in the sidebar navigation.

## Why Standalone?

mdbind does not require Node.js, Python, or complex build setups. A single Go binary compiles everything into static HTML and CSS.
`,
	}

	for path, body := range chapters {
		fullPath := filepath.Join(contentDir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "mdbind: error creating dir: %v\n", err)
			return 1
		}
		if err := os.WriteFile(fullPath, []byte(body), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "mdbind: error writing file %q: %v\n", fullPath, err)
			return 1
		}
	}

	fmt.Printf("✓ Initialized new mdbind book project in %s\n\n", targetDir)
	fmt.Println("Next steps:")
	if targetDir != "." {
		fmt.Printf("  cd %s\n", targetDir)
	}
	fmt.Println("  mdbind serve      # preview book at http://localhost:8080")
	fmt.Println("  mdbind build      # compile static site to dist/")
	return 0
}
