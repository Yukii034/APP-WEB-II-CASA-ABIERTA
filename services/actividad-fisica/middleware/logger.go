package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		

		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next(rec, r)
		
		log.Printf(
			"[INFO] %s %s - %d %s - %v",
			r.Method,
			r.URL.Path,
			rec.statusCode,
			http.StatusText(rec.statusCode),
			time.Since(start),
		)
	}
}