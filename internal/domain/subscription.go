package domain

import "time"

// Subscription доменная модель подписки (структура для бизнес-логики)
// Даты храним как первое число месяца в UTC.
type Subscription struct {
	ID          string
	ServiceName string
	Price       int        // рубли, целое
	UserID      string     // UUID
	StartDate   time.Time  // 1-е число месяца, UTC
	EndDate     *time.Time // nil = бессрочная
}

// MonthStart нормализует дату к первому дню месяца (00:00:00 UTC)
func MonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// Normalize приводит StartDate и  EndDate к первому дню месяца (00:00:00 UTC)
func (s *Subscription) Normalize() {
	s.StartDate = MonthStart(s.StartDate)
	if s.EndDate != nil {
		e := MonthStart(*s.EndDate)
		s.EndDate = &e
	}
}

// Validate проверяет базовые инварианты модели и используем ошибки из файла /internal/domain/errors.go
// Цена должна быть > 0 и дата начала должна быть до даты конца
func (s *Subscription) Validate() error {
	if s.Price <= 0 {
		return ErrInvalidPrice
	}
	if s.EndDate != nil && s.EndDate.Before(s.StartDate) {
		return ErrInvalidDates
	}
	return nil
}
