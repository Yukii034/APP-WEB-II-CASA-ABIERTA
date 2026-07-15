package handlers

import (
	"cuidabien/citas/logger"
	"cuidabien/citas/models"
	"cuidabien/citas/storage"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
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

func (h *Handlers) CrearCitaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	var req models.CrearCitaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if req.PacienteID == "" || req.DoctorID == "" || req.Fecha == "" || req.Hora == "" {
		writeError(w, http.StatusBadRequest, "Los campos paciente_id, doctor_id, fecha y hora son obligatorios")
		return
	}

	req.PacienteID = storage.Sanitizar(req.PacienteID)
	req.DoctorID = storage.Sanitizar(req.DoctorID)
	req.Motivo = storage.Sanitizar(req.Motivo)

	if h.Store.FindPacienteByID(req.PacienteID) == nil {
		writeError(w, http.StatusBadRequest, "El paciente no existe")
		return
	}
	if h.Store.FindDoctorByID(req.DoctorID) == nil {
		writeError(w, http.StatusBadRequest, "El doctor no existe")
		return
	}

	if req.Prioridad == "" {
		auto := storage.DetectarPrioridadAutomatica(req.Motivo)
		if auto != "" {
			req.Prioridad = auto
			logger.LogJSON("INFO", fmt.Sprintf("Prioridad auto-detectada: %s por motivo: %s", auto, req.Motivo), "auto_prioridad", r.URL.Path, "")
		} else {
			req.Prioridad = "normal"
		}
	}
	if req.Prioridad != "normal" && req.Prioridad != "urgente" && req.Prioridad != "control" {
		writeError(w, http.StatusBadRequest, "Prioridad invalida. Opciones: normal, urgente, control")
		return
	}

	if h.Store.EsFechaPasada(req.Fecha, req.Hora) {
		writeError(w, http.StatusBadRequest, "No se pueden crear citas en fechas pasadas")
		return
	}
	if h.Store.MedicoOcupado(req.DoctorID, req.Fecha, req.Hora, "") {
		writeError(w, http.StatusConflict, "El medico ya tiene una cita en ese horario")
		return
	}
	if h.Store.PacienteOcupado(req.PacienteID, req.Fecha, req.Hora, "") {
		writeError(w, http.StatusConflict, "El paciente ya tiene una cita en ese horario")
		return
	}

	nuevaCita := models.Cita{
		ID:          h.Store.GenerateID(),
		PacienteID:  req.PacienteID,
		DoctorID:    req.DoctorID,
		Fecha:       req.Fecha,
		Hora:        req.Hora,
		Estado:      "pendiente",
		Prioridad:   req.Prioridad,
		Motivo:      req.Motivo,
		NotasMedico: "",
	}

	h.Store.Citas = append(h.Store.Citas, nuevaCita)
	h.Store.RegistrarHistorial(nuevaCita.ID, "creacion", "", "pendiente", "")
	h.Store.ActualizarMetricas()
	h.Store.TodayCreated++

	logger.LogJSON("INFO", fmt.Sprintf("Cita %s creada para %s con %s", nuevaCita.ID, req.PacienteID, req.DoctorID), "crear_cita", r.URL.Path, "")
	writeJSON(w, http.StatusCreated, nuevaCita)
}

func (h *Handlers) ListarCitasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	q := r.URL.Query()

	filtroDoctor := q.Get("doctor_id")
	filtroPaciente := q.Get("paciente_id")
	filtroEstado := q.Get("estado")
	filtroPrioridad := q.Get("prioridad")
	filtroFechaInicio := q.Get("fecha_inicio")
	filtroFechaFin := q.Get("fecha_fin")
	filtroMotivo := q.Get("motivo")

	page := 1
	limit := 10
	if p := q.Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := q.Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var filtradas []models.Cita
	for _, c := range h.Store.Citas {
		if filtroDoctor != "" && c.DoctorID != filtroDoctor {
			continue
		}
		if filtroPaciente != "" && c.PacienteID != filtroPaciente {
			continue
		}
		if filtroEstado != "" && c.Estado != filtroEstado {
			continue
		}
		if filtroPrioridad != "" && c.Prioridad != filtroPrioridad {
			continue
		}
		if filtroFechaInicio != "" && c.Fecha < filtroFechaInicio {
			continue
		}
		if filtroFechaFin != "" && c.Fecha > filtroFechaFin {
			continue
		}
		if filtroMotivo != "" && !strings.Contains(strings.ToLower(c.Motivo), strings.ToLower(filtroMotivo)) {
			continue
		}
		filtradas = append(filtradas, c)
	}

	total := len(filtradas)
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Data:       filtradas[start:end],
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

func (h *Handlers) ObtenerCitaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cita-medica/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de cita requerido")
		return
	}

	cita := h.Store.FindCitaByID(id)
	if cita == nil {
		writeError(w, http.StatusNotFound, "Cita no encontrada")
		return
	}

	writeJSON(w, http.StatusOK, cita)
}

func (h *Handlers) DetalleCitaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cita-medica/")
	id = strings.TrimSuffix(id, "/detalle")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de cita requerido")
		return
	}

	cita := h.Store.FindCitaByID(id)
	if cita == nil {
		writeError(w, http.StatusNotFound, "Cita no encontrada")
		return
	}

	detalle := models.DetalleCita{
		Cita:     *cita,
		Paciente: h.Store.FindPacienteByID(cita.PacienteID),
		Doctor:   h.Store.FindDoctorByID(cita.DoctorID),
	}

	infoID := storage.InformacionSaludIDPorPaciente(cita.PacienteID)
	detalle.InformacionSaludID = infoID
	if infoID == "" {
		detalle.InformacionSaludAviso = "No existe mapeo hacia informacion-salud para este paciente"
		writeJSON(w, http.StatusOK, detalle)
		return
	}

	info, aviso := consultarInformacionSalud(infoID)
	detalle.InformacionSalud = info
	detalle.InformacionSaludAviso = aviso

	writeJSON(w, http.StatusOK, detalle)
}

func consultarInformacionSalud(infoID string) (*models.InformacionSalud, string) {
	baseURL := os.Getenv("INFORMACION_SALUD_URL")
	if baseURL == "" {
		return nil, "INFORMACION_SALUD_URL no configurada"
	}

	resp, err := http.Get(strings.TrimRight(baseURL, "/") + "/api/informacion-salud/" + infoID)
	if err != nil {
		return nil, "No se pudo contactar informacion-salud"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "No se pudo leer la respuesta de informacion-salud"
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Sprintf("informacion-salud respondio con estado %d", resp.StatusCode)
	}

	var info models.InformacionSalud
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, "Respuesta invalida de informacion-salud"
	}

	return &info, ""
}

func (h *Handlers) ActualizarCitaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cita-medica/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de cita requerido")
		return
	}

	cita := h.Store.FindCitaByID(id)
	if cita == nil {
		writeError(w, http.StatusNotFound, "Cita no encontrada")
		return
	}

	if !h.Store.EsEstadoModificable(cita.Estado) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("No se puede modificar una cita en estado '%s'", cita.Estado))
		return
	}

	var req models.ActualizarCitaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	nuevaFecha := cita.Fecha
	nuevaHora := cita.Hora

	if req.Fecha != "" {
		nuevaFecha = storage.Sanitizar(req.Fecha)
	}
	if req.Hora != "" {
		nuevaHora = storage.Sanitizar(req.Hora)
	}

	if req.Fecha != "" || req.Hora != "" {
		if h.Store.EsFechaPasada(nuevaFecha, nuevaHora) {
			writeError(w, http.StatusBadRequest, "No se pueden reprogramar citas a fechas pasadas")
			return
		}
		if h.Store.MedicoOcupado(cita.DoctorID, nuevaFecha, nuevaHora, id) {
			writeError(w, http.StatusConflict, "El medico ya tiene una cita en ese horario")
			return
		}
		if h.Store.PacienteOcupado(cita.PacienteID, nuevaFecha, nuevaHora, id) {
			writeError(w, http.StatusConflict, "El paciente ya tiene una cita en ese horario")
			return
		}
	}

	if req.Fecha != "" {
		cita.Fecha = nuevaFecha
	}
	if req.Hora != "" {
		cita.Hora = nuevaHora
	}
	if req.Prioridad != "" {
		if req.Prioridad != "normal" && req.Prioridad != "urgente" && req.Prioridad != "control" {
			writeError(w, http.StatusBadRequest, "Prioridad invalida. Opciones: normal, urgente, control")
			return
		}
		cita.Prioridad = req.Prioridad
	}
	if req.Motivo != "" {
		cita.Motivo = storage.Sanitizar(req.Motivo)
	}

	h.Store.RegistrarHistorial(id, "actualizacion", "", cita.Estado, "Cita reprogramada")
	logger.LogJSON("INFO", fmt.Sprintf("Cita %s actualizada", id), "actualizar_cita", r.URL.Path, "")
	writeJSON(w, http.StatusOK, cita)
}

func (h *Handlers) CancelarCitaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cita-medica/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de cita requerido")
		return
	}

	cita := h.Store.FindCitaByID(id)
	if cita == nil {
		writeError(w, http.StatusNotFound, "Cita no encontrada")
		return
	}

	if cita.Estado == "cancelada" {
		writeError(w, http.StatusBadRequest, "La cita ya esta cancelada")
		return
	}
	if cita.Estado == "completada" {
		writeError(w, http.StatusBadRequest, "No se puede cancelar una cita completada")
		return
	}

	estadoAnt := cita.Estado
	cita.Estado = "cancelada"
	h.Store.RegistrarHistorial(id, "cancelacion", estadoAnt, "cancelada", "")
	h.Store.ActualizarMetricas()
	h.Store.TodayCancelled++

	logger.LogJSON("INFO", fmt.Sprintf("Cita %s cancelada", id), "cancelar_cita", r.URL.Path, "")
	writeJSON(w, http.StatusOK, map[string]string{
		"mensaje": "Cita cancelada exitosamente",
		"id":      id,
	})
}

func (h *Handlers) ConfirmarCitaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cita-medica/")
	id = strings.TrimSuffix(id, "/confirmar")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de cita requerido")
		return
	}

	cita := h.Store.FindCitaByID(id)
	if cita == nil {
		writeError(w, http.StatusNotFound, "Cita no encontrada")
		return
	}

	if cita.Estado != "pendiente" {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Solo se pueden confirmar citas pendientes. Estado actual: '%s'", cita.Estado))
		return
	}

	cita.Estado = "confirmada"
	h.Store.RegistrarHistorial(id, "confirmacion", "pendiente", "confirmada", "")
	logger.LogJSON("INFO", fmt.Sprintf("Cita %s confirmada", id), "confirmar_cita", r.URL.Path, "")
	writeJSON(w, http.StatusOK, map[string]string{
		"mensaje": "Cita confirmada exitosamente",
		"id":      id,
	})
}

func (h *Handlers) CompletarCitaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cita-medica/")
	id = strings.TrimSuffix(id, "/completar")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de cita requerido")
		return
	}

	cita := h.Store.FindCitaByID(id)
	if cita == nil {
		writeError(w, http.StatusNotFound, "Cita no encontrada")
		return
	}

	if cita.Estado == "cancelada" {
		writeError(w, http.StatusBadRequest, "Una cita cancelada no puede marcarse como completada")
		return
	}
	if cita.Estado != "confirmada" {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Solo se pueden completar citas confirmadas. Estado actual: '%s'", cita.Estado))
		return
	}

	cita.Estado = "completada"
	h.Store.RegistrarHistorial(id, "completado", "confirmada", "completada", "")
	h.Store.ActualizarMetricas()
	h.Store.TodayCompleted++

	logger.LogJSON("INFO", fmt.Sprintf("Cita %s completada", id), "completar_cita", r.URL.Path, "")
	writeJSON(w, http.StatusOK, map[string]string{
		"mensaje": "Cita completada exitosamente",
		"id":      id,
	})
}

func (h *Handlers) NotasCitaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cita-medica/")
	id = strings.TrimSuffix(id, "/notas")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de cita requerido")
		return
	}

	cita := h.Store.FindCitaByID(id)
	if cita == nil {
		writeError(w, http.StatusNotFound, "Cita no encontrada")
		return
	}

	if cita.Estado == "cancelada" {
		writeError(w, http.StatusBadRequest, "No se pueden agregar notas a una cita cancelada")
		return
	}

	var req models.NotasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if req.NotasMedico == "" {
		writeError(w, http.StatusBadRequest, "El campo notas_medico es obligatorio")
		return
	}

	cita.NotasMedico = storage.Sanitizar(req.NotasMedico)
	h.Store.RegistrarHistorial(id, "nota_agregada", "", cita.Estado, req.NotasMedico)
	logger.LogJSON("INFO", fmt.Sprintf("Notas agregadas a cita %s", id), "notas_cita", r.URL.Path, "")
	writeJSON(w, http.StatusOK, map[string]string{
		"mensaje":      "Notas actualizadas exitosamente",
		"id":           id,
		"notas_medico": cita.NotasMedico,
	})
}

func (h *Handlers) CitasPorPacienteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/cita-medica/paciente/")
	pacienteID := path
	if pacienteID == "" {
		writeError(w, http.StatusBadRequest, "ID de paciente requerido")
		return
	}

	if h.Store.FindPacienteByID(pacienteID) == nil {
		writeError(w, http.StatusNotFound, "Paciente no encontrado")
		return
	}

	var result []models.Cita
	for _, c := range h.Store.Citas {
		if c.PacienteID == pacienteID {
			result = append(result, c)
		}
	}

	if result == nil {
		result = []models.Cita{}
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handlers) RecordatoriosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	now := time.Now()
	manana := now.Add(24 * time.Hour)
	fechaLimite := manana.Format("2006-01-02")

	var reminders []models.Reminder
	for _, c := range h.Store.Citas {
		if c.Estado == "cancelada" || c.Estado == "completada" {
			continue
		}
		if c.Fecha == fechaLimite {
			mensaje := fmt.Sprintf("Recordatorio: Citas manana %s a las %s", c.Fecha, c.Hora)
			if c.Prioridad == "urgente" {
				mensaje = fmt.Sprintf("RECORDATORIO URGENTE: Citas manana %s a las %s", c.Fecha, c.Hora)
			}
			reminders = append(reminders, models.Reminder{
				CitaID:     c.ID,
				PacienteID: c.PacienteID,
				DoctorID:   c.DoctorID,
				Fecha:      c.Fecha,
				Hora:       c.Hora,
				Mensaje:    mensaje,
				Prioridad:  c.Prioridad,
			})
		}
	}

	if reminders == nil {
		reminders = []models.Reminder{}
	}
	writeJSON(w, http.StatusOK, reminders)
}

func (h *Handlers) HistorialCitaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cita-medica/historial/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de cita requerido")
		return
	}

	var eventos []models.EventoHistorial
	for _, e := range h.Store.Historial {
		if e.CitaID == id {
			eventos = append(eventos, e)
		}
	}

	if eventos == nil {
		eventos = []models.EventoHistorial{}
	}
	writeJSON(w, http.StatusOK, eventos)
}

func (h *Handlers) MetricasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	h.Store.ActualizarMetricas()
	m := models.Metricas{
		CitasCreadasHoy:     h.Store.TodayCreated,
		CitasCanceladasHoy:  h.Store.TodayCancelled,
		CitasCompletadasHoy: h.Store.TodayCompleted,
		TotalCitasActivas:   h.Store.ContarCitasActivas(),
	}
	writeJSON(w, http.StatusOK, m)
}

func (h *Handlers) CitasRecurrentesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	var req models.RecurrenteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if req.PacienteID == "" || req.DoctorID == "" || req.FechaInicio == "" || req.Hora == "" {
		writeError(w, http.StatusBadRequest, "Los campos paciente_id, doctor_id, fecha_inicio y hora son obligatorios")
		return
	}

	if req.Intervalo < 1 {
		writeError(w, http.StatusBadRequest, "El intervalo_dias debe ser mayor a 0")
		return
	}
	if req.Cantidad < 1 || req.Cantidad > 12 {
		writeError(w, http.StatusBadRequest, "La cantidad debe ser entre 1 y 12")
		return
	}

	if h.Store.FindPacienteByID(req.PacienteID) == nil {
		writeError(w, http.StatusBadRequest, "El paciente no existe")
		return
	}
	if h.Store.FindDoctorByID(req.DoctorID) == nil {
		writeError(w, http.StatusBadRequest, "El doctor no existe")
		return
	}

	if req.Prioridad == "" {
		auto := storage.DetectarPrioridadAutomatica(req.Motivo)
		if auto != "" {
			req.Prioridad = auto
		} else {
			req.Prioridad = "normal"
		}
	}

	fechaBase, err := time.Parse("2006-01-02", req.FechaInicio)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Formato de fecha invalido. Usa YYYY-MM-DD")
		return
	}

	var citasCreadas []models.Cita
	for i := 0; i < req.Cantidad; i++ {
		fecha := fechaBase.AddDate(0, 0, i*req.Intervalo)
		fechaStr := fecha.Format("2006-01-02")

		if h.Store.EsFechaPasada(fechaStr, req.Hora) {
			continue
		}
		if h.Store.MedicoOcupado(req.DoctorID, fechaStr, req.Hora, "") {
			continue
		}
		if h.Store.PacienteOcupado(req.PacienteID, fechaStr, req.Hora, "") {
			continue
		}

		nuevaCita := models.Cita{
			ID:          h.Store.GenerateID(),
			PacienteID:  storage.Sanitizar(req.PacienteID),
			DoctorID:    storage.Sanitizar(req.DoctorID),
			Fecha:       fechaStr,
			Hora:        req.Hora,
			Estado:      "pendiente",
			Prioridad:   req.Prioridad,
			Motivo:      storage.Sanitizar(req.Motivo),
			NotasMedico: "",
		}
		h.Store.Citas = append(h.Store.Citas, nuevaCita)
		h.Store.RegistrarHistorial(nuevaCita.ID, "creacion_recurrente", "", "pendiente", "")
		citasCreadas = append(citasCreadas, nuevaCita)
	}

	h.Store.ActualizarMetricas()
	h.Store.TodayCreated += len(citasCreadas)

	logger.LogJSON("INFO", fmt.Sprintf("Creadas %d citas recurrentes para paciente %s", len(citasCreadas), req.PacienteID), "citas_recurrentes", r.URL.Path, "")
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"mensaje":       fmt.Sprintf("Se crearon %d citas recurrentes", len(citasCreadas)),
		"citas_creadas": citasCreadas,
	})
}

func (h *Handlers) ListarPacientesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}
	writeJSON(w, http.StatusOK, h.Store.Pacientes)
}

func (h *Handlers) ListarDoctoresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}
	writeJSON(w, http.StatusOK, h.Store.Doctores)
}
