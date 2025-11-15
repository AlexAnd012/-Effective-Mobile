package dto

// TotalCostQuery — параметры запроса для подсчёта суммы.
// from/to — обязательные месяцы в формате "MM-YYYY".
type TotalCostQuery struct {
	From        string  `query:"from" example:"01-2025"` // MM-YYYY
	To          string  `query:"to" example:"12-2025"`   // MM-YYYY
	UserID      *string `query:"user_id"`
	ServiceName *string `query:"service_name"`
}

// TotalCostResponse ответ по суммарной стоимости.
type TotalCostResponse struct {
	Total         int64  `json:"total" example:"5600"`
	Currency      string `json:"currency" example:"RUB"`
	MonthsCounted int    `json:"months_counted" example:"14"`
}
