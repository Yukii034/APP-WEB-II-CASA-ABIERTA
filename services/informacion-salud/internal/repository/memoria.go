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
// con un mutex para uso concurrente. Los datos se pierden si el
// contenedor se reinicia (ver docs/arquitectura.md - limitaciones).
type memoriaRepository struct {
	mu          sync.Mutex
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

func (r *memoriaRepository) Listar() []model.InformacionSalud {
	r.mu.Lock()
	defer r.mu.Unlock()

	resultado := make([]model.InformacionSalud, 0, len(r.registros))
	for _, reg := range r.registros {
		resultado = append(resultado, reg)
	}
	return resultado
}

func (r *memoriaRepository) Obtener(id string) (model.InformacionSalud, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	reg, ok := r.registros[id]
	return reg, ok
}

func (r *memoriaRepository) Guardar(registro model.InformacionSalud) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registros[registro.ID] = registro
}

func (r *memoriaRepository) SiguienteID() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := strconv.Itoa(r.siguienteID)
	r.siguienteID++
	return id
}
