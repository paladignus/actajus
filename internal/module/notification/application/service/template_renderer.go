// Package service
package service

import "context"

type TemplateRenderer interface {
	RenderHTML(ctx context.Context, templateName string, data any) (subject string, html string, err error)
}
