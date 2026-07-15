package storage

import (
	"cuidabien/contacto-emergencia/models"
	"fmt"
	"sort"
	"time"
)

type Store struct {
	Contactos      []models.Contacto
	Alertas        []models.Alerta
	Historial      []models.EventoHistorial
	NextContactoID int
	NextAlertaID   int
	TodayDate      string
	TodayActivadas int
	TodayAtendidas int
}

func NewStore() *Store {
	return &Store{
		Contactos: []models.Contacto{
			{ID: "C001", PacienteID: "P001", Nombre: "Sofia Garcia", Telefono: "555-0102", Parentesco: "Hija", Prioridad: 1, Principal: true},
			{ID: "C002", PacienteID: "P001", Nombre: "Pedro Garcia", Telefono: "555-0103", Parentesco: "Hijo", Prioridad: 2, Principal: false},
			{ID: "C003", PacienteID: "P002", Nombre: "Rosa Lopez", Telefono: "555-0202", Parentesco: "Esposa", Prioridad: 1, Principal: true},
			{ID: "C004", PacienteID: "P003", Nombre: "Diego Martinez", Telefono: "555-0302", Parentesco: "Cuidador", Prioridad: 1, Principal: true},
		},
		NextContactoID: 5,
		NextAlertaID:   1,
		TodayDate:      time.Now().Format("2006-01-02"),
	}
}

// --- Generacion de IDs ---

func (s *Store) GenerateContactoID() string {
	id := fmt.Sprintf("C%03d", s.NextContactoID)
	s.NextContactoID++
	return id
}

func (s *Store) GenerateAlertaID() string {
	id := fmt.Sprintf("A%03d", s.NextAlertaID)
	s.NextAlertaID++
	return id
}

// --- Busquedas ---

func (s *Store) FindContactoByID(id string) *models.Contacto {
	for i := range s.Contactos {
		if s.Contactos[i].ID == id {
			return &s.Contactos[i]
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

// ContactosPorPaciente devuelve los contactos de un paciente
// ordenados por prioridad (1 = primero a notificar).
func (s *Store) ContactosPorPaciente(pacienteID string) []models.Contacto {
	var resultado []models.Contacto
	for _, c := range s.Contactos {
		if c.PacienteID == pacienteID {
			resultado = append(resultado, c)
		}
	}
	sort.Slice(resultado, func(i, j int) bool {
		return resultado[i].Prioridad < resultado[j].Prioridad
	})
	return resultado
}

// EliminarContacto quita un contacto de la lista. Devuelve true si existia.
func (s *Store) EliminarContacto(id string) bool {
	for i := range s.Contactos {
		if s.Contactos[i].ID == id {
			s.Contactos = append(s.Contactos[:i], s.Contactos[i+1:]...)
			return true
		}
	}
	return false
}

// NotificarContactos simula el envio de una notificacion a cada contacto
// del paciente, respetando el orden de prioridad. Devuelve la lista de
// nombres notificados (para dejar registro en la alerta).
func (s *Store) NotificarContactos(pacienteID string) []string {
	contactos := s.ContactosPorPaciente(pacienteID)
	notificados := make([]string, 0, len(contactos))
	for _, c := range contactos {
		// Aqui iria la integracion real (SMS, llamada, push, etc.)
		// Por ahora es una simulacion en memoria.
		notificados = append(notificados, fmt.Sprintf("%s (%s) - %s", c.Nombre, c.Parentesco, c.Telefono))
	}
	return notificados
}

// --- Historial ---

func (s *Store) RegistrarHistorial(alertaID, accion, notas string) {
	h := models.EventoHistorial{
		AlertaID:  alertaID,
		Accion:    accion,
		Timestamp: time.Now().Format(time.RFC3339),
		Notas:     notas,
	}
	s.Historial = append(s.Historial, h)
}

func (s *Store) HistorialPorAlerta(alertaID string) []models.EventoHistorial {
	var resultado []models.EventoHistorial
	for _, h := range s.Historial {
		if h.AlertaID == alertaID {
			resultado = append(resultado, h)
		}
	}
	return resultado
}

// --- Metricas ---

func (s *Store) ActualizarContadoresDiarios() {
	hoy := time.Now().Format("2006-01-02")
	if s.TodayDate != hoy {
		s.TodayDate = hoy
		s.TodayActivadas = 0
		s.TodayAtendidas = 0
	}
}

func (s *Store) NivelValido(nivel string) bool {
	return nivel == "leve" || nivel == "moderado" || nivel == "critico"
}
