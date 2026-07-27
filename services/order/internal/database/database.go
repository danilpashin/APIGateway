package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pgxslogger "github.com/induzo/gocom/database/pgx-slog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

func NewPgxPool(connStr string, logger *slog.Logger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conn string: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute

	tracelogger := pgxslogger.NewLogger(logger)

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   tracelogger,
		LogLevel: tracelog.LogLevelInfo,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return pool, nil
}
