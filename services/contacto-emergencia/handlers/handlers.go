package handlers

import (
	"cuidabien/contacto-emergencia/models"
	"cuidabien/contacto-emergencia/storage"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Handlers struct {
	Store *storage.Store
}

func New(store *storage.Store) *Handlers {
	return &Handlers{Store: store}
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func sanitizar(s string) string {
	return strings.TrimSpace(s)
}

// --- Health ---

func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// --- Contactos.. ---

// GET /api/contacts (soporta ?paciente_id= para filtrar)
func (h *Handlers) ListarContactosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	pacienteID := r.URL.Query().Get("paciente_id")
	if pacienteID != "" {
		writeJSON(w, http.StatusOK, h.Store.ContactosPorPaciente(pacienteID))
		return
	}

	writeJSON(w, http.StatusOK, h.Store.Contactos)
}

// POST /api/contacts
func (h *Handlers) CrearContactoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	var req models.CrearContactoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	req.PacienteID = sanitizar(req.PacienteID)
	req.Nombre = sanitizar(req.Nombre)
	req.Telefono = sanitizar(req.Telefono)
	req.Parentesco = sanitizar(req.Parentesco)

	if req.PacienteID == "" || req.Nombre == "" || req.Telefono == "" || req.Parentesco == "" {
		writeError(w, http.StatusBadRequest, "paciente_id, nombre, telefono y parentesco son obligatorios")
		return
	}

	if req.Prioridad <= 0 {
		req.Prioridad = 1
	}

	contacto := models.Contacto{
		ID:         h.Store.GenerateContactoID(),
		PacienteID: req.PacienteID,
		Nombre:     req.Nombre,
		Telefono:   req.Telefono,
		Parentesco: req.Parentesco,
		Prioridad:  req.Prioridad,
		Principal:  req.Principal,
	}

	h.Store.Contactos = append(h.Store.Contactos, contacto)
	writeJSON(w, http.StatusCreated, contacto)
}

// GET /api/contacts/{id}
func (h *Handlers) ObtenerContactoHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/contacts/")
	contacto := h.Store.FindContactoByID(id)
	if contacto == nil {
		writeError(w, http.StatusNotFound, "Contacto no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, contacto)
}

// PUT /api/contacts/{id}
func (h *Handlers) ActualizarContactoHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/contacts/")
	contacto := h.Store.FindContactoByID(id)
	if contacto == nil {
		writeError(w, http.StatusNotFound, "Contacto no encontrado")
		return
	}

	var req models.ActualizarContactoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if sanitizar(req.Nombre) != "" {
		contacto.Nombre = sanitizar(req.Nombre)
	}
	if sanitizar(req.Telefono) != "" {
		contacto.Telefono = sanitizar(req.Telefono)
	}
	if sanitizar(req.Parentesco) != "" {
		contacto.Parentesco = sanitizar(req.Parentesco)
	}
	if req.Prioridad > 0 {
		contacto.Prioridad = req.Prioridad
	}
	if req.Principal != nil {
		contacto.Principal = *req.Principal
	}

	writeJSON(w, http.StatusOK, contacto)
}

// DELETE /api/contacts/{id}
func (h *Handlers) EliminarContactoHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/contacts/")
	if !h.Store.EliminarContacto(id) {
		writeError(w, http.StatusNotFound, "Contacto no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "Contacto eliminado"})
}

// --- Alertas ---

// POST /api/alerts -> crea la alerta y "notifica" (simulado) a los contactos del paciente
func (h *Handlers) CrearAlertaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	var req models.CrearAlertaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	req.PacienteID = sanitizar(req.PacienteID)
	req.Mensaje = sanitizar(req.Mensaje)
	req.Nivel = strings.ToLower(sanitizar(req.Nivel))

	if req.PacienteID == "" || req.Mensaje == "" {
		writeError(w, http.StatusBadRequest, "paciente_id y mensaje son obligatorios")
		return
	}
	if req.Nivel == "" {
		req.Nivel = "moderado"
	}
	if !h.Store.NivelValido(req.Nivel) {
		writeError(w, http.StatusBadRequest, "nivel debe ser: leve, moderado o critico")
		return
	}

	contactos := h.Store.ContactosPorPaciente(req.PacienteID)
	if len(contactos) == 0 {
		writeError(w, http.StatusBadRequest, "El paciente no tiene contactos de emergencia registrados")
		return
	}

	notificados := h.Store.NotificarContactos(req.PacienteID)

	alerta := models.Alerta{
		ID:                   h.Store.GenerateAlertaID(),
		PacienteID:           req.PacienteID,
		Mensaje:              req.Mensaje,
		Nivel:                req.Nivel,
		Estado:               "activa",
		Timestamp:            time.Now().Format(time.RFC3339),
		ContactosNotificados: notificados,
	}

	h.Store.Alertas = append(h.Store.Alertas, alerta)
	h.Store.RegistrarHistorial(alerta.ID, "creada", "Alerta activada y contactos notificados")

	h.Store.ActualizarContadoresDiarios()
	h.Store.TodayActivadas++

	writeJSON(w, http.StatusCreated, alerta)
}

// GET /api/alerts
func (h *Handlers) ListarAlertasHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.Alertas)
}

// GET /api/alerts/{id}
func (h *Handlers) ObtenerAlertaHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/alerts/")
	alerta := h.Store.FindAlertaByID(id)
	if alerta == nil {
		writeError(w, http.StatusNotFound, "Alerta no encontrada")
		return
	}
	writeJSON(w, http.StatusOK, alerta)
}

// PATCH /api/alerts/{id}/attend -> marca la alerta como atendida
func (h *Handlers) AtenderAlertaHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/alerts/"), "/attend")
	alerta := h.Store.FindAlertaByID(id)
	if alerta == nil {
		writeError(w, http.StatusNotFound, "Alerta no encontrada")
		return
	}
	if alerta.Estado != "activa" {
		writeError(w, http.StatusConflict, "Solo se puede atender una alerta activa")
		return
	}

	var req models.AtenderAlertaRequest
	json.NewDecoder(r.Body).Decode(&req) // notas es opcional, no es error si viene vacio

	alerta.Estado = "atendida"
	h.Store.RegistrarHistorial(alerta.ID, "atendida", req.Notas)

	h.Store.ActualizarContadoresDiarios()
	h.Store.TodayAtendidas++

	writeJSON(w, http.StatusOK, alerta)
}

// DELETE /api/alerts/{id} -> cancela una alerta activa
func (h *Handlers) CancelarAlertaHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/alerts/")
	alerta := h.Store.FindAlertaByID(id)
	if alerta == nil {
		writeError(w, http.StatusNotFound, "Alerta no encontrada")
		return
	}
	if alerta.Estado != "activa" {
		writeError(w, http.StatusConflict, "Solo se puede cancelar una alerta activa")
		return
	}

	alerta.Estado = "cancelada"
	h.Store.RegistrarHistorial(alerta.ID, "cancelada", "")

	writeJSON(w, http.StatusOK, alerta)
}

// GET /api/alerts/history/{id}
func (h *Handlers) HistorialAlertaHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/alerts/history/")
	if h.Store.FindAlertaByID(id) == nil {
		writeError(w, http.StatusNotFound, "Alerta no encontrada")
		return
	}
	writeJSON(w, http.StatusOK, h.Store.HistorialPorAlerta(id))
}

// --- Metricas ---

// GET /api/metrics
func (h *Handlers) MetricasHandler(w http.ResponseWriter, r *http.Request) {
	h.Store.ActualizarContadoresDiarios()

	m := models.Metricas{
		TotalContactos:      len(h.Store.Contactos),
		AlertasActivasHoy:   h.Store.TodayActivadas,
		AlertasAtendidasHoy: h.Store.TodayAtendidas,
		TotalAlertas:        len(h.Store.Alertas),
	}

	writeJSON(w, http.StatusOK, m)
}
