package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"cuidabien/cuidadores/internal/model"
	"cuidabien/cuidadores/internal/repository"
	"cuidabien/cuidadores/internal/service"
)

type CuidadorHandler struct {
	service *service.CuidadorService
}

func NuevoCuidadorHandler(s *service.CuidadorService) *CuidadorHandler {
	return &CuidadorHandler{service: s}
}

// Health es el endpoint obligatorio usado por el gateway y docker-compose.
func (h *CuidadorHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *CuidadorHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var c model.Cuidador
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		responderJSONError(w, http.StatusBadRequest, "json invalido")
		return
	}

	creado, err := h.service.Crear(c)
	if err != nil {
		responderJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(creado)
}

func (h *CuidadorHandler) Listar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.service.Listar())
}

func (h *CuidadorHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := h.service.ObtenerPorID(id)
	if err != nil {
		h.responderError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (h *CuidadorHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var c model.Cuidador
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		responderJSONError(w, http.StatusBadRequest, "json invalido")
		return
	}

	actualizado, err := h.service.Actualizar(id, c)
	if err != nil {
		h.responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actualizado)
}

func (h *CuidadorHandler) Eliminar(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.Eliminar(id); err != nil {
		h.responderError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CuidadorHandler) ListarPorPaciente(w http.ResponseWriter, r *http.Request) {
	pacienteID := r.PathValue("pacienteId")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.service.ObtenerPorPaciente(pacienteID))
}

func (h *CuidadorHandler) responderError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNoEncontrado) {
		responderJSONError(w, http.StatusNotFound, "cuidador no encontrado")
		return
	}
	responderJSONError(w, http.StatusBadRequest, err.Error())
}

func responderJSONError(w http.ResponseWriter, status int, mensaje string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": mensaje})
}
