package middleware

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Begin: %s, %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("Request ended: %s, %s. Time duration: %s", r.Method, r.URL.Path, duration)
	})
}
