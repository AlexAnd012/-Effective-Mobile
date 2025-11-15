package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New Создаем пул, на входе ctx, чтобы можно было отменить создание пула
func New(ctx context.Context, dsn string, maxConns, minConns int32, maxLife, maxIdle time.Duration) (*pgxpool.Pool, error) {
	// Парсим DSN в *pgxpool.Config
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Переопределяем дефолты, если передали значения
	if maxConns > 0 {
		cfg.MaxConns = maxConns
	}
	if minConns >= 0 {
		cfg.MinConns = minConns
	}
	if maxLife > 0 {
		cfg.MaxConnLifetime = maxLife
	}
	if maxIdle > 0 {
		cfg.MaxConnIdleTime = maxIdle
	}

	// Создаём пул с нашей конфигурацией, ctx позволяет оборвать создание
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Проверяем, можем ли установить соединение
	pctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := pool.Ping(pctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
