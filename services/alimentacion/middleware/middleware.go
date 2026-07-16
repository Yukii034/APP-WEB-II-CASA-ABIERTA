package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger registra método, ruta y duración de cada petición. No cambia
// el comportamiento de la petición, solo la envuelve.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(inicio))
	})
}
