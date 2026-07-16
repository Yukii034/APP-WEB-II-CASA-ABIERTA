package storage

import (
	"time"

	"monitoreo-signos-vitales/internal/models"
)

// Sembrar loads representative records for the Casa Abierta demonstration.
func Sembrar(repo SignosVitalesRepository) {
	temperaturaNormal, temperaturaAlta := 36.7, 38.5
	spo2Normal, spo2Bajo := 97, 87
	glucosaNormal, glucosaAlta := 108.0, 215.0
	dolorLeve, dolorAlto := 2, 7

	ejemplos := []models.EntradaSignosVitales{
		{IDAdultoMayor: "1", RegistradoPor: "Cuidador Ana", PresionSistolica: 122, PresionDiastolica: 78, FrecuenciaCardiaca: 72, Temperatura: &temperaturaNormal, SaturacionOxigeno: &spo2Normal, NivelGlucosa: &glucosaNormal, NivelDolor: &dolorLeve, Observaciones: "Control de rutina."},
		{IDAdultoMayor: "1", RegistradoPor: "Cuidador Ana", PresionSistolica: 145, PresionDiastolica: 92, FrecuenciaCardiaca: 98, Temperatura: &temperaturaNormal, SaturacionOxigeno: &spo2Normal, NivelGlucosa: &glucosaNormal, NivelDolor: &dolorLeve, Observaciones: "Revisar presión en próxima toma."},
		{IDAdultoMayor: "2", RegistradoPor: "Cuidador Luis", PresionSistolica: 165, PresionDiastolica: 104, FrecuenciaCardiaca: 122, Temperatura: &temperaturaAlta, SaturacionOxigeno: &spo2Bajo, NivelGlucosa: &glucosaAlta, NivelDolor: &dolorAlto, Observaciones: "Se recomienda contactar al profesional de salud."},
	}

	for indice, entrada := range ejemplos {
		registro := models.SignosVitales{ID: repo.SiguienteID(), IDAdultoMayor: entrada.IDAdultoMayor, RegistradoPor: entrada.RegistradoPor, FechaRegistro: time.Now().Add(-time.Duration(indice) * 24 * time.Hour), PresionSistolica: entrada.PresionSistolica, PresionDiastolica: entrada.PresionDiastolica, FrecuenciaCardiaca: entrada.FrecuenciaCardiaca, Temperatura: entrada.Temperatura, SaturacionOxigeno: entrada.SaturacionOxigeno, NivelGlucosa: entrada.NivelGlucosa, Peso: entrada.Peso, Altura: entrada.Altura, NivelDolor: entrada.NivelDolor, Observaciones: entrada.Observaciones}
		repo.Guardar(registro)
	}
}
