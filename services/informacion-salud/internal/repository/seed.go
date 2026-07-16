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
			ID:                   "P001",
			NombrePaciente:       "Maria Garcia",
			Diagnosticos:         []string{"hipertensión arterial"},
			Alergias:             []string{"penicilina"},
			EnfermedadesCronicas: []string{"diabetes tipo 2"},
			AntecedentesMedicos:  []string{"cirugía de cadera (2019)"},
		},
		{
			ID:                   "P002",
			NombrePaciente:       "Juan Lopez",
			Diagnosticos:         []string{"artritis"},
			Alergias:             []string{},
			EnfermedadesCronicas: []string{"hipotiroidismo"},
			AntecedentesMedicos:  []string{"marcapasos (2021)"},
		},
		{
			ID:                   "P003",
			NombrePaciente:       "Ana Martinez",
			Diagnosticos:         []string{},
			Alergias:             []string{"aspirina", "mariscos"},
			EnfermedadesCronicas: []string{},
			AntecedentesMedicos:  []string{"ninguno relevante"},
		},
	}

	for _, ejemplo := range ejemplos {
		ejemplo.ActualizadoEn = ahora
		repo.Guardar(ejemplo)
	}
}
