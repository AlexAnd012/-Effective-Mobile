package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexAnd012/-Effective-Mobile.git/internal/domain"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/dto"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/repo"

	"github.com/google/uuid"
)

// Service обёртка над репозиторием с бизнес-правилами
// проверяем входные данные
// парсим/нормализуем даты
// маппим доменные модели в DTO
// отдаём понятные ошибки наверх
type Service struct{ repo repo.SubscriptionRepository }

func New(r repo.SubscriptionRepository) *Service { return &Service{repo: r} }

// Create Валидируем поля, парсим даты, запсиываем в бд
func (s *Service) Create(ctx context.Context, in dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	// Валидация и парсинг
	if in.ServiceName == "" {
		return nil, fmt.Errorf("service_name is required")
	}
	if in.Price <= 0 {
		return nil, domain.ErrInvalidPrice
	}
	if _, err := uuid.Parse(in.UserID); err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}
	start, err := parseMonth(in.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	var end *time.Time
	if in.EndDate != nil && *in.EndDate != "" {
		e, err := parseMonth(*in.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		if e.Before(start) {
			return nil, domain.ErrInvalidDates
		}
		end = &e
	}

	// собираем domain. Subscription и вызываем repo. Create
	created, err := s.repo.Create(ctx, &domain.Subscription{
		ServiceName: in.ServiceName, Price: in.Price, UserID: in.UserID,
		StartDate: start, EndDate: end,
	})
	if err != nil {
		return nil, err
	}
	return toDTO(created), nil
}

// Get Вызываем repo. Get, преобразуем доменную модель в DTO
func (s *Service) Get(ctx context.Context, id string) (*dto.SubscriptionResponse, error) {
	out, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDTO(out), nil
}

// List Пробрасываем фильтры/лимиты в repo.List через repo.ListFilter.
// Переводим []domain.Subscription в []dto.SubscriptionResponse.
func (s *Service) List(ctx context.Context, q dto.ListQuery) ([]dto.SubscriptionResponse, error) {
	items, err := s.repo.List(ctx, repo.ListFilter{
		UserID: q.UserID, ServiceName: q.ServiceName, Limit: q.Limit, Offset: q.Offset,
	})
	if err != nil {
		return nil, err
	}
	res := make([]dto.SubscriptionResponse, 0, len(items))
	for i := range items {
		res = append(res, *toDTO(&items[i]))
	}
	return res, nil
}

// Update полная замена put, всё валидируем с нуля, формируем полную доменную модель и сохраняем
func (s *Service) Update(ctx context.Context, id string, in dto.UpdateSubscriptionRequest) error {

	if in.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}
	if in.Price <= 0 {
		return domain.ErrInvalidPrice
	}
	if _, err := uuid.Parse(in.UserID); err != nil {
		return fmt.Errorf("invalid user_id: %w", err)
	}

	start, err := parseMonth(in.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date: %w", err)
	}

	var end *time.Time
	if in.EndDate != nil { // nil = бессрочно
		if *in.EndDate == "" {
			return fmt.Errorf("end_date must be null or 'YYYY-MM'")
		}
		e, err := parseMonth(*in.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date: %w", err)
		}
		end = &e
	}

	if end != nil && end.Before(start) {
		return domain.ErrInvalidDates
	}

	return s.repo.Update(ctx, &domain.Subscription{
		ID:          id,
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      in.UserID,
		StartDate:   start,
		EndDate:     end,
	})
}

// Delete Выполняем repo.Delete
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// TotalCost  Парсим from и to как месяцы через parseMonth
// Выполняем repo.CalcTotal
func (s *Service) TotalCost(ctx context.Context, q dto.TotalCostQuery) (dto.TotalCostResponse, error) {
	from, err := parseMonth(q.From)
	if err != nil {
		return dto.TotalCostResponse{}, fmt.Errorf("invalid from: %w", err)
	}
	to, err := parseMonth(q.To)
	if err != nil {
		return dto.TotalCostResponse{}, fmt.Errorf("invalid to: %w", err)
	}
	total, months, err := s.repo.CalcTotal(ctx, from, to, q.UserID, q.ServiceName)
	if err != nil {
		return dto.TotalCostResponse{}, err
	}
	// Возвращаем DTO
	return dto.TotalCostResponse{Total: total, Currency: "RUB", MonthsCounted: months}, nil
}

// Вспомогательные функции

// Возвращаем время к первому числу месяца (UTC)
func parseMonth(s string) (time.Time, error) {
	if t, err := time.ParseInLocation("01-2006", s, time.UTC); err == nil {
		return domain.MonthStart(t), nil
	}
	return time.Time{}, fmt.Errorf("expected MM-YYYY")
}

// toDTO маппим доменную модель в ответ и форматируем месяцы
func toDTO(s *domain.Subscription) *dto.SubscriptionResponse {
	var end *string
	if s.EndDate != nil {
		v := s.EndDate.Format("01-2006")
		end = &v
	}
	return &dto.SubscriptionResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   s.StartDate.Format("01-2006"),
		EndDate:     end,
	}
}
