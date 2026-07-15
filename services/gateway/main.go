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
	http.HandleFunc("/api/medicamentos", proxyHandler("MEDICAMENTOS_URL"))
	http.HandleFunc("/api/appointments", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/appointments/", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/patients", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/doctors", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/reportes", multiProxyHandler([]string{"REPORTES_URL", "CITAS_URL", "MEDICAMENTOS_URL", "ALIMENTACION_URL"}))
	http.HandleFunc("/api/reportes/", proxyHandler("REPORTES_URL"))

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

// Proxy generico: reenvia la peticion al servicio indicado en la variable de entorno.
func proxyHandler(envVar string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceURL := os.Getenv(envVar)
		if serviceURL == "" {
			http.Error(w, envVar+" no configurada", http.StatusInternalServerError)
			return
		}

		destino := serviceURL + r.URL.Path
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
			http.Error(w, "Error al contactar el servicio", http.StatusBadGateway)
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
}

// Proxy multi-servicio: intenta cada servicio en orden, devuelve el primero que responda.
func multiProxyHandler(envVars []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, envVar := range envVars {
			serviceURL := os.Getenv(envVar)
			if serviceURL == "" {
				continue
			}

			destino := serviceURL + r.URL.Path
			if r.URL.RawQuery != "" {
				destino += "?" + r.URL.RawQuery
			}

			req, err := http.NewRequest(r.Method, destino, r.Body)
			if err != nil {
				continue
			}
			req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			w.Write(body)
			return
		}
		http.Error(w, "Ningun servicio disponible", http.StatusBadGateway)
	}
}
