package service

import (
	"errors"
	"strings"
	"time"

	"cuidabien/actividad-fisica/internal/model"
	"cuidabien/actividad-fisica/internal/repository"
)

var (
	ErrNombreRequerido    = errors.New("el nombre del paciente es obligatorio")
	ErrTipoRequerido      = errors.New("el tipo de actividad es obligatorio")
	ErrDuracionInvalida   = errors.New("la duración debe ser mayor que cero")
	ErrIntensidadInvalida = errors.New("la intensidad debe ser baja, moderada o alta")
	ErrFechaInvalida      = errors.New("la fecha debe tener el formato AAAA-MM-DD")
	ErrEstadoInvalido     = errors.New("el estado debe ser pendiente, completada o cancelada")
)

type ActividadFisicaService struct{ repo repository.Repository }

func NuevoActividadFisicaService(repo repository.Repository) *ActividadFisicaService {
	return &ActividadFisicaService{repo: repo}
}

func (s *ActividadFisicaService) Listar() []model.ActividadFisica { return s.repo.Listar() }
func (s *ActividadFisicaService) Obtener(id string) (model.ActividadFisica, bool) {
	return s.repo.Obtener(id)
}

func (s *ActividadFisicaService) Crear(entrada model.EntradaActividadFisica) (model.ActividadFisica, error) {
	entrada = normalizarEntrada(entrada)
	if err := validar(entrada); err != nil {
		return model.ActividadFisica{}, err
	}
	ahora := time.Now()
	actividad := construir(s.repo.SiguienteID(), entrada, ahora, ahora)
	s.repo.Guardar(actividad)
	return actividad, nil
}

func (s *ActividadFisicaService) Actualizar(id string, entrada model.EntradaActividadFisica) (model.ActividadFisica, error, bool) {
	existente, ok := s.repo.Obtener(id)
	if !ok {
		return model.ActividadFisica{}, nil, false
	}
	entrada = normalizarEntrada(entrada)
	if err := validar(entrada); err != nil {
		return model.ActividadFisica{}, err, true
	}
	actualizada := construir(id, entrada, existente.CreadoEn, time.Now())
	s.repo.Guardar(actualizada)
	return actualizada, nil, true
}

func (s *ActividadFisicaService) Eliminar(id string) bool { return s.repo.Eliminar(id) }

func normalizarEntrada(e model.EntradaActividadFisica) model.EntradaActividadFisica {
	e.NombrePaciente = strings.TrimSpace(e.NombrePaciente)
	e.TipoActividad = strings.TrimSpace(e.TipoActividad)
	e.Intensidad = strings.ToLower(strings.TrimSpace(e.Intensidad))
	e.Estado = strings.ToLower(strings.TrimSpace(e.Estado))
	e.Fecha = strings.TrimSpace(e.Fecha)
	e.Observaciones = strings.TrimSpace(e.Observaciones)
	return e
}

func validar(e model.EntradaActividadFisica) error {
	if e.NombrePaciente == "" {
		return ErrNombreRequerido
	}
	if e.TipoActividad == "" {
		return ErrTipoRequerido
	}
	if e.DuracionMinutos <= 0 {
		return ErrDuracionInvalida
	}
	if e.Intensidad != "baja" && e.Intensidad != "moderada" && e.Intensidad != "alta" {
		return ErrIntensidadInvalida
	}
	if _, err := time.Parse("2006-01-02", e.Fecha); err != nil {
		return ErrFechaInvalida
	}
	if e.Estado != "pendiente" && e.Estado != "completada" && e.Estado != "cancelada" {
		return ErrEstadoInvalido
	}
	return nil
}

func construir(id string, e model.EntradaActividadFisica, creado, actualizado time.Time) model.ActividadFisica {
	return model.ActividadFisica{
		ID: id, NombrePaciente: e.NombrePaciente, TipoActividad: e.TipoActividad,
		DuracionMinutos: e.DuracionMinutos, Intensidad: e.Intensidad, Fecha: e.Fecha,
		Estado: e.Estado, Observaciones: e.Observaciones,
		CaloriasEstimadas: calcularCalorias(e.DuracionMinutos, e.Intensidad),
		CreadoEn:          creado, ActualizadoEn: actualizado,
	}
}

func calcularCalorias(duracion int, intensidad string) int {
	factor := 3
	switch intensidad {
	case "moderada":
		factor = 5
	case "alta":
		factor = 8
	}
	return duracion * factor
}
