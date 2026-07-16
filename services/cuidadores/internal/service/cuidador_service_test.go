package service

import (
	"errors"
	"testing"

	"cuidabien/cuidadores/internal/model"
	"cuidabien/cuidadores/internal/repository"
)

func TestCrear_Valido(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()
	svc := NuevoCuidadorService(repo)

	c, err := svc.Crear(model.Cuidador{Nombre: "Jeremy", Telefono: "0999999999"})
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if c.ID == "" {
		t.Error("se esperaba que el repositorio asignara un ID")
	}
	if !c.Activo {
		t.Error("un cuidador nuevo debe quedar activo por defecto")
	}
}

func TestCrear_SinNombreFalla(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()
	svc := NuevoCuidadorService(repo)

	_, err := svc.Crear(model.Cuidador{Telefono: "0999999999"})
	if !errors.Is(err, ErrValidacion) {
		t.Fatalf("se esperaba ErrValidacion, se obtuvo: %v", err)
	}
}

func TestObtenerPorID_NoExiste(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()
	svc := NuevoCuidadorService(repo)

	_, err := svc.ObtenerPorID("no-existe")
	if !errors.Is(err, repository.ErrNoEncontrado) {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo: %v", err)
	}
}

func TestObtenerPorPaciente(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()
	svc := NuevoCuidadorService(repo)

	svc.Crear(model.Cuidador{Nombre: "Jeremy", Telefono: "099", Pacientes: []string{"p1"}})
	svc.Crear(model.Cuidador{Nombre: "Ana", Telefono: "098", Pacientes: []string{"p2"}})

	resultado := svc.ObtenerPorPaciente("p1")
	if len(resultado) != 1 || resultado[0].Nombre != "Jeremy" {
		t.Fatalf("se esperaba solo a Jeremy asociado a p1, se obtuvo: %+v", resultado)
	}
}

func TestEliminar(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()
	svc := NuevoCuidadorService(repo)

	c, _ := svc.Crear(model.Cuidador{Nombre: "Jeremy", Telefono: "099"})
	if err := svc.Eliminar(c.ID); err != nil {
		t.Fatalf("no se esperaba error al eliminar: %v", err)
	}
	if _, err := svc.ObtenerPorID(c.ID); !errors.Is(err, repository.ErrNoEncontrado) {
		t.Error("el cuidador debería haber sido eliminado")
	}
}
