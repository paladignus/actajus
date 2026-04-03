package webtemplate

import (
	"io"
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

	if err := renderer.RenderPage(io.Discard, "pages/home", Page{
		Title:     "Portal",
		Entry:     "src/main.ts",
		Bootstrap: map[string]string{"apiBaseURL": "/api"},
	}); err != nil {
		t.Fatalf("render page: %v", err)
	}
}
