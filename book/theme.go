package book

import (
	"html/template"
	"reflect"
	"strconv"
	"strings"
)

const defaultThemeKey = "krewire-theme"

const ThemeModeVarsCSS = `:root { --show-sun: block; --show-moon: none; }
:root[data-theme="dark"] { --show-sun: none; --show-moon: block; }
@media (prefers-color-scheme: dark) { :root:not([data-theme]) { --show-sun: none; --show-moon: block; } }`

const ThemeToggleCSS = `.theme-toggle { border: 1px solid var(--neutral); border-radius: 999px; background: transparent; color: var(--ink); padding: .35rem .5rem; line-height: 0; cursor: pointer; transition: border-color .2s ease, color .2s ease; }
.theme-toggle:hover { border-color: var(--primary); color: var(--primary); }
.theme-toggle svg { width: 1.05rem; height: 1.05rem; display: block; }
.theme-toggle .icon-sun { display: var(--show-sun); }
.theme-toggle .icon-moon { display: var(--show-moon); }`

type Color string

type Palette struct {
	Base1            Color `css:"base-1"`
	Base1Content     Color `css:"base-1-content"`
	Base2            Color `css:"base-2"`
	Base2Content     Color `css:"base-2-content"`
	Base3            Color `css:"base-3"`
	Base3Content     Color `css:"base-3-content"`
	Primary          Color `css:"primary"`
	PrimaryContent   Color `css:"primary-content"`
	Secondary        Color `css:"secondary"`
	SecondaryContent Color `css:"secondary-content"`
	Accent           Color `css:"accent"`
	AccentContent    Color `css:"accent-content"`
	Ghost            Color `css:"ghost"`
	GhostContent     Color `css:"ghost-content"`
	Neutral          Color `css:"neutral"`
	NeutralContent   Color `css:"neutral-content"`
	Success          Color `css:"success"`
	SuccessContent   Color `css:"success-content"`
	Info             Color `css:"info"`
	InfoContent      Color `css:"info-content"`
	Warning          Color `css:"warning"`
	WarningContent   Color `css:"warning-content"`
	Error            Color `css:"error"`
	ErrorContent     Color `css:"error-content"`
}

var DefaultLightPalette = Palette{
	Base1:            "#fdfdfb",
	Base1Content:     "#1a1a1a",
	Base2:            "#f5f2eb",
	Base2Content:     "#666",
	Base3:            "#efece4",
	Base3Content:     "#4a4a4a",
	Primary:          "#00c853",
	PrimaryContent:   "#ffffff",
	Secondary:        "#4a8f7a",
	SecondaryContent: "#ffffff",
	Accent:           "#f59e0b",
	AccentContent:    "#ffffff",
	Ghost:            "#e9e6dc",
	GhostContent:     "#666",
	Neutral:          "#e3e0d8",
	NeutralContent:   "#666",
	Success:          "#15803d",
	SuccessContent:   "#ffffff",
	Info:             "#0d9488",
	InfoContent:      "#ffffff",
	Warning:          "#f59e0b",
	WarningContent:   "#ffffff",
	Error:            "#b91c1c",
	ErrorContent:     "#ffffff",
}

var DefaultDarkPalette = Palette{
	Base1:            "#17171a",
	Base1Content:     "#e6e4dd",
	Base2:            "#202024",
	Base2Content:     "#9a968b",
	Base3:            "#2a2a30",
	Base3Content:     "#c6c2b8",
	Primary:          "#5cff8e",
	PrimaryContent:   "#06160d",
	Secondary:        "#7fc99d",
	SecondaryContent: "#07130b",
	Accent:           "#ffb020",
	AccentContent:    "#241a00",
	Ghost:            "#26262b",
	GhostContent:     "#9a968b",
	Neutral:          "#2e2e34",
	NeutralContent:   "#9a968b",
	Success:          "#4ade80",
	SuccessContent:   "#0b1f10",
	Info:             "#2dd4bf",
	InfoContent:      "#062a23",
	Warning:          "#ffb020",
	WarningContent:   "#241a00",
	Error:            "#f87171",
	ErrorContent:     "#260a0a",
}

func (p Palette) cssVars(defaults Palette) string {
	var b strings.Builder
	kiw, rdd := reflect.ValueOf(p), reflect.ValueOf(defaults)
	for i := 0; i < kiw.NumField(); i++ {
		v := kiw.Field(i).String()
		if v == "" {
			v = rdd.Field(i).String()
		}
		if v == "" {
			continue
		}
		b.WriteString("--")
		b.WriteString(kiw.Type().Field(i).Tag.Get("css"))
		b.WriteString(": ")
		b.WriteString(v)
		b.WriteString(";")
	}
	return b.String()
}

type Theme struct {
	StorageKey string
	Default    string
	Light      Palette
	Dark       Palette
}

func (t Theme) storageKey() string {
	if t.StorageKey == "" {
		return defaultThemeKey
	}
	return t.StorageKey
}

func (t Theme) defaultTheme() string {
	switch t.Default {
	case "light", "dark":
		return t.Default
	default:
		return "auto"
	}
}

func (t Theme) Style() template.CSS {
	light := ":root{" + t.Light.cssVars(DefaultLightPalette) + "}"
	dark := ":root[data-theme=\"dark\"]{" + t.Dark.cssVars(DefaultDarkPalette) + "}"
	fallback := "@media (prefers-color-scheme: dark){:root:not([data-theme]){" +
		t.Dark.cssVars(DefaultDarkPalette) + "}}"
	return template.CSS("<style>" + light + dark + fallback + "</style>")
}

func (t Theme) Script() template.HTML {
	return template.HTML(string(t.Style()) + "<script>" + themeJS(t.storageKey(), t.defaultTheme()) + "</script>")
}

func (t Theme) Button() template.HTML {
	return template.HTML(`<button type="button" class="theme-toggle" data-theme-toggle aria-label="Toggle light/dark theme" title="Toggle theme"><svg class="icon-sun" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg><svg class="icon-moon" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg></button>`)
}

func themeJS(key, def string) string {
	var b strings.Builder
	b.WriteString(`(function(){var k=` + strconv.Quote(key) + `,d=` + strconv.Quote(def) + `,t=d;try{t=localStorage.getItem(k)||d}catch(e){}var m=window.matchMedia('(prefers-color-scheme: dark)');function apply(){var dark=t==='dark'||(t==='auto'&&m.matches);var r=document.documentElement;r.dataset.theme=dark?'dark':'light';r.style.colorScheme=dark?'dark':'light';var meta=document.querySelector('meta[name="color-scheme"]');if(meta)meta.content=dark?'dark':'light'}apply();if(m.addEventListener){m.addEventListener('change',function(){if(t==='auto')apply()})}window.krewireTheme={get:function(){return t},set:function(x){t=x;try{localStorage.setItem(k,x)}catch(e){}apply()},toggle:function(){var dark=t==='dark'||(t==='auto'&&m.matches);window.krewireTheme.set(dark?'light':'dark')}};document.addEventListener('click',function(e){var el=e.target&&e.target.closest?e.target.closest('[data-theme-toggle]'):null;if(el)window.krewireTheme.toggle()})})();`)
	return b.String()
}
