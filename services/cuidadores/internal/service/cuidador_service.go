package service

import (
	"errors"

	"cuidabien/cuidadores/internal/model"
	"cuidabien/cuidadores/internal/repository"
)

// ErrValidacion se devuelve cuando los datos enviados no cumplen
// las reglas mínimas del negocio.
var ErrValidacion = errors.New("nombre y telefono son obligatorios")

type CuidadorService struct {
	repo repository.CuidadorRepository
}

func NuevoCuidadorService(repo repository.CuidadorRepository) *CuidadorService {
	return &CuidadorService{repo: repo}
}

func (s *CuidadorService) Crear(c model.Cuidador) (model.Cuidador, error) {
	if c.Nombre == "" || c.Telefono == "" {
		return model.Cuidador{}, ErrValidacion
	}
	c.Activo = true
	return s.repo.Crear(c), nil
}

func (s *CuidadorService) Listar() []model.Cuidador {
	return s.repo.Listar()
}

func (s *CuidadorService) ObtenerPorID(id string) (model.Cuidador, error) {
	return s.repo.ObtenerPorID(id)
}

func (s *CuidadorService) ObtenerPorPaciente(pacienteID string) []model.Cuidador {
	return s.repo.ObtenerPorPaciente(pacienteID)
}

func (s *CuidadorService) Actualizar(id string, c model.Cuidador) (model.Cuidador, error) {
	if c.Nombre == "" || c.Telefono == "" {
		return model.Cuidador{}, ErrValidacion
	}
	return s.repo.Actualizar(id, c)
}

func (s *CuidadorService) Eliminar(id string) error {
	return s.repo.Eliminar(id)
}
