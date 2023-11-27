package main

import (
	"log/slog"
	"net/http"
	"os"
	"urlC/internal/config"
	"urlC/internal/http-server/handlers/url/redirect"
	"urlC/internal/http-server/handlers/url/save"
	"urlC/internal/http-server/middleware/logger"
	"urlC/internal/lib/logger/handlers/slogcolory"
	"urlC/internal/lib/logger/sl"
	"urlC/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting URLcutter", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init db", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-cutter", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
	})

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	//TODO
	//router.Delete("/{alias}", delete.New(log, storage))

	log.Info("starting server", slog.String("adress", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupColoryLogger()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupColoryLogger() *slog.Logger {
	opts := slogcolory.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
