package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"monitoreo-signos-vitales/internal/models"
	"monitoreo-signos-vitales/internal/service"
)

type SignosVitalesHandler struct{ service *service.SignosVitalesService }

func NuevoSignosVitalesHandler(s *service.SignosVitalesService) *SignosVitalesHandler {
	return &SignosVitalesHandler{service: s}
}

func Health(w http.ResponseWriter, _ *http.Request) {
	responderJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
func (h *SignosVitalesHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var entrada models.EntradaSignosVitales
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, http.StatusBadRequest, "body inválido, se espera JSON")
		return
	}
	registro, err := h.service.Crear(entrada)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, registro)
}
func (h *SignosVitalesHandler) Historial(w http.ResponseWriter, r *http.Request) {
	id := segmento(r.URL.Path, 2)
	if id == "" {
		responderError(w, http.StatusBadRequest, "id_adulto_mayor es obligatorio")
		return
	}
	responderJSON(w, http.StatusOK, h.service.Historial(id))
}
func (h *SignosVitalesHandler) Ultimo(w http.ResponseWriter, r *http.Request) {
	registro, ok := h.service.Ultimo(segmento(r.URL.Path, 2))
	if !ok {
		responderError(w, http.StatusNotFound, "no existen registros para este adulto mayor")
		return
	}
	responderJSON(w, http.StatusOK, registro)
}
func (h *SignosVitalesHandler) Tendencia(w http.ResponseWriter, r *http.Request) {
	dias, _ := strconv.Atoi(r.URL.Query().Get("dias"))
	puntos, err := h.service.Tendencia(segmento(r.URL.Path, 2), r.URL.Query().Get("parametro"), dias)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, puntos)
}
func segmento(path string, indice int) string {
	partes := strings.Split(strings.Trim(path, "/"), "/")
	if len(partes) > indice {
		return partes[indice]
	}
	return ""
}
