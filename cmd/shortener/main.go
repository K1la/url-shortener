package main

import (
	"github.com/K1la/url-shortener/internal/api/handler"
	"github.com/K1la/url-shortener/internal/api/router"
	"github.com/K1la/url-shortener/internal/api/server"
	"github.com/K1la/url-shortener/internal/cache"
	"github.com/K1la/url-shortener/internal/config"
	"github.com/K1la/url-shortener/internal/repository"
	"github.com/K1la/url-shortener/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/zlog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	zlog.Init()

	cfg := config.Init()

	db := repository.NewDB(cfg)
	repo := repository.New(db)
	ch := cache.New(cfg.Redis)
	srvc := service.New(repo, ch)

	// sig channel to handle SIGINT and SIGTERM for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	val := validator.New()
	hndlr := handler.New(srvc, val)
	r := router.New(hndlr)
	s := server.New(cfg.HTTPServer.Address, r)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to start server")
		}
		zlog.Logger.Info().Msg("successfully started server on " + cfg.HTTPServer.Address)
	}()

	// Блокируем main горутину до получения сигнала завершения
	<-sigChan
	zlog.Logger.Info().Msg("Shutting down gracefully...")
}
