// Package site wires the shared WEB site bundle and presentation layer.
package site

import (
	"embed"
	"os"
	"strings"

	webapp "github.com/paladignus/actajus/internal/shared/presentation/web/app"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

//go:embed content/templates/layouts/*.gohtml content/templates/pages/*.gohtml content/templates/partials/*.gohtml content/templates/components/*.gohtml content/dist/.vite/manifest.json content/dist/assets/* content/dist/*.svg
var contentFS embed.FS

type Dependencies = webapp.Dependencies

type Module = webapp.Module

func NewModule(dep Dependencies) (Module, error) {
	devServerURL := strings.TrimSpace(os.Getenv("ACTAJUS_WEB_DEV_SERVER_URL"))
	if devServerURL == "" && strings.EqualFold(strings.TrimSpace(os.Getenv("ACTAJUS_WEB_ENV")), "development") {
		devServerURL = "http://localhost:5173"
	}

	dep.RendererSource = webtemplate.Source{
		TemplateFS: contentFS,
		TemplatePatterns: []string{
			"content/templates/layouts/*.gohtml",
			"content/templates/pages/*.gohtml",
			"content/templates/partials/*.gohtml",
			"content/templates/components/*.gohtml",
		},
		AssetFS:      contentFS,
		AssetRoot:    "content/dist",
		ManifestPath: "content/dist/.vite/manifest.json",
		PublicPath:   "/assets",
		DevServerURL: devServerURL,
	}

	return webapp.NewModule(dep)
}
