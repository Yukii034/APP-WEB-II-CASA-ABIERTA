package services

import (
	"fmt"

	"cuidabien/estado-animo/models"
)

func GenerarAlerta(registros []models.RegistroEstadoAnimo) models.Alerta {

	alerta := models.Alerta{
		GenerarAlerta: false,
		Mensaje:       "El estado de ánimo se encuentra estable.",
	}

	n := len(registros)

	if n >= 2 {

		ultimo := registros[n-1]
		penultimo := registros[n-2]

		if ultimo.Nivel <= 2 && penultimo.Nivel <= 2 {

			alerta.GenerarAlerta = true
			alerta.Mensaje = fmt.Sprintf(
				"ALERTA: Se ha detectado un estado de ánimo bajo persistente. Últimas emociones: '%s' y '%s'.",
				penultimo.Emocion,
				ultimo.Emocion,
			)
		}
	}

	return alerta
}