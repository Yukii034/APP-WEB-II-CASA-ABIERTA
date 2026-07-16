package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"cuidabien/recordatorios-medicamentos/internal/model"
	"cuidabien/recordatorios-medicamentos/internal/service"
)

// RecordatorioMedicamentoHandler traduce peticiones HTTP
// hacia el service. No contiene lógica de negocio.
type RecordatorioMedicamentoHandler struct {
	service *service.RecordatorioMedicamentoService
}

func NuevoRecordatorioMedicamentoHandler(
	s *service.RecordatorioMedicamentoService,
) *RecordatorioMedicamentoHandler {
	return &RecordatorioMedicamentoHandler{
		service: s,
	}
}

// Health responde el endpoint obligatorio.
func Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"status": "ok",
		},
	)
}

// GET /api/recordatorio-medicamentos
func (h *RecordatorioMedicamentoHandler) Listar(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		h.service.Listar(),
	)
}

// POST /api/recordatorio-medicamentos
func (h *RecordatorioMedicamentoHandler) Crear(
	w http.ResponseWriter,
	r *http.Request,
) {
	var entrada model.EntradaRecordatorioMedicamento

	if err := json.NewDecoder(
		r.Body,
	).Decode(&entrada); err != nil {
		http.Error(
			w,
			"Body inválido, se espera un JSON "+
				"con los datos del recordatorio",
			http.StatusBadRequest,
		)

		return
	}

	nuevo, err := h.service.Crear(entrada)

	if err != nil {
		responderErrorService(w, err)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(nuevo)
}

// GET /api/recordatorio-medicamentos/{id}
func (h *RecordatorioMedicamentoHandler) Obtener(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := chi.URLParam(r, "id")

	existente, ok := h.service.Obtener(id)

	if !ok {
		http.Error(
			w,
			"No existe un recordatorio para ese id",
			http.StatusNotFound,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(existente)
}

// PUT /api/recordatorio-medicamentos/{id}
func (h *RecordatorioMedicamentoHandler) Actualizar(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := chi.URLParam(r, "id")

	var entrada model.EntradaRecordatorioMedicamento

	if err := json.NewDecoder(
		r.Body,
	).Decode(&entrada); err != nil {
		http.Error(
			w,
			"Body inválido, se espera un JSON "+
				"con los datos del recordatorio",
			http.StatusBadRequest,
		)

		return
	}

	actualizado, ok, err :=
		h.service.Actualizar(id, entrada)

	if err != nil {
		responderErrorService(w, err)
		return
	}

	if !ok {
		http.Error(
			w,
			"No existe un recordatorio para ese id",
			http.StatusNotFound,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(actualizado)
}

// DELETE /api/recordatorio-medicamentos/{id}
func (h *RecordatorioMedicamentoHandler) Eliminar(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := chi.URLParam(r, "id")

	if !h.service.Eliminar(id) {
		http.Error(
			w,
			"No existe un recordatorio para ese id",
			http.StatusNotFound,
		)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /api/recordatorio-medicamentos/{id}/estado
func (h *RecordatorioMedicamentoHandler) CambiarEstado(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := chi.URLParam(r, "id")

	var entrada model.EntradaEstadoRecordatorio

	if err := json.NewDecoder(
		r.Body,
	).Decode(&entrada); err != nil ||
		entrada.Activo == nil {
		http.Error(
			w,
			"Body inválido, se espera {\"activo\": true}",
			http.StatusBadRequest,
		)

		return
	}

	actualizado, ok :=
		h.service.CambiarEstado(
			id,
			*entrada.Activo,
		)

	if !ok {
		http.Error(
			w,
			"No existe un recordatorio para ese id",
			http.StatusNotFound,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(actualizado)
}

// POST /api/recordatorio-medicamentos/verificar
func (h *RecordatorioMedicamentoHandler) Verificar(
	w http.ResponseWriter,
	r *http.Request,
) {
	var entrada model.EntradaVerificacion

	if err := json.NewDecoder(
		r.Body,
	).Decode(&entrada); err != nil {
		http.Error(
			w,
			"Body inválido, se espera "+
				"{\"hora\": \"08:00\"}",
			http.StatusBadRequest,
		)

		return
	}

	resultado, err :=
		h.service.VerificarHora(
			entrada.Hora,
		)

	if err != nil {
		responderErrorService(w, err)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(resultado)
}

func responderErrorService(
	w http.ResponseWriter,
	err error,
) {
	if errors.Is(
		err,
		service.ErrDatosInvalidos,
	) {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	http.Error(
		w,
		"Error interno del servidor",
		http.StatusInternalServerError,
	)
}
