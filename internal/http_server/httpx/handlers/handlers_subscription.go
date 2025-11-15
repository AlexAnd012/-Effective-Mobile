package handlers

import (
	"encoding/json"
	"errors"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"

	"github.com/AlexAnd012/-Effective-Mobile.git/internal/dto"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/http_server/httpx"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/service"
)

// SubHandlers Структура, в которой бизнес-логика из service. Service
type SubHandlers struct{ svc *service.Service }

func NewSubHandlers(s *service.Service) *SubHandlers { return &SubHandlers{svc: s} }

// Routes передаём сервис, получаем готовый набор хендлеров
func (h *SubHandlers) Routes(r chi.Router) {
	r.Post("/", h.create)
	r.Get("/", h.list)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update) // полное обновление записи
		r.Delete("/", h.delete)
	})
}

// @Summary      Create subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input  body  dto.CreateSubscriptionRequest  true  "Данные подписки"
// @Success      201    {object}  dto.SubscriptionResponse
// @Failure      400    {object}  httpx.ErrorResponse
// @Router       /subscriptions [post]

func (h *SubHandlers) create(w http.ResponseWriter, r *http.Request) {
	// Читаем JSON тела в dto.CreateSubscriptionRequest
	var req dto.CreateSubscriptionRequest
	if err := decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	// Вызываем бизнес-логику
	out, err := h.svc.Create(r.Context(), req)
	if err != nil {
		httpx.Error(w, statusByErr(err), err)
		return
	}
	// используем обертку вокруг encoding/json
	httpx.JSON(w, http.StatusCreated, out)
}

// @Summary      Get subscription by ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "ID подписки (UUID)"
// @Success      200  {object}  dto.SubscriptionResponse
// @Failure      404  {object}  httpx.ErrorResponse
// @Router       /subscriptions/{id} [get]
func (h *SubHandlers) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// Вызываем бизнес-логику
	out, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httpx.Error(w, statusByErr(err), err)
		return
	}
	// используем обертку вокруг encoding/json
	httpx.JSON(w, http.StatusOK, out)
}

// @Summary      List subscriptions
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query  string  false  "Фильтр по UUID пользователя"
// @Param        service_name  query  string  false  "Фильтр по названию сервиса (ILIKE)"
// @Param        limit         query  int     false  "Лимит, по умолчанию 50"
// @Param        offset        query  int     false  "Смещение, по умолчанию 0"
// @Success      200  {array}   dto.SubscriptionResponse
// @Router       /subscriptions [get]
func (h *SubHandlers) list(w http.ResponseWriter, r *http.Request) {
	// Читаем Query-параметры
	q := dto.ListQuery{
		Limit:  queryInt(r, "limit", 50),
		Offset: queryInt(r, "offset", 0),
	}
	if v := r.URL.Query().Get("user_id"); v != "" {
		q.UserID = &v
	}
	if v := r.URL.Query().Get("service_name"); v != "" {
		q.ServiceName = &v
	}
	// Вызываем бизнес-логику
	out, err := h.svc.List(r.Context(), q)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	// используем обертку вокруг encoding/json
	httpx.JSON(w, http.StatusOK, out)
}

// @Summary      Update subscription
// @Description  Частичное обновление. Пустая строка в end_date снимает дату окончания.
// @Tags         subscriptions
// @Accept       json
// @Param        id     path  string                         true  "ID подписки (UUID)"
// @Param        input  body  dto.UpdateSubscriptionRequest  true  "Поля для обновления"
// @Success      204
// @Failure      400  {object}  httpx.ErrorResponse
// @Failure      404  {object}  httpx.ErrorResponse
// @Router       /subscriptions/{id} [put]
func (h *SubHandlers) update(w http.ResponseWriter, r *http.Request) {
	// Читаем JSON тела в dto.UpdateSubscriptionRequest
	id := chi.URLParam(r, "id")
	var req dto.UpdateSubscriptionRequest
	if err := decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}
	// Вызываем бизнес-логику
	if err := h.svc.Update(r.Context(), id, req); err != nil {
		httpx.Error(w, statusByErr(err), err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Delete subscription
// @Tags         subscriptions
// @Param        id   path  string  true  "ID подписки (UUID)"
// @Success      204
// @Failure      404  {object}  httpx.ErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *SubHandlers) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// Вызываем бизнес-логику
	if err := h.svc.Delete(r.Context(), id); err != nil {
		httpx.Error(w, statusByErr(err), err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Читаем JSON из тела запроса, DisallowUnknownFields защита от лишних полей/опечаток
func decode(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// Безопасно парсим целое из query с дефолтом
func queryInt(r *http.Request, name string, def int) int {
	if v := r.URL.Query().Get(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// Маппим доменные ошибки в HTTP-коды, errors. Is для работы с обернутыми ошибками
func statusByErr(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidDates), errors.Is(err, domain.ErrInvalidPrice):
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
