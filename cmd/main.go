package main

import (
	"GoServise/internal/config"
	"GoServise/internal/http-server/handlers/redirect"
	"GoServise/internal/http-server/handlers/refactor"
	"GoServise/internal/http-server/handlers/remove"
	"GoServise/internal/http-server/handlers/save"
	mwLogger "GoServise/internal/http-server/middleware"
	"GoServise/internal/lib/logger"
	"GoServise/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	/*ssoClient, err := ssogrpc.New(context.Background(), log,
		cfg.Client.SSO.Address, cfg.Client.SSO.Timeout, cfg.Client.SSO.RetriesCount)
	if err != nil {
		log.Error("error creating storage", logger.Err(err))
		os.Exit(1)
	}

	_ = ssoClient*/

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error creating storage", logger.Err(err))
		os.Exit(1)
	}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.NewLogger(log)) // router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// авторизация
	router.Route("/admin", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", remove.New(log, storage))
		r.Patch("/", refactor.New(log, storage))
	})
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("env", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("error starting server", logger.Err(err))
	}
	log.Error("failed server")

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
