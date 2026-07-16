package middleware

import (
	"net/http"
	"time"

	"cuidabien/estimulacion-cognitiva/logger"
)

type responseWriterConEstado struct {
	http.ResponseWriter
	status int
}

func (w *responseWriterConEstado) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Logging(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inicio := time.Now()
			wrapped := &responseWriterConEstado{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			log.Info("%s %s -> %d (%s)", r.Method, r.URL.Path, wrapped.status, time.Since(inicio))
		})
	}
}

func Recover(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("panic recuperado en %s %s: %v", r.Method, r.URL.Path, err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"error interno del servidor"}`))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
