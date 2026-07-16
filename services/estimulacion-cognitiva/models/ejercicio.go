package models

import "time"

type Ejercicio struct {
	ID    string    `json:"id"`
	Tipo  string    `json:"tipo"`
	Fecha time.Time `json:"fecha"`
}

type Resumen struct {
	Ejercicios      []Ejercicio `json:"ejercicios"`
	Total           int         `json:"total"`
	EjerciciosHoy   int         `json:"ejercicios_hoy"`
	UltimoEjercicio *Ejercicio  `json:"ultimo_ejercicio,omitempty"`
	DiasDesdeUltimo int         `json:"dias_desde_ultimo"`
	HayAlerta       bool        `json:"hay_alerta"`
	Mensaje         string      `json:"mensaje"`
}
