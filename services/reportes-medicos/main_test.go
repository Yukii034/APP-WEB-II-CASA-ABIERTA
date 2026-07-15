package main

import "testing"

func TestCrearResumenConDatos(t *testing.T) {
	resumen := crearResumen(reportes)

	if resumen.PacientesEvaluados != 2 {
		t.Fatalf("se esperaban 2 pacientes, se obtuvo %d", resumen.PacientesEvaluados)
	}
	if resumen.AlertasTotales != 4 {
		t.Fatalf("se esperaban 4 alertas, se obtuvo %d", resumen.AlertasTotales)
	}
	if resumen.EstadoGeneral != "requiere seguimiento" {
		t.Fatalf("se esperaba estado requiere seguimiento, se obtuvo %s", resumen.EstadoGeneral)
	}
}

func TestCrearResumenSinDatos(t *testing.T) {
	resumen := crearResumen(nil)

	if resumen.PacientesEvaluados != 0 {
		t.Fatalf("se esperaban 0 pacientes, se obtuvo %d", resumen.PacientesEvaluados)
	}
	if resumen.EstadoGeneral != "sin datos" {
		t.Fatalf("se esperaba estado sin datos, se obtuvo %s", resumen.EstadoGeneral)
	}
}
