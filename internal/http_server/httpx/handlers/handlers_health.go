package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/AlexAnd012/-Effective-Mobile.git/internal/http_server/httpx"
)

// Pinger Интерфейс для pgxpool.Pool, для пинга BD
type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	DB    Pinger    // Для проверки готовности
	Start time.Time // Момент старта сервиса для расчёта аптайма
}

func NewHealth(db Pinger) *HealthHandler {
	return &HealthHandler{
		DB:    db,
		Start: time.Now().UTC(),
	}
}

// Liveness GET /healthz процесс жив
// параметр r типа *http. Request не используется, но он требуется для роутера
func (h *HealthHandler) Liveness(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"uptime": time.Since(h.Start).String(),
	})
}

// Readiness GET /readyz сервис готов и БД доступна
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		httpx.JSON(w, http.StatusServiceUnavailable, map[string]any{"status": "no db"})
		return
	}
	// Ставим предел ожидания 300мс и ждем бд, иначе отдаем ошибку
	ctx, cancel := context.WithTimeout(r.Context(), 300*time.Millisecond)
	// Освобождаем ресурсы таймера/контекста
	defer cancel()

	if err := h.DB.Ping(ctx); err != nil {
		httpx.JSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "db down",
			"error":  err.Error(),
		})
		return
	}
	// используем обертку вокруг encoding/json
	httpx.JSON(w, http.StatusOK, map[string]any{"status": "ready"})
}
