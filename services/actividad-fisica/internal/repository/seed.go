package repository

import (
	"time"

	"cuidabien/actividad-fisica/internal/model"
)

// Sembrar agrega registros de ejemplo para facilitar las pruebas iniciales.
func Sembrar(repo Repository) {
	ahora := time.Now()
	repo.Guardar(model.ActividadFisica{
		ID: "1", NombrePaciente: "María López", TipoActividad: "Caminata",
		DuracionMinutos: 30, Intensidad: "moderada", Fecha: ahora.Format("2006-01-02"),
		Estado: "completada", Observaciones: "Actividad realizada sin molestias",
		CaloriasEstimadas: 150, CreadoEn: ahora, ActualizadoEn: ahora,
	})
}
