package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"cuidabien/actividad-fisica/internal/repository"
	"cuidabien/actividad-fisica/internal/service"
)


func setup() *ActividadFisicaHandler {
	repo := repository.NuevaMemoriaRepository()
	svc := service.NuevoActividadFisicaService(repo)
	return NuevoActividadFisicaHandler(svc)
}

func TestListar(t *testing.T) {
	h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/actividad-fisica", nil)
	rec := httptest.NewRecorder()


	h.Listar(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Se esperaba status 200, se obtuvo %d", rec.Code)
	}
}

func TestCrear_BadRequest(t *testing.T) {
	h := setup()
	req := httptest.NewRequest(http.MethodPost, "/api/actividad-fisica", bytes.NewBuffer([]byte(`{invalid}`)))
	rec := httptest.NewRecorder()

	h.Crear(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Se esperaba status 400, se obtuvo %d", rec.Code)
	}
}

func TestObtener_NotFound(t *testing.T) {
	h := setup()
	rec := httptest.NewRecorder()

	h.Obtener(rec, "999")

	if rec.Code != http.StatusNotFound {
		t.Errorf("Se esperaba status 404, se obtuvo %d", rec.Code)
	}
}

func TestEliminar_NotFound(t *testing.T) {
	h := setup()
	rec := httptest.NewRecorder()

	h.Eliminar(rec, "999")

	if rec.Code != http.StatusNotFound {
		t.Errorf("Se esperaba status 404, se obtuvo %d", rec.Code)
	}
}