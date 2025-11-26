package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware registra solicitações com método, path, remote addr e duração
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Início: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		dur := time.Since(start)
		log.Printf("Fim: %s (demorou %v)", r.URL.Path, dur)
	})
}
