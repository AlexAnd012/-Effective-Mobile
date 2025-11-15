package domain

import "errors"

var (
	// ErrNotFound запись не найдена.
	ErrNotFound = errors.New("subscription not found")

	// ErrInvalidDates дата окончания раньше даты начала.
	ErrInvalidDates = errors.New("end_date before start_date")

	// ErrInvalidPrice цена должна быть > 0.
	ErrInvalidPrice = errors.New("price must be > 0")
)
