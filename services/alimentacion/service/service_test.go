package service

import (
	"testing"
	"time"

	"cuidabien/alimentacion/modelo"
)

func esperadas() []modelo.ComidaEsperada {
	return []modelo.ComidaEsperada{
		{Tipo: "desayuno", HoraLimite: "10:00"},
		{Tipo: "almuerzo", HoraLimite: "15:00"},
		{Tipo: "cena", HoraLimite: "21:00"},
	}
}

func TestResumenSinRegistrosAntesDeCualquierLimite(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 8, 0, 0, 0, time.Local)
	r := calcularResumen([]modelo.RegistroComida{}, esperadas(), ahora)

	if r.HaySaltadas {
		t.Errorf("a las 8am ninguna comida debería estar saltada todavía")
	}
	if r.ComidasHechas != 0 {
		t.Errorf("esperaba 0 comidas hechas, obtuve %d", r.ComidasHechas)
	}
}

func TestResumenDesayunoSaltado(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 11, 0, 0, 0, time.Local)
	r := calcularResumen([]modelo.RegistroComida{}, esperadas(), ahora)

	if !r.HaySaltadas {
		t.Errorf("a las 11am sin desayuno registrado debería marcar saltada")
	}
	if !r.Comidas[0].Saltada {
		t.Errorf("el desayuno debería estar marcado como saltado")
	}
}

func TestResumenComidaRegistradaNoCuentaComoSaltada(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 11, 0, 0, 0, time.Local)
	regs := []modelo.RegistroComida{
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

func TestNivelAlertaOkSinSaltadas(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 8, 0, 0, 0, time.Local)
	r := calcularResumen([]modelo.RegistroComida{}, esperadas(), ahora)

	if r.NivelAlerta != "ok" {
		t.Errorf("esperaba nivel_alerta 'ok', obtuve %q", r.NivelAlerta)
	}
}

func TestNivelAlertaAtencionConUnaSaltada(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 11, 0, 0, 0, time.Local)
	r := calcularResumen([]modelo.RegistroComida{}, esperadas(), ahora)

	if r.NivelAlerta != "atencion" {
		t.Errorf("esperaba nivel_alerta 'atencion' con 1 comida saltada, obtuve %q", r.NivelAlerta)
	}
}

func TestNivelAlertaUrgenteConDosSaltadas(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 16, 0, 0, 0, time.Local)
	r := calcularResumen([]modelo.RegistroComida{}, esperadas(), ahora)

	if r.NivelAlerta != "urgente" {
		t.Errorf("esperaba nivel_alerta 'urgente' con 2 comidas saltadas, obtuve %q", r.NivelAlerta)
	}
}
