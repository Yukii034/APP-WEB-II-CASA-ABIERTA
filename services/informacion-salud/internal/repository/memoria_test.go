package repository

import (
	"testing"

	"cuidabien/informacion-salud/internal/model"
)

func TestMemoriaRepositorySiguienteIDIncrementa(t *testing.T) {
	repo := NuevaMemoriaRepository()

	primero := repo.SiguienteID()
	segundo := repo.SiguienteID()

	if primero == segundo {
		t.Errorf("esperaba ids distintos, obtuve %s y %s", primero, segundo)
	}
}

func TestMemoriaRepositoryGuardarYObtener(t *testing.T) {
	repo := NuevaMemoriaRepository()
	registro := model.InformacionSalud{ID: "1", NombrePaciente: "Ana"}

	repo.Guardar(registro)

	obtenido, ok := repo.Obtener("1")
	if !ok {
		t.Fatal("esperaba encontrar el registro guardado")
	}
	if obtenido.NombrePaciente != "Ana" {
		t.Errorf("esperaba nombre 'Ana', obtuve %s", obtenido.NombrePaciente)
	}
}

func TestMemoriaRepositoryObtenerIDInexistente(t *testing.T) {
	repo := NuevaMemoriaRepository()

	_, ok := repo.Obtener("no-existe")
	if ok {
		t.Fatal("esperaba ok=false para un id que no existe")
	}
}

func TestMemoriaRepositoryListar(t *testing.T) {
	repo := NuevaMemoriaRepository()
	repo.Guardar(model.InformacionSalud{ID: "1"})
	repo.Guardar(model.InformacionSalud{ID: "2"})

	lista := repo.Listar()

	if len(lista) != 2 {
		t.Errorf("esperaba 2 registros, obtuve %d", len(lista))
	}
}

func TestSembrarAgregaRegistrosConID(t *testing.T) {
	repo := NuevaMemoriaRepository()

	Sembrar(repo)

	lista := repo.Listar()
	if len(lista) == 0 {
		t.Fatal("esperaba que Sembrar agregue al menos un registro")
	}
	for _, reg := range lista {
		if reg.ID == "" {
			t.Errorf("esperaba que cada registro sembrado tenga id, obtuve uno vacío: %+v", reg)
		}
		if reg.NombrePaciente == "" {
			t.Errorf("esperaba que cada registro sembrado tenga nombre, obtuve uno vacío: %+v", reg)
		}
	}
}
