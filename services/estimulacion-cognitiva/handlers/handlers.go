package handlers

import (
	"encoding/json"
	"net/http"

	"cuidabien/estimulacion-cognitiva/storage"
)

type Handlers struct {
	store *storage.Storage
}

func New(store *storage.Storage) *Handlers {
	return &Handlers{store: store}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) Listar(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.Listar())
}

type agregarRequest struct {
	Tipo string `json:"tipo"`
}

func (h *Handlers) Agregar(w http.ResponseWriter, r *http.Request) {
	var req agregarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido, se esperaba JSON"})
		return
	}

	ejercicio, err := h.store.Agregar(req.Tipo)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, ejercicio)
}

func (h *Handlers) Resumen(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.Resumen())
}

func (h *Handlers) Reset(w http.ResponseWriter, r *http.Request) {
	h.store.Reset()
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
