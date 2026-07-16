package models

type Paciente struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}

type Medicamento struct {
	ID          string   `json:"id"`
	PacienteID  string   `json:"paciente_id"`
	Nombre      string   `json:"nombre"`
	Dosis       string   `json:"dosis"`
	Frecuencia  string   `json:"frecuencia"`
	Horarios    []string `json:"horarios"`
	FechaInicio string   `json:"fecha_inicio"`
	FechaFin    string   `json:"fecha_fin,omitempty"`
	Estado      string   `json:"estado"`
	Notas       string   `json:"notas,omitempty"`
}

type Toma struct {
	ID                    string `json:"id"`
	MedicamentoID         string `json:"medicamento_id"`
	PacienteID            string `json:"paciente_id"`
	FechaHoraProgramada   string `json:"fecha_hora_programada"`
	Estado                string `json:"estado"`
	FechaHoraReal         string `json:"fecha_hora_real,omitempty"`
	Notas                 string `json:"notas,omitempty"`
}

type Alerta struct {
	ID            string `json:"id"`
	PacienteID    string `json:"paciente_id"`
	MedicamentoID string `json:"medicamento_id"`
	Tipo          string `json:"tipo"`
	Mensaje       string `json:"mensaje"`
	FechaCreacion string `json:"fecha_creacion"`
	Leida         bool   `json:"leida"`
}

type Interaccion struct {
	ID           string `json:"id"`
	MedicamentoA string `json:"medicamento_a"`
	MedicamentoB string `json:"medicamento_b"`
	Gravedad     string `json:"gravedad"`
	Descripcion  string `json:"descripcion"`
}

type Adherencia struct {
	PacienteID      string  `json:"paciente_id"`
	TotalTomas      int     `json:"total_tomas"`
	TomasCumplidas  int     `json:"tomas_cumplidas"`
	TomasNoCumplidas int    `json:"tomas_no_cumplidas"`
	Porcentaje      float64 `json:"porcentaje"`
}

type CrearMedicamentoRequest struct {
	PacienteID string   `json:"paciente_id"`
	Nombre     string   `json:"nombre"`
	Dosis      string   `json:"dosis"`
	Frecuencia string   `json:"frecuencia"`
	Horarios   []string `json:"horarios"`
	FechaInicio string  `json:"fecha_inicio"`
	FechaFin   string   `json:"fecha_fin,omitempty"`
	Notas      string   `json:"notas,omitempty"`
}

type ActualizarMedicamentoRequest struct {
	Nombre     string   `json:"nombre,omitempty"`
	Dosis      string   `json:"dosis,omitempty"`
	Frecuencia string   `json:"frecuencia,omitempty"`
	Horarios   []string `json:"horarios,omitempty"`
	FechaFin   string   `json:"fecha_fin,omitempty"`
	Notas      string   `json:"notas,omitempty"`
}

type RegistrarTomaRequest struct {
	Estado string `json:"estado"`
	Notas  string `json:"notas,omitempty"`
}

type VerificarInteraccionesRequest struct {
	PacienteID   string `json:"paciente_id"`
	MedicamentoA string `json:"medicamento_a"`
	MedicamentoB string `json:"medicamento_b"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
