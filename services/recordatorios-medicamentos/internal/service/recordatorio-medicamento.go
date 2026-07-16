package service

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"cuidabien/recordatorios-medicamentos/internal/model"
	"cuidabien/recordatorios-medicamentos/internal/repository"
)

// ErrDatosInvalidos permite al handler distinguir
// los errores producidos por validaciones.
var ErrDatosInvalidos = errors.New(
	"datos del recordatorio inválidos",
)

// RecordatorioMedicamentoService contiene la lógica de negocio.
type RecordatorioMedicamentoService struct {
	repo   repository.Repository
	salida io.Writer
}

func NuevoRecordatorioMedicamentoService(
	repo repository.Repository,
	salida io.Writer,
) *RecordatorioMedicamentoService {
	if salida == nil {
		salida = io.Discard
	}

	return &RecordatorioMedicamentoService{
		repo:   repo,
		salida: salida,
	}
}

// Listar devuelve todos los recordatorios.
func (s *RecordatorioMedicamentoService) Listar() []model.RecordatorioMedicamento {
	return s.repo.Listar()
}

// Crear registra un nuevo recordatorio después
// de validar los datos recibidos.
func (s *RecordatorioMedicamentoService) Crear(
	entrada model.EntradaRecordatorioMedicamento,
) (model.RecordatorioMedicamento, error) {
	entrada = normalizarEntrada(entrada)

	if err := validarCreacion(entrada); err != nil {
		return model.RecordatorioMedicamento{}, err
	}

	activo := true

	if entrada.Activo != nil {
		activo = *entrada.Activo
	}

	ahora := time.Now()

	nuevo := model.RecordatorioMedicamento{
		ID:             s.repo.SiguienteID(),
		AdultoMayorID:  entrada.AdultoMayorID,
		NombrePaciente: entrada.NombrePaciente,
		Medicamento:    entrada.Medicamento,
		Dosis:          entrada.Dosis,
		Hora:           entrada.Hora,
		Frecuencia:     entrada.Frecuencia,
		Activo:         activo,
		CreadoEn:       ahora,
		ActualizadoEn:  ahora,
	}

	s.repo.Guardar(nuevo)

	return nuevo, nil
}

// Obtener busca un recordatorio por id.
func (s *RecordatorioMedicamentoService) Obtener(
	id string,
) (model.RecordatorioMedicamento, bool) {
	return s.repo.Obtener(id)
}

// Actualizar aplica una actualización parcial.
// Los campos vacíos conservan su valor anterior.
func (s *RecordatorioMedicamentoService) Actualizar(
	id string,
	entrada model.EntradaRecordatorioMedicamento,
) (model.RecordatorioMedicamento, bool, error) {
	existente, ok := s.repo.Obtener(id)

	if !ok {
		return model.RecordatorioMedicamento{}, false, nil
	}

	entrada = normalizarEntrada(entrada)

	actualizado := actualizarRegistro(
		existente,
		entrada,
		time.Now(),
	)

	if err := validarRegistro(actualizado); err != nil {
		return model.RecordatorioMedicamento{}, true, err
	}

	s.repo.Guardar(actualizado)

	return actualizado, true, nil
}

// Eliminar borra un recordatorio por id.
func (s *RecordatorioMedicamentoService) Eliminar(
	id string,
) bool {
	return s.repo.Eliminar(id)
}

// CambiarEstado activa o desactiva un recordatorio.
func (s *RecordatorioMedicamentoService) CambiarEstado(
	id string,
	activo bool,
) (model.RecordatorioMedicamento, bool) {
	existente, ok := s.repo.Obtener(id)

	if !ok {
		return model.RecordatorioMedicamento{}, false
	}

	existente.Activo = activo
	existente.ActualizadoEn = time.Now()

	s.repo.Guardar(existente)

	return existente, true
}

// VerificarHora busca recordatorios activos y
// genera alertas simuladas en la consola.
func (s *RecordatorioMedicamentoService) VerificarHora(
	hora string,
) (model.ResultadoVerificacion, error) {
	hora = strings.TrimSpace(hora)

	if !horaValida(hora) {
		return model.ResultadoVerificacion{},
			fmt.Errorf(
				"%w: la hora debe usar el formato HH:MM",
				ErrDatosInvalidos,
			)
	}

	registros := s.repo.BuscarActivosPorHora(hora)

	alertas := make(
		[]model.AlertaMedicamento,
		0,
		len(registros),
	)

	for _, registro := range registros {
		mensaje := fmt.Sprintf(
			"%s, es hora de tomar %s. Dosis: %s",
			registro.NombrePaciente,
			registro.Medicamento,
			registro.Dosis,
		)

		fmt.Fprintf(
			s.salida,
			"ALERTA | Paciente: %s | Medicamento: %s | "+
				"Dosis: %s | Hora: %s\n",
			registro.NombrePaciente,
			registro.Medicamento,
			registro.Dosis,
			registro.Hora,
		)

		alertas = append(
			alertas,
			model.AlertaMedicamento{
				RecordatorioID: registro.ID,
				AdultoMayorID:  registro.AdultoMayorID,
				NombrePaciente: registro.NombrePaciente,
				Medicamento:    registro.Medicamento,
				Dosis:          registro.Dosis,
				Hora:           registro.Hora,
				Mensaje:        mensaje,
			},
		)
	}

	return model.ResultadoVerificacion{
		Hora:     hora,
		Cantidad: len(alertas),
		Alertas:  alertas,
	}, nil
}

func validarCreacion(
	entrada model.EntradaRecordatorioMedicamento,
) error {
	if entrada.AdultoMayorID == "" ||
		entrada.NombrePaciente == "" ||
		entrada.Medicamento == "" ||
		entrada.Dosis == "" ||
		entrada.Frecuencia == "" {
		return fmt.Errorf(
			"%w: adulto_mayor_id, nombre_paciente, "+
				"medicamento, dosis y frecuencia son obligatorios",
			ErrDatosInvalidos,
		)
	}

	if !horaValida(entrada.Hora) {
		return fmt.Errorf(
			"%w: la hora debe usar el formato HH:MM",
			ErrDatosInvalidos,
		)
	}

	return nil
}

func validarRegistro(
	registro model.RecordatorioMedicamento,
) error {
	entrada := model.EntradaRecordatorioMedicamento{
		AdultoMayorID:  registro.AdultoMayorID,
		NombrePaciente: registro.NombrePaciente,
		Medicamento:    registro.Medicamento,
		Dosis:          registro.Dosis,
		Hora:           registro.Hora,
		Frecuencia:     registro.Frecuencia,
	}

	return validarCreacion(entrada)
}

func actualizarRegistro(
	existente model.RecordatorioMedicamento,
	entrada model.EntradaRecordatorioMedicamento,
	ahora time.Time,
) model.RecordatorioMedicamento {
	if entrada.AdultoMayorID != "" {
		existente.AdultoMayorID = entrada.AdultoMayorID
	}

	if entrada.NombrePaciente != "" {
		existente.NombrePaciente = entrada.NombrePaciente
	}

	if entrada.Medicamento != "" {
		existente.Medicamento = entrada.Medicamento
	}

	if entrada.Dosis != "" {
		existente.Dosis = entrada.Dosis
	}

	if entrada.Hora != "" {
		existente.Hora = entrada.Hora
	}

	if entrada.Frecuencia != "" {
		existente.Frecuencia = entrada.Frecuencia
	}

	if entrada.Activo != nil {
		existente.Activo = *entrada.Activo
	}

	existente.ActualizadoEn = ahora

	return existente
}

func normalizarEntrada(
	entrada model.EntradaRecordatorioMedicamento,
) model.EntradaRecordatorioMedicamento {
	entrada.AdultoMayorID =
		strings.TrimSpace(entrada.AdultoMayorID)

	entrada.NombrePaciente =
		strings.TrimSpace(entrada.NombrePaciente)

	entrada.Medicamento =
		strings.TrimSpace(entrada.Medicamento)

	entrada.Dosis =
		strings.TrimSpace(entrada.Dosis)

	entrada.Hora =
		strings.TrimSpace(entrada.Hora)

	entrada.Frecuencia =
		strings.TrimSpace(entrada.Frecuencia)

	return entrada
}

func horaValida(hora string) bool {
	valor, err := time.Parse("15:04", hora)

	return err == nil &&
		valor.Format("15:04") == hora
}
