package book

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Source is a loaded manuscript file.
type source struct {
	file  string
	title string
	body  []byte
}

// item is one top-level manuscript unit: a chapter file or a directory that
// becomes a chapter with subchapters.
type item struct {
	// name is the entry name used for ordering and slugging.
	name string
	// chapter holds the chapter body; for directories this is the optional
	// index.md/_index.md, empty when the chapter page is auto-generated.
	chapter source
	// subs holds the subchapter sources for a directory unit.
	subs []source
}

// Load reads a manuscript directory into a Book in reading order, served from
// the site root. Default include/exclude rules apply (README notes skipped).
func Load(input, title, author string) (*Book, error) {
	return LoadWithBase(input, title, author, "/")
}

// LoadWithRules reads a content directory into a Book applying explicit
// include/exclude glob rules (slash-separated patterns relative to input;
// ** crosses segments). nil selects the defaults — include "**/*.md",
// exclude README/readme developer notes.
func LoadWithRules(input, title, author, base string, include, exclude []string) (*Book, error) {
	return loadWithRules(input, title, author, base, include, exclude)
}

// LoadWithBase reads a manuscript directory into a Book in reading order. Top
// level Markdown files become chapters; a directory becomes a chapter whose
// subchapters are the Markdown files inside it, ordered by numeric prefix. An
// index.md or _index.md inside a directory becomes the chapter page body;
// without one, the chapter page lists its subchapters automatically. The base
// is the URL prefix the site will be served under, e.g. "/guide/".
// YAML frontmatter (leading --- block) is tolerated and stripped so content
// files stay compatible with framework/web/ssg collections.
func LoadWithBase(input, title, author, base string) (*Book, error) {
	return loadWithRules(input, title, author, base, nil, nil)
}

func loadWithRules(input, title, author, base string, include, exclude []string) (*Book, error) {
	if input == "" {
		input = "."
	}
	entries, err := os.ReadDir(input)
	if err != nil {
		return nil, fmt.Errorf("book: read manuscript %s: %w", input, err)
	}

	var items []item
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			it := item{name: name}
			subEntries, err := os.ReadDir(filepath.Join(input, name))
			if err != nil {
				return nil, fmt.Errorf("book: read %s: %w", name, err)
			}
			for _, se := range subEntries {
				sn := se.Name()
				if se.IsDir() || !strings.HasSuffix(sn, ".md") {
					continue
				}
				if !accepts(include, exclude, name+"/"+sn) {
					continue
				}
				data, err := os.ReadFile(filepath.Join(input, name, sn))
				if err != nil {
					return nil, fmt.Errorf("book: read %s: %w", filepath.Join(name, sn), err)
				}
				if isIndex(sn) {
					it.chapter = source{file: sn, body: stripFrontmatter(data)}
					continue
				}
				it.subs = append(it.subs, source{file: sn, body: stripFrontmatter(data)})
			}
			sort.Slice(it.subs, func(i, j int) bool {
				ni, si := splitOrder(it.subs[i].file)
				nj, sj := splitOrder(it.subs[j].file)
				if ni != nj {
					return ni < nj
				}
				return si < sj
			})
			if it.chapter.file != "" || len(it.subs) > 0 {
				items = append(items, it)
			}
			continue
		}
		if !strings.HasSuffix(name, ".md") || isIndex(name) {
			continue
		}
		if !accepts(include, exclude, name) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(input, name))
		if err != nil {
			return nil, fmt.Errorf("book: read %s: %w", name, err)
		}
		items = append(items, item{name: name, chapter: source{file: name, body: stripFrontmatter(data)}})
	}

	sort.Slice(items, func(i, j int) bool {
		ni, si := splitOrder(items[i].name)
		nj, sj := splitOrder(items[j].name)
		if ni != nj {
			return ni < nj
		}
		return si < sj
	})

	base = normalizeBase(base)
	book := &Book{Title: title, Author: author, base: base, Chapters: make([]Chapter, 0, len(items))}
	for i, it := range items {
		slug := slugFor(it.name)
		book.Chapters = append(book.Chapters, Chapter{
			Number: i + 1,
			Slug:   slug,
			Title:  titleFor(it.chapter.body, slug),
		})
		ch := &book.Chapters[len(book.Chapters)-1]
		if it.chapter.file != "" {
			html, err := renderMarkdown(it.chapter.body, base)
			if err != nil {
				return nil, fmt.Errorf("book: render %s: %w", it.chapter.file, err)
			}
			ch.Body = template.HTML(html)
		} else {
			ch.Body = template.HTML(book.autoSubList(ch, it.subs))
		}
		for j, s := range it.subs {
			slug := slugFor(s.file)
			sub := Chapter{
				Number: i + 1,
				Sub:    j + 1,
				Slug:   slug,
				Title:  titleFor(s.body, slug),
				Parent: ch,
			}
			html, err := renderMarkdown(s.body, base)
			if err != nil {
				return nil, fmt.Errorf("book: render %s: %w", s.file, err)
			}
			sub.Body = template.HTML(html)
			ch.Subs = append(ch.Subs, &sub)
		}
	}
	flat := book.flattened()
	for i, c := range flat {
		if i > 0 {
			c.Prev = flat[i-1]
		}
		if i < len(flat)-1 {
			c.Next = flat[i+1]
		}
	}
	return book, nil
}

// isIndex reports whether name is a directory chapter index file.
func isIndex(name string) bool {
	stem := strings.TrimSuffix(name, filepath.Ext(name))
	return stem == "index" || stem == "_index"
}

// stripFrontmatter removes a leading UTF-8 BOM and YAML frontmatter block
// (--- ... ---) so content files shared with framework/web/ssg collections
// render clean.
func stripFrontmatter(body []byte) []byte {
	body = bytes.TrimPrefix(body, []byte{0xEF, 0xBB, 0xBF})
	text := strings.TrimLeft(string(body), " \t\r\n")
	if !strings.HasPrefix(text, "---") {
		return body
	}
	rest := text[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return body
	}
	after := strings.TrimLeft(rest[idx+4:], "\r\n")
	if after == "" {
		return []byte("\n")
	}
	return []byte(after)
}
