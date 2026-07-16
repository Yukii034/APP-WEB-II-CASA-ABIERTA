package repository

import (
	"strconv"
	"sync"
	"testing"

	"cuidabien/informacion-salud/internal/model"
)

func TestMemoriaRepositorySiguienteIDDevuelveSecuencia(t *testing.T) {
	repo := NuevaMemoriaRepository()

	for i := 1; i <= 3; i++ {
		esperado := strconv.Itoa(i)
		obtenido := repo.SiguienteID()
		if obtenido != esperado {
			t.Errorf("esperaba id %q en la llamada %d, obtuve %q", esperado, i, obtenido)
		}
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

func TestMemoriaRepositoryGuardarSobrescribeRegistroExistente(t *testing.T) {
	repo := NuevaMemoriaRepository()
	repo.Guardar(model.InformacionSalud{ID: "1", NombrePaciente: "Ana"})

	// Guardar con el mismo ID debe actualizar, no duplicar.
	repo.Guardar(model.InformacionSalud{ID: "1", NombrePaciente: "Ana María"})

	obtenido, ok := repo.Obtener("1")
	if !ok {
		t.Fatal("esperaba encontrar el registro")
	}
	if obtenido.NombrePaciente != "Ana María" {
		t.Errorf("esperaba que se sobrescribiera a 'Ana María', obtuve %s", obtenido.NombrePaciente)
	}
	if len(repo.Listar()) != 1 {
		t.Errorf("esperaba que siguiera habiendo 1 solo registro, obtuve %d", len(repo.Listar()))
	}
}

func TestMemoriaRepositoryObtenerIDInexistente(t *testing.T) {
	repo := NuevaMemoriaRepository()

	_, ok := repo.Obtener("no-existe")
	if ok {
		t.Fatal("esperaba ok=false para un id que no existe")
	}
}

func TestMemoriaRepositoryListarVacioDevuelveListaVaciaNoNil(t *testing.T) {
	repo := NuevaMemoriaRepository()

	lista := repo.Listar()

	if lista == nil {
		t.Fatal("esperaba una lista vacía, no nil (afecta la serialización a JSON)")
	}
	if len(lista) != 0 {
		t.Errorf("esperaba 0 registros, obtuve %d", len(lista))
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

func TestMemoriaRepositoryConcurrencia(t *testing.T) {
	// Corre este test con -race para que tenga sentido:
	//   go test ./internal/repository/... -race
	repo := NuevaMemoriaRepository()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := repo.SiguienteID()
			repo.Guardar(model.InformacionSalud{ID: id, NombrePaciente: "Paciente " + id})
		}()
	}
	wg.Wait()

	if len(repo.Listar()) != 50 {
		t.Errorf("esperaba 50 registros tras el uso concurrente, obtuve %d", len(repo.Listar()))
	}
}

func TestSembrarAgregaExactamenteLosRegistrosEsperados(t *testing.T) {
	repo := NuevaMemoriaRepository()

	Sembrar(repo)

	lista := repo.Listar()
	if len(lista) != 3 {
		t.Fatalf("esperaba exactamente 3 registros sembrados, obtuve %d", len(lista))
	}

	nombres := map[string]bool{}
	for _, reg := range lista {
		if reg.ID == "" {
			t.Errorf("esperaba que cada registro sembrado tenga id, obtuve uno vacío: %+v", reg)
		}
		nombres[reg.NombrePaciente] = true
	}

	esperados := []string{"María Pérez", "José Ramírez", "Carmen Torres"}
	for _, nombre := range esperados {
		if !nombres[nombre] {
			t.Errorf("esperaba encontrar a %q entre los sembrados, no apareció", nombre)
		}
	}
}

func TestObtenerDevuelveCopiaIndependiente(t *testing.T) {
	repo := NuevaMemoriaRepository()
	repo.Guardar(model.InformacionSalud{ID: "1", Alergias: []string{"penicilina"}})

	obtenido, _ := repo.Obtener("1")
	obtenido.Alergias[0] = "modificado desde afuera"

	otraLectura, _ := repo.Obtener("1")
	if otraLectura.Alergias[0] != "penicilina" {
		t.Errorf("mutar el slice devuelto no debería afectar el dato interno, obtuve %v", otraLectura.Alergias)
	}
}
