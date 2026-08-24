package book

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/krewire/framework/ui"
	"github.com/krewire/framework/web"
)

// Pages builds the site's web.Page tree: the root landing page and one page
// per chapter. Exported page paths are always root-relative; only the links
// emitted into the rendered HTML carry the book's base path.
func (b *Book) Pages() ([]web.Page, error) {
	tpls, err := parseTemplates()
	if err != nil {
		return nil, err
	}
	toc := b.chapterList(nil)
	root := web.Page{
		Path: "/",
		Render: func(w io.Writer) error {
			data := b.pageData(nil, toc)
			if flat := b.flattened(); len(flat) > 0 {
				data.First = &linkRef{Title: flat[0].Title, Link: b.link(flat[0].Path())}
				data.Pages = len(flat)
			} else {
				data.Pages = 0
			}
			return tpls.ExecuteTemplate(w, "index.tmpl", data)
		},
	}
	pages := []web.Page{root}
	for _, ch := range b.flattened() {
		data := b.pageData(ch, b.chapterList(ch))
		data.Crumbs = b.crumbs(ch)
		if ch.Prev != nil {
			data.Prev = &linkRef{Title: ch.Prev.Title, Link: b.link(ch.Prev.Path())}
		}
		if ch.Next != nil {
			data.Next = &linkRef{Title: ch.Next.Title, Link: b.link(ch.Next.Path())}
		}
		pages = append(pages, web.Page{
			Path: ch.Path(),
			Render: func(w io.Writer) error {
				return tpls.ExecuteTemplate(w, "chapter.tmpl", data)
			},
		})
	}
	return pages, nil
}

// crumbs builds the breadcrumb trail for a chapter page: Home, the containing
// chapter, and the page itself.
func (b *Book) crumbs(ch *Chapter) []crumb {
	out := []crumb{{Title: "Home", Link: b.base}}
	if ch.Parent != nil {
		out = append(out, crumb{Title: ch.Parent.Title, Link: b.link(ch.Parent.Path())})
	}
	if ch != nil {
		out = append(out, crumb{Title: ch.Title, Link: b.link(ch.Path()), Current: true})
	}
	return out
}

// chapterList builds the sidebar chapter list, marking active as the current
// page. The subs of the current (or current's parent) section are expanded
// inline with indented rows.
func (b *Book) chapterList(active *Chapter) []tocEntry {
	toc := make([]tocEntry, 0, len(b.Chapters))
	for i := range b.Chapters {
		ch := &b.Chapters[i]
		expanded := ch == active || (active != nil && active.Parent == ch)
		entry := tocEntry{
			Label:  ch.Label() + ".",
			Title:  ch.Title,
			Link:   b.link(ch.Path()),
			Active: ch == active,
		}
		if expanded {
			for _, s := range ch.Subs {
				entry.Subs = append(entry.Subs, tocEntry{
					Label:  s.Label(),
					Title:  s.Title,
					Link:   b.link(s.Path()),
					Active: s == active,
				})
			}
		}
		toc = append(toc, entry)
	}
	return toc
}

// pageData assembles the shared template payload for a page.
func (b *Book) pageData(ch *Chapter, toc []tocEntry) pageData {
	return pageData{
		BasePath:   b.base,
		Book:       b,
		Chapter:    ch,
		Toc:        toc,
		NavLinks:   b.navLinks,
		FooterText: b.footerText,
		Theme:      b.theme,
	}
}

// autoSubList renders the auto-generated chapter page for a directory without
// an index file: a list of its subchapters with numeric labels.
func (b *Book) autoSubList(ch *Chapter, subs []source) string {
	var sb strings.Builder
	sb.WriteString("<p class=\"lede\">This chapter has the following sections.</p>\n<ul class=\"sub-toc\">\n")
	for i, s := range subs {
		slug := slugFor(s.file)
		title := titleFor(s.body, slug)
		link := b.link("/" + ch.Slug + "/" + slug)
		fmt.Fprintf(&sb, "<li><a href=\"%s\"><strong>%d.%d</strong> %s</a></li>\n", link, ch.Number, i+1, template.HTMLEscapeString(title))
	}
	sb.WriteString("</ul>\n")
	return sb.String()
}

// exportPages writes each page under outDir as the sibling .html file its
// extensionless URL maps to: "/" becomes index.html and "/a/b" becomes
// outDir/a/b.html.
func exportPages(outDir string, pages []web.Page) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	for _, p := range pages {
		name := "index.html"
		if rel := strings.Trim(p.Path, "/"); rel != "" {
			name = rel + ".html"
		}
		out := filepath.Join(outDir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		f, err := os.Create(out)
		if err != nil {
			return err
		}
		if err := p.Render(f); err != nil {
			f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

// createdPaths computes the files the build will produce for the book,
// mirroring the export shape for deterministic reporting.
func (b *Book) createdPaths(outDir string) []string {
	paths := []string{filepath.Join(outDir, "index.html"), filepath.Join(outDir, "assets", "mdbind.css")}
	for _, ch := range b.flattened() {
		rel := strings.TrimPrefix(ch.Path(), "/")
		paths = append(paths, filepath.Join(outDir, filepath.FromSlash(rel)+".html"))
	}
	sort.Strings(paths)
	return paths
}

// pageData is the payload for the index and chapter page templates.
type pageData struct {
	BasePath   string
	Book       *Book
	Chapter    *Chapter
	Toc        []tocEntry
	Crumbs     []crumb
	Pages      int
	First      *linkRef
	Prev       *linkRef
	Next       *linkRef
	NavLinks   []Link
	FooterText string
	Theme      *ui.Theme
}

// tocEntry is one row in the sidebar chapter list. Subs holds the expanded
// subchapter rows of the active section.
type tocEntry struct {
	Label  string
	Title  string
	Link   string
	Active bool
	Subs   []tocEntry
}

// crumb is one step in a page's breadcrumb trail.
type crumb struct {
	Title   string
	Link    string
	Current bool
}

// linkRef is a prev/next navigation target.
type linkRef struct {
	Title string
	Link  string
}
