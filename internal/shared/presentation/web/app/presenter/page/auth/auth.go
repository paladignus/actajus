package authpage

import webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"

func LoginPage(email, message string) webtemplate.Page {
	return webtemplate.Page{
		Title:       "Entrar",
		Description: "Acesso WEB do Actajus.",
		NavKey:      "login",
		Entry:       "src/main.ts",
		Data: map[string]any{
			"Email": email,
			"Error": message,
		},
	}
}
