package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"cuidabien/recordatorios-medicamentos/internal/model"
	"cuidabien/recordatorios-medicamentos/internal/repository"
	"cuidabien/recordatorios-medicamentos/internal/service"
)

func nuevoRouterDePrueba() http.Handler {
	repo := repository.NuevaMemoriaRepository()

	svc := service.NuevoRecordatorioMedicamentoService(
		repo,
		nil,
	)

	h := NuevoRecordatorioMedicamentoHandler(svc)

	r := chi.NewRouter()

	r.Get("/health", Health)

	r.Route(
		"/api/recordatorio-medicamentos",
		func(r chi.Router) {
			r.Get("/", h.Listar)
			r.Post("/", h.Crear)
			r.Post("/verificar", h.Verificar)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Obtener)
				r.Put("/", h.Actualizar)
				r.Delete("/", h.Eliminar)
				r.Patch("/estado", h.CambiarEstado)
			})
		},
	)

	return r
}

func TestHealthResponde200(t *testing.T) {
	router := nuevoRouterDePrueba()

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(
			"esperaba status 200, obtuve %d",
			rr.Code,
		)
	}
}

func TestCrearYLuegoObtener(t *testing.T) {
	router := nuevoRouterDePrueba()

	body := `{
		"adulto_mayor_id": "AM-001",
		"nombre_paciente": "María Pérez",
		"medicamento": "Losartán",
		"dosis": "1 tableta",
		"hora": "08:00",
		"frecuencia": "diaria"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/recordatorio-medicamentos",
		strings.NewReader(body),
	)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf(
			"esperaba 201, obtuve %d: %s",
			rr.Code,
			rr.Body.String(),
		)
	}

	var creado model.RecordatorioMedicamento

	if err := json.NewDecoder(
		rr.Body,
	).Decode(&creado); err != nil {
		t.Fatalf(
			"no se pudo decodificar: %v",
			err,
		)
	}

	req = httptest.NewRequest(
		http.MethodGet,
		"/api/recordatorio-medicamentos/"+creado.ID,
		nil,
	)

	rr = httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf(
			"esperaba 200, obtuve %d: %s",
			rr.Code,
			rr.Body.String(),
		)
	}
}

func TestCrearBodyInvalidoDevuelve400(
	t *testing.T,
) {
	router := nuevoRouterDePrueba()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/recordatorio-medicamentos",
		strings.NewReader("no-json"),
	)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf(
			"esperaba 400, obtuve %d",
			rr.Code,
		)
	}
}

func TestObtenerInexistenteDevuelve404(
	t *testing.T,
) {
	router := nuevoRouterDePrueba()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/recordatorio-medicamentos/no-existe",
		nil,
	)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf(
			"esperaba 404, obtuve %d",
			rr.Code,
		)
	}
}

func TestVerificarHoraDevuelve200(
	t *testing.T,
) {
	router := nuevoRouterDePrueba()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/recordatorio-medicamentos/verificar",
		strings.NewReader(
			`{"hora":"08:00"}`,
		),
	)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(
			"esperaba 200, obtuve %d: %s",
			rr.Code,
			rr.Body.String(),
		)
	}
}
