package main

import (
	"GoServise/internal/config"
	"GoServise/internal/lib/logger/sl"
	"GoServise/internal/storage/sqlite"
	"log/slog"
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

	// log = log.With(slog.String("env", cfg.Env)) для того, чтобы каждый лог был с таким параметром

	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug logging enabled")

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("error creating storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router: chi, "chi render"

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
