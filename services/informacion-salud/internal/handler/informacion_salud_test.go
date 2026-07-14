package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"cuidabien/informacion-salud/internal/model"
	"cuidabien/informacion-salud/internal/repository"
	"cuidabien/informacion-salud/internal/service"
)

func nuevoRouterDePrueba() http.Handler {
	repo := repository.NuevaMemoriaRepository()
	svc := service.NuevoInformacionSaludService(repo)
	h := NuevoInformacionSaludHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", Health)
	r.Route("/api/informacion-salud", func(r chi.Router) {
		r.Get("/", h.Listar)
		r.Post("/", h.Crear)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.Obtener)
			r.Put("/", h.Actualizar)
		})
	})
	return r
}

func TestHealthResponde200(t *testing.T) {
	router := nuevoRouterDePrueba()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("esperaba status 200, obtuve %d", rr.Code)
	}
}

func TestCrearYLuegoObtener(t *testing.T) {
	router := nuevoRouterDePrueba()

	body := `{"nombre_paciente": "Juan", "diagnosticos": ["gripe"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/informacion-salud", strings.NewReader(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("esperaba status 201 al crear, obtuve %d: %s", rr.Code, rr.Body.String())
	}

	var creado model.InformacionSalud
	if err := json.NewDecoder(rr.Body).Decode(&creado); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/informacion-salud/"+creado.ID, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("esperaba status 200 al consultar, obtuve %d: %s", rr.Code, rr.Body.String())
	}
}

func TestObtenerIDInexistenteDevuelve404(t *testing.T) {
	router := nuevoRouterDePrueba()

	req := httptest.NewRequest(http.MethodGet, "/api/informacion-salud/no-existe", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("esperaba status 404, obtuve %d", rr.Code)
	}
}

func TestCrearConBodyInvalidoDevuelve400(t *testing.T) {
	router := nuevoRouterDePrueba()

	req := httptest.NewRequest(http.MethodPost, "/api/informacion-salud", strings.NewReader("esto no es json"))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("esperaba status 400, obtuve %d", rr.Code)
	}
}
