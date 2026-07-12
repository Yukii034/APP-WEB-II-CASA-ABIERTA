package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// RegistroComida representa una comida registrada.
type RegistroComida struct {
	ID          string    `json:"id"`
	TipoComida  string    `json:"tipo_comida"` // desayuno | almuerzo | cena | merienda
	Descripcion string    `json:"descripcion,omitempty"`
	Hora        time.Time `json:"hora"`
}

// EstadoComida indica si una comida del día ya se registró o si se saltó.
type EstadoComida struct {
	TipoComida string `json:"tipo_comida"`
	Registrada bool   `json:"registrada"`
	Saltada    bool   `json:"saltada"`
	HoraLimite string `json:"hora_limite"`
}

// Resumen es el estado del día: qué comidas van, cuáles faltan y si hay
// alguna que ya se saltó (pasó la hora límite y no se registró).
type Resumen struct {
	Comidas       []EstadoComida `json:"comidas"`
	ComidasHechas int            `json:"comidas_hechas"`
	ComidasTotal  int            `json:"comidas_total"`
	HaySaltadas   bool           `json:"hay_saltadas"`
	Mensaje       string         `json:"mensaje,omitempty"`
}

type comidaEsperada struct {
	tipo       string
	horaLimite string // formato "HH:MM", hora del día tras la cual se considera saltada
}

var (
	mu          sync.Mutex
	registros   []RegistroComida
	siguienteID = 1

	comidasEsperadas = []comidaEsperada{
		{tipo: "desayuno", horaLimite: "10:00"},
		{tipo: "almuerzo", horaLimite: "15:00"},
		{tipo: "cena", horaLimite: "21:00"},
	}
)

func main() {
	if v := os.Getenv("DESAYUNO_HASTA"); v != "" {
		comidasEsperadas[0].horaLimite = v
	}
	if v := os.Getenv("ALMUERZO_HASTA"); v != "" {
		comidasEsperadas[1].horaLimite = v
	}
	if v := os.Getenv("CENA_HASTA"); v != "" {
		comidasEsperadas[2].horaLimite = v
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/alimentacion", alimentacionHandler)
	http.HandleFunc("/api/alimentacion/resumen", resumenHandler)
	http.HandleFunc("/api/alimentacion/reset", resetHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servicio de alimentación corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /api/alimentacion -> lista las comidas registradas hoy
// POST /api/alimentacion -> registra una comida { "tipo_comida": "almuerzo", "descripcion": "sopa y pollo" }
func alimentacionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(registrosDeHoy())

	case http.MethodPost:
		var entrada struct {
			TipoComida  string `json:"tipo_comida"`
			Descripcion string `json:"descripcion"`
		}
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil || entrada.TipoComida == "" {
			http.Error(w, "Body inválido, se espera {\"tipo_comida\": \"almuerzo\", \"descripcion\": \"...\"}", http.StatusBadRequest)
			return
		}

		mu.Lock()
		nuevo := RegistroComida{
			ID:          strconv.Itoa(siguienteID),
			TipoComida:  entrada.TipoComida,
			Descripcion: entrada.Descripcion,
			Hora:        time.Now(),
		}
		siguienteID++
		registros = append(registros, nuevo)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(nuevo)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /api/alimentacion/resumen -> estado de las comidas del día
func resumenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calcularResumen(registrosDeHoy(), comidasEsperadas, time.Now()))
}

// POST /api/alimentacion/reset -> limpia los registros de hoy (útil para demos)
func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	mu.Lock()
	registros = nil
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "reiniciado"})
}

func registrosDeHoy() []RegistroComida {
	mu.Lock()
	defer mu.Unlock()

	hoy := time.Now().Format("2006-01-02")
	var resultado []RegistroComida
	for _, reg := range registros {
		if reg.Hora.Format("2006-01-02") == hoy {
			resultado = append(resultado, reg)
		}
	}
	if resultado == nil {
		resultado = []RegistroComida{}
	}
	return resultado
}

// calcularResumen es una función pura (sin estado global) para que sea
// fácil de probar con go test.
func calcularResumen(regs []RegistroComida, esperadas []comidaEsperada, ahora time.Time) Resumen {
	registradas := map[string]bool{}
	for _, reg := range regs {
		registradas[reg.TipoComida] = true
	}

	var comidas []EstadoComida
	comidasHechas := 0
	haySaltadas := false

	for _, esperada := range esperadas {
		hecha := registradas[esperada.tipo]
		saltada := false

		if !hecha {
			limite, err := time.ParseInLocation("15:04", esperada.horaLimite, ahora.Location())
			if err == nil {
				limiteHoy := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), limite.Hour(), limite.Minute(), 0, 0, ahora.Location())
				if ahora.After(limiteHoy) {
					saltada = true
					haySaltadas = true
				}
			}
		} else {
			comidasHechas++
		}

		comidas = append(comidas, EstadoComida{
			TipoComida: esperada.tipo,
			Registrada: hecha,
			Saltada:    saltada,
			HoraLimite: esperada.horaLimite,
		})
	}

	mensaje := ""
	if haySaltadas {
		mensaje = "Hay una o más comidas que no se registraron a tiempo hoy."
	}

	return Resumen{
		Comidas:       comidas,
		ComidasHechas: comidasHechas,
		ComidasTotal:  len(esperadas),
		HaySaltadas:   haySaltadas,
		Mensaje:       mensaje,
	}
}
