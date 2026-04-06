// Package webtemplate provides HTML rendering primitives for the WEB presentation layer.
package webtemplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

type Renderer struct {
	templates  *template.Template
	manifest   Manifest
	publicPath string
	assetFS    fs.FS
	devServer  string
}

type Source struct {
	TemplateFS       fs.FS
	TemplatePatterns []string
	AssetFS          fs.FS
	AssetRoot        string
	ManifestPath     string
	PublicPath       string
	DevServerURL     string
}

func NewRenderer(source Source) (*Renderer, error) {
	if source.TemplateFS == nil {
		return nil, fmt.Errorf("template fs is required")
	}
	if len(source.TemplatePatterns) == 0 {
		return nil, fmt.Errorf("template patterns are required")
	}
	funcs := template.FuncMap{
		"safeHTML": func(v string) template.HTML { return template.HTML(v) },
	}
	tpl, err := template.New("web").Funcs(funcs).ParseFS(source.TemplateFS, source.TemplatePatterns...)
	if err != nil {
		return nil, fmt.Errorf("parse web templates: %w", err)
	}
	publicPath := strings.TrimRight(source.PublicPath, "/")
	if publicPath == "" {
		publicPath = "/assets"
	}

	renderer := &Renderer{
		templates:  tpl,
		publicPath: publicPath,
		assetFS:    source.AssetFS,
		devServer:  strings.TrimRight(strings.TrimSpace(source.DevServerURL), "/"),
	}
	if renderer.devServer == "" && source.AssetFS != nil && source.ManifestPath != "" {
		manifest, err := LoadManifest(source.AssetFS, source.ManifestPath)
		if err != nil {
			return nil, err
		}
		renderer.manifest = manifest
	}
	return renderer, nil
}

func (r *Renderer) RenderPage(w io.Writer, templateName string, page Page) error {
	assets, err := r.resolveAssets(page.Entry)
	if err != nil {
		return fmt.Errorf("resolve entry assets: %w", err)
	}
	bootstrapJSON, err := marshalBootstrap(page.Bootstrap)
	if err != nil {
		return fmt.Errorf("marshal bootstrap payload: %w", err)
	}
	vm := ViewModel{
		Page:          withPageDefaults(page),
		Assets:        assets,
		BootstrapJSON: bootstrapJSON,
		Data:          page.Data,
	}
	if err := r.templates.ExecuteTemplate(w, templateName, vm); err != nil {
		return fmt.Errorf("render template %q: %w", templateName, err)
	}
	return nil
}

func (r *Renderer) resolveAssets(entry string) (EntryAssets, error) {
	if r.devServer != "" {
		return EntryAssets{
			Scripts: []string{
				r.devServer + "/@vite/client",
				r.devServer + "/" + strings.TrimLeft(entry, "/"),
			},
		}, nil
	}
	return r.manifest.Resolve(entry, r.publicPath)
}

func (r *Renderer) RenderHTTP(w http.ResponseWriter, status int, templateName string, page Page) error {
	var buf bytes.Buffer
	if err := r.RenderPage(&buf, templateName, page); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, err := io.Copy(w, &buf)
	return err
}

func (r *Renderer) AssetsHandler(root string) (http.Handler, error) {
	if r.assetFS == nil {
		return nil, fmt.Errorf("asset fs is required")
	}
	servedFS := r.assetFS
	if root != "" && root != "." {
		var err error
		servedFS, err = fs.Sub(r.assetFS, root)
		if err != nil {
			return nil, fmt.Errorf("sub asset fs %q: %w", root, err)
		}
	}
	return http.FileServerFS(servedFS), nil
}

func withPageDefaults(page Page) Page {
	if strings.TrimSpace(page.Title) == "" {
		page.Title = "Actajus"
	}
	if strings.TrimSpace(page.Lang) == "" {
		page.Lang = "pt-BR"
	}
	if len(page.Sidebar.Groups) == 0 && strings.TrimSpace(page.NavKey) != "" && page.NavKey != "login" {
		page.Sidebar = defaultSidebar(page.NavKey)
	}
	return page
}

func marshalBootstrap(value any) (template.JS, error) {
	if value == nil {
		return template.JS("{}"), nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return template.JS(raw), nil
}
