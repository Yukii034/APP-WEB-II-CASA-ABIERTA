package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Estructura de ejemplo, cada equipo la cambia según lo que necesite
type Item struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}

// Datos de ejemplo en memoria (luego pueden usar una base de datos si quieren)
var items = []Item{
	{ID: "1", Nombre: "Ejemplo 1"},
	{ID: "2", Nombre: "Ejemplo 2"},
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/items", itemsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // puerto por defecto dentro del contenedor
	}

	log.Printf("Servicio corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Endpoint obligatorio: usado por el gateway y por docker-compose
// para saber si el servicio está vivo
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Endpoint de ejemplo, cada equipo lo reemplaza por su lógica real
func itemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
