package service

import (
	"cuidabien/actividad-fisica/internal/model"
	"cuidabien/actividad-fisica/internal/repository"
	"testing"
)

func TestCrearActividad(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()
	svc := NuevoActividadFisicaService(repo)
	actividad, err := svc.Crear(model.EntradaActividadFisica{
		NombrePaciente: "Ana", TipoActividad: "Caminata", DuracionMinutos: 30,
		Intensidad: "moderada", Fecha: "2026-07-15", Estado: "completada",
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actividad.CaloriasEstimadas != 150 {
		t.Fatalf("se esperaban 150 calorías, se obtuvo %d", actividad.CaloriasEstimadas)
	}
}

func TestCrearActividadConIntensidadInvalida(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()
	svc := NuevoActividadFisicaService(repo)
	_, err := svc.Crear(model.EntradaActividadFisica{
		NombrePaciente: "Ana", TipoActividad: "Caminata", DuracionMinutos: 30,
		Intensidad: "extrema", Fecha: "2026-07-15", Estado: "completada",
	})
	if err != ErrIntensidadInvalida {
		t.Fatalf("se esperaba ErrIntensidadInvalida, se obtuvo %v", err)
	}
}
