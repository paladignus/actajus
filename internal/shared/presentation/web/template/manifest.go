// Package webtemplate provides HTML rendering primitives for the WEB presentation layer.
package webtemplate

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"
)

type Manifest map[string]ManifestEntry

type ManifestEntry struct {
	File    string   `json:"file"`
	Name    string   `json:"name"`
	Src     string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
	CSS     []string `json:"css"`
	Imports []string `json:"imports"`
}

type EntryAssets struct {
	Scripts       []string
	Styles        []string
	ModulePreload []string
}

func LoadManifest(assetFS fs.FS, manifestPath string) (Manifest, error) {
	if assetFS == nil {
		return nil, fmt.Errorf("asset fs is required")
	}
	f, err := assetFS.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("open manifest %q: %w", manifestPath, err)
	}
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read manifest %q: %w", manifestPath, err)
	}
	var manifest Manifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, fmt.Errorf("decode manifest %q: %w", manifestPath, err)
	}
	return manifest, nil
}

func (m Manifest) Resolve(entry string, publicPath string) (EntryAssets, error) {
	if len(m) == 0 {
		return EntryAssets{}, fmt.Errorf("manifest is empty")
	}
	seenScripts := map[string]struct{}{}
	seenStyles := map[string]struct{}{}
	seenPreloads := map[string]struct{}{}
	var assets EntryAssets

	var walk func(string, bool) error
	walk = func(name string, root bool) error {
		chunk, ok := m[name]
		if !ok {
			return fmt.Errorf("entry %q not found in manifest", name)
		}
		if root {
			appendUnique(&assets.Scripts, joinPublicPath(publicPath, chunk.File), seenScripts)
		} else {
			appendUnique(&assets.ModulePreload, joinPublicPath(publicPath, chunk.File), seenPreloads)
		}
		for _, css := range chunk.CSS {
			appendUnique(&assets.Styles, joinPublicPath(publicPath, css), seenStyles)
		}
		for _, imported := range chunk.Imports {
			if err := walk(imported, false); err != nil {
				return err
			}
		}
		return nil
	}

	if err := walk(entry, true); err != nil {
		return EntryAssets{}, err
	}
	sort.Strings(assets.Styles)
	sort.Strings(assets.ModulePreload)
	return assets, nil
}

func joinPublicPath(basePath, file string) string {
	file = strings.TrimSpace(file)
	if file == "" {
		if basePath == "" {
			return "/"
		}
		return path.Clean(basePath)
	}
	file = strings.TrimPrefix(path.Clean(file), "/")
	if basePath == "" || basePath == "/" {
		return "/" + file
	}
	basePath = path.Clean(basePath)
	baseName := path.Base(basePath)
	if prefix, ok := strings.CutPrefix(file, baseName+"/"); ok && baseName != "." && baseName != "/" {
		file = prefix
	}
	return basePath + "/" + file
}

func appendUnique(dst *[]string, value string, seen map[string]struct{}) {
	if _, ok := seen[value]; ok {
		return
	}
	seen[value] = struct{}{}
	*dst = append(*dst, value)
}
