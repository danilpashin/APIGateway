package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pkg/env"
	"product/internal/config"
	"product/internal/database"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	env.LoadEnv()
	cfg := config.Load()
	db := initDB()

	runServer(cfg, db)
}

func runServer(cfg *config.Config, db *sql.DB) {
	r := newRouter()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: r,
	}
	log.Print("Server starting on", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed: ", err)
		}
	}()

	gracefulShutdown(&srv, db)
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/product", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Write([]byte("Checking graceful shutdown"))
	}))

	return r
}

func gracefulShutdown(srv *http.Server, db *sql.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("DB close error: %v", err)
	}

	log.Print("Server exit")
}

func initDB() *sql.DB {
	connStr := env.GetEnv("APP_DB_URL")
	if connStr == "" {
		log.Fatalf("APP_DB_URL is required")
	}

	db, err := database.NewDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	log.Print("Database connected")

	return db
}
