package dto

// CreateSubscriptionRequest тело запроса на создание подписки
// example чтобы на swagger были примеры
type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" example:"Yandex Plus"`
	Price       int     `json:"price" example:"400"`
	UserID      string  `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string  `json:"start_date" example:"07-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"01-2026"`
}

// UpdateSubscriptionRequest частичное обновление
// Все поля опциональны, пустая строка в EndDate удаляет дату окончания
type UpdateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

// SubscriptionResponse объект, который отдаем наружу
type SubscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

// ListQuery параметры фильтрации/пагинации для списка
// Парсим их из r.URL.Query() в хендлере
type ListQuery struct {
	UserID      *string `query:"user_id"`
	ServiceName *string `query:"service_name"`
	Limit       int     `query:"limit"`
	Offset      int     `query:"offset"`
}
