package book

import (
	"log/slog"
	"net/http"
)

// Handler returns an http.Handler that serves the book's pages and assets
// through a stdlib ServeMux. It is intentionally not dependent on
// framework/web so the book can be served without pulling the framework.
func (b *Book) Handler() (http.Handler, error) {
	pages, err := b.Pages()
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	for _, p := range pages {
		pp := p
		route := pp.Path
		if b.mount != "" {
			route = "/" + b.mount + pp.Path
		}
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != route {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_ = pp.Render(w)
		})
	}
	css := Assets()["assets/mdbind.css"]
	mux.HandleFunc(b.mount+"/assets/mdbind.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write([]byte(css))
	})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, pattern := mux.Handler(r)
		if pattern == "" {
			http.Error(w, "404 — page not found", http.StatusNotFound)
			return
		}
		h.ServeHTTP(w, r)
	})
	return handler, nil
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
