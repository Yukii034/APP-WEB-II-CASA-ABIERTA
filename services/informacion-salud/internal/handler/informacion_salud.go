package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"cuidabien/informacion-salud/internal/model"
	"cuidabien/informacion-salud/internal/service"
)

// InformacionSaludHandler traduce peticiones HTTP hacia el service y
// las respuestas del service de vuelta a JSON. No contiene lógica de
// negocio propia.
type InformacionSaludHandler struct {
	service *service.InformacionSaludService
}

// NuevoInformacionSaludHandler construye el adaptador HTTP del servicio.
func NuevoInformacionSaludHandler(s *service.InformacionSaludService) *InformacionSaludHandler {
	return &InformacionSaludHandler{service: s}
}

// Health responde el endpoint obligatorio usado por el gateway y por
// docker-compose para saber si el servicio está vivo.
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /api/informacion-salud -> lista todas las fichas registradas
func (h *InformacionSaludHandler) Listar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.service.Listar())
}

// POST /api/informacion-salud -> crea una nueva ficha
// Body esperado: {"nombre_paciente": "...", "diagnosticos": [...], "alergias": [...],
//                 "enfermedades_cronicas": [...], "antecedentes_medicos": [...]}
func (h *InformacionSaludHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var entrada model.EntradaInformacionSalud
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "Body inválido, se espera un JSON con la información de salud", http.StatusBadRequest)
		return
	}

	nuevo := h.service.Crear(entrada)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevo)
}

// GET /api/informacion-salud/{id} -> consulta la ficha de un paciente
func (h *InformacionSaludHandler) Obtener(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	existente, ok := h.service.Obtener(id)
	if !ok {
		http.Error(w, "No existe información de salud para ese id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existente)
}

// PUT /api/informacion-salud/{id} -> actualiza la ficha de un paciente
func (h *InformacionSaludHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var entrada model.EntradaInformacionSalud
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "Body inválido, se espera un JSON con la información de salud", http.StatusBadRequest)
		return
	}

	actualizado, ok := h.service.Actualizar(id, entrada)
	if !ok {
		http.Error(w, "No existe información de salud para ese id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actualizado)
}
