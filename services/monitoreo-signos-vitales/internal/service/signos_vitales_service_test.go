package service

import (
	"testing"

	"monitoreo-signos-vitales/internal/models"
	"monitoreo-signos-vitales/internal/storage"
)

func TestCrearEvaluaRegistroCritico(t *testing.T) {
	servicio := NuevoSignosVitalesService(storage.NuevaMemoriaRepository())
	registro, err := servicio.Crear(models.EntradaSignosVitales{
		IDAdultoMayor: "1", PresionSistolica: 165, PresionDiastolica: 80, FrecuenciaCardiaca: 70,
	})
	if err != nil {
		t.Fatalf("Crear devolvió error: %v", err)
	}
	if registro.Evaluacion.EstadoGeneral != models.EstadoCritico {
		t.Errorf("estado general = %q; se esperaba critico", registro.Evaluacion.EstadoGeneral)
	}
}

func TestCrearRechazaCamposObligatorios(t *testing.T) {
	servicio := NuevoSignosVitalesService(storage.NuevaMemoriaRepository())
	_, err := servicio.Crear(models.EntradaSignosVitales{IDAdultoMayor: "1", PresionSistolica: 120, PresionDiastolica: 80})
	if err == nil {
		t.Fatal("Crear aceptó una frecuencia cardiaca inválida")
	}
}
