package storage

import "cuidabien/reportes-medicos/models"

type Store struct {
	Reportes []models.ReportePaciente
}

func NewStore() *Store {
	return &Store{
		Reportes: []models.ReportePaciente{
			{
				PacienteID:          "P001",
				Nombre:              "Maria Garcia",
				Periodo:             "semana actual",
				CitasProgramadas:    3,
				CitasCompletadas:    2,
				ComidasRegistradas:  18,
				AlertasSalud:        1,
				AdherenciaMedicinas: 92,
				EstadoGeneral:       "estable",
				Recomendaciones: []string{
					"Mantener controles medicos programados",
					"Continuar con horarios de alimentacion",
				},
			},
			{
				PacienteID:          "P002",
				Nombre:              "Juan Lopez",
				Periodo:             "semana actual",
				CitasProgramadas:    2,
				CitasCompletadas:    1,
				ComidasRegistradas:  12,
				AlertasSalud:        3,
				AdherenciaMedicinas: 76,
				EstadoGeneral:       "requiere seguimiento",
				Recomendaciones: []string{
					"Revisar alertas de salud con el cuidador",
					"Mejorar cumplimiento de medicacion",
				},
			},
		},
	}
}

func (s *Store) ListarReportes() []models.ReportePaciente {
	if s.Reportes == nil {
		return []models.ReportePaciente{}
	}
	return s.Reportes
}

func (s *Store) BuscarPorPaciente(id string) *models.ReportePaciente {
	for i := range s.Reportes {
		if s.Reportes[i].PacienteID == id {
			return &s.Reportes[i]
		}
	}
	return nil
}

func (s *Store) CrearResumen() models.ResumenGeneral {
	return CrearResumen(s.Reportes)
}

func CrearResumen(data []models.ReportePaciente) models.ResumenGeneral {
	resumen := models.ResumenGeneral{
		PacientesEvaluados: len(data),
		EstadoGeneral:      "estable",
		Promedios:          map[string]int{},
		Pacientes:          data,
	}

	if len(data) == 0 {
		resumen.EstadoGeneral = "sin datos"
		return resumen
	}

	var totalAlertas, totalAdherencia, totalComidas int
	for _, reporte := range data {
		totalAlertas += reporte.AlertasSalud
		totalAdherencia += reporte.AdherenciaMedicinas
		totalComidas += reporte.ComidasRegistradas
	}

	resumen.AlertasTotales = totalAlertas
	resumen.Promedios["adherencia_medicinas"] = totalAdherencia / len(data)
	resumen.Promedios["comidas_registradas"] = totalComidas / len(data)

	if totalAlertas >= 3 || resumen.Promedios["adherencia_medicinas"] < 80 {
		resumen.EstadoGeneral = "requiere seguimiento"
	}

	return resumen
}
