package httperror

import "net/http"

func BadRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func Forbidden(w http.ResponseWriter) {
	http.Error(w, "forbidden", http.StatusForbidden)
}

func Render(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}
