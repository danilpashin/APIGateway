package mw

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"
)

const requestIDKey string = "request_id"

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger := slog.Default()

				logger.ErrorContext(r.Context(), "panic recovered",
					slog.Any("error", err),
					slog.String("path", r.URL.Path),
					slog.String("stack", string(debug.Stack())),
				)

				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)

		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
