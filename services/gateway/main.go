package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/medicamentos", medicamentosHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Este handler NO tiene lógica propia: solo redirige (proxy) la petición
// al microservicio de medicamentos y devuelve su respuesta tal cual.
// Este es el patrón que usarían para conectar cualquier servicio nuevo.
func medicamentosHandler(w http.ResponseWriter, r *http.Request) {
	// La URL del servicio se lee de una variable de entorno, nunca hardcodeada
	medicamentosURL := os.Getenv("MEDICAMENTOS_URL")
	if medicamentosURL == "" {
		http.Error(w, "MEDICAMENTOS_URL no configurada", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get(medicamentosURL + "/api/items")
	if err != nil {
		http.Error(w, "Error al contactar el servicio de medicamentos", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error al leer la respuesta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
