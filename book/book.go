// Package book assembles Markdown manuscripts into book-shaped websites using
// the Krewire web framework. A manuscript is a directory of ordered Markdown
// files and subdirectories; the builder renders it into pages and assets that
// can be exported to disk or served over HTTP.
//
// File responsibilities: chapter.go (chapter model), load.go (manuscript
// loading), naming.go (ordering/slugs/titles), links.go (base-path links),
// pages.go (page-tree rendering), build.go (export orchestration),
// serve.go (HTTP serving).
package book

import (
	"github.com/krewire/framework/ui"
)

// Book is a manuscript, loaded into an ordered set of chapters.
type Book struct {
	// Title is the book title.
	Title string
	// Author is the book author.
	Author string
	// Chapters holds the top-level chapters in reading order.
	Chapters []Chapter
	// base is the URL base prefix applied to generated links, e.g. "/guide/".
	base string
	// navLinks holds optional navbar links.
	navLinks []Link
	// footerText holds optional footer text.
	footerText string
	// theme, when set, enables the light/dark theme switcher.
	theme *ui.Theme
}

// Link is a named navigation target shown in the navbar.
type Link struct {
	Text string
	URL  string
}

// flattened returns all pages in reading order: each chapter page followed by
// its subchapters.
func (b *Book) flattened() []*Chapter {
	var flat []*Chapter
	for i := range b.Chapters {
		flat = append(flat, &b.Chapters[i])
		flat = append(flat, b.Chapters[i].Subs...)
	}
	return flat
}
