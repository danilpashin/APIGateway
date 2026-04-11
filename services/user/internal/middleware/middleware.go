package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeReq := time.Now()

		next.ServeHTTP(w, r)

		timeSince := time.Since(timeReq)
		log.Print("Request time: ", timeSince)
	})
}
