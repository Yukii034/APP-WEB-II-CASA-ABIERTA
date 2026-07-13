package service

import (
	"time"

	"cuidabien/informacion-salud/internal/model"
	"cuidabien/informacion-salud/internal/repository"
)

// InformacionSaludService contiene la lógica de negocio del servicio.
// No sabe si los datos están en memoria o en una base de datos: solo
// conoce la interfaz repository.Repository.
type InformacionSaludService struct {
	repo repository.Repository
}

func NuevoInformacionSaludService(repo repository.Repository) *InformacionSaludService {
	return &InformacionSaludService{repo: repo}
}

// Listar devuelve todas las fichas registradas.
func (s *InformacionSaludService) Listar() []model.InformacionSalud {
	return s.repo.Listar()
}

// Crear registra una nueva ficha de salud.
func (s *InformacionSaludService) Crear(entrada model.EntradaInformacionSalud) model.InformacionSalud {
	nuevo := nuevoRegistro(s.repo.SiguienteID(), entrada, time.Now())
	s.repo.Guardar(nuevo)
	return nuevo
}

// Obtener busca la ficha de un paciente por id.
func (s *InformacionSaludService) Obtener(id string) (model.InformacionSalud, bool) {
	return s.repo.Obtener(id)
}

// Actualizar aplica una actualización parcial sobre la ficha de un
// paciente existente. Devuelve ok=false si el id no existe.
func (s *InformacionSaludService) Actualizar(id string, entrada model.EntradaInformacionSalud) (model.InformacionSalud, bool) {
	existente, ok := s.repo.Obtener(id)
	if !ok {
		return model.InformacionSalud{}, false
	}

	actualizado := actualizarRegistro(existente, entrada, time.Now())
	s.repo.Guardar(actualizado)
	return actualizado, true
}

// nuevoRegistro es una función pura (sin estado) para que sea fácil de
// probar con go test, sin necesidad de mockear el repository.
func nuevoRegistro(id string, entrada model.EntradaInformacionSalud, ahora time.Time) model.InformacionSalud {
	return model.InformacionSalud{
		ID:                   id,
		NombrePaciente:       entrada.NombrePaciente,
		Diagnosticos:         normalizar(entrada.Diagnosticos),
		Alergias:             normalizar(entrada.Alergias),
		EnfermedadesCronicas: normalizar(entrada.EnfermedadesCronicas),
		AntecedentesMedicos:  normalizar(entrada.AntecedentesMedicos),
		ActualizadoEn:        ahora,
	}
}

// actualizarRegistro combina la ficha existente con los campos
// enviados, sin borrar datos que no vinieron en la petición
// (actualización parcial).
func actualizarRegistro(existente model.InformacionSalud, entrada model.EntradaInformacionSalud, ahora time.Time) model.InformacionSalud {
	if entrada.NombrePaciente != "" {
		existente.NombrePaciente = entrada.NombrePaciente
	}
	if entrada.Diagnosticos != nil {
		existente.Diagnosticos = normalizar(entrada.Diagnosticos)
	}
	if entrada.Alergias != nil {
		existente.Alergias = normalizar(entrada.Alergias)
	}
	if entrada.EnfermedadesCronicas != nil {
		existente.EnfermedadesCronicas = normalizar(entrada.EnfermedadesCronicas)
	}
	if entrada.AntecedentesMedicos != nil {
		existente.AntecedentesMedicos = normalizar(entrada.AntecedentesMedicos)
	}
	existente.ActualizadoEn = ahora
	return existente
}

// normalizar evita que el JSON de salida muestre "null" en vez de una
// lista vacía cuando no se envía ese campo.
func normalizar(valores []string) []string {
	if valores == nil {
		return []string{}
	}
	return valores
}
