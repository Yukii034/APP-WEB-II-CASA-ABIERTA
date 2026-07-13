package main

import (
	"testing"
	"time"
)

func TestNuevoRegistroNormalizaListasNulas(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC)
	entrada := entradaInformacionSalud{
		NombrePaciente: "María Pérez",
		Diagnosticos:   []string{"hipertensión"},
	}

	r := nuevoRegistro("1", entrada, ahora)

	if r.ID != "1" {
		t.Errorf("esperaba id 1, obtuve %s", r.ID)
	}
	if r.NombrePaciente != "María Pérez" {
		t.Errorf("esperaba nombre 'María Pérez', obtuve %s", r.NombrePaciente)
	}
	if len(r.Diagnosticos) != 1 || r.Diagnosticos[0] != "hipertensión" {
		t.Errorf("diagnósticos no se copiaron correctamente: %v", r.Diagnosticos)
	}
	if r.Alergias == nil || len(r.Alergias) != 0 {
		t.Errorf("esperaba alergias como lista vacía, obtuve %v", r.Alergias)
	}
	if !r.ActualizadoEn.Equal(ahora) {
		t.Errorf("esperaba fecha %v, obtuve %v", ahora, r.ActualizadoEn)
	}
}

func TestActualizarRegistroMantieneCamposNoEnviados(t *testing.T) {
	original := InformacionSalud{
		ID:                   "1",
		NombrePaciente:       "María Pérez",
		Diagnosticos:         []string{"hipertensión"},
		Alergias:             []string{"penicilina"},
		EnfermedadesCronicas: []string{"diabetes tipo 2"},
		AntecedentesMedicos:  []string{"cirugía de cadera 2019"},
		ActualizadoEn:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Solo se envía una actualización de alergias, el resto no debería borrarse.
	entrada := entradaInformacionSalud{
		Alergias: []string{"penicilina", "aspirina"},
	}
	ahora := time.Date(2026, 7, 12, 11, 0, 0, 0, time.UTC)

	actualizado := actualizarRegistro(original, entrada, ahora)

	if actualizado.NombrePaciente != "María Pérez" {
		t.Errorf("el nombre no debería cambiar, obtuve %s", actualizado.NombrePaciente)
	}
	if len(actualizado.Diagnosticos) != 1 || actualizado.Diagnosticos[0] != "hipertensión" {
		t.Errorf("los diagnósticos no deberían cambiar, obtuve %v", actualizado.Diagnosticos)
	}
	if len(actualizado.Alergias) != 2 {
		t.Errorf("esperaba 2 alergias tras la actualización, obtuve %v", actualizado.Alergias)
	}
	if !actualizado.ActualizadoEn.Equal(ahora) {
		t.Errorf("esperaba que se actualizara la fecha a %v, obtuve %v", ahora, actualizado.ActualizadoEn)
	}
}

func TestNormalizarConvierteNilEnListaVacia(t *testing.T) {
	resultado := normalizar(nil)
	if resultado == nil {
		t.Fatal("normalizar no debería devolver nil")
	}
	if len(resultado) != 0 {
		t.Errorf("esperaba lista vacía, obtuve %v", resultado)
	}
}
