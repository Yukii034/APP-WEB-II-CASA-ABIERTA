package model

import "time"

// RecordatorioMedicamento representa un medicamento programado
// para un adulto mayor.
type RecordatorioMedicamento struct {
	ID             string    `json:"id"`
	AdultoMayorID  string    `json:"adulto_mayor_id"`
	NombrePaciente string    `json:"nombre_paciente"`
	Medicamento    string    `json:"medicamento"`
	Dosis          string    `json:"dosis"`
	Hora           string    `json:"hora"`
	Frecuencia     string    `json:"frecuencia"`
	Activo         bool      `json:"activo"`
	CreadoEn       time.Time `json:"creado_en"`
	ActualizadoEn  time.Time `json:"actualizado_en"`
}

// EntradaRecordatorioMedicamento es el body esperado
// para crear o actualizar un recordatorio.
type EntradaRecordatorioMedicamento struct {
	AdultoMayorID  string `json:"adulto_mayor_id"`
	NombrePaciente string `json:"nombre_paciente"`
	Medicamento    string `json:"medicamento"`
	Dosis          string `json:"dosis"`
	Hora           string `json:"hora"`
	Frecuencia     string `json:"frecuencia"`
	Activo         *bool  `json:"activo,omitempty"`
}

// EntradaEstadoRecordatorio es el body esperado
// para activar o desactivar un recordatorio.
type EntradaEstadoRecordatorio struct {
	Activo *bool `json:"activo"`
}

// EntradaVerificacion es el body esperado
// para comprobar los medicamentos de una hora.
type EntradaVerificacion struct {
	Hora string `json:"hora"`
}

// AlertaMedicamento representa una notificación simulada.
type AlertaMedicamento struct {
	RecordatorioID string `json:"recordatorio_id"`
	AdultoMayorID  string `json:"adulto_mayor_id"`
	NombrePaciente string `json:"nombre_paciente"`
	Medicamento    string `json:"medicamento"`
	Dosis          string `json:"dosis"`
	Hora           string `json:"hora"`
	Mensaje        string `json:"mensaje"`
}

// ResultadoVerificacion contiene las alertas
// encontradas para una hora determinada.
type ResultadoVerificacion struct {
	Hora     string              `json:"hora"`
	Cantidad int                 `json:"cantidad_alertas"`
	Alertas  []AlertaMedicamento `json:"alertas"`
}
