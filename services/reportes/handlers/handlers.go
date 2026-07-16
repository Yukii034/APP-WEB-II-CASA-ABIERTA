package handlers

import (
	"cuidabien/reportes/logger"
	"cuidabien/reportes/storage"
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

// GET /api/reportes/semanal/{paciente_id}
func (h *Handlers) ReporteSemanalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/reportes/semanal/")
	pacienteID := path
	if pacienteID == "" {
		writeError(w, http.StatusBadRequest, "ID de paciente requerido")
		return
	}

	if h.Store.FindPacienteByID(pacienteID) == nil {
		writeError(w, http.StatusNotFound, "Paciente no encontrado")
		return
	}

	reporte := h.Store.GenerarReporteSemanal(pacienteID)
	logger.LogJSON("INFO", "Reporte semanal generado para "+pacienteID, "reporte_semanal", r.URL.Path, "")
	writeJSON(w, http.StatusOK, reporte)
}

// GET /api/reportes/paciente/{paciente_id}
func (h *Handlers) ReportePacienteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/reportes/paciente/")
	pacienteID := path
	if pacienteID == "" {
		writeError(w, http.StatusBadRequest, "ID de paciente requerido")
		return
	}

	if h.Store.FindPacienteByID(pacienteID) == nil {
		writeError(w, http.StatusNotFound, "Paciente no encontrado")
		return
	}

	reporte := h.Store.GenerarReportePaciente(pacienteID)
	logger.LogJSON("INFO", "Reporte de paciente generado para "+pacienteID, "reporte_paciente", r.URL.Path, "")
	writeJSON(w, http.StatusOK, reporte)
}

// GET /api/reportes/resumen
func (h *Handlers) ResumenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	dashboard := h.Store.GenerarDashboard()
	logger.LogJSON("INFO", "Dashboard generado", "dashboard", r.URL.Path, "")
	writeJSON(w, http.StatusOK, dashboard)
}

// GET /api/reportes/todos
func (h *Handlers) TodosReportesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	var reportes []interface{}
	for _, p := range h.Store.Pacientes {
		reportes = append(reportes, h.Store.GenerarReporteSemanal(p.ID))
	}

	if reportes == nil {
		reportes = []interface{}{}
	}
	writeJSON(w, http.StatusOK, reportes)
}

// GET /api/patients
func (h *Handlers) ListarPacientesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}
	writeJSON(w, http.StatusOK, h.Store.Pacientes)
}
