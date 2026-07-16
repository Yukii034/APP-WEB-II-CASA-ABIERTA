package router

import (
	"cuidabien/reportes/handlers"
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

		if path == "/api/reportes/resumen" {
			h.ResumenHandler(w, r)
			return
		}

		if path == "/api/reportes/todos" {
			h.TodosReportesHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/reportes/semanal/") {
			h.ReporteSemanalHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/reportes/paciente/") {
			h.ReportePacienteHandler(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, `{"error":"Ruta no encontrada"}`, http.StatusNotFound)
	}
}
