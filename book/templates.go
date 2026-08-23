package book

import (
	"embed"
	"html/template"

	"github.com/krewire/framework/ui"
)

// content embeds the site templates and shipped assets.
//
//go:embed templates assets
var content embed.FS

// parseTemplates loads the page templates from the embedded filesystem.
func parseTemplates() (*template.Template, error) {
	return template.ParseFS(content, "templates/*.tmpl")
}

// Assets returns the site's shipped assets: name to file body. The theme-mode
// variables and toggle button styles come from the ui framework and are
// appended to the reader stylesheet so consumers never hardcode them.
func Assets() map[string]string {
	body, err := content.ReadFile("assets/mdbind.css")
	if err != nil {
		return map[string]string{"assets/mdbind.css": "/* mdbind default stylesheet */"}
	}
	css := string(body) + "\n" + ui.ThemeModeVarsCSS + "\n" + ui.ThemeToggleCSS
	return map[string]string{"assets/mdbind.css": css}
}
