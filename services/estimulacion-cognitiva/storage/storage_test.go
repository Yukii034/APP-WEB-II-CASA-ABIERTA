package storage

import (
	"testing"
	"time"
)

func TestAgregarTipoRequerido(t *testing.T) {
	s := New()

	if _, err := s.Agregar(""); err == nil {
		t.Fatal("se esperaba error al no enviar tipo")
	}
	if _, err := s.Agregar("   "); err == nil {
		t.Fatal("se esperaba error al enviar solo espacios en blanco")
	}
}

func TestAgregarValido(t *testing.T) {
	s := New()

	e, err := s.Agregar("memoria")
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if e.ID == "" {
		t.Fatal("se esperaba un ID asignado")
	}
	if e.Tipo != "memoria" {
		t.Fatalf("tipo esperado 'memoria', obtenido '%s'", e.Tipo)
	}
}

func TestResumenSinEjerciciosAlertaActiva(t *testing.T) {
	s := New()

	resumen := s.Resumen()
	if resumen.Total != 0 {
		t.Fatalf("total esperado 0, obtenido %d", resumen.Total)
	}
	if !resumen.HayAlerta {
		t.Fatal("se esperaba alerta activa cuando no hay ejercicios registrados")
	}
}

func TestResumenConEjercicioHoyNoHayAlerta(t *testing.T) {
	s := New()
	s.Agregar("trivia")

	resumen := s.Resumen()
	if resumen.EjerciciosHoy != 1 {
		t.Fatalf("se esperaba 1 ejercicio hoy, se obtuvo %d", resumen.EjerciciosHoy)
	}
	if resumen.DiasDesdeUltimo != 0 {
		t.Fatalf("dias_desde_ultimo esperado 0, obtenido %d", resumen.DiasDesdeUltimo)
	}
	if resumen.HayAlerta {
		t.Fatal("no se esperaba alerta si hay un ejercicio hecho hoy")
	}
}

func TestResumenConInactividadActivaAlerta(t *testing.T) {
	s := New()
	s.agregarConFecha("memoria", time.Now().AddDate(0, 0, -umbralDiasAlerta))

	resumen := s.Resumen()
	if resumen.DiasDesdeUltimo != umbralDiasAlerta {
		t.Fatalf("dias_desde_ultimo esperado %d, obtenido %d", umbralDiasAlerta, resumen.DiasDesdeUltimo)
	}
	if !resumen.HayAlerta {
		t.Fatal("se esperaba alerta activa tras varios días sin actividad")
	}
}

func TestResumenTomaElEjercicioMasReciente(t *testing.T) {
	s := New()
	s.agregarConFecha("trivia", time.Now().AddDate(0, 0, -5))
	s.agregarConFecha("memoria", time.Now())
	s.agregarConFecha("sopa_letras", time.Now().AddDate(0, 0, -3))

	resumen := s.Resumen()
	if resumen.UltimoEjercicio == nil {
		t.Fatal("se esperaba un último ejercicio calculado")
	}
	if resumen.UltimoEjercicio.Tipo != "memoria" {
		t.Fatalf("se esperaba que el último ejercicio fuera 'memoria', fue '%s'", resumen.UltimoEjercicio.Tipo)
	}
	if resumen.DiasDesdeUltimo != 0 {
		t.Fatalf("dias_desde_ultimo esperado 0, obtenido %d", resumen.DiasDesdeUltimo)
	}
}

func TestResetBorraEjercicios(t *testing.T) {
	s := New()
	s.Agregar("memoria")
	s.Reset()

	if len(s.Listar()) != 0 {
		t.Fatal("se esperaba que no quedaran ejercicios tras el reset")
	}
}

func TestSeedDejaAlertaActiva(t *testing.T) {
	s := New()
	s.Seed()

	resumen := s.Resumen()
	if resumen.Total != 2 {
		t.Fatalf("se esperaban 2 ejercicios de ejemplo, se obtuvieron %d", resumen.Total)
	}
	if !resumen.HayAlerta {
		t.Fatal("se esperaba que el seed dejara la alerta activa a propósito")
	}
	if resumen.DiasDesdeUltimo != umbralDiasAlerta {
		t.Fatalf("dias_desde_ultimo esperado %d tras el seed, obtenido %d", umbralDiasAlerta, resumen.DiasDesdeUltimo)
	}
}

func TestSeedLuegoRegistrarHoyResuelveAlerta(t *testing.T) {
	s := New()
	s.Seed()
	s.Agregar("rompecabezas")
	resumen := s.Resumen()
	if resumen.HayAlerta {
		t.Fatal("la alerta debería resolverse tras registrar un ejercicio hoy")
	}
	if resumen.EjerciciosHoy != 1 {
		t.Fatalf("se esperaba 1 ejercicio hoy, se obtuvo %d", resumen.EjerciciosHoy)
	}
}
