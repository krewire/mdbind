package book

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Page is a renderable page for export or serve. Path is the extensionless
// URL (e.g. "/getting-started"); Render writes the HTML to w.
type Page struct {
	Path   string
	Render func(w io.Writer) error
}

// Pages builds the site's Page tree: the root landing page and one page
// per chapter. Exported page paths are always root-relative; only the links
// emitted into the rendered HTML carry the book's base path.
func (b *Book) Pages() ([]Page, error) {
	tpls, err := parseTemplates()
	if err != nil {
		return nil, err
	}
	toc := b.chapterList(nil)
	pages := make([]Page, 0, len(b.Chapters)+1)
	if !b.noRootTOC {
		root := Page{
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
		pages = append(pages, root)
	}
	for _, ch := range b.flattened() {
		data := b.pageData(ch, b.chapterList(ch))
		data.Crumbs = b.crumbs(ch)
		if ch.Prev != nil {
			data.Prev = &linkRef{Title: ch.Prev.Title, Link: b.link(ch.Prev.Path())}
		}
		if ch.Next != nil {
			data.Next = &linkRef{Title: ch.Next.Title, Link: b.link(ch.Next.Path())}
		}
		pages = append(pages, Page{
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

// chapterList builds the sidebar chapter list for responsive collapsible nav.
// Every chapter with subs gets a collapsible group; Open is true when the
// group contains the active page (itself or any sub). All groups are rendered
// so the client can toggle them without a round-trip — CSS handles the
// collapsed state, JS persists it.
func (b *Book) chapterList(active *Chapter) []tocEntry {
	toc := make([]tocEntry, 0, len(b.Chapters))
	for i := range b.Chapters {
		ch := &b.Chapters[i]
		hasSubs := len(ch.Subs) > 0
		activeSelf := ch == active
		activeInSubs := false
		subs := make([]tocEntry, 0, len(ch.Subs))
		for _, s := range ch.Subs {
			isActive := s == active
			if isActive {
				activeInSubs = true
			}
			subs = append(subs, tocEntry{
				Label:  s.Label(),
				Title:  s.Title,
				Link:   b.link(s.Path()),
				Active: isActive,
			})
		}
		open := activeSelf || activeInSubs
		toc = append(toc, tocEntry{
			Label:   ch.Label() + ".",
			Title:   ch.Title,
			Link:    b.link(ch.Path()),
			Active:  activeSelf,
			HasSubs: hasSubs,
			Open:    open,
			Subs:    subs,
		})
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
		Version:    b.version,
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
// outDir/a/b.html. A non-empty mount sub-path relocates the files so a book
// can share one output with an ssg site.
func exportPages(outDir string, pages []Page, mount string) error {
	root := outDir
	if mount != "" {
		root = filepath.Join(outDir, filepath.FromSlash(mount))
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	for _, p := range pages {
		name := "index.html"
		if rel := strings.Trim(p.Path, "/"); rel != "" {
			name = rel + ".html"
		}
		out := filepath.Join(root, filepath.FromSlash(name))
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
	baseOut := outDir
	if b.mount != "" {
		baseOut = filepath.Join(outDir, filepath.FromSlash(b.mount))
	}
	var paths []string
	if !b.noRootTOC {
		paths = append(paths, filepath.Join(baseOut, "index.html"))
	}
	paths = append(paths, filepath.Join(baseOut, "assets", "mdbind.css"))
	for _, ch := range b.flattened() {
		rel := strings.TrimPrefix(ch.Path(), "/")
		paths = append(paths, filepath.Join(baseOut, filepath.FromSlash(rel)+".html"))
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
	Version    string
	Theme      *Theme
}

// tocEntry is one row in the sidebar chapter list. HasSubs indicates a
// collapsible group; Open indicates it should be expanded initially (active
// path). Subs holds all subchapters for that group.
type tocEntry struct {
	Label   string
	Title   string
	Link    string
	Active  bool
	HasSubs bool
	Open    bool
	Subs    []tocEntry
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
