package httpserver

import (
	"net/http"
	"strings"

	"monitoreo-signos-vitales/internal/handlers"
)

// NuevoRouter registers the public, static API of the microservice.
func NuevoRouter(h *handlers.SignosVitalesHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("POST /api/vitales", h.Crear)
	mux.HandleFunc("POST /api/vitales/", h.Crear)
	mux.HandleFunc("GET /api/vitales/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/ultimo") {
			h.Ultimo(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/tendencia") {
			h.Tendencia(w, r)
			return
		}
		h.Historial(w, r)
	})
	return cors(mux)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
