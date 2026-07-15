package handler

import (
	"encoding/json"
	"net/http"

	"cuidabien/actividad-fisica/internal/model"
	"cuidabien/actividad-fisica/internal/service"
)

type ActividadFisicaHandler struct {
	service *service.ActividadFisicaService
}

func NuevoActividadFisicaHandler(s *service.ActividadFisicaService) *ActividadFisicaHandler {
	return &ActividadFisicaHandler{service: s}
}

func Health(w http.ResponseWriter, _ *http.Request) {
	responderJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "actividad-fisica"})
}

func (h *ActividadFisicaHandler) Listar(w http.ResponseWriter, _ *http.Request) {
	responderJSON(w, http.StatusOK, h.service.Listar())
}

func (h *ActividadFisicaHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var entrada model.EntradaActividadFisica
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, http.StatusBadRequest, "Body inválido: se espera un JSON")
		return
	}
	actividad, err := h.service.Crear(entrada)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, actividad)
}

func (h *ActividadFisicaHandler) Obtener(w http.ResponseWriter, id string) {
	actividad, ok := h.service.Obtener(id)
	if !ok {
		responderError(w, http.StatusNotFound, "No existe una actividad con ese id")
		return
	}
	responderJSON(w, http.StatusOK, actividad)
}

func (h *ActividadFisicaHandler) Actualizar(w http.ResponseWriter, r *http.Request, id string) {
	var entrada model.EntradaActividadFisica
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, http.StatusBadRequest, "Body inválido: se espera un JSON")
		return
	}
	actividad, err, existe := h.service.Actualizar(id, entrada)
	if !existe {
		responderError(w, http.StatusNotFound, "No existe una actividad con ese id")
		return
	}
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, actividad)
}

func (h *ActividadFisicaHandler) Eliminar(w http.ResponseWriter, id string) {
	if !h.service.Eliminar(id) {
		responderError(w, http.StatusNotFound, "No existe una actividad con ese id")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func responderJSON(w http.ResponseWriter, estado int, datos any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(estado)
	_ = json.NewEncoder(w).Encode(datos)
}

func responderError(w http.ResponseWriter, estado int, mensaje string) {
	responderJSON(w, estado, map[string]string{"error": mensaje})
}
