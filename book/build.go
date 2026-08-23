package book

import (
	"github.com/krewire/framework/ui"
	"github.com/krewire/framework/web"
)

// Config configures a build.
type Config struct {
	// Title is the book title.
	Title string
	// Author is the book author.
	Author string
	// Input is the manuscript directory.
	Input string
	// Output is the destination directory for the exported site.
	Output string
	// BasePath is the URL base the site will be served under, e.g. "/guide/".
	// Generated links and asset references are prefixed accordingly. Defaults
	// to "/".
	BasePath string
	// NavLinks are optional extra links rendered in the navbar.
	NavLinks []Link
	// FooterText is optional text shown in the footer; when empty the footer
	// shows the author copyright instead.
	FooterText string
	// Theme, when non-nil, enables the light/dark theme switcher on every page.
	Theme *ui.Theme
}

// Build renders the manuscript in cfg.Input into a static website in
// cfg.Output and returns the list of created paths, sorted for determinism.
func Build(cfg Config) ([]string, error) {
	if cfg.Input == "" {
		cfg.Input = "."
	}
	if cfg.Output == "" {
		cfg.Output = "site"
	}
	b, err := LoadWithBase(cfg.Input, cfg.Title, cfg.Author, cfg.BasePath)
	if err != nil {
		return nil, err
	}
	b.WithChrome(cfg.NavLinks, cfg.FooterText)
	if cfg.Theme != nil {
		b.WithTheme(cfg.Theme)
	}
	pages, err := b.Pages()
	if err != nil {
		return nil, err
	}
	if err := web.Export(cfg.Output, pages, Assets()); err != nil {
		return nil, err
	}
	return b.createdPaths(cfg.Output), nil
}

// WithChrome sets the book's navigation chrome: optional navbar links and
// footer text. Returns the Book for chaining.
func (b *Book) WithChrome(nav []Link, footer string) *Book {
	b.navLinks = append([]Link(nil), nav...)
	b.footerText = footer
	return b
}

// WithTheme enables the light/dark theme switcher. Returns the Book for
// chaining.
func (b *Book) WithTheme(t *ui.Theme) *Book {
	if t != nil {
		b.theme = t
	}
	return b
}
