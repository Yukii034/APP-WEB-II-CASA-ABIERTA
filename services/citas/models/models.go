package models

type Paciente struct {
	ID                 string   `json:"id"`
	Nombre             string   `json:"nombre"`
	Telefono           string   `json:"telefono"`
	ContactoEmergencia string   `json:"contacto_emergencia"`
	Alergias           []string `json:"alergias"`
}

type Doctor struct {
	ID           string `json:"id"`
	Nombre       string `json:"nombre"`
	Especialidad string `json:"especialidad"`
}

type Cita struct {
	ID          string `json:"id"`
	PacienteID  string `json:"paciente_id"`
	DoctorID    string `json:"doctor_id"`
	Fecha       string `json:"fecha"`
	Hora        string `json:"hora"`
	Estado      string `json:"estado"`
	Prioridad   string `json:"prioridad"`
	Motivo      string `json:"motivo"`
	NotasMedico string `json:"notas_medico"`
}

type EventoHistorial struct {
	CitaID    string `json:"cita_id"`
	Accion    string `json:"accion"`
	EstadoAnt string `json:"estado_anterior,omitempty"`
	EstadoNue string `json:"estado_nuevo,omitempty"`
	Timestamp string `json:"timestamp"`
	Notas     string `json:"notas,omitempty"`
}

type Metricas struct {
	CitasCreadasHoy     int `json:"citas_creadas_hoy"`
	CitasCanceladasHoy  int `json:"citas_canceladas_hoy"`
	CitasCompletadasHoy int `json:"citas_completadas_hoy"`
	TotalCitasActivas   int `json:"total_citas_activas"`
}

type Reminder struct {
	CitaID     string `json:"cita_id"`
	PacienteID string `json:"paciente_id"`
	DoctorID   string `json:"doctor_id"`
	Fecha      string `json:"fecha"`
	Hora       string `json:"hora"`
	Mensaje    string `json:"mensaje"`
	Prioridad  string `json:"prioridad"`
}

type CrearCitaRequest struct {
	PacienteID string `json:"paciente_id"`
	DoctorID   string `json:"doctor_id"`
	Fecha      string `json:"fecha"`
	Hora       string `json:"hora"`
	Prioridad  string `json:"prioridad"`
	Motivo     string `json:"motivo"`
}

type ActualizarCitaRequest struct {
	Fecha     string `json:"fecha,omitempty"`
	Hora      string `json:"hora,omitempty"`
	Prioridad string `json:"prioridad,omitempty"`
	Motivo    string `json:"motivo,omitempty"`
}

type NotasRequest struct {
	NotasMedico string `json:"notas_medico"`
}

type RecurrenteRequest struct {
	PacienteID  string `json:"paciente_id"`
	DoctorID    string `json:"doctor_id"`
	FechaInicio string `json:"fecha_inicio"`
	Hora        string `json:"hora"`
	Prioridad   string `json:"prioridad"`
	Motivo      string `json:"motivo"`
	Intervalo   int    `json:"intervalo_dias"`
	Cantidad    int    `json:"cantidad"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
