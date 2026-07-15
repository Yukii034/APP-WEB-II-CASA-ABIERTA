package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cuidabien/estado-animo/models"
	"cuidabien/estado-animo/services"
	"cuidabien/estado-animo/storage"
)

// HealthHandler verifica que el servicio esté activo
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// EstadoAnimoHandler maneja GET y POST
func EstadoAnimoHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {

	case http.MethodGet:

		storage.Mutex.Lock()
		defer storage.Mutex.Unlock()

		json.NewEncoder(w).Encode(storage.BaseDeDatos)

	case http.MethodPost:

		var nuevo models.RegistroEstadoAnimo

		if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		if nuevo.Nivel < 1 || nuevo.Nivel > 5 {
			http.Error(w, "El nivel debe estar entre 1 y 5", http.StatusBadRequest)
			return
		}

		storage.Mutex.Lock()
		defer storage.Mutex.Unlock()

		storage.UltimoID++
		nuevo.ID = fmt.Sprintf("%d", storage.UltimoID)

		if nuevo.Fecha == "" {
			nuevo.Fecha = time.Now().Format("2006-01-02")
		}

		storage.BaseDeDatos = append(storage.BaseDeDatos, nuevo)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(nuevo)

	default:

		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// AlertasHandler devuelve una alerta si detecta ánimo bajo persistente
func AlertasHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	alerta := services.GenerarAlerta(storage.BaseDeDatos)

	if alerta.GenerarAlerta {
		log.Printf("[ALERTA] %s", alerta.Mensaje)
	}

	json.NewEncoder(w).Encode(alerta)
}