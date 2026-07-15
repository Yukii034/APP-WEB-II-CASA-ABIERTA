package storage

import (
	"cuidabien/reportes-medicos/models"
	"testing"
)

func TestCrearResumenConDatos(t *testing.T) {
	store := NewStore()
	resumen := store.CrearResumen()

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
	resumen := CrearResumen(nil)

	if resumen.PacientesEvaluados != 0 {
		t.Fatalf("se esperaban 0 pacientes, se obtuvo %d", resumen.PacientesEvaluados)
	}
	if resumen.EstadoGeneral != "sin datos" {
		t.Fatalf("se esperaba estado sin datos, se obtuvo %s", resumen.EstadoGeneral)
	}
}

func TestBuscarPorPaciente(t *testing.T) {
	store := NewStore()
	reporte := store.BuscarPorPaciente("P001")

	if reporte == nil {
		t.Fatal("se esperaba encontrar reporte para P001")
	}
	if reporte.Nombre != "Maria Garcia" {
		t.Fatalf("se esperaba Maria Garcia, se obtuvo %s", reporte.Nombre)
	}
}

func TestListarReportesNil(t *testing.T) {
	store := &Store{Reportes: nil}
	reportes := store.ListarReportes()

	if reportes == nil {
		t.Fatal("se esperaba slice vacio, no nil")
	}
	if len(reportes) != 0 {
		t.Fatalf("se esperaban 0 reportes, se obtuvo %d", len(reportes))
	}
}

func TestCrearResumenManual(t *testing.T) {
	resumen := CrearResumen([]models.ReportePaciente{{PacienteID: "P001", AdherenciaMedicinas: 100, ComidasRegistradas: 3}})

	if resumen.EstadoGeneral != "estable" {
		t.Fatalf("se esperaba estable, se obtuvo %s", resumen.EstadoGeneral)
	}
}
