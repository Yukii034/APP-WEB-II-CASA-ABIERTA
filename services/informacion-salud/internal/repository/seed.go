package repository

import (
	"time"

	"cuidabien/informacion-salud/internal/model"
)

// Sembrar agrega fichas de ejemplo al Repository. Sirve para tener
// datos visibles al levantar el servicio (demos, casa abierta, pruebas
// manuales en Postman) sin tener que crear todo a mano primero.
func Sembrar(repo Repository) {
	ahora := time.Now()

	ejemplos := []model.InformacionSalud{
		{
			NombrePaciente:       "María Pérez",
			Diagnosticos:         []string{"hipertensión arterial"},
			Alergias:             []string{"penicilina"},
			EnfermedadesCronicas: []string{"diabetes tipo 2"},
			AntecedentesMedicos:  []string{"cirugía de cadera (2019)"},
		},
		{
			NombrePaciente:       "José Ramírez",
			Diagnosticos:         []string{"artritis"},
			Alergias:             []string{},
			EnfermedadesCronicas: []string{"hipotiroidismo"},
			AntecedentesMedicos:  []string{"marcapasos (2021)"},
		},
		{
			NombrePaciente:       "Carmen Torres",
			Diagnosticos:         []string{},
			Alergias:             []string{"aspirina", "mariscos"},
			EnfermedadesCronicas: []string{},
			AntecedentesMedicos:  []string{"ninguno relevante"},
		},
	}

	for _, ejemplo := range ejemplos {
		ejemplo.ID = repo.SiguienteID()
		ejemplo.ActualizadoEn = ahora
		repo.Guardar(ejemplo)
	}
}
