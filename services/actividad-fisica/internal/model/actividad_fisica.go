package model

import "time"

// ActividadFisica representa una actividad registrada para un adulto mayor.
type ActividadFisica struct {
	ID                string    `json:"id"`
	NombrePaciente    string    `json:"nombre_paciente"`
	TipoActividad     string    `json:"tipo_actividad"`
	DuracionMinutos   int       `json:"duracion_minutos"`
	Intensidad        string    `json:"intensidad"`
	Fecha             string    `json:"fecha"`
	Estado            string    `json:"estado"`
	Observaciones     string    `json:"observaciones,omitempty"`
	CaloriasEstimadas int       `json:"calorias_estimadas"`
	CreadoEn          time.Time `json:"creado_en"`
	ActualizadoEn     time.Time `json:"actualizado_en"`
}

// EntradaActividadFisica es el body esperado para crear y actualizar.
type EntradaActividadFisica struct {
	NombrePaciente  string `json:"nombre_paciente"`
	TipoActividad   string `json:"tipo_actividad"`
	DuracionMinutos int    `json:"duracion_minutos"`
	Intensidad      string `json:"intensidad"`
	Fecha           string `json:"fecha"`
	Estado          string `json:"estado"`
	Observaciones   string `json:"observaciones"`
}
