package repo

import (
	"context"
	"errors"
	"time"

	"github.com/AlexAnd012/-Effective-Mobile.git/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ListFilter фильтры и пагинация для списка.
type ListFilter struct {
	UserID      *string
	ServiceName *string
	Limit       int
	Offset      int
}

// SubscriptionRepository CRUD + сумма за период
type SubscriptionRepository interface {
	Create(ctx context.Context, s *domain.Subscription) (*domain.Subscription, error)
	Get(ctx context.Context, id string) (*domain.Subscription, error)
	List(ctx context.Context, f ListFilter) ([]domain.Subscription, error)
	Update(ctx context.Context, s *domain.Subscription) error
	Delete(ctx context.Context, id string) error
	CalcTotal(ctx context.Context, from, to time.Time, userID *string, serviceName *string) (int64, int, error)
}

type PGRepo struct{ db *pgxpool.Pool }

func NewPGRepo(db *pgxpool.Pool) *PGRepo { return &PGRepo{db: db} }

// Create Вставляем запись и сразу возвращаем все нужные поля
// Параметры передаются через плейсхолдеры
func (r *PGRepo) Create(ctx context.Context, s *domain.Subscription) (*domain.Subscription, error) {
	const q = `
insert into subscriptions(service_name, price, user_id, start_date, end_date)
values ($1,$2,$3,$4,$5)
returning id, service_name, price, user_id, start_date, end_date`
	row := r.db.QueryRow(ctx, q, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate)

	// Создаем доменную модель для бизнес-логики
	out := new(domain.Subscription)
	// scanSub хелпер для Scan
	if err := scanSub(row, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get Читаем по id
func (r *PGRepo) Get(ctx context.Context, id string) (*domain.Subscription, error) {
	const q = `
select id, service_name, price, user_id, start_date, end_date from subscriptions where id=$1`

	row := r.db.QueryRow(ctx, q, id)
	// Создаем доменную модель для бизнес-логики
	out := new(domain.Subscription)
	// scanSub хелпер для Scan
	if err := scanSub(row, out); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// маппим в ErrNotFound, чтобы HTTP-слой отдал ошибку
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return out, nil
}

func (r *PGRepo) List(ctx context.Context, f ListFilter) ([]domain.Subscription, error) {
	limit := f.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 { // чтобы не уронить БД случайным запросом
		limit = 200
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}
	// service_name
	var servName *string
	if f.ServiceName != nil && *f.ServiceName != "" {
		like := "%" + *f.ServiceName + "%"
		servName = &like
	}

	const q = `
select id, service_name, price, user_id, start_date, end_date
from subscriptions
where ($1::uuid is null or user_id = $1::uuid)
  and ($2::text is null or service_name ilike $2)
order by start_date desc, id desc
limit $3 offset $4;`

	rows, err := r.db.Query(ctx, q, f.UserID, servName, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// срез с capacity=16, чтобы уменьшить количество реаллокаций при небольшом ответе
	res := make([]domain.Subscription, 0, 16)
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, rows.Err()
}

// Update Полное обновление всех полей, если строка не найдена, возвращаем ошибку
func (r *PGRepo) Update(ctx context.Context, s *domain.Subscription) error {
	const q = `
update subscriptions
set service_name=$2, price=$3, user_id=$4, start_date=$5, end_date=$6
where id=$1`
	ct, err := r.db.Exec(ctx, q, s.ID, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Delete Удаляем по id, если строка не найдена, возвращаем ошибку
func (r *PGRepo) Delete(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `delete from subscriptions where id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PGRepo) CalcTotal(ctx context.Context, from, to time.Time, userID *string, serviceName *string) (int64, int, error) {
	// $1 user_id
	// $2 service_name
	// $3 from
	// $4 to

	const q = `
-- Если $1 или $2 = NULL, то условие даёт TRUE и не сужает выборку

with filtered as (
  select price, start_date, end_date from subscriptions
  where ($1::uuid is null or user_id = $1::uuid)
    and ($2::text is null or service_name ilike $2)
),

-- обрезаем подписку рамками периода
-- приводим дату к первому числу месяца
-- COALESCE(end_date, $4) — бессрочные подписки

clamped as (
  select
    price,
    greatest(date_trunc('month', start_date), date_trunc('month', $3::date)) as s,
    least(date_trunc('month', coalesce(end_date, $4::date)), date_trunc('month', $4::date)) as e
  from filtered
),

-- считаем месяцы включительно

counts AS (
  SELECT
    price,
    ( ((EXTRACT(YEAR FROM e) * 12 + EXTRACT(MONTH FROM e))::int)
    -  ((EXTRACT(YEAR FROM s) * 12 + EXTRACT(MONTH FROM s))::int) + 1 ) AS months
  FROM clamped
  WHERE e >= s
)

-- Соединяем , COALESCE даёт нули, если ничего не нашлось
select
  coalesce(sum((price * months)::bigint), 0) as total,
  coalesce(sum(months), 0) as months_counted
from counts;`
	var total int64
	var months int
	err := r.db.QueryRow(ctx, q, userID, serviceName, from, to).Scan(&total, &months)
	return total, months, err
}

// scanSub хелпер для Scan
func scanSub(r pgx.Row, s *domain.Subscription) error {
	return r.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate)
}
