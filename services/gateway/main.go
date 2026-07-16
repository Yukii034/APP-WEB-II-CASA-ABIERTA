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
	http.HandleFunc("/api/vitales", proxyHandler("MONITOREO_SIGNOS_VITALES_URL"))
	http.HandleFunc("/api/vitales/", proxyHandler("MONITOREO_SIGNOS_VITALES_URL"))
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

// Proxy al microservicio de información de salud.
// Reenvia la peticion completa (metodo, body, path, query) al servicio.
func informacionSaludHandler(w http.ResponseWriter, r *http.Request) {
	informacionSaludURL := os.Getenv("INFORMACION_SALUD_URL")
	if informacionSaludURL == "" {
		http.Error(w, "INFORMACION_SALUD_URL no configurada", http.StatusInternalServerError)
		return
	}

	destino := informacionSaludURL + r.URL.Path
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
		http.Error(w, "Error al contactar el servicio de información de salud", http.StatusBadGateway)
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
