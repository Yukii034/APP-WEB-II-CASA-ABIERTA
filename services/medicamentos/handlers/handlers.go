package handlers

import (
	"cuidabien/medicamentos/logger"
	"cuidabien/medicamentos/models"
	"cuidabien/medicamentos/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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

// POST /api/medications
func (h *Handlers) CrearMedicamentoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	var req models.CrearMedicamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if req.PacienteID == "" || req.Nombre == "" || req.Dosis == "" {
		writeError(w, http.StatusBadRequest, "Los campos paciente_id, nombre y dosis son obligatorios")
		return
	}

	req.Nombre = storage.Sanitizar(req.Nombre)
	req.Dosis = storage.Sanitizar(req.Dosis)
	req.Frecuencia = storage.Sanitizar(req.Frecuencia)

	if h.Store.FindPacienteByID(req.PacienteID) == nil {
		writeError(w, http.StatusBadRequest, "El paciente no existe")
		return
	}

	if len(req.Horarios) > 0 {
		for _, horario := range req.Horarios {
			if !storage.ValidarHorario(horario) {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("Horario invalido: '%s'. Formato requerido: HH:MM (ej: 08:00)", horario))
				return
			}
		}
	}

	if req.FechaInicio == "" {
		req.FechaInicio = time.Now().Format("2006-01-02")
	}

	med := models.Medicamento{
		ID:          h.Store.GenerateMedID(),
		PacienteID:  req.PacienteID,
		Nombre:      req.Nombre,
		Dosis:       req.Dosis,
		Frecuencia:  req.Frecuencia,
		Horarios:    req.Horarios,
		FechaInicio: req.FechaInicio,
		FechaFin:    req.FechaFin,
		Estado:      "activo",
		Notas:       req.Notas,
	}

	h.Store.Medicamentos = append(h.Store.Medicamentos, med)
	h.Store.VerificarInteraccionesActivas(req.PacienteID)

	logger.LogJSON("INFO", fmt.Sprintf("Medicamento %s creado para paciente %s", med.ID, req.PacienteID), "crear_medicamento", r.URL.Path, "")
	writeJSON(w, http.StatusCreated, med)
}

// GET /api/medications
func (h *Handlers) ListarMedicamentosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	q := r.URL.Query()
	filtroPaciente := q.Get("paciente_id")
	filtroEstado := q.Get("estado")
	filtroNombre := q.Get("nombre")

	var filtrados []models.Medicamento
	for _, m := range h.Store.Medicamentos {
		if filtroPaciente != "" && m.PacienteID != filtroPaciente {
			continue
		}
		if filtroEstado != "" && m.Estado != filtroEstado {
			continue
		}
		if filtroNombre != "" && !strings.Contains(strings.ToLower(m.Nombre), strings.ToLower(filtroNombre)) {
			continue
		}
		filtrados = append(filtrados, m)
	}

	if filtrados == nil {
		filtrados = []models.Medicamento{}
	}
	writeJSON(w, http.StatusOK, filtrados)
}

// GET /api/medications/{id}
func (h *Handlers) ObtenerMedicamentoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/medications/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de medicamento requerido")
		return
	}

	med := h.Store.FindMedicamentoByID(id)
	if med == nil {
		writeError(w, http.StatusNotFound, "Medicamento no encontrado")
		return
	}

	writeJSON(w, http.StatusOK, med)
}

// PUT /api/medications/{id}
func (h *Handlers) ActualizarMedicamentoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/medications/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de medicamento requerido")
		return
	}

	med := h.Store.FindMedicamentoByID(id)
	if med == nil {
		writeError(w, http.StatusNotFound, "Medicamento no encontrado")
		return
	}

	if med.Estado == "completado" {
		writeError(w, http.StatusBadRequest, "No se puede modificar un medicamento completado")
		return
	}

	var req models.ActualizarMedicamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if req.Nombre != "" {
		med.Nombre = storage.Sanitizar(req.Nombre)
	}
	if req.Dosis != "" {
		med.Dosis = storage.Sanitizar(req.Dosis)
	}
	if req.Frecuencia != "" {
		med.Frecuencia = storage.Sanitizar(req.Frecuencia)
	}
	if len(req.Horarios) > 0 {
		for _, horario := range req.Horarios {
			if !storage.ValidarHorario(horario) {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("Horario invalido: '%s'", horario))
				return
			}
		}
		med.Horarios = req.Horarios
	}
	if req.FechaFin != "" {
		med.FechaFin = req.FechaFin
	}
	if req.Notas != "" {
		med.Notas = storage.Sanitizar(req.Notas)
	}

	logger.LogJSON("INFO", fmt.Sprintf("Medicamento %s actualizado", id), "actualizar_medicamento", r.URL.Path, "")
	writeJSON(w, http.StatusOK, med)
}

// DELETE /api/medications/{id}
func (h *Handlers) EliminarMedicamentoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/medications/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de medicamento requerido")
		return
	}

	med := h.Store.FindMedicamentoByID(id)
	if med == nil {
		writeError(w, http.StatusNotFound, "Medicamento no encontrado")
		return
	}

	if med.Estado == "suspendido" {
		writeError(w, http.StatusBadRequest, "El medicamento ya esta suspendido")
		return
	}

	med.Estado = "suspendido"
	logger.LogJSON("INFO", fmt.Sprintf("Medicamento %s suspendido", id), "suspender_medicamento", r.URL.Path, "")
	writeJSON(w, http.StatusOK, map[string]string{
		"mensaje": "Medicamento suspendido exitosamente",
		"id":      id,
	})
}

// POST /api/medications/{id}/take
func (h *Handlers) RegistrarTomaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/medications/")
	id = strings.TrimSuffix(id, "/take")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de medicamento requerido")
		return
	}

	med := h.Store.FindMedicamentoByID(id)
	if med == nil {
		writeError(w, http.StatusNotFound, "Medicamento no encontrado")
		return
	}

	if med.Estado != "activo" {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Solo se pueden registrar tomas de medicamentos activos. Estado actual: '%s'", med.Estado))
		return
	}

	var req models.RegistrarTomaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if req.Estado == "" {
		req.Estado = "cumplida"
	}
	if req.Estado != "cumplida" && req.Estado != "no_cumplida" && req.Estado != "saltada" {
		writeError(w, http.StatusBadRequest, "Estado invalido. Opciones: cumplida, no_cumplida, saltada")
		return
	}

	now := time.Now()
	horaAhora := now.Format("15:04")
	fechaAhora := now.Format("2006-01-02")
	fechaHoraProgramada := fechaAhora + " " + horaAhora

	toma := models.Toma{
		ID:                  h.Store.GenerateTomaID(),
		MedicamentoID:       id,
		PacienteID:          med.PacienteID,
		FechaHoraProgramada: fechaHoraProgramada,
		Estado:              req.Estado,
		FechaHoraReal:       now.Format(time.RFC3339),
		Notas:               req.Notas,
	}

	h.Store.Tomas = append(h.Store.Tomas, toma)
	logger.LogJSON("INFO", fmt.Sprintf("Toma %s registrada para medicamento %s (%s)", toma.ID, id, req.Estado), "registrar_toma", r.URL.Path, "")
	writeJSON(w, http.StatusCreated, toma)
}

// GET /api/medications/{id}/history
func (h *Handlers) HistorialTomasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/medications/")
	id := strings.TrimSuffix(path, "/history")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de medicamento requerido")
		return
	}

	if h.Store.FindMedicamentoByID(id) == nil {
		writeError(w, http.StatusNotFound, "Medicamento no encontrado")
		return
	}

	tomas := h.Store.TomasPorMedicamento(id)
	if tomas == nil {
		tomas = []models.Toma{}
	}
	writeJSON(w, http.StatusOK, tomas)
}

// GET /api/medications/adherence/{patient_id}
func (h *Handlers) AdherenciaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/medications/adherence/")
	pacienteID := path
	if pacienteID == "" {
		writeError(w, http.StatusBadRequest, "ID de paciente requerido")
		return
	}

	if h.Store.FindPacienteByID(pacienteID) == nil {
		writeError(w, http.StatusNotFound, "Paciente no encontrado")
		return
	}

	adherencia := h.Store.CalcularAdherencia(pacienteID)
	writeJSON(w, http.StatusOK, adherencia)
}

// GET /api/medications/alerts
func (h *Handlers) ListarAlertasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	h.Store.GenerarAlertasVencimiento()

	q := r.URL.Query()
	filtroPaciente := q.Get("paciente_id")

	var filtradas []models.Alerta
	for _, a := range h.Store.Alertas {
		if filtroPaciente != "" && a.PacienteID != filtroPaciente {
			continue
		}
		filtradas = append(filtradas, a)
	}

	if filtradas == nil {
		filtradas = []models.Alerta{}
	}
	writeJSON(w, http.StatusOK, filtradas)
}

// GET /api/medications/alerts/{patient_id}
func (h *Handlers) AlertasPacienteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/medications/alerts/")
	pacienteID := path
	if pacienteID == "" {
		writeError(w, http.StatusBadRequest, "ID de paciente requerido")
		return
	}

	h.Store.GenerarAlertasVencimiento()

	alertas := h.Store.AlertasPorPaciente(pacienteID)
	if alertas == nil {
		alertas = []models.Alerta{}
	}
	writeJSON(w, http.StatusOK, alertas)
}

// PATCH /api/medications/alerts/{id}/read
func (h *Handlers) MarcarAlertaLeidaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/medications/alerts/")
	id = strings.TrimSuffix(id, "/read")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de alerta requerido")
		return
	}

	alerta := h.Store.FindAlertaByID(id)
	if alerta == nil {
		writeError(w, http.StatusNotFound, "Alerta no encontrada")
		return
	}

	alerta.Leida = true
	writeJSON(w, http.StatusOK, map[string]string{
		"mensaje": "Alerta marcada como leida",
		"id":      id,
	})
}

// GET /api/medications/interactions
func (h *Handlers) VerificarInteraccionesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	q := r.URL.Query()
	pacienteID := q.Get("paciente_id")

	if pacienteID != "" {
		interacciones := h.Store.VerificarInteraccionesActivas(pacienteID)
		if interacciones == nil {
			interacciones = []models.Interaccion{}
		}
		writeJSON(w, http.StatusOK, interacciones)
		return
	}

	medA := q.Get("medicamento_a")
	medB := q.Get("medicamento_b")
	if medA != "" && medB != "" {
		inter := h.Store.BuscarInteraccion(medA, medB)
		if inter == nil {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"interaccion_encontrada": false,
				"mensaje":               "No se encontraron interacciones entre estos medicamentos",
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"interaccion_encontrada": true,
			"interaccion":            inter,
		})
		return
	}

	writeError(w, http.StatusBadRequest, "Se requiere paciente_id o medicamento_a + medicamento_b")
}

// GET /api/patients
func (h *Handlers) ListarPacientesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}
	writeJSON(w, http.StatusOK, h.Store.Pacientes)
}
