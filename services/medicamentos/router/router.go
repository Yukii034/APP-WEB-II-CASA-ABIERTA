package router

import (
	"cuidabien/medicamentos/handlers"
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

		if path == "/api/patients" {
			h.ListarPacientesHandler(w, r)
			return
		}

		if path == "/api/medications/alerts" {
			h.ListarAlertasHandler(w, r)
			return
		}

		if path == "/api/medications/interactions" {
			h.VerificarInteraccionesHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/medications/adherence/") {
			h.AdherenciaHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/medications/alerts/") {
			if strings.HasSuffix(path, "/read") {
				h.MarcarAlertaLeidaHandler(w, r)
				return
			}
			h.AlertasPacienteHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/medications/") {
			rest := strings.TrimPrefix(path, "/api/medications/")
			if rest == "" {
				if r.Method == http.MethodPost {
					h.CrearMedicamentoHandler(w, r)
				} else {
					h.ListarMedicamentosHandler(w, r)
				}
				return
			}

			parts := strings.Split(rest, "/")
			if len(parts) == 1 {
				switch r.Method {
				case http.MethodGet:
					h.ObtenerMedicamentoHandler(w, r)
				case http.MethodPut:
					h.ActualizarMedicamentoHandler(w, r)
				case http.MethodDelete:
					h.EliminarMedicamentoHandler(w, r)
				default:
					http.Error(w, `{"error":"Metodo no permitido"}`, http.StatusMethodNotAllowed)
				}
				return
			}

			if parts[1] == "take" {
				h.RegistrarTomaHandler(w, r)
				return
			}
			if parts[1] == "history" {
				h.HistorialTomasHandler(w, r)
				return
			}

			http.Error(w, `{"error":"Ruta no encontrada"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, `{"error":"Ruta no encontrada"}`, http.StatusNotFound)
	}
}
