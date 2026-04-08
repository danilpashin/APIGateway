package main

import (
	"fmt"
	"log"
	"net/http"
	"pkg/env"
	"product/internal/config"

	"github.com/go-chi/chi/v5"
)

func main() {
	env.LoadEnv()
	cfg := config.Load()

	runServer(cfg)
}

func runServer(cfg *config.Config) {
	r := newRouter()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: r,
	}
	log.Print("Server starting on: ", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server failed: ", err)
	}
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/product", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Write([]byte("First handler in this project"))
	}))

	return r
}
