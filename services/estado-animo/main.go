package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// RegistroEstadoAnimo representa la estructura de datos que guardaremos
type RegistroEstadoAnimo struct {
	ID         string `json:"id"`
	Fecha      string `json:"fecha"`      // Formato "YYYY-MM-DD"
	Nivel      int    `json:"nivel"`      // Escala del 1 al 5 (1: Muy mal, 5: Muy bien)
	Emocion    string `json:"emocion"`    // Ej: "Feliz", "Triste", "Ansioso", "Tranquilo"
	Comentario string `json:"comentario"` // Nota opcional
}

// Alerta representa una advertencia si detectamos cambios negativos
type Alerta struct {
	GenerarAlerta bool   `json:"generar_alerta"`
	Mensaje       string `json:"mensaje"`
}

// Usamos un Mutex para evitar colisiones al escribir/leer datos concurrentemente en memoria
var (
	baseDeDatos []RegistroEstadoAnimo
	ultimoID    int
	mutex       sync.Mutex
)

func main() {
	// Inicializamos con un par de datos de prueba en memoria
	baseDeDatos = []RegistroEstadoAnimo{
		{
			ID:         "1",
			Fecha:      time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			Nivel:      4,
			Emocion:    "Tranquilo",
			Comentario: "Pasé una tarde agradable.",
		},
		{
			ID:         "2",
			Fecha:      time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			Nivel:      5,
			Emocion:    "Feliz",
			Comentario: "Me visitaron mis nietos.",
		},
	}
	ultimoID = 2

	// Registramos los endpoints de nuestro servicio
	http.HandleFunc("/health", healthHandler)                    // Requerido por el gateway/Docker
	http.HandleFunc("/api/estado-animo", estadoAnimoHandler)     // Listar (GET) y Registrar (POST)
	http.HandleFunc("/api/estado-animo/alertas", alertasHandler) // Obtener alertas de desánimo

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // puerto por defecto dentro del contenedor
	}

	log.Printf("Servicio de Estado de Ánimo corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Endpoint obligatorio: usado por el gateway y por docker-compose para saber si el servicio está vivo
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Handler principal que reemplaza a 'itemsHandler' para manejar el estado de ánimo
func estadoAnimoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// CORS básico en caso de que lo consuman directamente
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		mutex.Lock()
		json.NewEncoder(w).Encode(baseDeDatos)
		mutex.Unlock()

	case http.MethodPost:
		var nuevoRegistro RegistroEstadoAnimo
		err := json.NewDecoder(r.Body).Decode(&nuevoRegistro)
		if err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		// Validamos que el nivel sea correcto (entre 1 y 5)
		if nuevoRegistro.Nivel < 1 || nuevoRegistro.Nivel > 5 {
			http.Error(w, "El nivel de estado de ánimo debe estar entre 1 y 5", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		ultimoID++
		nuevoRegistro.ID = fmt.Sprintf("%d", ultimoID)

		// Si no envían fecha, asignamos la de hoy automáticamente
		if nuevoRegistro.Fecha == "" {
			nuevoRegistro.Fecha = time.Now().Format("2006-01-02")
		}

		baseDeDatos = append(baseDeDatos, nuevoRegistro)
		mutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(nuevoRegistro)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Handler adicional para procesar y alertar desánimo continuo
func alertasHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	alerta := Alerta{
		GenerarAlerta: false,
		Mensaje:       "El estado de ánimo se encuentra estable.",
	}

	n := len(baseDeDatos)
	// Si hay por lo menos 2 registros en el historial
	if n >= 2 {
		ultimo := baseDeDatos[n-1]
		penultimo := baseDeDatos[n-2]

		// Si los últimos dos días reportó un ánimo decaído (1 o 2)
		if ultimo.Nivel <= 2 && penultimo.Nivel <= 2 {
			alerta.GenerarAlerta = true
			alerta.Mensaje = fmt.Sprintf(
				"ALERTA: Se ha detectado un estado de ánimo bajo persistente. Últimas emociones registradas: '%s' y '%s'. Se sugiere que un cuidador contacte al adulto mayor.",
				penultimo.Emocion, ultimo.Emocion,
			)
			log.Printf("[ALERTA EMITIDA] %s", alerta.Mensaje)
		}
	}

	json.NewEncoder(w).Encode(alerta)
}
