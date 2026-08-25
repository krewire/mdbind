// Package book assembles Markdown manuscripts into book-shaped websites.
// A manuscript is a directory of ordered Markdown files and subdirectories;
// the builder renders it into pages and assets that can be exported to disk
// or served over HTTP. It is intentionally lightweight and depends only on
// libs/markdown and stdlib — not on framework/web or framework/ssg — so it
// can be depended on alongside framework progressively.
//
// File responsibilities: chapter.go (chapter model), load.go (manuscript
// loading), naming.go (ordering/slugs/titles), links.go (base-path links),
// pages.go (page-tree rendering), build.go (export orchestration),
// serve.go (HTTP serving), theme.go (local light/dark theming).
package book

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
	// mount is the export sub-path ("", "docs", …) prefixed to every exported
	// page path; it mirrors base without the trailing slash.
	mount string
	// navLinks holds optional navbar links.
	navLinks []Link
	// footerText holds optional footer text.
	footerText string
	// version holds the product version shown in the credit line.
	version string
	// theme, when set, enables the light/dark theme switcher.
	theme *Theme
	// noRootTOC suppresses the generated root TOC page so a book can share
	// one output with an ssg site whose landing page owns "/". Set via
	// Config.NoRootTOC.
	noRootTOC bool
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
