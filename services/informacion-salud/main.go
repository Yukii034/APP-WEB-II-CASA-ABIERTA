package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// InformacionSalud representa la ficha de salud de un adulto mayor.
type InformacionSalud struct {
	ID                   string    `json:"id"`
	NombrePaciente       string    `json:"nombre_paciente,omitempty"`
	Diagnosticos         []string  `json:"diagnosticos"`
	Alergias             []string  `json:"alergias"`
	EnfermedadesCronicas []string  `json:"enfermedades_cronicas"`
	AntecedentesMedicos  []string  `json:"antecedentes_medicos"`
	ActualizadoEn        time.Time `json:"actualizado_en"`
}

// entradaInformacionSalud es el body esperado en POST y PUT.
type entradaInformacionSalud struct {
	NombrePaciente       string   `json:"nombre_paciente"`
	Diagnosticos         []string `json:"diagnosticos"`
	Alergias             []string `json:"alergias"`
	EnfermedadesCronicas []string `json:"enfermedades_cronicas"`
	AntecedentesMedicos  []string `json:"antecedentes_medicos"`
}

var (
	mu          sync.Mutex
	registros   = map[string]InformacionSalud{}
	siguienteID = 1
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", healthHandler)

	r.Route("/api/informacion-salud", func(r chi.Router) {
		r.Get("/", listarHandler)
		r.Post("/", crearHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", obtenerHandler)
			r.Put("/", actualizarHandler)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // puerto por defecto dentro del contenedor
	}

	log.Printf("Servicio de información de salud corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Endpoint obligatorio: usado por el gateway y por docker-compose
// para saber si el servicio está vivo
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /api/informacion-salud -> lista todas las fichas registradas
func listarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listarRegistros())
}

// POST /api/informacion-salud -> crea una nueva ficha
// Body esperado: {"nombre_paciente": "...", "diagnosticos": [...], "alergias": [...],
//
//	"enfermedades_cronicas": [...], "antecedentes_medicos": [...]}
func crearHandler(w http.ResponseWriter, r *http.Request) {
	var entrada entradaInformacionSalud
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "Body inválido, se espera un JSON con la información de salud", http.StatusBadRequest)
		return
	}

	mu.Lock()
	nuevo := nuevoRegistro(strconv.Itoa(siguienteID), entrada, time.Now())
	siguienteID++
	registros[nuevo.ID] = nuevo
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevo)
}

// GET /api/informacion-salud/{id} -> consulta la ficha de un paciente
func obtenerHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mu.Lock()
	existente, ok := registros[id]
	mu.Unlock()
	if !ok {
		http.Error(w, "No existe información de salud para ese id", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existente)
}

// PUT /api/informacion-salud/{id} -> actualiza la ficha de un paciente
func actualizarHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var entrada entradaInformacionSalud
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "Body inválido, se espera un JSON con la información de salud", http.StatusBadRequest)
		return
	}

	mu.Lock()
	existente, ok := registros[id]
	if !ok {
		mu.Unlock()
		http.Error(w, "No existe información de salud para ese id", http.StatusNotFound)
		return
	}
	actualizado := actualizarRegistro(existente, entrada, time.Now())
	registros[id] = actualizado
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actualizado)
}

func listarRegistros() []InformacionSalud {
	mu.Lock()
	defer mu.Unlock()

	resultado := make([]InformacionSalud, 0, len(registros))
	for _, reg := range registros {
		resultado = append(resultado, reg)
	}
	return resultado
}

// nuevoRegistro es una función pura (sin estado global) para que sea
// fácil de probar con go test.
func nuevoRegistro(id string, entrada entradaInformacionSalud, ahora time.Time) InformacionSalud {
	return InformacionSalud{
		ID:                   id,
		NombrePaciente:       entrada.NombrePaciente,
		Diagnosticos:         normalizar(entrada.Diagnosticos),
		Alergias:             normalizar(entrada.Alergias),
		EnfermedadesCronicas: normalizar(entrada.EnfermedadesCronicas),
		AntecedentesMedicos:  normalizar(entrada.AntecedentesMedicos),
		ActualizadoEn:        ahora,
	}
}

// actualizarRegistro combina la ficha existente con los campos enviados,
// sin borrar datos que no vinieron en la petición (actualización parcial).
func actualizarRegistro(existente InformacionSalud, entrada entradaInformacionSalud, ahora time.Time) InformacionSalud {
	if entrada.NombrePaciente != "" {
		existente.NombrePaciente = entrada.NombrePaciente
	}
	if entrada.Diagnosticos != nil {
		existente.Diagnosticos = normalizar(entrada.Diagnosticos)
	}
	if entrada.Alergias != nil {
		existente.Alergias = normalizar(entrada.Alergias)
	}
	if entrada.EnfermedadesCronicas != nil {
		existente.EnfermedadesCronicas = normalizar(entrada.EnfermedadesCronicas)
	}
	if entrada.AntecedentesMedicos != nil {
		existente.AntecedentesMedicos = normalizar(entrada.AntecedentesMedicos)
	}
	existente.ActualizadoEn = ahora
	return existente
}

// normalizar evita que el JSON de salida muestre "null" en vez de una
// lista vacía cuando no se envía ese campo.
func normalizar(valores []string) []string {
	if valores == nil {
		return []string{}
	}
	return valores
}
