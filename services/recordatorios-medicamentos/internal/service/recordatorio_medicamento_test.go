package service

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"cuidabien/recordatorios-medicamentos/internal/model"
	"cuidabien/recordatorios-medicamentos/internal/repository"
)

func TestCrearRecordatorioValido(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()

	svc := NuevoRecordatorioMedicamentoService(
		repo,
		nil,
	)

	creado, err := svc.Crear(
		model.EntradaRecordatorioMedicamento{
			AdultoMayorID:  "AM-001",
			NombrePaciente: "María Pérez",
			Medicamento:    "Losartán",
			Dosis:          "1 tableta",
			Hora:           "08:00",
			Frecuencia:     "diaria",
		},
	)

	if err != nil {
		t.Fatalf(
			"no se esperaba error: %v",
			err,
		)
	}

	if creado.ID == "" || !creado.Activo {
		t.Fatalf(
			"recordatorio creado incorrectamente: %+v",
			creado,
		)
	}
}

func TestCrearRechazaHoraInvalida(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()

	svc := NuevoRecordatorioMedicamentoService(
		repo,
		nil,
	)

	_, err := svc.Crear(
		model.EntradaRecordatorioMedicamento{
			AdultoMayorID:  "AM-001",
			NombrePaciente: "María Pérez",
			Medicamento:    "Losartán",
			Dosis:          "1 tableta",
			Hora:           "25:90",
			Frecuencia:     "diaria",
		},
	)

	if !errors.Is(err, ErrDatosInvalidos) {
		t.Fatalf(
			"esperaba ErrDatosInvalidos, obtuve %v",
			err,
		)
	}
}

func TestActualizarMantieneCamposNoEnviados(
	t *testing.T,
) {
	repo := repository.NuevaMemoriaRepository()

	svc := NuevoRecordatorioMedicamentoService(
		repo,
		nil,
	)

	creado, _ := svc.Crear(
		model.EntradaRecordatorioMedicamento{
			AdultoMayorID:  "AM-001",
			NombrePaciente: "María Pérez",
			Medicamento:    "Losartán",
			Dosis:          "1 tableta",
			Hora:           "08:00",
			Frecuencia:     "diaria",
		},
	)

	actualizado, ok, err := svc.Actualizar(
		creado.ID,
		model.EntradaRecordatorioMedicamento{
			Dosis: "2 tabletas",
		},
	)

	if err != nil || !ok {
		t.Fatalf(
			"no se pudo actualizar: ok=%v err=%v",
			ok,
			err,
		)
	}

	if actualizado.Medicamento != "Losartán" ||
		actualizado.Dosis != "2 tabletas" {
		t.Fatalf(
			"actualización incorrecta: %+v",
			actualizado,
		)
	}
}

func TestVerificarHoraGeneraAlerta(t *testing.T) {
	repo := repository.NuevaMemoriaRepository()

	var salida bytes.Buffer

	svc := NuevoRecordatorioMedicamentoService(
		repo,
		&salida,
	)

	_, _ = svc.Crear(
		model.EntradaRecordatorioMedicamento{
			AdultoMayorID:  "AM-001",
			NombrePaciente: "María Pérez",
			Medicamento:    "Metformina",
			Dosis:          "500 mg",
			Hora:           "14:30",
			Frecuencia:     "diaria",
		},
	)

	resultado, err := svc.VerificarHora("14:30")

	if err != nil {
		t.Fatalf(
			"no se esperaba error: %v",
			err,
		)
	}

	if resultado.Cantidad != 1 {
		t.Fatalf(
			"esperaba 1 alerta, obtuve %d",
			resultado.Cantidad,
		)
	}

	if !strings.Contains(
		salida.String(),
		"ALERTA",
	) {
		t.Fatalf(
			"no se imprimió la alerta: %q",
			salida.String(),
		)
	}
}

func TestActualizarRegistroConservaCreadoEn(
	t *testing.T,
) {
	creadoEn := time.Date(
		2026,
		7,
		15,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	existente := model.RecordatorioMedicamento{
		ID:             "1",
		AdultoMayorID:  "AM-001",
		NombrePaciente: "María",
		Medicamento:    "Losartán",
		Dosis:          "1 tableta",
		Hora:           "08:00",
		Frecuencia:     "diaria",
		Activo:         true,
		CreadoEn:       creadoEn,
	}

	actualizado := actualizarRegistro(
		existente,
		model.EntradaRecordatorioMedicamento{
			Dosis: "2 tabletas",
		},
		creadoEn.Add(time.Hour),
	)

	if !actualizado.CreadoEn.Equal(creadoEn) {
		t.Fatal("CreadoEn no debía cambiar")
	}
}
