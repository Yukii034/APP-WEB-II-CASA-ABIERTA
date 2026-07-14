package router

import (
	"cuidabien/citas/handlers"
	"net/http"
	"strings"
)

func New(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/health" {
			h.HealthHandler(w, r)
			return
		}

		if path == "/api/appointments/reminders" {
			h.RecordatoriosHandler(w, r)
			return
		}

		if path == "/api/appointments/metrics" {
			h.MetricasHandler(w, r)
			return
		}

		if path == "/api/appointments/recurring" {
			h.CitasRecurrentesHandler(w, r)
			return
		}

		if path == "/api/patients" {
			h.ListarPacientesHandler(w, r)
			return
		}
		if path == "/api/doctors" {
			h.ListarDoctoresHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/appointments/patient/") {
			h.CitasPorPacienteHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/appointments/history/") {
			h.HistorialCitaHandler(w, r)
			return
		}

		if strings.HasSuffix(path, "/confirm") && strings.HasPrefix(path, "/api/appointments/") {
			h.ConfirmarCitaHandler(w, r)
			return
		}

		if strings.HasSuffix(path, "/complete") && strings.HasPrefix(path, "/api/appointments/") {
			h.CompletarCitaHandler(w, r)
			return
		}

		if strings.HasSuffix(path, "/notes") && strings.HasPrefix(path, "/api/appointments/") {
			h.NotasCitaHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/appointments/") && path != "/api/appointments/" {
			switch r.Method {
			case http.MethodGet:
				h.ObtenerCitaHandler(w, r)
			case http.MethodPut:
				h.ActualizarCitaHandler(w, r)
			case http.MethodDelete:
				h.CancelarCitaHandler(w, r)
			default:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(`{"error":"Metodo no permitido"}`))
			}
			return
		}

		if path == "/api/appointments" && r.Method == http.MethodPost {
			h.CrearCitaHandler(w, r)
			return
		}

		if path == "/api/appointments" {
			h.ListarCitasHandler(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, `{"error":"Ruta no encontrada"}`, http.StatusNotFound)
	}
}
