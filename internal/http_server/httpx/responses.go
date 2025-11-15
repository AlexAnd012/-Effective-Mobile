package httpx

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse Формат ошибок для клиента
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// JSON Обертка для json
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Error собираем ErrorResponse и вызываем JSON
func Error(w http.ResponseWriter, status int, err error) {
	JSON(w, status, ErrorResponse{Error: http.StatusText(status), Message: err.Error()})
}
