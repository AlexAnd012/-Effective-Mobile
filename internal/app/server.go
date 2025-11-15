package app

import (
	"context"
	"net/http"
	"time"
)

type HTTPServer struct {
	srv *http.Server
}

// NewHTTPServer обертка над стандартным http. Server
func NewHTTPServer(addr string, handler http.Handler, readTO, writeTO time.Duration) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  readTO,
			WriteTimeout: writeTO,
		},
	}
}

// Start запускаем HTTP-сервер, блокируем горутину, пока сервер работает
func (s *HTTPServer) Start() error {
	return s.srv.ListenAndServe()
}

// Shutdown перестаём принимать новые коннекты, ждём завершения текущих запросов
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
