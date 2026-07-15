package storage

import (
	"cuidabien/citas/models"
	"testing"
	"time"
)

func newTestStore() *Store {
	s := NewStore()
	s.Citas = []models.Cita{}
	s.Historial = []models.EventoHistorial{}
	s.NextID = 1
	s.TodayCreated = 0
	s.TodayCancelled = 0
	s.TodayCompleted = 0
	return s
}

// ==================== Tests: Validacion de fechas ====================

func TestEsFechaPasada_FechaPasada(t *testing.T) {
	s := newTestStore()
	if !s.EsFechaPasada("2020-01-01", "10:00") {
		t.Error("Se esperaba que 2020-01-01 fuera fecha pasada")
	}
}

func TestEsFechaPasada_FechaFutura(t *testing.T) {
	s := newTestStore()
	futura := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	if s.EsFechaPasada(futura, "10:00") {
		t.Errorf("Se esperaba que %s fuera fecha valida", futura)
	}
}

func TestEsFechaPasada_FormatoInvalido(t *testing.T) {
	s := newTestStore()
	if !s.EsFechaPasada("fecha-mala", "hora-mala") {
		t.Error("Se esperaba que formato invalido devolviera true")
	}
}

// ==================== Tests: Deteccion de prioridad ====================

func TestDetectarPrioridad_Dolor(t *testing.T) {
	resultado := DetectarPrioridadAutomatica("Tengo mucho dolor de cabeza")
	if resultado != "urgente" {
		t.Errorf("Se esperaba 'urgente', se obtuvo '%s'", resultado)
	}
}

func TestDetectarPrioridad_Fiebre(t *testing.T) {
	resultado := DetectarPrioridadAutomatica("Fiebre alta desde ayer")
	if resultado != "urgente" {
		t.Errorf("Se esperaba 'urgente', se obtuvo '%s'", resultado)
	}
}

func TestDetectarPrioridad_Sangrado(t *testing.T) {
	resultado := DetectarPrioridadAutomatica("Sangrado nasal persistente")
	if resultado != "urgente" {
		t.Errorf("Se esperaba 'urgente', se obtuvo '%s'", resultado)
	}
}

func TestDetectarPrioridad_Control(t *testing.T) {
	resultado := DetectarPrioridadAutomatica("Control mensual de presion")
	if resultado != "" {
		t.Errorf("Se esperaba prioridad normal, se obtuvo '%s'", resultado)
	}
}

// ==================== Tests: Sanitizacion ====================

func TestSanitizar_EspaciosExtra(t *testing.T) {
	resultado := Sanitizar("  hola   mundo  ")
	if resultado != "hola mundo" {
		t.Errorf("Se esperaba 'hola mundo', se obtuvo '%s'", resultado)
	}
}

func TestSanitizar_StringVacio(t *testing.T) {
	resultado := Sanitizar("")
	if resultado != "" {
		t.Errorf("Se esperaba string vacio, se obtuvo '%s'", resultado)
	}
}

// ==================== Tests: Reglas de negocio ====================

func TestMedicoOcupado_Disponible(t *testing.T) {
	s := newTestStore()
	if s.MedicoOcupado("D001", "2030-01-01", "10:00", "") {
		t.Error("El medico deberia estar disponible")
	}
}

func TestMedicoOcupado_Ocupado(t *testing.T) {
	s := newTestStore()
	s.Citas = append(s.Citas, models.Cita{
		ID: "C001", DoctorID: "D001", PacienteID: "P001",
		Fecha: "2030-01-01", Hora: "10:00", Estado: "pendiente",
	})
	if !s.MedicoOcupado("D001", "2030-01-01", "10:00", "") {
		t.Error("El medico deberia estar ocupado")
	}
}

func TestMedicoOcupado_CitaCancelada(t *testing.T) {
	s := newTestStore()
	s.Citas = append(s.Citas, models.Cita{
		ID: "C001", DoctorID: "D001", PacienteID: "P001",
		Fecha: "2030-01-01", Hora: "10:00", Estado: "cancelada",
	})
	if s.MedicoOcupado("D001", "2030-01-01", "10:00", "") {
		t.Error("El medico deberia estar disponible si la cita esta cancelada")
	}
}

func TestMedicoOcupado_ExcluirCitaActual(t *testing.T) {
	s := newTestStore()
	s.Citas = append(s.Citas, models.Cita{
		ID: "C001", DoctorID: "D001", PacienteID: "P001",
		Fecha: "2030-01-01", Hora: "10:00", Estado: "pendiente",
	})
	if s.MedicoOcupado("D001", "2030-01-01", "10:00", "C001") {
		t.Error("Deberia excluir la cita actual al verificar")
	}
}

func TestPacienteOcupado_Disponible(t *testing.T) {
	s := newTestStore()
	if s.PacienteOcupado("P001", "2030-01-01", "10:00", "") {
		t.Error("El paciente deberia estar disponible")
	}
}

func TestPacienteOcupado_Ocupado(t *testing.T) {
	s := newTestStore()
	s.Citas = append(s.Citas, models.Cita{
		ID: "C001", DoctorID: "D001", PacienteID: "P001",
		Fecha: "2030-01-01", Hora: "10:00", Estado: "pendiente",
	})
	if !s.PacienteOcupado("P001", "2030-01-01", "10:00", "") {
		t.Error("El paciente deberia estar ocupado")
	}
}

// ==================== Tests: Estados ====================

func TestEsEstadoModificable_Pendiente(t *testing.T) {
	s := newTestStore()
	if !s.EsEstadoModificable("pendiente") {
		t.Error("El estado pendiente deberia ser modificable")
	}
}

func TestEsEstadoModificable_Confirmada(t *testing.T) {
	s := newTestStore()
	if !s.EsEstadoModificable("confirmada") {
		t.Error("El estado confirmada deberia ser modificable")
	}
}

func TestEsEstadoModificable_Cancelada(t *testing.T) {
	s := newTestStore()
	if s.EsEstadoModificable("cancelada") {
		t.Error("El estado cancelada NO deberia ser modificable")
	}
}

func TestEsEstadoModificable_Completada(t *testing.T) {
	s := newTestStore()
	if s.EsEstadoModificable("completada") {
		t.Error("El estado completada NO deberia ser modificable")
	}
}

// ==================== Tests: Metricas ====================

func TestActualizarMetricas_CambioDeDia(t *testing.T) {
	s := newTestStore()
	s.TodayDate = "2000-01-01"
	s.ActualizarMetricas()
	hoy := time.Now().Format("2006-01-02")
	if s.TodayDate != hoy {
		t.Errorf("Se esperaba fecha '%s', se obtuvo '%s'", hoy, s.TodayDate)
	}
	if s.TodayCreated != 0 || s.TodayCancelled != 0 || s.TodayCompleted != 0 {
		t.Error("Las metricas deberian resetearse al cambiar de dia")
	}
}

func TestContarCitasActivas(t *testing.T) {
	s := newTestStore()
	s.Citas = []models.Cita{
		{ID: "C001", Estado: "pendiente"},
		{ID: "C002", Estado: "confirmada"},
		{ID: "C003", Estado: "cancelada"},
		{ID: "C004", Estado: "completada"},
	}
	if s.ContarCitasActivas() != 2 {
		t.Errorf("Se esperaban 2 citas activas, se obtuvieron %d", s.ContarCitasActivas())
	}
}

// ==================== Tests: Historial ====================

func TestRegistrarHistorial(t *testing.T) {
	s := newTestStore()
	s.RegistrarHistorial("C001", "creacion", "", "pendiente", "")
	if len(s.Historial) != 1 {
		t.Errorf("Se esperaba 1 evento, se obtuvieron %d", len(s.Historial))
	}
	if s.Historial[0].CitaID != "C001" {
		t.Errorf("Se esperaba cita_id 'C001', se obtuvo '%s'", s.Historial[0].CitaID)
	}
	if s.Historial[0].Accion != "creacion" {
		t.Errorf("Se esperaba accion 'creacion', se obtuvo '%s'", s.Historial[0].Accion)
	}
}

func TestRegistrarHistorial_ConNotas(t *testing.T) {
	s := newTestStore()
	s.RegistrarHistorial("C001", "nota_agregada", "", "confirmada", "Paciente estable")
	if s.Historial[0].Notas != "Paciente estable" {
		t.Errorf("Se esperaban notas 'Paciente estable', se obtuvo '%s'", s.Historial[0].Notas)
	}
}
