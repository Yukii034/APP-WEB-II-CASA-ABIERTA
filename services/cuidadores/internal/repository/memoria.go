package repository

import (
	"errors"
	"strconv"
	"sync"

	"cuidabien/cuidadores/internal/model"
)

// ErrNoEncontrado se devuelve cuando no existe un cuidador con el ID pedido.
var ErrNoEncontrado = errors.New("cuidador no encontrado")

// CuidadorRepository define el contrato de almacenamiento.
// Empezamos en memoria; si luego se necesita una base de datos,
// solo hay que crear otra implementación de esta interfaz.
type CuidadorRepository interface {
	Crear(c model.Cuidador) model.Cuidador
	Listar() []model.Cuidador
	ObtenerPorID(id string) (model.Cuidador, error)
	ObtenerPorPaciente(pacienteID string) []model.Cuidador
	Actualizar(id string, c model.Cuidador) (model.Cuidador, error)
	Eliminar(id string) error
}

// MemoriaRepository es la implementación en memoria, segura para uso
// concurrente mediante un RWMutex.
type MemoriaRepository struct {
	mu          sync.RWMutex
	data        map[string]model.Cuidador
	consecutivo int
}

func NuevaMemoriaRepository() *MemoriaRepository {
	return &MemoriaRepository{data: make(map[string]model.Cuidador)}
}

func (r *MemoriaRepository) Crear(c model.Cuidador) model.Cuidador {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.consecutivo++
	c.ID = strconv.Itoa(r.consecutivo)
	r.data[c.ID] = c
	return c
}

func (r *MemoriaRepository) Listar() []model.Cuidador {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lista := make([]model.Cuidador, 0, len(r.data))
	for _, c := range r.data {
		lista = append(lista, c)
	}
	return lista
}

func (r *MemoriaRepository) ObtenerPorID(id string) (model.Cuidador, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.data[id]
	if !ok {
		return model.Cuidador{}, ErrNoEncontrado
	}
	return c, nil
}

func (r *MemoriaRepository) ObtenerPorPaciente(pacienteID string) []model.Cuidador {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var lista []model.Cuidador
	for _, c := range r.data {
		for _, p := range c.Pacientes {
			if p == pacienteID {
				lista = append(lista, c)
				break
			}
		}
	}
	return lista
}

func (r *MemoriaRepository) Actualizar(id string, c model.Cuidador) (model.Cuidador, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return model.Cuidador{}, ErrNoEncontrado
	}
	c.ID = id
	r.data[id] = c
	return c, nil
}

func (r *MemoriaRepository) Eliminar(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrNoEncontrado
	}
	delete(r.data, id)
	return nil
}
