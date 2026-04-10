package main

import (
	"apigateway/services/product/internal/config"
	"apigateway/services/product/internal/database"
	"apigateway/services/product/internal/handler"
	"apigateway/services/product/internal/middleware"
	"apigateway/services/product/internal/repository/postgres"
	"apigateway/services/product/internal/service"
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

	if handleCLI() {
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
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed: ", err)
		}
	}()

	gracefulShutdown(&srv, db)
}

func handleCLI() bool {
	var cmd string
	var version int

	flag.StringVar(&cmd, "cmd", "", "migration command: up, down, force, version")
	flag.IntVar(&version, "version", version, "version for force or migrate command")
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
	productRepo := postgres.NewPostgresProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(*productService)

	r := chi.NewRouter()
	r.Use(middleware.PanicRecoveryMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Post("/products", productHandler.CreateProductHandler)
	r.Put("/products/{id}", productHandler.UpdateProductHandler)
	r.Get("/products/{id}", productHandler.GetProductHandler)
	r.Get("/products", productHandler.ListProductsHandler)
	r.Delete("/products/{id}", productHandler.DeleteProductHandler)

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
