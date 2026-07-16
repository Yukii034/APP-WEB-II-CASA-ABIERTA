package service

import (
	"strconv"
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

func TestActualizarRegistro(t *testing.T) {
	original := model.InformacionSalud{
		ID:                   "1",
		NombrePaciente:       "María Pérez",
		Diagnosticos:         []string{"hipertensión"},
		Alergias:             []string{"penicilina"},
		EnfermedadesCronicas: []string{"diabetes tipo 2"},
		AntecedentesMedicos:  []string{"cirugía de cadera 2019"},
		ActualizadoEn:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	ahora := time.Date(2026, 7, 12, 11, 0, 0, 0, time.UTC)

	t.Run("mantiene los campos que no se envían", func(t *testing.T) {
		// Solo se envía una actualización de alergias, el resto no debería borrarse.
		entrada := model.EntradaInformacionSalud{
			Alergias: []string{"penicilina", "aspirina"},
		}

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
	})

	t.Run("una lista vacía explícita sí borra el campo", func(t *testing.T) {
		// A diferencia de no enviar el campo (nil), enviar [] es una
		// instrucción explícita de "ya no tiene diagnósticos".
		entrada := model.EntradaInformacionSalud{
			Diagnosticos: []string{},
		}

		actualizado := actualizarRegistro(original, entrada, ahora)

		if len(actualizado.Diagnosticos) != 0 {
			t.Errorf("esperaba que una lista vacía explícita borre el campo, obtuve %v", actualizado.Diagnosticos)
		}
		// Los campos no tocados por esta petición siguen intactos.
		if len(actualizado.Alergias) != 1 || actualizado.Alergias[0] != "penicilina" {
			t.Errorf("las alergias no deberían cambiar, obtuve %v", actualizado.Alergias)
		}
	})

	t.Run("no envía nada y no cambia ningún campo, solo la fecha", func(t *testing.T) {
		entrada := model.EntradaInformacionSalud{}

		actualizado := actualizarRegistro(original, entrada, ahora)

		if actualizado.NombrePaciente != original.NombrePaciente {
			t.Errorf("el nombre no debería cambiar, obtuve %s", actualizado.NombrePaciente)
		}
		if len(actualizado.Diagnosticos) != len(original.Diagnosticos) {
			t.Errorf("los diagnósticos no deberían cambiar, obtuve %v", actualizado.Diagnosticos)
		}
		if !actualizado.ActualizadoEn.Equal(ahora) {
			t.Errorf("esperaba que igual se actualice la fecha a %v, obtuve %v", ahora, actualizado.ActualizadoEn)
		}
	})

	t.Run("actualiza los demás campos enviados", func(t *testing.T) {
		entrada := model.EntradaInformacionSalud{
			NombrePaciente:       "María González",
			EnfermedadesCronicas: []string{"asma"},
			AntecedentesMedicos:  []string{"cirugía de rodilla 2020"},
		}

		actualizado := actualizarRegistro(original, entrada, ahora)

		if actualizado.NombrePaciente != "María González" {
			t.Errorf("esperaba actualizar el nombre, obtuve %s", actualizado.NombrePaciente)
		}
		if len(actualizado.EnfermedadesCronicas) != 1 || actualizado.EnfermedadesCronicas[0] != "asma" {
			t.Errorf("esperaba actualizar enfermedades crónicas, obtuve %v", actualizado.EnfermedadesCronicas)
		}
		if len(actualizado.AntecedentesMedicos) != 1 || actualizado.AntecedentesMedicos[0] != "cirugía de rodilla 2020" {
			t.Errorf("esperaba actualizar antecedentes médicos, obtuve %v", actualizado.AntecedentesMedicos)
		}
	})
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

func TestNormalizarConservaListaConValores(t *testing.T) {
	resultado := normalizar([]string{"a", "b"})
	if len(resultado) != 2 {
		t.Errorf("esperaba conservar los 2 valores, obtuve %v", resultado)
	}
}

func TestServiceCrearYObtener(t *testing.T) {
	svc := NuevoInformacionSaludService(nuevoRepositorioFake())

	creado := svc.Crear(model.EntradaInformacionSalud{
		NombrePaciente: "Juan",
		Diagnosticos:   []string{"gripe"},
	})
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
	if len(obtenido.Diagnosticos) != 1 || obtenido.Diagnosticos[0] != "gripe" {
		t.Errorf("esperaba diagnóstico 'gripe', obtuve %v", obtenido.Diagnosticos)
	}
}

func TestServiceCrearAsignaIDsDistintos(t *testing.T) {
	svc := NuevoInformacionSaludService(nuevoRepositorioFake())

	a := svc.Crear(model.EntradaInformacionSalud{NombrePaciente: "Juan"})
	b := svc.Crear(model.EntradaInformacionSalud{NombrePaciente: "Ana"})

	if a.ID == b.ID {
		t.Errorf("esperaba ids distintos para registros distintos, ambos fueron %s", a.ID)
	}
}

func TestServiceListarDevuelveTodosLosRegistros(t *testing.T) {
	svc := NuevoInformacionSaludService(nuevoRepositorioFake())
	svc.Crear(model.EntradaInformacionSalud{NombrePaciente: "Juan"})
	svc.Crear(model.EntradaInformacionSalud{NombrePaciente: "Ana"})

	lista := svc.Listar()

	if len(lista) != 2 {
		t.Errorf("esperaba 2 registros, obtuve %d", len(lista))
	}
}

func TestServiceObtenerIDInexistente(t *testing.T) {
	svc := NuevoInformacionSaludService(nuevoRepositorioFake())

	_, ok := svc.Obtener("no-existe")
	if ok {
		t.Fatal("esperaba ok=false al consultar un id que no existe")
	}
}

func TestServiceActualizarIDInexistente(t *testing.T) {
	svc := NuevoInformacionSaludService(nuevoRepositorioFake())

	_, ok := svc.Actualizar("no-existe", model.EntradaInformacionSalud{})
	if ok {
		t.Fatal("esperaba ok=false al actualizar un id que no existe")
	}
}

func TestServiceActualizarPersisteElCambio(t *testing.T) {
	svc := NuevoInformacionSaludService(nuevoRepositorioFake())
	creado := svc.Crear(model.EntradaInformacionSalud{NombrePaciente: "Juan"})

	_, ok := svc.Actualizar(creado.ID, model.EntradaInformacionSalud{
		Alergias: []string{"penicilina"},
	})
	if !ok {
		t.Fatal("esperaba poder actualizar un registro existente")
	}

	obtenido, _ := svc.Obtener(creado.ID)
	if len(obtenido.Alergias) != 1 || obtenido.Alergias[0] != "penicilina" {
		t.Errorf("esperaba que la actualización quedara guardada, obtuve %v", obtenido.Alergias)
	}
}

// repositorioFake es un repository.Repository mínimo para probar el
// service sin depender de la implementación en memoria real.
type repositorioFake struct {
	registros   map[string]model.InformacionSalud
	siguienteID int
}

func nuevoRepositorioFake() *repositorioFake {
	return &repositorioFake{registros: map[string]model.InformacionSalud{}}
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
	return "fake-" + strconv.Itoa(r.siguienteID)
}
