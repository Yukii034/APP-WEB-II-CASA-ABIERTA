package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type ReportePaciente struct {
	PacienteID          string   `json:"paciente_id"`
	Nombre              string   `json:"nombre"`
	Periodo             string   `json:"periodo"`
	CitasProgramadas    int      `json:"citas_programadas"`
	CitasCompletadas    int      `json:"citas_completadas"`
	ComidasRegistradas  int      `json:"comidas_registradas"`
	AlertasSalud        int      `json:"alertas_salud"`
	AdherenciaMedicinas int      `json:"adherencia_medicinas"`
	EstadoGeneral       string   `json:"estado_general"`
	Recomendaciones     []string `json:"recomendaciones"`
}

type ResumenGeneral struct {
	PacientesEvaluados int               `json:"pacientes_evaluados"`
	EstadoGeneral      string            `json:"estado_general"`
	AlertasTotales     int               `json:"alertas_totales"`
	Promedios          map[string]int    `json:"promedios"`
	Pacientes          []ReportePaciente `json:"pacientes"`
}

var reportes = []ReportePaciente{
	{
		PacienteID:          "P001",
		Nombre:              "Maria Garcia",
		Periodo:             "semana actual",
		CitasProgramadas:    3,
		CitasCompletadas:    2,
		ComidasRegistradas:  18,
		AlertasSalud:        1,
		AdherenciaMedicinas: 92,
		EstadoGeneral:       "estable",
		Recomendaciones: []string{
			"Mantener controles medicos programados",
			"Continuar con horarios de alimentacion",
		},
	},
	{
		PacienteID:          "P002",
		Nombre:              "Juan Lopez",
		Periodo:             "semana actual",
		CitasProgramadas:    2,
		CitasCompletadas:    1,
		ComidasRegistradas:  12,
		AlertasSalud:        3,
		AdherenciaMedicinas: 76,
		EstadoGeneral:       "requiere seguimiento",
		Recomendaciones: []string{
			"Revisar alertas de salud con el cuidador",
			"Mejorar cumplimiento de medicacion",
		},
	},
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/reportes/resumen", resumenHandler)
	http.HandleFunc("/api/reportes/semanal", semanalHandler)
	http.HandleFunc("/api/reportes/paciente/", pacienteHandler)
	http.HandleFunc("/api/reportes", semanalHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servicio de reportes medicos corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func resumenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	writeJSON(w, http.StatusOK, crearResumen(reportes))
}

func semanalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	writeJSON(w, http.StatusOK, reportes)
}

func pacienteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/reportes/paciente/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de paciente requerido")
		return
	}

	for _, reporte := range reportes {
		if reporte.PacienteID == id {
			writeJSON(w, http.StatusOK, reporte)
			return
		}
	}

	writeError(w, http.StatusNotFound, "Reporte no encontrado")
}

func crearResumen(data []ReportePaciente) ResumenGeneral {
	resumen := ResumenGeneral{
		PacientesEvaluados: len(data),
		EstadoGeneral:      "estable",
		Promedios:          map[string]int{},
		Pacientes:          data,
	}

	if len(data) == 0 {
		resumen.EstadoGeneral = "sin datos"
		return resumen
	}

	var totalAlertas, totalAdherencia, totalComidas int
	for _, reporte := range data {
		totalAlertas += reporte.AlertasSalud
		totalAdherencia += reporte.AdherenciaMedicinas
		totalComidas += reporte.ComidasRegistradas
	}

	resumen.AlertasTotales = totalAlertas
	resumen.Promedios["adherencia_medicinas"] = totalAdherencia / len(data)
	resumen.Promedios["comidas_registradas"] = totalComidas / len(data)

	if totalAlertas >= 3 || resumen.Promedios["adherencia_medicinas"] < 80 {
		resumen.EstadoGeneral = "requiere seguimiento"
	}

	return resumen
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
