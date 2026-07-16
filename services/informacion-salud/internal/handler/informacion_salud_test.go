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

// crearRegistro es un helper para no repetir el POST en cada test que
// necesita partir de un registro ya existente.
func crearRegistro(t *testing.T, router http.Handler, body string) model.InformacionSalud {
	t.Helper()

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
	return creado
}

func TestHealthResponde200(t *testing.T) {
	router := nuevoRouterDePrueba()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("esperaba status 200, obtuve %d", rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("esperaba status 'ok' en el body, obtuve %v", body)
	}
}

func TestCrearYLuegoObtener(t *testing.T) {
	router := nuevoRouterDePrueba()

	creado := crearRegistro(t, router, `{"nombre_paciente": "Juan", "diagnosticos": ["gripe"]}`)

	req := httptest.NewRequest(http.MethodGet, "/api/informacion-salud/"+creado.ID, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("esperaba status 200 al consultar, obtuve %d: %s", rr.Code, rr.Body.String())
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("esperaba Content-Type application/json, obtuve %s", ct)
	}

	var obtenido model.InformacionSalud
	if err := json.NewDecoder(rr.Body).Decode(&obtenido); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if obtenido.NombrePaciente != "Juan" {
		t.Errorf("esperaba nombre 'Juan', obtuve %s", obtenido.NombrePaciente)
	}
	if len(obtenido.Diagnosticos) != 1 || obtenido.Diagnosticos[0] != "gripe" {
		t.Errorf("esperaba diagnóstico 'gripe', obtuve %v", obtenido.Diagnosticos)
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

func TestListar(t *testing.T) {
	t.Run("lista vacía al inicio devuelve [] no null", func(t *testing.T) {
		router := nuevoRouterDePrueba()

		req := httptest.NewRequest(http.MethodGet, "/api/informacion-salud", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("esperaba status 200, obtuve %d", rr.Code)
		}
		if body := strings.TrimSpace(rr.Body.String()); body != "[]" {
			t.Errorf("esperaba '[]' con la lista vacía, obtuve %q", body)
		}
	})

	t.Run("devuelve los registros creados", func(t *testing.T) {
		router := nuevoRouterDePrueba()
		crearRegistro(t, router, `{"nombre_paciente": "Juan"}`)
		crearRegistro(t, router, `{"nombre_paciente": "Ana"}`)

		req := httptest.NewRequest(http.MethodGet, "/api/informacion-salud", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		var lista []model.InformacionSalud
		if err := json.NewDecoder(rr.Body).Decode(&lista); err != nil {
			t.Fatalf("no se pudo decodificar la respuesta: %v", err)
		}
		if len(lista) != 2 {
			t.Errorf("esperaba 2 registros, obtuve %d", len(lista))
		}
	})
}

func TestActualizar(t *testing.T) {
	t.Run("actualiza y persiste el cambio", func(t *testing.T) {
		router := nuevoRouterDePrueba()
		creado := crearRegistro(t, router, `{"nombre_paciente": "Juan"}`)

		body := `{"alergias": ["penicilina"]}`
		req := httptest.NewRequest(http.MethodPut, "/api/informacion-salud/"+creado.ID, strings.NewReader(body))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("esperaba status 200 al actualizar, obtuve %d: %s", rr.Code, rr.Body.String())
		}

		var actualizado model.InformacionSalud
		if err := json.NewDecoder(rr.Body).Decode(&actualizado); err != nil {
			t.Fatalf("no se pudo decodificar la respuesta: %v", err)
		}
		if len(actualizado.Alergias) != 1 || actualizado.Alergias[0] != "penicilina" {
			t.Errorf("esperaba alergia 'penicilina', obtuve %v", actualizado.Alergias)
		}
		// El nombre no se envió en el PUT, así que debe mantenerse.
		if actualizado.NombrePaciente != "Juan" {
			t.Errorf("esperaba que el nombre se mantuviera 'Juan', obtuve %s", actualizado.NombrePaciente)
		}

		// Confirma que el cambio quedó guardado, no solo en la respuesta del PUT.
		req = httptest.NewRequest(http.MethodGet, "/api/informacion-salud/"+creado.ID, nil)
		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		var obtenido model.InformacionSalud
		json.NewDecoder(rr.Body).Decode(&obtenido)
		if len(obtenido.Alergias) != 1 {
			t.Errorf("esperaba que el cambio persistiera tras el PUT, obtuve %v", obtenido.Alergias)
		}
	})

	t.Run("id inexistente devuelve 404", func(t *testing.T) {
		router := nuevoRouterDePrueba()

		req := httptest.NewRequest(http.MethodPut, "/api/informacion-salud/no-existe", strings.NewReader(`{}`))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("esperaba status 404, obtuve %d", rr.Code)
		}
	})

	t.Run("body invalido devuelve 400", func(t *testing.T) {
		router := nuevoRouterDePrueba()
		creado := crearRegistro(t, router, `{"nombre_paciente": "Juan"}`)

		req := httptest.NewRequest(http.MethodPut, "/api/informacion-salud/"+creado.ID, strings.NewReader("esto no es json"))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("esperaba status 400, obtuve %d", rr.Code)
		}
	})
}
