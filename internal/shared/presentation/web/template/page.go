// Package webtemplate provides HTML rendering primitives for the WEB presentation layer.
package webtemplate

import "html/template"

type Page struct {
	Title       string
	Lang        string
	Description string
	NavKey      string
	Sidebar     Sidebar
	Entry       string
	Bootstrap   any
	Data        any
}

type Sidebar struct {
	Groups []SidebarGroup
}

type SidebarGroup struct {
	Key      string
	Label    string
	Badge    string
	Items    []SidebarItem
	Active   bool
	Expanded bool
}

type SidebarItem struct {
	Key      string
	Label    string
	Href     string
	Active   bool
	Disabled bool
}

type ViewModel struct {
	Page          Page
	Assets        EntryAssets
	BootstrapJSON template.JS
	Data          any
}

func defaultSidebar(activeKey string) Sidebar {
	sidebar := Sidebar{
		Groups: []SidebarGroup{
			{
				Key:   "identity",
				Label: "Identity",
				Badge: "Core",
				Items: []SidebarItem{
					{Key: "identity.sessions", Label: "Minhas sessoes", Href: "/sessions"},
					{Key: "identity.users", Label: "Usuarios", Href: "/users"},
					{Key: "identity.roles", Label: "Roles", Href: "/roles"},
					{Key: "identity.permissions", Label: "Permissoes", Href: "/permissions"},
				},
			},
			{
				Key:   "companies",
				Label: "Companies",
				Badge: "CRM",
				Items: []SidebarItem{
					{Key: "companies.list", Label: "Listar companies", Href: "/companies"},
					{Key: "companies.new", Label: "Nova company", Href: "/companies/new"},
				},
			},
		},
	}

	for groupIndex := range sidebar.Groups {
		group := &sidebar.Groups[groupIndex]
		for itemIndex := range group.Items {
			item := &group.Items[itemIndex]
			item.Active = item.Key == activeKey
			if item.Active {
				group.Active = true
				group.Expanded = true
			}
		}
	}

	return sidebar
}
