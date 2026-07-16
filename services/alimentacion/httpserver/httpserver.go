package httpserver

import (
	"net/http"

	"cuidabien/alimentacion/handlers"
)

// New arma el mux HTTP conectando cada ruta con su handler. Es la
// única capa que conoce las rutas del servicio.
func New(h *handlers.Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HealthHandler)
	mux.HandleFunc("/api/alimentacion", h.AlimentacionHandler)
	mux.HandleFunc("/api/alimentacion/resumen", h.ResumenHandler)
	mux.HandleFunc("/api/alimentacion/reset", h.ResetHandler)
	mux.HandleFunc("/api/alimentacion/historial", h.HistorialHandler)
	mux.HandleFunc("/api/alimentacion/hidratacion", h.HidratacionHandler)
	mux.HandleFunc("/api/alimentacion/restricciones", h.RestriccionesHandler)

	return mux
}
