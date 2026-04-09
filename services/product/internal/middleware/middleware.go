package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		next.ServeHTTP(w, r)
		timeSince := time.Since(timeStart)
		log.Print("Request time: ", timeSince)
	})
}
