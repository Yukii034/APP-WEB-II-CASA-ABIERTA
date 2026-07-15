package handlers

import (
	"cuidabien/reportes-medicos/models"
	"cuidabien/reportes-medicos/storage"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type citaMedica struct {
	ID         string `json:"id"`
	PacienteID string `json:"paciente_id"`
	Estado     string `json:"estado"`
}

type citasPaginadas struct {
	Data []citaMedica `json:"data"`
}

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

	reportes := h.reportesConCitas()
	writeJSON(w, http.StatusOK, storage.CrearResumen(reportes))
}

func (h *Handlers) SemanalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	writeJSON(w, http.StatusOK, h.reportesConCitas())
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

	reporteConCitas := h.aplicarCitas(*reporte)
	writeJSON(w, http.StatusOK, reporteConCitas)
}

func (h *Handlers) reportesConCitas() []models.ReportePaciente {
	reportes := h.Store.ListarReportes()
	resultado := make([]models.ReportePaciente, 0, len(reportes))
	for _, reporte := range reportes {
		resultado = append(resultado, h.aplicarCitas(reporte))
	}
	return resultado
}

func (h *Handlers) aplicarCitas(reporte models.ReportePaciente) models.ReportePaciente {
	citas, err := consultarCitasPaciente(reporte.PacienteID)
	if err != nil {
		return reporte
	}

	reporte.CitasProgramadas = len(citas)
	reporte.CitasCompletadas = 0
	for _, cita := range citas {
		if cita.Estado == "completada" {
			reporte.CitasCompletadas++
		}
	}

	return reporte
}

func consultarCitasPaciente(pacienteID string) ([]citaMedica, error) {
	baseURL := os.Getenv("CITAS_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("CITAS_URL no configurada")
	}

	paths := []string{"/api/cita-medica", "/api/appointments"}
	var ultimoErr error
	for _, path := range paths {
		citas, err := consultarCitasEnRuta(baseURL, path, pacienteID)
		if err == nil {
			return citas, nil
		}
		ultimoErr = err
	}

	return nil, ultimoErr
}

func consultarCitasEnRuta(baseURL, path, pacienteID string) ([]citaMedica, error) {
	endpoint := strings.TrimRight(baseURL, "/") + path + "?paciente_id=" + url.QueryEscape(pacienteID)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("citas respondio con estado %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paginadas citasPaginadas
	if err := json.Unmarshal(body, &paginadas); err == nil && paginadas.Data != nil {
		return paginadas.Data, nil
	}

	var citas []citaMedica
	if err := json.Unmarshal(body, &citas); err != nil {
		return nil, err
	}
	return citas, nil
}
