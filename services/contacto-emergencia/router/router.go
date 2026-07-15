package router

import (
	"cuidabien/contacto-emergencia/handlers"
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

		if path == "/api/metrics" {
			h.MetricasHandler(w, r)
			return
		}

		// --- Alertas ---

		if strings.HasPrefix(path, "/api/alerts/history/") {
			h.HistorialAlertaHandler(w, r)
			return
		}

		if strings.HasSuffix(path, "/attend") && strings.HasPrefix(path, "/api/alerts/") {
			h.AtenderAlertaHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/alerts/") && path != "/api/alerts/" {
			switch r.Method {
			case http.MethodGet:
				h.ObtenerAlertaHandler(w, r)
			case http.MethodDelete:
				h.CancelarAlertaHandler(w, r)
			default:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(`{"error":"Metodo no permitido"}`))
			}
			return
		}

		if path == "/api/alerts" {
			switch r.Method {
			case http.MethodPost:
				h.CrearAlertaHandler(w, r)
			case http.MethodGet:
				h.ListarAlertasHandler(w, r)
			default:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(`{"error":"Metodo no permitido"}`))
			}
			return
		}

		// --- Contactos ---

		if strings.HasPrefix(path, "/api/contacts/") && path != "/api/contacts/" {
			switch r.Method {
			case http.MethodGet:
				h.ObtenerContactoHandler(w, r)
			case http.MethodPut:
				h.ActualizarContactoHandler(w, r)
			case http.MethodDelete:
				h.EliminarContactoHandler(w, r)
			default:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(`{"error":"Metodo no permitido"}`))
			}
			return
		}

		if path == "/api/contacts" {
			switch r.Method {
			case http.MethodPost:
				h.CrearContactoHandler(w, r)
			case http.MethodGet:
				h.ListarContactosHandler(w, r)
			default:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(`{"error":"Metodo no permitido"}`))
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, `{"error":"Ruta no encontrada"}`, http.StatusNotFound)
	}
}
