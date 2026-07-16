package model

import "time"

// InformacionSalud representa la ficha de salud de un adulto mayor.
type InformacionSalud struct {
	// ID identifica de forma única la ficha dentro del servicio.
	ID string `json:"id"`
	// Los campos clínicos se representan como listas para admitir varios valores.
	NombrePaciente       string    `json:"nombre_paciente,omitempty"`
	Diagnosticos         []string  `json:"diagnosticos"`
	Alergias             []string  `json:"alergias"`
	EnfermedadesCronicas []string  `json:"enfermedades_cronicas"`
	AntecedentesMedicos  []string  `json:"antecedentes_medicos"`
	ActualizadoEn        time.Time `json:"actualizado_en"`
}

// EntradaInformacionSalud es el body esperado en POST y PUT. En un PUT,
// una lista nil indica que el campo no se debe modificar.
type EntradaInformacionSalud struct {
	PacienteID           string   `json:"paciente_id"`
	NombrePaciente       string   `json:"nombre_paciente"`
	Diagnosticos         []string `json:"diagnosticos"`
	Alergias             []string `json:"alergias"`
	EnfermedadesCronicas []string `json:"enfermedades_cronicas"`
	AntecedentesMedicos  []string `json:"antecedentes_medicos"`
}
