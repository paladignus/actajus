package webtemplate

import (
	"io"
	"strings"
	"testing"
	"testing/fstest"
)

func TestRendererRenderPage(t *testing.T) {
	templateFS := fstest.MapFS{
		"templates/base.gohtml": {Data: []byte(`{{define "pages/home"}}<html lang="{{.Page.Lang}}"><head><title>{{.Page.Title}}</title>{{range .Assets.Styles}}<link rel="stylesheet" href="{{.}}">{{end}}</head><body><script>window.__BOOTSTRAP__={{.BootstrapJSON}};</script>{{range .Assets.Scripts}}<script type="module" src="{{.}}"></script>{{end}}</body></html>{{end}}`)},
	}
	assetFS := fstest.MapFS{
		"dist/.vite/manifest.json": {Data: []byte(`{"src/main.ts":{"file":"assets/app.js","isEntry":true,"css":["assets/app.css"]}}`)},
	}

	renderer, err := NewRenderer(Source{
		TemplateFS:       templateFS,
		TemplatePatterns: []string{"templates/*.gohtml"},
		AssetFS:          assetFS,
		ManifestPath:     "dist/.vite/manifest.json",
		PublicPath:       "/assets",
	})
	if err != nil {
		t.Fatalf("new renderer: %v", err)
	}

	var out strings.Builder
	if err := renderer.RenderPage(&out, "pages/home", Page{
		Title:     "Portal",
		Entry:     "src/main.ts",
		Bootstrap: map[string]string{"apiBaseURL": "/api"},
	}); err != nil {
		t.Fatalf("render page: %v", err)
	}
	rendered := out.String()
	if !strings.Contains(rendered, `href="/assets/app.css"`) {
		t.Fatalf("expected css asset path without duplication, got %q", rendered)
	}
	if !strings.Contains(rendered, `src="/assets/app.js"`) {
		t.Fatalf("expected js asset path without duplication, got %q", rendered)
	}
}

func TestRendererRenderPageInDevMode(t *testing.T) {
	templateFS := fstest.MapFS{
		"templates/base.gohtml": {Data: []byte(`{{define "pages/home"}}{{range .Assets.Scripts}}<script type="module" src="{{.}}"></script>{{end}}{{end}}`)},
	}

	renderer, err := NewRenderer(Source{
		TemplateFS:       templateFS,
		TemplatePatterns: []string{"templates/*.gohtml"},
		DevServerURL:     "http://localhost:5173",
	})
	if err != nil {
		t.Fatalf("new renderer: %v", err)
	}

	if err := renderer.RenderPage(io.Discard, "pages/home", Page{
		Entry: "src/main.ts",
	}); err != nil {
		t.Fatalf("render page in dev mode: %v", err)
	}
}
