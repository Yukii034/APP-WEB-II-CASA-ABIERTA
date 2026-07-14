package service

import (
	"testing"
	"time"

	"cuidabien/informacion-salud/internal/model"
)

func TestNuevoRegistroNormalizaListasNulas(t *testing.T) {
	ahora := time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC)
	entrada := model.EntradaInformacionSalud{
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
	original := model.InformacionSalud{
		ID:                   "1",
		NombrePaciente:       "María Pérez",
		Diagnosticos:         []string{"hipertensión"},
		Alergias:             []string{"penicilina"},
		EnfermedadesCronicas: []string{"diabetes tipo 2"},
		AntecedentesMedicos:  []string{"cirugía de cadera 2019"},
		ActualizadoEn:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Solo se envía una actualización de alergias, el resto no debería borrarse.
	entrada := model.EntradaInformacionSalud{
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

func TestServiceCrearYObtener(t *testing.T) {
	svc := NuevoInformacionSaludService(&repositorioFake{registros: map[string]model.InformacionSalud{}})

	creado := svc.Crear(model.EntradaInformacionSalud{NombrePaciente: "Juan"})
	if creado.ID == "" {
		t.Fatal("esperaba que Crear asigne un id")
	}

	obtenido, ok := svc.Obtener(creado.ID)
	if !ok {
		t.Fatal("esperaba encontrar el registro recién creado")
	}
	if obtenido.NombrePaciente != "Juan" {
		t.Errorf("esperaba nombre 'Juan', obtuve %s", obtenido.NombrePaciente)
	}
}

func TestServiceActualizarIDInexistente(t *testing.T) {
	svc := NuevoInformacionSaludService(&repositorioFake{registros: map[string]model.InformacionSalud{}})

	_, ok := svc.Actualizar("no-existe", model.EntradaInformacionSalud{})
	if ok {
		t.Fatal("esperaba ok=false al actualizar un id que no existe")
	}
}

// repositorioFake es un repository.Repository mínimo para probar el
// service sin depender de la implementación en memoria real.
type repositorioFake struct {
	registros   map[string]model.InformacionSalud
	siguienteID int
}

func (r *repositorioFake) Listar() []model.InformacionSalud {
	resultado := make([]model.InformacionSalud, 0, len(r.registros))
	for _, reg := range r.registros {
		resultado = append(resultado, reg)
	}
	return resultado
}

func (r *repositorioFake) Obtener(id string) (model.InformacionSalud, bool) {
	reg, ok := r.registros[id]
	return reg, ok
}

func (r *repositorioFake) Guardar(registro model.InformacionSalud) {
	r.registros[registro.ID] = registro
}

func (r *repositorioFake) SiguienteID() string {
	r.siguienteID++
	return "fake-" + string(rune('0'+r.siguienteID))
}
