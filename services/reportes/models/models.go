package models

type Paciente struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}

type ReporteSemanal struct {
	PacienteID          string             `json:"paciente_id"`
	PacienteNombre      string             `json:"paciente_nombre"`
	FechaInicio         string             `json:"fecha_inicio"`
	FechaFin            string             `json:"fecha_fin"`
	ResumenCitas        ResumenCitas       `json:"resumen_citas"`
	ResumenMedicamentos ResumenMedicamentos `json:"resumen_medicamentos"`
	ResumenAlimentacion ResumenAlimentacion `json:"resumen_alimentacion"`
	ResumenSalud        ResumenSalud       `json:"resumen_salud"`
	EstadoGeneral       string             `json:"estado_general"`
	Recomendacion       string             `json:"recomendacion"`
}

type ResumenCitas struct {
	TotalProgramadas  int    `json:"total_programadas"`
	Completadas       int    `json:"completadas"`
	Canceladas        int    `json:"canceladas"`
	Pendientes        int    `json:"pendientes"`
	ProximaCita       string `json:"proxima_cita,omitempty"`
	ProximoDoctor     string `json:"proximo_doctor,omitempty"`
}

type ResumenMedicamentos struct {
	TotalActivos      int     `json:"total_activos"`
	TomasRegistradas  int     `json:"tomas_registradas"`
	TomasCumplidas    int     `json:"tomas_cumplidas"`
	TomasNoCumplidas  int     `json:"tomas_no_cumplidas"`
	PorcentajeAdherencia float64 `json:"porcentaje_adherencia"`
	AlertasActivas    int     `json:"alertas_activas"`
}

type ResumenAlimentacion struct {
	ComidasRegistradas int    `json:"comidas_registradas"`
	ComidasEsperadas   int    `json:"comidas_esperadas"`
	PorcentajeCumplido float64 `json:"porcentaje_cumplido"`
	ComidasSaltadas    int    `json:"comidas_saltadas"`
	UltimaComida       string `json:"ultima_comida,omitempty"`
}

type ResumenSalud struct {
	AlertasSalud       int    `json:"alertas_salud"`
	SignosVitalesOK    bool   `json:"signos_vitales_ok"`
	UltimoControl      string `json:"ultimo_control,omitempty"`
}

type ReportePaciente struct {
	PacienteID          string             `json:"paciente_id"`
	PacienteNombre      string             `json:"paciente_nombre"`
	TotalCitas          int                `json:"total_citas"`
	CitasCompletadas    int                `json:"citas_completadas"`
	TotalMedicamentos   int                `json:"total_medicamentos"`
	AdherenciaMedicacion float64           `json:"adherencia_medicacion"`
	ComidasRegistradas  int                `json:"comidas_registradas"`
	AlertasActivas      int                `json:"alertas_activas"`
	EstadoGeneral       string             `json:"estado_general"`
	HistorialCitas      []ResumenCita      `json:"historial_citas"`
	HistorialMedicamentos []ResumenMedicamento `json:"historial_medicamentos"`
}

type ResumenCita struct {
	ID        string `json:"id"`
	Fecha     string `json:"fecha"`
	Hora      string `json:"hora"`
	Doctor    string `json:"doctor"`
	Estado    string `json:"estado"`
	Prioridad string `json:"prioridad"`
}

type ResumenMedicamento struct {
	Nombre            string  `json:"nombre"`
	Dosis             string  `json:"dosis"`
	Estado            string  `json:"estado"`
	Adherencia        float64 `json:"adherencia"`
}

type ResumenGeneral struct {
	TotalPacientes        int     `json:"total_pacientes"`
	TotalCitasHoy         int     `json:"total_citas_hoy"`
	TotalMedicamentosActivos int `json:"total_medicamentos_activos"`
	TotalAlertasPendientes int   `json:"total_alertas_pendientes"`
	PromedioAdherencia    float64 `json:"promedio_adherencia"`
	PacientesConAlertas   int     `json:"pacientes_con_alertas"`
}

type DashboardData struct {
	ResumenGeneral ResumenGeneral  `json:"resumen_general"`
	Pacientes      []PacienteResumen `json:"pacientes"`
	AlertasRecientes []AlertaResumen `json:"alertas_recientes"`
}

type PacienteResumen struct {
	ID                string  `json:"id"`
	Nombre            string  `json:"nombre"`
	CitasProximas     int     `json:"citas_proximas"`
	MedicamentosActivos int   `json:"medicamentos_activos"`
	Adherencia        float64 `json:"adherencia"`
	Estado            string  `json:"estado"`
}

type AlertaResumen struct {
	Tipo      string `json:"tipo"`
	Mensaje   string `json:"mensaje"`
	Fuente    string `json:"fuente"`
	PacienteID string `json:"paciente_id"`
	Fecha     string `json:"fecha"`
}
