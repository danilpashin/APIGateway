package main

import (
	"apigateway/services/user/internal/config"
	"apigateway/services/user/internal/database"
	"apigateway/services/user/internal/handler"
	"apigateway/services/user/internal/middleware"
	"apigateway/services/user/internal/repository/postgres"
	"apigateway/services/user/internal/service"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pkg/env"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	env.LoadEnv()
	cfg := config.Load()

	if migrateCLI() {
		return
	}

	db := initDB()
	defer db.Close()

	runServer(cfg, db)
}

func runServer(cfg *config.Config, db *sql.DB) {
	r := newRouter(db)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: r,
	}
	log.Print("Server starting on", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	gracefulShutdown(&srv, db)
}

func migrateCLI() bool {
	var cmd string
	var version int

	flag.StringVar(&cmd, "cmd", "", "migration command: up, down, force, version")
	flag.IntVar(&version, "version", 0, "current version of migrations")
	flag.Parse()

	if cmd == "" {
		return false
	}

	if err := database.RunMigrations(cmd, version); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	log.Print("Migration completed")

	return true
}

func newRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	userRepo := postgres.NewUserRepository(db)
	userService := service.NewUserService(*userRepo)
	userHandler := handler.NewUserHandler(*userService)

	r.Use(middleware.LoggingMiddleware)
	r.Get("/users", userHandler.CheckHandler)

	return r
}

func initDB() *sql.DB {
	connStr := env.GetEnv("APP_DB_URL")
	if connStr == "" {
		log.Fatal("APP_DB_URL is required")
	}

	db, err := database.NewDB(connStr)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed ping to DB: %v", err)
	}
	log.Print("Database connected")

	return db
}

func gracefulShutdown(srv *http.Server, db *sql.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("DB close error: %v", err)
	}

	log.Print("Server exit")
}
