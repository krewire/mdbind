package book

import (
	"log/slog"
	"net/http"

	"github.com/krewire/framework/web"
)

// Handler returns an http.Handler that serves the book's pages and assets
// through the web router.
func (b *Book) Handler() (*web.Router, error) {
	pages, err := b.Pages()
	if err != nil {
		return nil, err
	}
	byPath := make(map[string]web.Page, len(pages))
	for _, p := range pages {
		byPath[p.Path] = p
	}
	r := web.NewRouter()
	r.Get("/", renderPage(byPath["/"]))
	for _, p := range pages {
		if p.Path != "/" {
			r.Get(p.Path, renderPage(p))
		}
	}
	r.Get("/assets/mdbind.css", func(w http.ResponseWriter, _ *http.Request, _ web.Params) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write([]byte(Assets()["assets/mdbind.css"]))
	})
	r.NotFound = func(w http.ResponseWriter, _ *http.Request, _ web.Params) {
		web.HTML(w, http.StatusNotFound, "404 — page not found")
	}
	return r, nil
}

// renderPage adapts a web.Page into a route handler.
func renderPage(p web.Page) web.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request, _ web.Params) {
		_ = p.Render(w)
	}
}

// Serve serves the book over HTTP on addr, blocking until the server stops.
func Serve(b *Book, addr string) error {
	h, err := b.Handler()
	if err != nil {
		return err
	}
	slog.Info("serving book", "addr", addr)
	return http.ListenAndServe(addr, h)
}
