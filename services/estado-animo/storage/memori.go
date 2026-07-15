package storage

import (
	"sync"
	"time"

	"cuidabien/estado-animo/models"
)

var (
	BaseDeDatos []models.RegistroEstadoAnimo
	UltimoID    int
	Mutex       sync.Mutex
)

func InicializarDatos() {
	BaseDeDatos = []models.RegistroEstadoAnimo{
		{
			ID:         "1",
			Fecha:      time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			Nivel:      4,
			Emocion:    "Tranquilo",
			Comentario: "Pasé una tarde agradable.",
		},
		{
			ID:         "2",
			Fecha:      time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			Nivel:      5,
			Emocion:    "Feliz",
			Comentario: "Me visitaron mis nietos.",
		},
	}

	UltimoID = 2
}