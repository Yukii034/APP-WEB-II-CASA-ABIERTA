package repository

import (
	"testing"

	"cuidabien/recordatorios-medicamentos/internal/model"
)

func TestMemoriaRepositorySiguienteIDIncrementa(t *testing.T) {
	repo := NuevaMemoriaRepository()

	primero := repo.SiguienteID()
	segundo := repo.SiguienteID()

	if primero == segundo {
		t.Errorf(
			"esperaba ids distintos, obtuve %s y %s",
			primero,
			segundo,
		)
	}
}

func TestMemoriaRepositoryGuardarYObtener(t *testing.T) {
	repo := NuevaMemoriaRepository()

	registro := model.RecordatorioMedicamento{
		ID:          "1",
		Medicamento: "Losartán",
	}

	repo.Guardar(registro)

	obtenido, ok := repo.Obtener("1")

	if !ok {
		t.Fatal(
			"esperaba encontrar el recordatorio guardado",
		)
	}

	if obtenido.Medicamento != "Losartán" {
		t.Errorf(
			"esperaba Losartán, obtuve %s",
			obtenido.Medicamento,
		)
	}
}

func TestMemoriaRepositoryObtenerIDInexistente(
	t *testing.T,
) {
	repo := NuevaMemoriaRepository()

	_, ok := repo.Obtener("no-existe")

	if ok {
		t.Fatal(
			"esperaba ok=false para un id que no existe",
		)
	}
}

func TestMemoriaRepositoryListar(t *testing.T) {
	repo := NuevaMemoriaRepository()

	repo.Guardar(
		model.RecordatorioMedicamento{
			ID: "1",
		},
	)

	repo.Guardar(
		model.RecordatorioMedicamento{
			ID: "2",
		},
	)

	lista := repo.Listar()

	if len(lista) != 2 {
		t.Errorf(
			"esperaba 2 recordatorios, obtuve %d",
			len(lista),
		)
	}
}

func TestMemoriaRepositoryEliminar(t *testing.T) {
	repo := NuevaMemoriaRepository()

	repo.Guardar(
		model.RecordatorioMedicamento{
			ID: "1",
		},
	)

	if !repo.Eliminar("1") {
		t.Fatal(
			"esperaba eliminar el recordatorio",
		)
	}

	if _, ok := repo.Obtener("1"); ok {
		t.Fatal(
			"el recordatorio todavía existe",
		)
	}
}

func TestBuscarActivosPorHoraIgnoraInactivos(
	t *testing.T,
) {
	repo := NuevaMemoriaRepository()

	repo.Guardar(
		model.RecordatorioMedicamento{
			ID:     "1",
			Hora:   "08:00",
			Activo: true,
		},
	)

	repo.Guardar(
		model.RecordatorioMedicamento{
			ID:     "2",
			Hora:   "08:00",
			Activo: false,
		},
	)

	repo.Guardar(
		model.RecordatorioMedicamento{
			ID:     "3",
			Hora:   "14:30",
			Activo: true,
		},
	)

	resultado := repo.BuscarActivosPorHora(
		"08:00",
	)

	if len(resultado) != 1 {
		t.Errorf(
			"esperaba 1 recordatorio activo, obtuve %d",
			len(resultado),
		)
	}
}

func TestSembrarAgregaRegistrosConID(t *testing.T) {
	repo := NuevaMemoriaRepository()

	Sembrar(repo)

	lista := repo.Listar()

	if len(lista) == 0 {
		t.Fatal(
			"esperaba que Sembrar agregue al menos un recordatorio",
		)
	}

	for _, registro := range lista {
		if registro.ID == "" {
			t.Errorf(
				"esperaba un id, obtuve uno vacío: %+v",
				registro,
			)
		}

		if registro.NombrePaciente == "" {
			t.Errorf(
				"esperaba un nombre de paciente: %+v",
				registro,
			)
		}

		if registro.CreadoEn.IsZero() ||
			registro.ActualizadoEn.IsZero() {
			t.Errorf(
				"esperaba fechas válidas: %+v",
				registro,
			)
		}
	}
}
