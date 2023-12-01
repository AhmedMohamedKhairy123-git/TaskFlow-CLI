package web

import (
	"log"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		
		next(w, r)
		
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	}
}

func (s *Server) recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				respondError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next(w, r)
	}
}

func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

func (s *Server) rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	limiter := make(chan struct{}, 10)
	
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case limiter <- struct{}{}:
			defer func() { <-limiter }()
			next(w, r)
		default:
			respondError(w, http.StatusTooManyRequests, "Rate limit exceeded")
		}
	}
}