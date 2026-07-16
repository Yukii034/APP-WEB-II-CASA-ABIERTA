package models

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
