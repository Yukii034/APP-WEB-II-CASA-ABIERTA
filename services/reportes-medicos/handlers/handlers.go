package handlers

import (
	"cuidabien/reportes-medicos/storage"
	"encoding/json"
	"net/http"
	"strings"
)

type Handlers struct {
	Store *storage.Store
}

func New(s *storage.Store) *Handlers {
	return &Handlers{Store: s}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) ResumenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	writeJSON(w, http.StatusOK, h.Store.CrearResumen())
}

func (h *Handlers) SemanalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	writeJSON(w, http.StatusOK, h.Store.ListarReportes())
}

func (h *Handlers) PacienteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/reportes-medicos/paciente/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de paciente requerido")
		return
	}

	reporte := h.Store.BuscarPorPaciente(id)
	if reporte == nil {
		writeError(w, http.StatusNotFound, "Reporte no encontrado")
		return
	}

	writeJSON(w, http.StatusOK, reporte)
}
