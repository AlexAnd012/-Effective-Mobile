package handlers

import (
	"net/http"

	"github.com/AlexAnd012/-Effective-Mobile.git/internal/dto"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/http_server/httpx"
)

// @Summary      Total cost
// @Description  Сумма стоимостей всех подписок за период (включительно), с фильтрами
// @Tags         cost
// @Produce      json
// @Param        from          query  string  true   "Начало периода, MM-YYYY"
// @Param        to            query  string  true   "Конец периода, MM-YYYY"
// @Param        user_id       query  string  false  "Фильтр по UUID пользователя"
// @Param        service_name  query  string  false  "Фильтр по названию сервиса"
// @Success      200  {object}  dto.TotalCostResponse
// @Failure      400  {object}  httpx.ErrorResponse
// @Router       /cost/total [get]
func (h *SubHandlers) TotalCost(w http.ResponseWriter, r *http.Request) {
	// Разбор query-параметров
	// from/to обязательны, user_id/service_name опциональны
	q := dto.TotalCostQuery{
		From: r.URL.Query().Get("from"),
		To:   r.URL.Query().Get("to"),
	}
	if v := r.URL.Query().Get("user_id"); v != "" {
		q.UserID = &v
	}
	if v := r.URL.Query().Get("service_name"); v != "" {
		q.ServiceName = &v
	}

	// Вызов бизнес-логики из service\subscription и ответ
	res, err := h.svc.TotalCost(r.Context(), q)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	// используем обертку вокруг encoding/json
	httpx.JSON(w, http.StatusOK, res)
}
