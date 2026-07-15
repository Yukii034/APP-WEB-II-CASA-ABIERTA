package models

// Contacto representa un familiar o cuidador que puede ser notificado
// en caso de emergencia de un paciente/adulto mayor.
type Contacto struct {
	ID         string `json:"id"`
	PacienteID string `json:"paciente_id"`
	Nombre     string `json:"nombre"`
	Telefono   string `json:"telefono"`
	Parentesco string `json:"parentesco"`
	Prioridad  int    `json:"prioridad"` // 1 = primero a notificar, 2 = segundo, etc.
	Principal  bool   `json:"principal"` // true si es el contacto principal
}

// Alerta representa una emergencia activada para un paciente.
// Al crearse, se notifica (simulado) a los contactos segun su prioridad.
type Alerta struct {
	ID                   string   `json:"id"`
	PacienteID           string   `json:"paciente_id"`
	Mensaje              string   `json:"mensaje"`
	Nivel                string   `json:"nivel"`  // leve, moderado, critico
	Estado               string   `json:"estado"` // activa, atendida, cancelada
	Timestamp            string   `json:"timestamp"`
	ContactosNotificados []string `json:"contactos_notificados"`
}

// EventoHistorial guarda un registro de cambios sobre una alerta.
type EventoHistorial struct {
	AlertaID  string `json:"alerta_id"`
	Accion    string `json:"accion"`
	Timestamp string `json:"timestamp"`
	Notas     string `json:"notas,omitempty"`
}

// Metricas expone un resumen rapido del estado del servicio.
type Metricas struct {
	TotalContactos      int `json:"total_contactos"`
	AlertasActivasHoy   int `json:"alertas_activas_hoy"`
	AlertasAtendidasHoy int `json:"alertas_atendidas_hoy"`
	TotalAlertas        int `json:"total_alertas"`
}

// --- Requests de entrada ---

type CrearContactoRequest struct {
	PacienteID string `json:"paciente_id"`
	Nombre     string `json:"nombre"`
	Telefono   string `json:"telefono"`
	Parentesco string `json:"parentesco"`
	Prioridad  int    `json:"prioridad"`
	Principal  bool   `json:"principal"`
}

type ActualizarContactoRequest struct {
	Nombre     string `json:"nombre,omitempty"`
	Telefono   string `json:"telefono,omitempty"`
	Parentesco string `json:"parentesco,omitempty"`
	Prioridad  int    `json:"prioridad,omitempty"`
	Principal  *bool  `json:"principal,omitempty"`
}

type CrearAlertaRequest struct {
	PacienteID string `json:"paciente_id"`
	Mensaje    string `json:"mensaje"`
	Nivel      string `json:"nivel"`
}

type AtenderAlertaRequest struct {
	Notas string `json:"notas,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
