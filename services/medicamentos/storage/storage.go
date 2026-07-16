package storage

import (
	"cuidabien/medicamentos/models"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

type Store struct {
	Medicamentos []models.Medicamento
	Tomas        []models.Toma
	Alertas      []models.Alerta
	Interacciones []models.Interaccion
	Pacientes    []models.Paciente
	NextMedID    int
	NextTomaID   int
	NextAlertaID int
	NextInterID  int
}

func NewStore() *Store {
	s := &Store{
		Pacientes: []models.Paciente{
			{ID: "P001", Nombre: "Maria Garcia"},
			{ID: "P002", Nombre: "Juan Lopez"},
			{ID: "P003", Nombre: "Ana Martinez"},
		},
		Interacciones: []models.Interaccion{
			{ID: "I001", MedicamentoA: "Ibuprofeno", MedicamentoB: "Aspirina", Gravedad: "severa", Descripcion: "Riesgo aumentado de sangrado gastrointestinal"},
			{ID: "I002", MedicamentoA: "Losartan", MedicamentoB: "Ibuprofeno", Gravedad: "moderada", Descripcion: "Puede reducir el efecto antihipertensivo"},
			{ID: "I003", MedicamentoA: "Metformina", MedicamentoB: "Alcohol", Gravedad: "severa", Descripcion: "Riesgo de acidosis lactica"},
		},
		NextMedID:    1,
		NextTomaID:   1,
		NextAlertaID: 1,
		NextInterID:  4,
	}
	s.seedMedicamentos()
	return s
}

func (s *Store) seedMedicamentos() {
	hoy := time.Now().Format("2006-01-02")
	semanaAntes := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	vencido := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	porVencer := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	s.Medicamentos = []models.Medicamento{
		{ID: "M001", PacienteID: "P001", Nombre: "Losartan", Dosis: "50mg", Frecuencia: "1 vez al dia", Horarios: []string{"08:00"}, FechaInicio: semanaAntes, Estado: "activo", Notas: "Tomar con agua"},
		{ID: "M002", PacienteID: "P001", Nombre: "Metformina", Dosis: "850mg", Frecuencia: "2 veces al dia", Horarios: []string{"08:00", "20:00"}, FechaInicio: semanaAntes, Estado: "activo"},
		{ID: "M003", PacienteID: "P002", Nombre: "Aspirina", Dosis: "100mg", Frecuencia: "1 vez al dia", Horarios: []string{"09:00"}, FechaInicio: semanaAntes, Estado: "activo"},
		{ID: "M004", PacienteID: "P002", Nombre: "Ibuprofeno", Dosis: "400mg", Frecuencia: "cada 8 horas", Horarios: []string{"08:00", "16:00"}, FechaInicio: semanaAntes, FechaFin: vencido, Estado: "activo"},
		{ID: "M005", PacienteID: "P003", Nombre: "Paracetamol", Dosis: "500mg", Frecuencia: "cada 6 horas", Horarios: []string{"08:00", "14:00", "20:00"}, FechaInicio: semanaAntes, FechaFin: porVencer, Estado: "activo"},
		{ID: "M006", PacienteID: "P003", Nombre: "Omeprazol", Dosis: "20mg", Frecuencia: "1 vez al dia", Horarios: []string{"07:00"}, FechaInicio: hoy, Estado: "suspendido", Notas: "Suspender por indicacion medica"},
	}
	s.NextMedID = 7

	s.Tomas = []models.Toma{
		{ID: "T001", MedicamentoID: "M001", PacienteID: "P001", FechaHoraProgramada: semanaAntes + " 08:00", Estado: "cumplida", FechaHoraReal: semanaAntes + " 08:15"},
		{ID: "T002", MedicamentoID: "M001", PacienteID: "P001", FechaHoraProgramada: time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 08:00", Estado: "cumplida", FechaHoraReal: time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 08:05"},
		{ID: "T003", MedicamentoID: "M001", PacienteID: "P001", FechaHoraProgramada: hoy + " 08:00", Estado: "cumplida", FechaHoraReal: hoy + " 08:10"},
		{ID: "T004", MedicamentoID: "M002", PacienteID: "P001", FechaHoraProgramada: time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 20:00", Estado: "no_cumplida"},
		{ID: "T005", MedicamentoID: "M003", PacienteID: "P002", FechaHoraProgramada: hoy + " 09:00", Estado: "cumplida", FechaHoraReal: hoy + " 09:20"},
	}
	s.NextTomaID = 6
}

func (s *Store) GenerateMedID() string {
	id := fmt.Sprintf("M%03d", s.NextMedID)
	s.NextMedID++
	return id
}

func (s *Store) GenerateTomaID() string {
	id := fmt.Sprintf("T%03d", s.NextTomaID)
	s.NextTomaID++
	return id
}

func (s *Store) GenerateAlertaID() string {
	id := fmt.Sprintf("A%03d", s.NextAlertaID)
	s.NextAlertaID++
	return id
}

func (s *Store) FindMedicamentoByID(id string) *models.Medicamento {
	for i := range s.Medicamentos {
		if s.Medicamentos[i].ID == id {
			return &s.Medicamentos[i]
		}
	}
	return nil
}

func (s *Store) FindPacienteByID(id string) *models.Paciente {
	for i := range s.Pacientes {
		if s.Pacientes[i].ID == id {
			return &s.Pacientes[i]
		}
	}
	return nil
}

func (s *Store) FindAlertaByID(id string) *models.Alerta {
	for i := range s.Alertas {
		if s.Alertas[i].ID == id {
			return &s.Alertas[i]
		}
	}
	return nil
}

func Sanitizar(str string) string {
	str = strings.TrimSpace(str)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(str, " ")
}

func ValidarHorario(h string) bool {
	re := regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)
	return re.MatchString(h)
}

func (s *Store) MedicamentosPorPaciente(pacienteID string) []models.Medicamento {
	var result []models.Medicamento
	for _, m := range s.Medicamentos {
		if m.PacienteID == pacienteID {
			result = append(result, m)
		}
	}
	return result
}

func (s *Store) MedicamentosActivosPorPaciente(pacienteID string) []models.Medicamento {
	var result []models.Medicamento
	for _, m := range s.Medicamentos {
		if m.PacienteID == pacienteID && m.Estado == "activo" {
			result = append(result, m)
		}
	}
	return result
}

func (s *Store) TomasPorMedicamento(medicamentoID string) []models.Toma {
	var result []models.Toma
	for _, t := range s.Tomas {
		if t.MedicamentoID == medicamentoID {
			result = append(result, t)
		}
	}
	return result
}

func (s *Store) TomasPorPaciente(pacienteID string) []models.Toma {
	var result []models.Toma
	for _, t := range s.Tomas {
		if t.PacienteID == pacienteID {
			result = append(result, t)
		}
	}
	return result
}

func (s *Store) AlertasPorPaciente(pacienteID string) []models.Alerta {
	var result []models.Alerta
	for _, a := range s.Alertas {
		if a.PacienteID == pacienteID {
			result = append(result, a)
		}
	}
	return result
}

func (s *Store) BuscarInteraccion(nombreA, nombreB string) *models.Interaccion {
	aLower := strings.ToLower(nombreA)
	bLower := strings.ToLower(nombreB)
	for i := range s.Interacciones {
		intA := strings.ToLower(s.Interacciones[i].MedicamentoA)
		intB := strings.ToLower(s.Interacciones[i].MedicamentoB)
		if (aLower == intA && bLower == intB) || (aLower == intB && bLower == intA) {
			return &s.Interacciones[i]
		}
	}
	return nil
}

func (s *Store) CalcularAdherencia(pacienteID string) models.Adherencia {
	tomas := s.TomasPorPaciente(pacienteID)
	cumplidas := 0
	noCumplidas := 0
	for _, t := range tomas {
		if t.Estado == "cumplida" {
			cumplidas++
		} else if t.Estado == "no_cumplida" {
			noCumplidas++
		}
	}
	total := cumplidas + noCumplidas
	porcentaje := 0.0
	if total > 0 {
		porcentaje = math.Round(float64(cumplidas)/float64(total)*100*100) / 100
	}
	return models.Adherencia{
		PacienteID:       pacienteID,
		TotalTomas:       total,
		TomasCumplidas:   cumplidas,
		TomasNoCumplidas: noCumplidas,
		Porcentaje:       porcentaje,
	}
}

func (s *Store) GenerarAlertasVencimiento() {
	hoy := time.Now()
	limite := hoy.AddDate(0, 0, 7)
	for _, m := range s.Medicamentos {
		if m.Estado != "activo" || m.FechaFin == "" {
			continue
		}
		vencimiento, err := time.Parse("2006-01-02", m.FechaFin)
		if err != nil {
			continue
		}
		if vencimiento.Before(hoy) {
			s.crearAlertaSiNoExiste(m.PacienteID, m.ID, "vencido",
				fmt.Sprintf("El medicamento '%s' esta vencido desde %s", m.Nombre, m.FechaFin))
		} else if vencimiento.Before(limite) {
			s.crearAlertaSiNoExiste(m.PacienteID, m.ID, "por_vencer",
				fmt.Sprintf("El medicamento '%s' vence el %s", m.Nombre, m.FechaFin))
		}
	}
}

func (s *Store) crearAlertaSiNoExiste(pacienteID, medicamentoID, tipo, mensaje string) {
	for _, a := range s.Alertas {
		if a.PacienteID == pacienteID && a.MedicamentoID == medicamentoID && a.Tipo == tipo && !a.Leida {
			return
		}
	}
	s.Alertas = append(s.Alertas, models.Alerta{
		ID:            s.GenerateAlertaID(),
		PacienteID:    pacienteID,
		MedicamentoID: medicamentoID,
		Tipo:          tipo,
		Mensaje:       mensaje,
		FechaCreacion: time.Now().Format(time.RFC3339),
		Leida:         false,
	})
}

func (s *Store) VerificarInteraccionesActivas(pacienteID string) []models.Interaccion {
	activos := s.MedicamentosActivosPorPaciente(pacienteID)
	var encontradas []models.Interaccion
	for i := 0; i < len(activos); i++ {
		for j := i + 1; j < len(activos); j++ {
			if inter := s.BuscarInteraccion(activos[i].Nombre, activos[j].Nombre); inter != nil {
				encontradas = append(encontradas, *inter)
			}
		}
	}
	return encontradas
}
