package storage

import (
	"cuidabien/reportes/models"
	"testing"
)

func newTestStore() *Store {
	s := NewStore()
	s.CitasURL = ""
	s.MedicamentosURL = ""
	s.AlimentacionURL = ""
	return s
}

// ==================== Tests: Pacientes ====================

func TestFindPacienteByID_Encontrado(t *testing.T) {
	s := NewStore()
	p := s.FindPacienteByID("P001")
	if p == nil || p.Nombre != "Maria Garcia" {
		t.Error("Deberia encontrar paciente P001")
	}
}

func TestFindPacienteByID_NoEncontrado(t *testing.T) {
	s := NewStore()
	p := s.FindPacienteByID("P999")
	if p != nil {
		t.Error("No deberia encontrar P999")
	}
}

// ==================== Tests: Datos simulados ====================

func TestCitasSimuladas_P001(t *testing.T) {
	s := newTestStore()
	citas := s.citasSimuladas("P001")
	if len(citas) != 3 {
		t.Errorf("Se esperaban 3 citas para P001, se obtuvieron %d", len(citas))
	}
}

func TestCitasSimuladas_PacienteNoExiste(t *testing.T) {
	s := newTestStore()
	citas := s.citasSimuladas("P999")
	if len(citas) != 0 {
		t.Errorf("Se esperaban 0 citas, se obtuvieron %d", len(citas))
	}
}

func TestMedicamentosSimulados_P001(t *testing.T) {
	s := newTestStore()
	meds := s.medicamentosSimulados("P001")
	if len(meds) != 2 {
		t.Errorf("Se esperaban 2 medicamentos, se obtuvieron %d", len(meds))
	}
}

func TestAdherenciaSimulada_P001(t *testing.T) {
	s := newTestStore()
	adh := s.adherenciaSimulada("P001")
	if adh.Porcentaje != 85.7 {
		t.Errorf("Se esperaba 85.7, se obtuvo %.1f", adh.Porcentaje)
	}
}

func TestAlimentacionSimulada(t *testing.T) {
	s := newTestStore()
	res := s.alimentacionSimulada()
	if res.ComidasHechas != 2 {
		t.Errorf("Se esperaban 2 comidas, se obtuvieron %d", res.ComidasHechas)
	}
	if !res.HaySaltadas {
		t.Error("Deberia haber comidas saltadas")
	}
}

// ==================== Tests: Estado general ====================

func TestCalcularEstadoGeneral_Excelente(t *testing.T) {
	s := newTestStore()
	citas := models.ResumenCitas{TotalProgramadas: 2, Completadas: 2, Canceladas: 0, Pendientes: 0}
	meds := models.ResumenMedicamentos{PorcentajeAdherencia: 80.0, AlertasActivas: 0}
	alim := models.ResumenAlimentacion{PorcentajeCumplido: 90.0}
	estado := s.calcularEstadoGeneral(citas, meds, alim)
	if estado != "excelente" && estado != "estable" {
		t.Errorf("Se esperaba excelente/estable, se obtuvo '%s'", estado)
	}
}

func TestCalcularEstadoGeneral_Critico(t *testing.T) {
	s := newTestStore()
	citas := models.ResumenCitas{TotalProgramadas: 4, Completadas: 1, Canceladas: 3, Pendientes: 0}
	meds := models.ResumenMedicamentos{PorcentajeAdherencia: 20.0, AlertasActivas: 3}
	alim := models.ResumenAlimentacion{PorcentajeCumplido: 10.0, ComidasSaltadas: 5}
	estado := s.calcularEstadoGeneral(citas, meds, alim)
	if estado != "requiere_atencion" && estado != "critico" {
		t.Errorf("Se esperaba requiere_atencion/critico, se obtuvo '%s'", estado)
	}
}

func TestCalcularEstadoGeneral_Estable(t *testing.T) {
	s := newTestStore()
	citas := models.ResumenCitas{TotalProgramadas: 3, Completadas: 2, Canceladas: 0, Pendientes: 1}
	meds := models.ResumenMedicamentos{PorcentajeAdherencia: 75.0, AlertasActivas: 0}
	alim := models.ResumenAlimentacion{PorcentajeCumplido: 70.0}
	estado := s.calcularEstadoGeneral(citas, meds, alim)
	if estado != "estable" && estado != "requiere_atencion" {
		t.Errorf("Se esperaba estable/requiere_atencion, se obtuvo '%s'", estado)
	}
}

// ==================== Tests: Recomendaciones ====================

func TestGenerarRecomendacion_Buena(t *testing.T) {
	s := newTestStore()
	citas := models.ResumenCitas{TotalProgramadas: 2, Completadas: 2, Pendientes: 0}
	meds := models.ResumenMedicamentos{PorcentajeAdherencia: 95.0, AlertasActivas: 0}
	alim := models.ResumenAlimentacion{PorcentajeCumplido: 100.0, ComidasSaltadas: 0}
	rec := s.generarRecomendacion(citas, meds, alim)
	if rec == "" {
		t.Error("La recomendacion no deberia estar vacia")
	}
}

func TestGenerarRecomendacion_Mala(t *testing.T) {
	s := newTestStore()
	citas := models.ResumenCitas{TotalProgramadas: 4, Completadas: 1, Pendientes: 3}
	meds := models.ResumenMedicamentos{PorcentajeAdherencia: 30.0, AlertasActivas: 2}
	alim := models.ResumenAlimentacion{PorcentajeCumplido: 40.0, ComidasSaltadas: 3}
	rec := s.generarRecomendacion(citas, meds, alim)
	if rec == "" {
		t.Error("La recomendacion no deberia estar vacia")
	}
}

// ==================== Tests: Dashboard ====================

func TestGenerarDashboard(t *testing.T) {
	s := newTestStore()
	dash := s.GenerarDashboard()
	if dash.ResumenGeneral.TotalPacientes != 3 {
		t.Errorf("Se esperaban 3 pacientes, se obtuvieron %d", dash.ResumenGeneral.TotalPacientes)
	}
	if len(dash.Pacientes) != 3 {
		t.Errorf("Se esperaban 3 resumenes de pacientes, se obtuvieron %d", len(dash.Pacientes))
	}
}

func TestGenerarReporteSemanal(t *testing.T) {
	s := newTestStore()
	reporte := s.GenerarReporteSemanal("P001")
	if reporte.PacienteID != "P001" {
		t.Errorf("Se esperaba paciente P001, se obtuvo '%s'", reporte.PacienteID)
	}
	if reporte.PacienteNombre != "Maria Garcia" {
		t.Errorf("Se esperaba nombre 'Maria Garcia', se obtuvo '%s'", reporte.PacienteNombre)
	}
	if reporte.EstadoGeneral == "" {
		t.Error("El estado general no deberia estar vacio")
	}
}

func TestGenerarReportePaciente(t *testing.T) {
	s := newTestStore()
	reporte := s.GenerarReportePaciente("P002")
	if reporte.PacienteID != "P002" {
		t.Errorf("Se esperaba paciente P002, se obtuvo '%s'", reporte.PacienteID)
	}
}
