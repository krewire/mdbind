package book

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// gold parses and renders Markdown to HTML.
var gold = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// absLinkRe matches absolute href/src references like href="/getting-started".
var absLinkRe = regexp.MustCompile(`(href|src)="/([^"]*)"`)

// renderMarkdown converts a Markdown body to an HTML fragment, prefixing
// absolute links with the site's base path ("/guide/" etc.) when it is not the
// site root.
func renderMarkdown(body []byte, base string) (string, error) {
	var buf bytes.Buffer
	if err := gold.Convert(body, &buf); err != nil {
		return "", err
	}
	return prefixLinks(buf.String(), base), nil
}

// prefixLinks rewrites absolute href/src links so they resolve under the given
// base path, e.g. "/getting-started" -> "/guide/getting-started/". Page links
// gain a trailing slash so static hosts serve directory indexes without a
// redirect.
func prefixLinks(html, base string) string {
	prefix := ""
	if base != "" && base != "/" {
		prefix = strings.Trim(base, "/")
	}
	return absLinkRe.ReplaceAllStringFunc(html, func(m string) string {
		sub := absLinkRe.FindStringSubmatch(m)
		attr, rest := sub[1], sub[2]
		if strings.HasPrefix(rest, "/") || (prefix != "" && (strings.HasPrefix(rest, prefix+"/") || rest == prefix)) {
			return m
		}
		out := attr + `="/`
		if prefix != "" {
			out += prefix + "/"
		}
		out += rest
		if isPagePath("/"+rest) && !strings.HasSuffix(rest, "/") {
			out += "/"
		}
		return out + `"`
	})
}
