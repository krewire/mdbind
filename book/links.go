package book

import (
	"strings"
)

// link prefixes a site path with the book's base path and normalizes chapter
// links to trailing-slash form for static hosting.
func (b *Book) link(path string) string {
	path = pageLink(path)
	if b.base == "" || b.base == "/" {
		return path
	}
	return strings.TrimRight(b.base, "/") + path
}

// pageLink appends a trailing slash to chapter page paths so static hosts
// serve directory indexes without a redirect.
func pageLink(path string) string {
	if isPagePath(path) && !strings.HasSuffix(path, "/") {
		return path + "/"
	}
	return path
}

// isPagePath reports whether path points to a rendered page (a directory
// index) rather than the site root or an asset file.
func isPagePath(path string) bool {
	if path == "" || path == "/" {
		return false
	}
	last := path
	if i := strings.LastIndex(path, "/"); i >= 0 {
		last = path[i+1:]
	}
	return last != "" && !strings.Contains(last, ".")
}

// normalizeBase normalizes a base URL into leading-and-trailing-slash form
// ("/", "/guide/", …).
func normalizeBase(base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return "/"
	}
	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	if len(base) > 1 {
		base = strings.TrimRight(base, "/") + "/"
	}
	return base
}
