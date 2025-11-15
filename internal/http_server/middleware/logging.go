package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func AccessLog(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			//Обёртка вокруг ResponseWriter, чтобы
			//перехватить WriteHeader и узнать статус ответа,
			//посчитать байты записанного тела
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			//Пускаем запрос дальше
			next.ServeHTTP(ww, r)
			//Структурно логируем
			log.Info("http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"dur", time.Since(start).String(),
				// в middleware/requited.go мы добавили X-Request-ID в ответ и положили его в контекст
				"req_id", chimiddleware.GetReqID(r.Context()),
			)
		})
	}
}
