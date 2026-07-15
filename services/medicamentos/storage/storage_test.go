package storage

import (
	"cuidabien/medicamentos/models"
	"testing"
	"time"
)

func newTestStore() *Store {
	s := NewStore()
	s.Medicamentos = []models.Medicamento{}
	s.Tomas = []models.Toma{}
	s.Alertas = []models.Alerta{}
	s.NextMedID = 1
	s.NextTomaID = 1
	s.NextAlertaID = 1
	return s
}

// ==================== Tests: Validaciones ====================

func TestValidarHorario_Correcto(t *testing.T) {
	if !ValidarHorario("08:00") {
		t.Error("08:00 deberia ser un horario valido")
	}
	if !ValidarHorario("23:59") {
		t.Error("23:59 deberia ser un horario valido")
	}
}

func TestValidarHorario_Invalido(t *testing.T) {
	if ValidarHorario("25:00") {
		t.Error("25:00 no deberia ser valido")
	}
	if ValidarHorario("12:60") {
		t.Error("12:60 no deberia ser valido")
	}
	if ValidarHorario("abc") {
		t.Error("abc no deberia ser valido")
	}
	if ValidarHorario("8:00") {
		t.Error("8:00 sin cero adelante no deberia ser valido")
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

// ==================== Tests: Busqueda ====================

func TestFindMedicamentoByID_Encontrado(t *testing.T) {
	s := newTestStore()
	s.Medicamentos = append(s.Medicamentos, models.Medicamento{ID: "M001", Nombre: "Test"})
	med := s.FindMedicamentoByID("M001")
	if med == nil || med.Nombre != "Test" {
		t.Error("Deberia encontrar el medicamento M001")
	}
}

func TestFindMedicamentoByID_NoEncontrado(t *testing.T) {
	s := newTestStore()
	med := s.FindMedicamentoByID("M999")
	if med != nil {
		t.Error("No deberia encontrar M999")
	}
}

func TestFindPacienteByID_Encontrado(t *testing.T) {
	s := NewStore()
	p := s.FindPacienteByID("P001")
	if p == nil || p.Nombre != "Maria Garcia" {
		t.Error("Deberia encontrar paciente P001")
	}
}

// ==================== Tests: Interacciones ====================

func TestBuscarInteraccion_Encontrada(t *testing.T) {
	s := NewStore()
	inter := s.BuscarInteraccion("Ibuprofeno", "Aspirina")
	if inter == nil {
		t.Fatal("Deberia encontrar interaccion entre Ibuprofeno y Aspirina")
	}
	if inter.Gravedad != "severa" {
		t.Errorf("Se esperaba gravedad 'severa', se obtuvo '%s'", inter.Gravedad)
	}
}

func TestBuscarInteraccion_Inversa(t *testing.T) {
	s := NewStore()
	inter := s.BuscarInteraccion("Aspirina", "Ibuprofeno")
	if inter == nil {
		t.Error("Deberia encontrar la interaccion en orden inverso")
	}
}

func TestBuscarInteraccion_NoEncontrada(t *testing.T) {
	s := NewStore()
	inter := s.BuscarInteraccion("Paracetamol", "Omeprazol")
	if inter != nil {
		t.Error("No deberia encontrar interaccion entre estos medicamentos")
	}
}

func TestBuscarInteraccion_CaseInsensitive(t *testing.T) {
	s := NewStore()
	inter := s.BuscarInteraccion("ibuprofeno", "aspirina")
	if inter == nil {
		t.Error("Deberia encontrar interaccion sin importar mayusculas")
	}
}

// ==================== Tests: Adherencia ====================

func TestCalcularAdherencia_TodasCumplidas(t *testing.T) {
	s := newTestStore()
	s.Tomas = []models.Toma{
		{ID: "T001", PacienteID: "P001", Estado: "cumplida"},
		{ID: "T002", PacienteID: "P001", Estado: "cumplida"},
	}
	adh := s.CalcularAdherencia("P001")
	if adh.Porcentaje != 100 {
		t.Errorf("Se esperaba 100%%, se obtuvo %.2f", adh.Porcentaje)
	}
	if adh.TomasCumplidas != 2 {
		t.Errorf("Se esperaban 2 tomas cumplidas, se obtuvieron %d", adh.TomasCumplidas)
	}
}

func TestCalcularAdherencia_NingunaCumplida(t *testing.T) {
	s := newTestStore()
	s.Tomas = []models.Toma{
		{ID: "T001", PacienteID: "P001", Estado: "no_cumplida"},
		{ID: "T002", PacienteID: "P001", Estado: "no_cumplida"},
	}
	adh := s.CalcularAdherencia("P001")
	if adh.Porcentaje != 0 {
		t.Errorf("Se esperaba 0%%, se obtuvo %.2f", adh.Porcentaje)
	}
}

func TestCalcularAdherencia_Mixta(t *testing.T) {
	s := newTestStore()
	s.Tomas = []models.Toma{
		{ID: "T001", PacienteID: "P001", Estado: "cumplida"},
		{ID: "T002", PacienteID: "P001", Estado: "cumplida"},
		{ID: "T003", PacienteID: "P001", Estado: "cumplida"},
		{ID: "T004", PacienteID: "P001", Estado: "no_cumplida"},
	}
	adh := s.CalcularAdherencia("P001")
	if adh.Porcentaje != 75 {
		t.Errorf("Se esperaba 75%%, se obtuvo %.2f", adh.Porcentaje)
	}
}

func TestCalcularAdherencia_SinTomas(t *testing.T) {
	s := newTestStore()
	adh := s.CalcularAdherencia("P001")
	if adh.Porcentaje != 0 {
		t.Errorf("Se esperaba 0%% sin tomas, se obtuvo %.2f", adh.Porcentaje)
	}
}

// ==================== Tests: Alertas ====================

func TestGenerarAlertasVencimiento_MedicamentoVencido(t *testing.T) {
	s := newTestStore()
	vencido := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	s.Medicamentos = append(s.Medicamentos, models.Medicamento{
		ID: "M001", PacienteID: "P001", Nombre: "TestMed",
		Estado: "activo", FechaFin: vencido,
	})
	s.GenerarAlertasVencimiento()
	if len(s.Alertas) != 1 {
		t.Fatalf("Se esperaba 1 alerta, se obtuvieron %d", len(s.Alertas))
	}
	if s.Alertas[0].Tipo != "vencido" {
		t.Errorf("Se esperaba tipo 'vencido', se obtuvo '%s'", s.Alertas[0].Tipo)
	}
}

func TestGenerarAlertasVencimiento_MedicamentoPorVencer(t *testing.T) {
	s := newTestStore()
	porVencer := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	s.Medicamentos = append(s.Medicamentos, models.Medicamento{
		ID: "M001", PacienteID: "P001", Nombre: "TestMed",
		Estado: "activo", FechaFin: porVencer,
	})
	s.GenerarAlertasVencimiento()
	if len(s.Alertas) != 1 {
		t.Fatalf("Se esperaba 1 alerta, se obtuvieron %d", len(s.Alertas))
	}
	if s.Alertas[0].Tipo != "por_vencer" {
		t.Errorf("Se esperaba tipo 'por_vencer', se obtuvo '%s'", s.Alertas[0].Tipo)
	}
}

func TestGenerarAlertasVencimiento_NoDuplicar(t *testing.T) {
	s := newTestStore()
	vencido := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	s.Medicamentos = append(s.Medicamentos, models.Medicamento{
		ID: "M001", PacienteID: "P001", Nombre: "TestMed",
		Estado: "activo", FechaFin: vencido,
	})
	s.GenerarAlertasVencimiento()
	s.GenerarAlertasVencimiento()
	if len(s.Alertas) != 1 {
		t.Errorf("No deberia duplicar alertas, se obtuvieron %d", len(s.Alertas))
	}
}

// ==================== Tests: Generar IDs ====================

func TestGenerateMedID(t *testing.T) {
	s := newTestStore()
	id1 := s.GenerateMedID()
	id2 := s.GenerateMedID()
	if id1 == id2 {
		t.Error("Los IDs deberian ser unicos")
	}
}

func TestGenerateTomaID(t *testing.T) {
	s := newTestStore()
	id1 := s.GenerateTomaID()
	id2 := s.GenerateTomaID()
	if id1 == id2 {
		t.Error("Los IDs deberian ser unicos")
	}
}

// ==================== Tests: Filtros ====================

func TestMedicamentosPorPaciente(t *testing.T) {
	s := newTestStore()
	s.Medicamentos = []models.Medicamento{
		{ID: "M001", PacienteID: "P001", Nombre: "A"},
		{ID: "M002", PacienteID: "P001", Nombre: "B"},
		{ID: "M003", PacienteID: "P002", Nombre: "C"},
	}
	result := s.MedicamentosPorPaciente("P001")
	if len(result) != 2 {
		t.Errorf("Se esperaban 2 medicamentos, se obtuvieron %d", len(result))
	}
}

func TestMedicamentosActivosPorPaciente(t *testing.T) {
	s := newTestStore()
	s.Medicamentos = []models.Medicamento{
		{ID: "M001", PacienteID: "P001", Nombre: "A", Estado: "activo"},
		{ID: "M002", PacienteID: "P001", Nombre: "B", Estado: "suspendido"},
		{ID: "M003", PacienteID: "P001", Nombre: "C", Estado: "activo"},
	}
	result := s.MedicamentosActivosPorPaciente("P001")
	if len(result) != 2 {
		t.Errorf("Se esperaban 2 activos, se obtuvieron %d", len(result))
	}
}
