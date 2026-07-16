package services

import (
	"cuidabien/estado-animo/models"
	"testing"
)

// 1. Test para verificar que NO se genera alerta con registros normales
func TestGenerarAlerta_Estable(t *testing.T) {
	registros := []models.RegistroEstadoAnimo{
		{Nivel: 4, Emocion: "Tranquilo"},
		{Nivel: 5, Emocion: "Feliz"},
	}
	alerta := GenerarAlerta(registros)
	if alerta.GenerarAlerta {
		t.Error("No debería generar alerta en estado estable")
	}
}

// 2. Test para verificar que SÍ se genera alerta con registros bajos persistentes
func TestGenerarAlerta_BajoPersistente(t *testing.T) {
	registros := []models.RegistroEstadoAnimo{
		{Nivel: 2, Emocion: "Triste"},
		{Nivel: 1, Emocion: "Ansioso"},
	}
	alerta := GenerarAlerta(registros)
	if !alerta.GenerarAlerta {
		t.Error("Debería generar alerta al detectar dos niveles bajos consecutivos")
	}
}

// 3. Test para verificar que NO se genera alerta con datos insuficientes
func TestGenerarAlerta_DatosInsuficientes(t *testing.T) {
	registros := []models.RegistroEstadoAnimo{
		{Nivel: 1, Emocion: "Triste"},
	}
	alerta := GenerarAlerta(registros)
	if alerta.GenerarAlerta {
		t.Error("No debería generar alerta si hay menos de dos registros")
	}
}
