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

		if path == "/api/cita-medica/recordatorios" {
			h.RecordatoriosHandler(w, r)
			return
		}

		if path == "/api/cita-medica/metricas" {
			h.MetricasHandler(w, r)
			return
		}

		if path == "/api/cita-medica/recurrentes" {
			h.CitasRecurrentesHandler(w, r)
			return
		}

		if path == "/api/cita-medica/pacientes" {
			h.ListarPacientesHandler(w, r)
			return
		}
		if path == "/api/cita-medica/doctores" {
			h.ListarDoctoresHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/cita-medica/paciente/") {
			h.CitasPorPacienteHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/cita-medica/historial/") {
			h.HistorialCitaHandler(w, r)
			return
		}

		if strings.HasSuffix(path, "/confirmar") && strings.HasPrefix(path, "/api/cita-medica/") {
			h.ConfirmarCitaHandler(w, r)
			return
		}

		if strings.HasSuffix(path, "/completar") && strings.HasPrefix(path, "/api/cita-medica/") {
			h.CompletarCitaHandler(w, r)
			return
		}

		if strings.HasSuffix(path, "/notas") && strings.HasPrefix(path, "/api/cita-medica/") {
			h.NotasCitaHandler(w, r)
			return
		}

		if strings.HasSuffix(path, "/detalle") && strings.HasPrefix(path, "/api/cita-medica/") {
			h.DetalleCitaHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/cita-medica/") && path != "/api/cita-medica/" {
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

		if path == "/api/cita-medica" && r.Method == http.MethodPost {
			h.CrearCitaHandler(w, r)
			return
		}

		if path == "/api/cita-medica" {
			h.ListarCitasHandler(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"Ruta no encontrada"}`, http.StatusNotFound)
	}
}
