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

func RegisterPage(data map[string]any) webtemplate.Page {
	if data == nil {
		data = map[string]any{}
	}
	if _, ok := data["GenderID"]; !ok {
		data["GenderID"] = uint(0)
	}
	return webtemplate.Page{
		Title:       "Criar conta",
		Description: "Cadastro WEB do Actajus.",
		NavKey:      "register",
		Entry:       "src/main.ts",
		Data:        data,
	}
}

func VerifyEmailPage(message string, success bool) webtemplate.Page {
	return webtemplate.Page{
		Title:       "Validar email",
		Description: "Confirmacao de email do Actajus.",
		NavKey:      "verify-email",
		Entry:       "src/main.ts",
		Data: map[string]any{
			"Message": message,
			"Success": success,
		},
	}
}

func VerificationPendingPage(email, message string) webtemplate.Page {
	return webtemplate.Page{
		Title:       "Verificacao pendente",
		Description: "Verificacao de email pendente no Actajus.",
		NavKey:      "verification-pending",
		Entry:       "src/main.ts",
		Data: map[string]any{
			"Email":   email,
			"Message": message,
		},
	}
}
