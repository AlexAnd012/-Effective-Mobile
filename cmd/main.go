// @title           Subscriptions API
// @version         1.0
// @description     REST API для управления подписками и расчёта суммарной стоимости.
// @BasePath        /api/v1
package main

import (
	"context"
	"errors"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/repo"
	pgxboot "github.com/AlexAnd012/-Effective-Mobile.git/internal/repo/postgres"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/service"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexAnd012/-Effective-Mobile.git/internal/app"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/config"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/http_server/httpx/handlers"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/http_server/middleware"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/http_server/router"
	"github.com/AlexAnd012/-Effective-Mobile.git/internal/logging"

	_ "github.com/AlexAnd012/-Effective-Mobile.git/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// 1) Конфиг
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// 2) Логгер
	log := logging.New(cfg.Log.Level)

	// 3) БД (pgxpool)
	pool, err := pgxboot.New(
		context.Background(),
		cfg.DB.DSN,
		cfg.DB.MaxConns,
		cfg.DB.MinConns,
		cfg.DB.ConnMaxLifetime,
		cfg.DB.ConnMaxIdleTime,
	)
	if err != nil {
		log.Error("db connect failed", slog.Any("err", err))
		os.Exit(1)
	}
	defer pool.Close()

	// 4) Сервисный слой и хендлеры
	rp := repo.NewPGRepo(pool)
	svc := service.New(rp)

	health := handlers.NewHealth(pool)
	subs := handlers.NewSubHandlers(svc)

	// 5) Роутер
	root := chi.NewRouter()

	// middleware до любых маршрутов
	root.Use(middleware.RequestID())
	root.Use(middleware.Recovery(log))
	root.Use(middleware.AccessLog(log))

	// Swagger
	root.Get("/swagger/*", httpSwagger.WrapHandler)

	// редирект с корня на swagger
	root.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})

	// API с /healthz, /readyz, /api/v1/...
	api := router.New(router.Handlers{Health: health, Subs: subs})
	root.Mount("/", api)

	// root передаём в сервер
	srv := app.NewHTTPServer(cfg.Server.Addr, root, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)

	if cfg.Server.EnableSwagger {
		root.Get("/swagger/*", httpSwagger.WrapHandler)
	}

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", slog.Any("err", err))
		}
	}()
	log.Info("listening", "addr", cfg.Server.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Info("stopped")
}
