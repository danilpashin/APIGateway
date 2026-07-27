module apigateway/services/order

go 1.26.3

require (
	github.com/go-chi/chi/v5 v5.3.1
	github.com/go-chi/httplog/v3 v3.4.0
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/induzo/gocom/database/pgx-slog v1.0.42
	github.com/jackc/pgx/v5 v5.10.0
	github.com/redis/go-redis/v9 v9.21.0
	pkg v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/text v0.38.0 // indirect
)

replace pkg => ../../pkg
