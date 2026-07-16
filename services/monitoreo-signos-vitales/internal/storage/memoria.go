package storage

import (
	"sort"
	"strconv"
	"sync"

	"monitoreo-signos-vitales/internal/models"
)

// MemoriaRepository stores records only while the service is running.
type MemoriaRepository struct {
	mu          sync.RWMutex
	registros   []models.SignosVitales
	siguienteID int
}

func NuevaMemoriaRepository() *MemoriaRepository {
	return &MemoriaRepository{siguienteID: 1}
}

func (r *MemoriaRepository) Guardar(registro models.SignosVitales) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registros = append(r.registros, registro)
}

func (r *MemoriaRepository) PorAdultoMayor(idAdultoMayor string) []models.SignosVitales {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resultado := make([]models.SignosVitales, 0)
	for _, registro := range r.registros {
		if registro.IDAdultoMayor == idAdultoMayor {
			resultado = append(resultado, registro)
		}
	}
	sort.Slice(resultado, func(i, j int) bool { return resultado[i].FechaRegistro.After(resultado[j].FechaRegistro) })
	return resultado
}

func (r *MemoriaRepository) SiguienteID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := strconv.Itoa(r.siguienteID)
	r.siguienteID++
	return id
}
