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

	// Monitoreo de signos vitales
	http.HandleFunc("/api/vitales", proxyHandler("MONITOREO_SIGNOS_VITALES_URL"))
	http.HandleFunc("/api/vitales/", proxyHandler("MONITOREO_SIGNOS_VITALES_URL"))

	// Citas medicas
	http.HandleFunc("/api/cita-medica", proxyHandler("CITAS_URL"))
	http.HandleFunc("/api/cita-medica/", proxyHandler("CITAS_URL"))

	// Reportes
	http.HandleFunc("/api/reportes", proxyHandler("REPORTES_URL"))
	http.HandleFunc("/api/reportes/", proxyHandler("REPORTES_URL"))

	// Reportes medicos
	http.HandleFunc("/api/reportes-medicos", proxyHandler("REPORTES_MEDICOS_URL"))
	http.HandleFunc("/api/reportes-medicos/", proxyHandler("REPORTES_MEDICOS_URL"))

	// Estado de animo
	http.HandleFunc("/api/estado-animo", proxyHandler("ESTADO_ANIMO_URL"))
	http.HandleFunc("/api/estado-animo/", proxyHandler("ESTADO_ANIMO_URL"))

	// Informacion de salud
	http.HandleFunc("/api/informacion-salud", proxyHandler("INFORMACION_SALUD_URL"))
	http.HandleFunc("/api/informacion-salud/", proxyHandler("INFORMACION_SALUD_URL"))

	// Contacto de emergencia
	http.HandleFunc("/api/contacto-emergencia", proxyHandler("CONTACTO_EMERGENCIA_URL"))
	http.HandleFunc("/api/contacto-emergencia/", proxyHandler("CONTACTO_EMERGENCIA_URL"))

	// Estimulacion cognitiva
	http.HandleFunc("/api/ejercicios", proxyHandler("ESTIMULACION_COGNITIVA_URL"))
	http.HandleFunc("/api/ejercicios/", proxyHandler("ESTIMULACION_COGNITIVA_URL"))

	// Alimentacion
	http.HandleFunc("/api/alimentacion", proxyHandler("ALIMENTACION_URL"))
	http.HandleFunc("/api/alimentacion/", proxyHandler("ALIMENTACION_URL"))

	// Recordatorios de medicamentos
	http.HandleFunc("/api/recordatorios-medicamentos", proxyHandler("RECORDATORIOS_MEDICAMENTOS_URL"))
	http.HandleFunc("/api/recordatorios-medicamentos/", proxyHandler("RECORDATORIOS_MEDICAMENTOS_URL"))

	// Actividad fisica
	http.HandleFunc("/api/actividad-fisica", proxyHandler("ACTIVIDAD_FISICA_URL"))
	http.HandleFunc("/api/actividad-fisica/", proxyHandler("ACTIVIDAD_FISICA_URL"))

	// Cuidadores
	http.HandleFunc("/api/cuidadores", proxyHandler("CUIDADORES_URL"))
	http.HandleFunc("/api/cuidadores/", proxyHandler("CUIDADORES_URL"))

	// Rutas de contacto-emergencia
	http.HandleFunc("/api/contacts", proxyHandler("CONTACTO_EMERGENCIA_URL"))
	http.HandleFunc("/api/contacts/", proxyHandler("CONTACTO_EMERGENCIA_URL"))
	http.HandleFunc("/api/alerts", proxyHandler("CONTACTO_EMERGENCIA_URL"))
	http.HandleFunc("/api/alerts/", proxyHandler("CONTACTO_EMERGENCIA_URL"))
	http.HandleFunc("/api/metrics", proxyHandler("CONTACTO_EMERGENCIA_URL"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

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
