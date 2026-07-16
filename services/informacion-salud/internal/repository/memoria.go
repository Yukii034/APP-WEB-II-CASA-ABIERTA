package repository

import (
	"strconv"
	"sync"

	"cuidabien/informacion-salud/internal/model"
)

// Repository define las operaciones de persistencia que necesita el
// service, sin importarle cómo están implementadas. Hoy solo existe
// memoriaRepository, pero mañana se podría agregar una implementación
// con base de datos sin tocar el service ni los handlers.
type Repository interface {
	Listar() []model.InformacionSalud
	Obtener(id string) (model.InformacionSalud, bool)
	Guardar(registro model.InformacionSalud)
	SiguienteID() string
}

// memoriaRepository guarda las fichas en un mapa en memoria, protegido
// con un RWMutex para uso concurrente (múltiples lecturas en paralelo,
// una sola escritura a la vez). Los datos se pierden si el contenedor
// se reinicia (ver docs/arquitectura.md - limitaciones).
type memoriaRepository struct {
	mu          sync.RWMutex
	registros   map[string]model.InformacionSalud
	siguienteID int
}

// NuevaMemoriaRepository crea un Repository respaldado en memoria.
func NuevaMemoriaRepository() Repository {
	return &memoriaRepository{
		registros:   map[string]model.InformacionSalud{},
		siguienteID: 1,
	}
}

// Listar devuelve copias de las fichas para preservar el estado interno.
func (r *memoriaRepository) Listar() []model.InformacionSalud {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resultado := make([]model.InformacionSalud, 0, len(r.registros))
	for _, reg := range r.registros {
		resultado = append(resultado, copiarRegistro(reg))
	}
	return resultado
}

// Obtener busca una ficha y también devuelve una copia independiente.
func (r *memoriaRepository) Obtener(id string) (model.InformacionSalud, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reg, ok := r.registros[id]
	if !ok {
		return model.InformacionSalud{}, false
	}
	return copiarRegistro(reg), true
}

// Guardar almacena una copia para evitar que el llamador altere los datos guardados.
func (r *memoriaRepository) Guardar(registro model.InformacionSalud) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registros[registro.ID] = copiarRegistro(registro)
}

// SiguienteID genera identificadores consecutivos de forma segura.
func (r *memoriaRepository) SiguienteID() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := strconv.Itoa(r.siguienteID)
	r.siguienteID++
	return id
}

// copiarRegistro devuelve una copia profunda de las listas del
// registro, para que quien reciba el valor (de Listar, Obtener o
// Guardar) no pueda mutar por accidente el estado interno del
// repository a través de sus slices.
func copiarRegistro(reg model.InformacionSalud) model.InformacionSalud {
	copia := reg
	copia.Diagnosticos = append([]string(nil), reg.Diagnosticos...)
	copia.Alergias = append([]string(nil), reg.Alergias...)
	copia.EnfermedadesCronicas = append([]string(nil), reg.EnfermedadesCronicas...)
	copia.AntecedentesMedicos = append([]string(nil), reg.AntecedentesMedicos...)
	return copia
}
