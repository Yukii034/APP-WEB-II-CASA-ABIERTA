package repository

import (
	"strconv"
	"sync"

	"cuidabien/recordatorios-medicamentos/internal/model"
)

// Repository define las operaciones de persistencia
// que necesita la capa service.
type Repository interface {
	Listar() []model.RecordatorioMedicamento
	Obtener(id string) (model.RecordatorioMedicamento, bool)
	Guardar(registro model.RecordatorioMedicamento)
	Eliminar(id string) bool
	SiguienteID() string
	BuscarActivosPorHora(hora string) []model.RecordatorioMedicamento
}

// memoriaRepository guarda los recordatorios en memoria,
// protegidos con un mutex para permitir peticiones concurrentes.
type memoriaRepository struct {
	mu          sync.Mutex
	registros   map[string]model.RecordatorioMedicamento
	siguienteID int
}

// NuevaMemoriaRepository crea un Repository
// respaldado por almacenamiento en memoria.
func NuevaMemoriaRepository() Repository {
	return &memoriaRepository{
		registros:   map[string]model.RecordatorioMedicamento{},
		siguienteID: 1,
	}
}

// Listar devuelve todos los recordatorios almacenados.
func (r *memoriaRepository) Listar() []model.RecordatorioMedicamento {
	r.mu.Lock()
	defer r.mu.Unlock()

	resultado := make(
		[]model.RecordatorioMedicamento,
		0,
		len(r.registros),
	)

	for _, registro := range r.registros {
		resultado = append(resultado, registro)
	}

	return resultado
}

// Obtener busca un recordatorio mediante su ID.
func (r *memoriaRepository) Obtener(
	id string,
) (model.RecordatorioMedicamento, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	registro, ok := r.registros[id]

	return registro, ok
}

// Guardar crea o reemplaza un recordatorio.
func (r *memoriaRepository) Guardar(
	registro model.RecordatorioMedicamento,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registros[registro.ID] = registro
}

// Eliminar borra un recordatorio por su ID.
func (r *memoriaRepository) Eliminar(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.registros[id]; !ok {
		return false
	}

	delete(r.registros, id)

	return true
}

// SiguienteID genera identificadores consecutivos.
func (r *memoriaRepository) SiguienteID() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := strconv.Itoa(r.siguienteID)
	r.siguienteID++

	return id
}

// BuscarActivosPorHora devuelve los recordatorios activos
// que coinciden con la hora recibida.
func (r *memoriaRepository) BuscarActivosPorHora(
	hora string,
) []model.RecordatorioMedicamento {
	r.mu.Lock()
	defer r.mu.Unlock()

	resultado := make(
		[]model.RecordatorioMedicamento,
		0,
	)

	for _, registro := range r.registros {
		if registro.Activo && registro.Hora == hora {
			resultado = append(resultado, registro)
		}
	}

	return resultado
}
