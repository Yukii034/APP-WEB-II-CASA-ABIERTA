package repository

import (
	"time"

	"cuidabien/recordatorios-medicamentos/internal/model"
)

// Sembrar agrega recordatorios de ejemplo para demostraciones
// y pruebas manuales sin tener que crearlos desde cero.
func Sembrar(repo Repository) {
	ahora := time.Now()

	ejemplos := []model.RecordatorioMedicamento{
		{
			AdultoMayorID:  "AM-001",
			NombrePaciente: "María Pérez",
			Medicamento:    "Losartán",
			Dosis:          "1 tableta",
			Hora:           "08:00",
			Frecuencia:     "diaria",
			Activo:         true,
		},
		{
			AdultoMayorID:  "AM-002",
			NombrePaciente: "José Ramírez",
			Medicamento:    "Metformina",
			Dosis:          "500 mg",
			Hora:           "14:30",
			Frecuencia:     "diaria",
			Activo:         true,
		},
	}

	for _, ejemplo := range ejemplos {
		ejemplo.ID = repo.SiguienteID()
		ejemplo.CreadoEn = ahora
		ejemplo.ActualizadoEn = ahora

		repo.Guardar(ejemplo)
	}
}
