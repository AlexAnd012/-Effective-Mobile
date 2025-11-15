package router

import (
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/http_server/httpx/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handlers контейнер для роутера
type Handlers struct {
	Health *handlers.HealthHandler
	Subs   *handlers.SubHandlers
}

func New(d Handlers, mws ...func(http.Handler) http.Handler) *chi.Mux {
	// Применяем все middleware
	r := chi.NewRouter()
	for _, mw := range mws {
		if mw != nil {
			r.Use(mw)
		}
	}
	// Здоровье и БД
	r.Get("/healthz", d.Health.Liveness)
	r.Get("/readyz", d.Health.Readiness)
	// создаем дочерний роутер с префиксом /api/v1
	r.Route("/api/v1", func(r chi.Router) {
		// Регистрируем пути в handlers/handlers_subscription
		r.Route("/subscriptions", d.Subs.Routes)
		// Ручка расчёта суммы
		r.Get("/cost/total", d.Subs.TotalCost)
	})
	return r
}
