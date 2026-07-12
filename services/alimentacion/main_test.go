package main

import (
	"testing"
	"time"
)

func esperadas() []comidaEsperada {
	return []comidaEsperada{
		{tipo: "desayuno", horaLimite: "10:00"},
		{tipo: "almuerzo", horaLimite: "15:00"},
		{tipo: "cena", horaLimite: "21:00"},
	}
}

func TestResumenSinRegistrosAntesDeCualquierLimite(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 8, 0, 0, 0, time.Local)
	r := calcularResumen([]RegistroComida{}, esperadas(), ahora)

	if r.HaySaltadas {
		t.Errorf("a las 8am ninguna comida debería estar saltada todavía")
	}
	if r.ComidasHechas != 0 {
		t.Errorf("esperaba 0 comidas hechas, obtuve %d", r.ComidasHechas)
	}
}

func TestResumenDesayunoSaltado(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 11, 0, 0, 0, time.Local)
	r := calcularResumen([]RegistroComida{}, esperadas(), ahora)

	if !r.HaySaltadas {
		t.Errorf("a las 11am sin desayuno registrado debería marcar saltada")
	}
	if !r.Comidas[0].Saltada {
		t.Errorf("el desayuno debería estar marcado como saltado")
	}
}

func TestResumenComidaRegistradaNoCuentaComoSaltada(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 11, 0, 0, 0, time.Local)
	regs := []RegistroComida{
		{ID: "1", TipoComida: "desayuno", Hora: ahora.Add(-2 * time.Hour)},
	}
	r := calcularResumen(regs, esperadas(), ahora)

	if r.Comidas[0].Saltada {
		t.Errorf("el desayuno ya registrado no debería estar saltado")
	}
	if r.ComidasHechas != 1 {
		t.Errorf("esperaba 1 comida hecha, obtuve %d", r.ComidasHechas)
	}
}
