package repository

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"cuidabien/actividad-fisica/internal/model"
)

// Repository define las operaciones de persistencia del microservicio.
type Repository interface {
	Listar() []model.ActividadFisica
	Obtener(id string) (model.ActividadFisica, bool)
	Guardar(actividad model.ActividadFisica)
	Eliminar(id string) bool
	SiguienteID() string
}

type MemoriaRepository struct {
	mu       sync.RWMutex
	datos    map[string]model.ActividadFisica
	ultimoID int
}

func NuevaMemoriaRepository() *MemoriaRepository {
	return &MemoriaRepository{datos: make(map[string]model.ActividadFisica)}
}

func (r *MemoriaRepository) Listar() []model.ActividadFisica {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resultado := make([]model.ActividadFisica, 0, len(r.datos))
	for _, actividad := range r.datos {
		resultado = append(resultado, actividad)
	}
	sort.Slice(resultado, func(i, j int) bool {
		a, _ := strconv.Atoi(resultado[i].ID)
		b, _ := strconv.Atoi(resultado[j].ID)
		return a < b
	})
	return resultado
}

func (r *MemoriaRepository) Obtener(id string) (model.ActividadFisica, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	actividad, ok := r.datos[id]
	return actividad, ok
}

func (r *MemoriaRepository) Guardar(actividad model.ActividadFisica) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.datos[actividad.ID] = actividad
	if id, err := strconv.Atoi(actividad.ID); err == nil && id > r.ultimoID {
		r.ultimoID = id
	}
}

func (r *MemoriaRepository) Eliminar(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.datos[id]; !ok {
		return false
	}
	delete(r.datos, id)
	return true
}

func (r *MemoriaRepository) SiguienteID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ultimoID++
	return fmt.Sprintf("%d", r.ultimoID)
}
