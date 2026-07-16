// Package models contains the domain structures exposed by the service.
package models

import "time"

// SignosVitales is one vital-sign measurement for an older adult.
type SignosVitales struct {
	ID                 string             `json:"id"`
	IDAdultoMayor      string             `json:"id_adulto_mayor"`
	RegistradoPor      string             `json:"registrado_por"`
	FechaRegistro      time.Time          `json:"fecha_registro"`
	PresionSistolica   int                `json:"presion_sistolica"`
	PresionDiastolica  int                `json:"presion_diastolica"`
	FrecuenciaCardiaca int                `json:"frecuencia_cardiaca"`
	Temperatura        *float64           `json:"temperatura,omitempty"`
	SaturacionOxigeno  *int               `json:"saturacion_oxigeno,omitempty"`
	NivelGlucosa       *float64           `json:"nivel_glucosa,omitempty"`
	Peso               *float64           `json:"peso,omitempty"`
	Altura             *float64           `json:"altura,omitempty"`
	NivelDolor         *int               `json:"nivel_dolor,omitempty"`
	Observaciones      string             `json:"observaciones,omitempty"`
	Evaluacion         EvaluacionRegistro `json:"evaluacion"`
}

// EntradaSignosVitales is the JSON body accepted by POST requests.
type EntradaSignosVitales struct {
	IDAdultoMayor      string   `json:"id_adulto_mayor"`
	RegistradoPor      string   `json:"registrado_por"`
	PresionSistolica   int      `json:"presion_sistolica"`
	PresionDiastolica  int      `json:"presion_diastolica"`
	FrecuenciaCardiaca int      `json:"frecuencia_cardiaca"`
	Temperatura        *float64 `json:"temperatura,omitempty"`
	SaturacionOxigeno  *int     `json:"saturacion_oxigeno,omitempty"`
	NivelGlucosa       *float64 `json:"nivel_glucosa,omitempty"`
	Peso               *float64 `json:"peso,omitempty"`
	Altura             *float64 `json:"altura,omitempty"`
	NivelDolor         *int     `json:"nivel_dolor,omitempty"`
	Observaciones      string   `json:"observaciones,omitempty"`
}

// Estado represents the clinical severity for a value or complete record.
type Estado string

const (
	EstadoNormal      Estado = "normal"
	EstadoBajo        Estado = "bajo"
	EstadoAlto        Estado = "alto"
	EstadoCritico     Estado = "critico"
	EstadoAdvertencia Estado = "advertencia"
)

// EvaluacionValor describes the clinical result for one measured parameter.
type EvaluacionValor struct {
	Parametro string `json:"parametro"`
	Estado    Estado `json:"estado"`
}

// EvaluacionRegistro groups the per-value results and the overall status.
type EvaluacionRegistro struct {
	EstadoGeneral Estado            `json:"estado_general"`
	Valores       []EvaluacionValor `json:"valores"`
}

// PuntoTendencia is a compact result for charts.
type PuntoTendencia struct {
	Fecha  time.Time `json:"fecha"`
	Valor  float64   `json:"valor"`
	Estado Estado    `json:"estado"`
}
