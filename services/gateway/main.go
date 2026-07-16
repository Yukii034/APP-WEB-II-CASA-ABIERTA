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
	http.HandleFunc("/api/informacion-salud", informacionSaludHandler)
	http.HandleFunc("/api/informacion-salud/", informacionSaludHandler)
	http.HandleFunc("/api/recordatorios-medicamentos", recordatoriosMedicamentosHandler)
	http.HandleFunc("/api/recordatorios-medicamentos/", recordatoriosMedicamentosHandler)
	http.HandleFunc("/api/medicamentos", proxyHandler("MEDICAMENTOS_URL"))
	http.HandleFunc("/api/appointments", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/appointments/", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/patients", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/doctors", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/cita-medica", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/cita-medica/", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/reportes", proxyHandler("REPORTES_URL"))
	http.HandleFunc("/api/reportes/", proxyHandler("REPORTES_URL"))
	http.HandleFunc("/api/reportes-medicos", proxyHandler("REPORTES_MEDICOS_URL"))
	http.HandleFunc("/api/reportes-medicos/", proxyHandler("REPORTES_MEDICOS_URL"))
	http.HandleFunc("/api/estado-animo", proxyHandler("ESTADO_ANIMO_URL"))
	http.HandleFunc("/api/estado-animo/", proxyHandler("ESTADO_ANIMO_URL"))
	http.HandleFunc("/api/informacion-salud", proxyHandler("INFORMACION_SALUD_URL"))
	http.HandleFunc("/api/informacion-salud/", proxyHandler("INFORMACION_SALUD_URL"))
	http.HandleFunc("/api/contacto-emergencia", proxyHandler("CONTACTO_EMERGENCIA_URL"))
	http.HandleFunc("/api/contacto-emergencia/", proxyHandler("CONTACTO_EMERGENCIA_URL"))

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

// Proxy al microservicio de recordatorio de medicamentos.
// Reenvía método, body, path y query sin agregar lógica propia.
func recordatoriosMedicamentosHandler(w http.ResponseWriter, r *http.Request) {
	recordatorioURL := os.Getenv("RECORDATORIOS_MEDICAMENTOS_URL")
	if recordatorioURL == "" {
		http.Error(w, "RECORDATORIOS_MEDICAMENTOS_URL no configurada", http.StatusInternalServerError)
		return
	}

	destino := recordatorioURL + r.URL.Path
	if r.URL.RawQuery != "" {
		destino += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequest(r.Method, destino, r.Body)

	if err != nil {
		http.Error(w, "Error al crear la petición", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error al contactar el servicio "+"de recordatorio de medicamentos", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error al leer la respuesta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
