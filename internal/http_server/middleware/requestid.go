package middleware

import (
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// RequestID добавляет X-Request-ID в ответ и кладёт его в контекст
// Если заголовок X-Request-ID уже пришёл то chi создаёт новый
func RequestID() func(http.Handler) http.Handler {
	return chimw.RequestID
}
