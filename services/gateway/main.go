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
	http.HandleFunc("/api/reportes-medicos/", reportesHandler)
	http.HandleFunc("/api/cita-medica", citasHandler)
	http.HandleFunc("/api/cita-medica/", citasHandler)

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

// Proxy al microservicio de citas medicas.
// Reenvia la peticion completa (metodo, body, path) al servicio de citas.
func citasHandler(w http.ResponseWriter, r *http.Request) {
	citasURL := os.Getenv("CITAS_URL")
	if citasURL == "" {
		http.Error(w, "CITAS_URL no configurada", http.StatusInternalServerError)
		return
	}

	destino := citasURL + r.URL.Path
	if r.URL.RawQuery != "" {
		destino += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequest(r.Method, destino, r.Body)
	if err != nil {
		http.Error(w, "Error al crear la peticion", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error al contactar el servicio de citas", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error al leer la respuesta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// Proxy al microservicio de reportes medicos.
func reportesHandler(w http.ResponseWriter, r *http.Request) {
	reportesURL := os.Getenv("REPORTES_URL")
	if reportesURL == "" {
		http.Error(w, "REPORTES_URL no configurada", http.StatusInternalServerError)
		return
	}

	destino := reportesURL + r.URL.Path
	if r.URL.RawQuery != "" {
		destino += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequest(r.Method, destino, r.Body)
	if err != nil {
		http.Error(w, "Error al crear la peticion", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error al contactar el servicio de reportes", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error al leer la respuesta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
