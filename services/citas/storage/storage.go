package storage

import (
	"cuidabien/citas/models"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Store struct {
	Citas          []models.Cita
	Pacientes      []models.Paciente
	Doctores       []models.Doctor
	Historial      []models.EventoHistorial
	NextID         int
	TodayCreated   int
	TodayCancelled int
	TodayCompleted int
	TodayDate      string
}

func NewStore() *Store {
	return &Store{
		Pacientes: []models.Paciente{
			{ID: "P001", Nombre: "Maria Garcia", Telefono: "555-0101", ContactoEmergencia: "555-0102", Alergias: []string{"Penicilina"}},
			{ID: "P002", Nombre: "Juan Lopez", Telefono: "555-0201", ContactoEmergencia: "555-0202", Alergias: []string{}},
			{ID: "P003", Nombre: "Ana Martinez", Telefono: "555-0301", ContactoEmergencia: "555-0302", Alergias: []string{"Ibuprofeno", "Sulfa"}},
		},
		Doctores: []models.Doctor{
			{ID: "D001", Nombre: "Dr. Carlos Ruiz", Especialidad: "Medicina General"},
			{ID: "D002", Nombre: "Dra. Laura Fernandez", Especialidad: "Cardiologia"},
			{ID: "D003", Nombre: "Dr. Pedro Sanchez", Especialidad: "Gerontologia"},
		},
		NextID:    1,
		TodayDate: time.Now().Format("2006-01-02"),
	}
}

func (s *Store) GenerateID() string {
	id := fmt.Sprintf("C%03d", s.NextID)
	s.NextID++
	return id
}

func (s *Store) FindCitaByID(id string) *models.Cita {
	for i := range s.Citas {
		if s.Citas[i].ID == id {
			return &s.Citas[i]
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

func (s *Store) FindDoctorByID(id string) *models.Doctor {
	for i := range s.Doctores {
		if s.Doctores[i].ID == id {
			return &s.Doctores[i]
		}
	}
	return nil
}

func InformacionSaludIDPorPaciente(pacienteID string) string {
	mapeo := map[string]string{
		"P001": "1",
		"P002": "2",
		"P003": "3",
	}
	return mapeo[pacienteID]
}

func (s *Store) EsFechaPasada(fecha, hora string) bool {
	fechaHoraStr := fecha + " " + hora
	t, err := time.Parse("2006-01-02 15:04", fechaHoraStr)
	if err != nil {
		return true
	}
	return t.Before(time.Now())
}

func (s *Store) MedicoOcupado(doctorID, fecha, hora, excludeID string) bool {
	for _, c := range s.Citas {
		if c.DoctorID == doctorID && c.Fecha == fecha && c.Hora == hora && c.ID != excludeID {
			if c.Estado != "cancelada" {
				return true
			}
		}
	}
	return false
}

func (s *Store) PacienteOcupado(pacienteID, fecha, hora, excludeID string) bool {
	for _, c := range s.Citas {
		if c.PacienteID == pacienteID && c.Fecha == fecha && c.Hora == hora && c.ID != excludeID {
			if c.Estado != "cancelada" {
				return true
			}
		}
	}
	return false
}

func (s *Store) EsEstadoModificable(estado string) bool {
	return estado == "pendiente" || estado == "confirmada"
}

func Sanitizar(s string) string {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}

var PalabrasUrgentes = []string{"dolor", "fiebre", "sangrado", "dificultad para respirar", "mareo", "desmayo", "accidente"}

func DetectarPrioridadAutomatica(motivo string) string {
	motivoLower := strings.ToLower(motivo)
	for _, palabra := range PalabrasUrgentes {
		if strings.Contains(motivoLower, palabra) {
			return "urgente"
		}
	}
	return ""
}

func (s *Store) RegistrarHistorial(citaID, accion, estadoAnt, estadoNue, notas string) {
	h := models.EventoHistorial{
		CitaID:    citaID,
		Accion:    accion,
		EstadoAnt: estadoAnt,
		EstadoNue: estadoNue,
		Timestamp: time.Now().Format(time.RFC3339),
		Notas:     notas,
	}
	s.Historial = append(s.Historial, h)
}

func (s *Store) ActualizarMetricas() {
	hoy := time.Now().Format("2006-01-02")
	if s.TodayDate != hoy {
		s.TodayDate = hoy
		s.TodayCreated = 0
		s.TodayCancelled = 0
		s.TodayCompleted = 0
	}
}

func (s *Store) ContarCitasActivas() int {
	count := 0
	for _, c := range s.Citas {
		if c.Estado == "pendiente" || c.Estado == "confirmada" {
			count++
		}
	}
	return count
}
