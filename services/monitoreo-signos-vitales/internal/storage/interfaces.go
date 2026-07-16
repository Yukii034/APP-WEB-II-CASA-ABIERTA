package storage

import "monitoreo-signos-vitales/internal/models"

// SignosVitalesRepository isolates persistence from business rules.
type SignosVitalesRepository interface {
	Guardar(models.SignosVitales)
	PorAdultoMayor(idAdultoMayor string) []models.SignosVitales
	SiguienteID() string
}
