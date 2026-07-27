package main

import (
	"apigateway/services/order/internal/config"
	"apigateway/services/order/internal/database"
	"apigateway/services/order/internal/handler"
	"apigateway/services/order/internal/repository/postgres"
	"apigateway/services/order/internal/service"
	"apigateway/services/order/internal/sl"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"pkg/env"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Application struct {
	logger    *slog.Logger
	logFormat *httplog.Schema
}

func main() {
	env.LoadEnv()
	cfg := config.Load()

	logHandler, logFormat := newLogger()

	logger := slog.New(&sl.ContextHandler{Handler: logHandler})
	slog.SetDefault(logger)

	app := &Application{
		logger:    logger,
		logFormat: logFormat,
	}

	// if app.migrateCLI() {
	// 	return
	// }

	pool, err := app.initPool()
	if err != nil {
		logger.Error("Critical initialization error", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	app.runServer(cfg, pool)
}

func (app *Application) runServer(cfg *config.Config, pool *pgxpool.Pool) {
	r := app.newRouter(pool)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: r,
	}
	app.logger.Info("Server starting", "address", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("Server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	app.gracefulShutdown(&srv, pool)
}

func (app *Application) newRouter(pool *pgxpool.Pool) *chi.Mux {
	orderRepo := postgres.NewOrderRepository(pool)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	r := chi.NewRouter()
	// r.Use(mw.RequestIDMiddleware)
	r.Use(httplog.RequestLogger(app.logger, &httplog.Options{
		Schema:        app.logFormat,
		RecoverPanics: true,
	}))
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("тестовая паника бизнес-логики")
	})
	r.Get("/health", app.healthHandler(pool))
	r.Post("/orders", orderHandler.AddProductToCart)

	return r
}

// func (app *Application) migrateCLI() bool {
// 	var cmd string
// 	var version int

// 	flag.StringVar(&cmd, "cmd", "", "migration command: up, down, force, version")
// 	flag.IntVar(&version, "version", 0, "current version of migrations")
// 	flag.Parse()

// 	if cmd == "" {
// 		return false
// 	}

// 	if err := database.RunMigrations(cmd, version); err != nil {
// 		app.logger.Error("Failed to migrate", slog.Any("error", err))
// 		os.Exit(1)
// 	}
// 	app.logger.Info("Migration completed")

// 	return true
// }

func (app *Application) initPool() (*pgxpool.Pool, error) {
	connStr := env.GetEnv("ORDER_DB_URL")
	if connStr == "" {
		return nil, fmt.Errorf("ORDER_DB_URL is required")
	}

	pool, err := database.NewPgxPool(connStr, app.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB")
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed ping to DB")
	}
	app.logger.Info("Database connected")

	return pool, nil
}

func (app *Application) gracefulShutdown(srv *http.Server, pool *pgxpool.Pool) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.logger.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.logger.Error("Forced shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	pool.Close()

	app.logger.Info("Server exited gracefully")
}

func (app *Application) healthHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := pool.Ping(r.Context()); err != nil {
			app.logger.ErrorContext(r.Context(), "DB health check failed", slog.Any("error", err))
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func newLogger() (slog.Handler, *httplog.Schema) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production"
	}

	var logHandler slog.Handler
	var logFormat *httplog.Schema

	if env == "local" {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		logFormat = httplog.SchemaECS.Concise(false)
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
		logFormat = httplog.SchemaECS
	}

	return logHandler, logFormat
}
