package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery middleware на паники и ошибки в хендлерах
func Recovery(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Сработает даже если внутри хендлера panic
			defer func() {
				// Ловим Panic и логируем ошибку и стек
				if rec := recover(); rec != nil {
					log.Error("panic recovered",
						"panic", rec,
						"stack", string(debug.Stack()),
						"method", r.Method,
						"path", r.URL.Path,
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
