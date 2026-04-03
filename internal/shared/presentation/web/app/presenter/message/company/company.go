package companymessage

import webflash "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/flash"

func DeleteFailed() webflash.Message {
	return webflash.Message{Source: "companies.delete", Code: "companies_delete_failed", Kind: "error", Message: "Nao foi possivel excluir a company."}
}

func DeleteSucceeded() webflash.Message {
	return webflash.Message{Source: "companies.delete", Code: "companies_delete_succeeded", Kind: "success", Message: "Company excluida com sucesso."}
}
