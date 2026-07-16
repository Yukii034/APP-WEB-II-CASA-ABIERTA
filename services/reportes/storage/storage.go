package storage

import (
	"cuidabien/reportes/models"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

type Store struct {
	Pacientes     []models.Paciente
	NextReportID  int
	CitasURL      string
	MedicamentosURL string
	AlimentacionURL string
}

func NewStore() *Store {
	return &Store{
		Pacientes: []models.Paciente{
			{ID: "P001", Nombre: "Maria Garcia"},
			{ID: "P002", Nombre: "Juan Lopez"},
			{ID: "P003", Nombre: "Ana Martinez"},
		},
		NextReportID:    1,
		CitasURL:        os.Getenv("CITAS_URL"),
		MedicamentosURL: os.Getenv("MEDICAMENTOS_URL"),
		AlimentacionURL: os.Getenv("ALIMENTACION_URL"),
	}
}

func (s *Store) FindPacienteByID(id string) *models.Paciente {
	for i := range s.Pacientes {
		if s.Pacientes[i].ID == id {
			return &s.Pacientes[i]
		}
	}
	return nil
}

// ==================== Clientes HTTP para otros servicios ====================

func (s *Store) fetchJSON(url string, target interface{}) error {
	if url == "" {
		return fmt.Errorf("URL no configurada")
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// ==================== Citas ====================

type citaAPI struct {
	ID        string `json:"id"`
	PacienteID string `json:"paciente_id"`
	DoctorID  string `json:"doctor_id"`
	Fecha     string `json:"fecha"`
	Hora      string `json:"hora"`
	Estado    string `json:"estado"`
	Prioridad string `json:"prioridad"`
}

func (s *Store) ObtenerCitas(pacienteID string) []citaAPI {
	var citas []citaAPI
	url := s.CitasURL + "/api/appointments?paciente_id=" + pacienteID
	if err := s.fetchJSON(url, &citas); err != nil {
		return s.citasSimuladas(pacienteID)
	}
	return citas
}

func (s *Store) ObtenerTodasCitas() []citaAPI {
	var citas []citaAPI
	url := s.CitasURL + "/api/appointments"
	if err := s.fetchJSON(url, &citas); err != nil {
		return s.todasCitasSimuladas()
	}
	return citas
}

func (s *Store) citasSimuladas(pacienteID string) []citaAPI {
	hoy := time.Now().Format("2006-01-02")
	manana := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
AYER := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	semana := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	switch pacienteID {
	case "P001":
		return []citaAPI{
			{ID: "C001", PacienteID: "P001", DoctorID: "D001", Fecha: AYER, Hora: "10:00", Estado: "completada", Prioridad: "normal"},
			{ID: "C002", PacienteID: "P001", DoctorID: "D002", Fecha: manana, Hora: "09:00", Estado: "pendiente", Prioridad: "urgente"},
			{ID: "C003", PacienteID: "P001", DoctorID: "D001", Fecha: semana, Hora: "10:00", Estado: "pendiente", Prioridad: "control"},
		}
	case "P002":
		return []citaAPI{
			{ID: "C004", PacienteID: "P002", DoctorID: "D003", Fecha: hoy, Hora: "11:00", Estado: "confirmada", Prioridad: "normal"},
			{ID: "C005", PacienteID: "P002", DoctorID: "D001", Fecha: AYER, Hora: "08:00", Estado: "cancelada", Prioridad: "normal"},
		}
	case "P003":
		return []citaAPI{
			{ID: "C006", PacienteID: "P003", DoctorID: "D002", Fecha: semana, Hora: "14:00", Estado: "pendiente", Prioridad: "control"},
		}
	}
	return []citaAPI{}
}

func (s *Store) todasCitasSimuladas() []citaAPI {
	var todas []citaAPI
	for _, p := range s.Pacientes {
		todas = append(todas, s.citasSimuladas(p.ID)...)
	}
	return todas
}

// ==================== Medicamentos ====================

type medicamentoAPI struct {
	ID         string   `json:"id"`
	PacienteID string   `json:"paciente_id"`
	Nombre     string   `json:"nombre"`
	Dosis      string   `json:"dosis"`
	Estado     string   `json:"estado"`
}

type tomaAPI struct {
	ID            string `json:"id"`
	MedicamentoID string `json:"medicamento_id"`
	Estado        string `json:"estado"`
}

type adherenciaAPI struct {
	PacienteID       string  `json:"paciente_id"`
	TotalTomas       int     `json:"total_tomas"`
	TomasCumplidas   int     `json:"tomas_cumplidas"`
	Porcentaje       float64 `json:"porcentaje"`
}

type alertaAPI struct {
	ID       string `json:"id"`
	Tipo     string `json:"tipo"`
	Mensaje  string `json:"mensaje"`
	Leida    bool   `json:"leida"`
}

func (s *Store) ObtenerMedicamentos(pacienteID string) []medicamentoAPI {
	var meds []medicamentoAPI
	url := s.MedicamentosURL + "/api/medications?paciente_id=" + pacienteID
	if err := s.fetchJSON(url, &meds); err != nil {
		return s.medicamentosSimulados(pacienteID)
	}
	return meds
}

func (s *Store) ObtenerAdherencia(pacienteID string) adherenciaAPI {
	var resultado adherenciaAPI
	url := s.MedicamentosURL + "/api/medications/adherence/" + pacienteID
	if err := s.fetchJSON(url, &resultado); err != nil {
		return s.adherenciaSimulada(pacienteID)
	}
	return resultado
}

func (s *Store) ObtenerAlertasMedicamentos(pacienteID string) []alertaAPI {
	var alertas []alertaAPI
	url := s.MedicamentosURL + "/api/medications/alerts?paciente_id=" + pacienteID
	if err := s.fetchJSON(url, &alertas); err != nil {
		return s.alertasMedSimuladas(pacienteID)
	}
	return alertas
}

func (s *Store) ContarTodasAlertasMedicamentos() int {
	var alertas []alertaAPI
	url := s.MedicamentosURL + "/api/medications/alerts"
	if err := s.fetchJSON(url, &alertas); err != nil {
		return 3
	}
	count := 0
	for _, a := range alertas {
		if !a.Leida {
			count++
		}
	}
	return count
}

func (s *Store) medicamentosSimulados(pacienteID string) []medicamentoAPI {
	switch pacienteID {
	case "P001":
		return []medicamentoAPI{
			{ID: "M001", PacienteID: "P001", Nombre: "Losartan", Dosis: "50mg", Estado: "activo"},
			{ID: "M002", PacienteID: "P001", Nombre: "Metformina", Dosis: "850mg", Estado: "activo"},
		}
	case "P002":
		return []medicamentoAPI{
			{ID: "M003", PacienteID: "P002", Nombre: "Aspirina", Dosis: "100mg", Estado: "activo"},
			{ID: "M004", PacienteID: "P002", Nombre: "Ibuprofeno", Dosis: "400mg", Estado: "activo"},
		}
	case "P003":
		return []medicamentoAPI{
			{ID: "M005", PacienteID: "P003", Nombre: "Paracetamol", Dosis: "500mg", Estado: "activo"},
			{ID: "M006", PacienteID: "P003", Nombre: "Omeprazol", Dosis: "20mg", Estado: "suspendido"},
		}
	}
	return []medicamentoAPI{}
}

func (s *Store) adherenciaSimulada(pacienteID string) adherenciaAPI {
	switch pacienteID {
	case "P001":
		return adherenciaAPI{PacienteID: "P001", TotalTomas: 14, TomasCumplidas: 12, Porcentaje: 85.7}
	case "P002":
		return adherenciaAPI{PacienteID: "P002", TotalTomas: 10, TomasCumplidas: 8, Porcentaje: 80.0}
	case "P003":
		return adherenciaAPI{PacienteID: "P003", TotalTomas: 21, TomasCumplidas: 21, Porcentaje: 100.0}
	}
	return adherenciaAPI{PacienteID: pacienteID, TotalTomas: 0, TomasCumplidas: 0, Porcentaje: 0}
}

func (s *Store) alertasMedSimuladas(pacienteID string) []alertaAPI {
	if pacienteID == "P002" {
		return []alertaAPI{
			{ID: "A001", Tipo: "por_vencer", Mensaje: "Ibuprofeno vence en 3 dias", Leida: false},
		}
	}
	return []alertaAPI{}
}

// ==================== Alimentacion ====================

type comidaAPI struct {
	TipoComida  string `json:"tipo_comida"`
	Registrada  bool   `json:"registrada"`
	Saltada     bool   `json:"saltada"`
}

type resumenAlimAPI struct {
	Comidas       []comidaAPI `json:"comidas"`
	ComidasHechas int         `json:"comidas_hechas"`
	ComidasTotal  int         `json:"comidas_total"`
	HaySaltadas   bool        `json:"hay_saltadas"`
}

func (s *Store) ObtenerResumenAlimentacion() resumenAlimAPI {
	var res resumenAlimAPI
	url := s.AlimentacionURL + "/api/alimentacion/resumen"
	if err := s.fetchJSON(url, &res); err != nil {
		return s.alimentacionSimulada()
	}
	return res
}

func (s *Store) ContarComidasPaciente(pacienteID string) int {
	res := s.ObtenerResumenAlimentacion()
	return res.ComidasHechas
}

func (s *Store) alimentacionSimulada() resumenAlimAPI {
	return resumenAlimAPI{
		Comidas: []comidaAPI{
			{TipoComida: "desayuno", Registrada: true, Saltada: false},
			{TipoComida: "almuerzo", Registrada: true, Saltada: false},
			{TipoComida: "cena", Registrada: false, Saltada: true},
		},
		ComidasHechas: 2,
		ComidasTotal:  3,
		HaySaltadas:   true,
	}
}

// ==================== Generación de Reportes ====================

func (s *Store) GenerarReporteSemanal(pacienteID string) models.ReporteSemanal {
	paciente := s.FindPacienteByID(pacienteID)
	nombre := "Desconocido"
	if paciente != nil {
		nombre = paciente.Nombre
	}

	hoy := time.Now()
	fechaInicio := hoy.AddDate(0, 0, -7).Format("2006-01-02")
	fechaFin := hoy.Format("2006-01-02")

	citas := s.ObtenerCitas(pacienteID)
	resCitas := s.agregarCitas(citas, fechaInicio, fechaFin)

	meds := s.ObtenerMedicamentos(pacienteID)
	adherencia := s.ObtenerAdherencia(pacienteID)
	alertasMed := s.ObtenerAlertasMedicamentos(pacienteID)
	resMeds := s.agregarMedicamentos(meds, adherencia, alertasMed)

	resAlim := s.agregarAlimentacion(pacienteID)

	estado := s.calcularEstadoGeneral(resCitas, resMeds, resAlim)
	recomendacion := s.generarRecomendacion(resCitas, resMeds, resAlim)

	return models.ReporteSemanal{
		PacienteID:          pacienteID,
		PacienteNombre:      nombre,
		FechaInicio:         fechaInicio,
		FechaFin:            fechaFin,
		ResumenCitas:        resCitas,
		ResumenMedicamentos: resMeds,
		ResumenAlimentacion: resAlim,
		ResumenSalud:        models.ResumenSalud{AlertasSalud: len(alertasMed), SignosVitalesOK: true},
		EstadoGeneral:       estado,
		Recomendacion:       recomendacion,
	}
}

func (s *Store) agregarCitas(citas []citaAPI, fechaInicio, fechaFin string) models.ResumenCitas {
	res := models.ResumenCitas{}
	proximaEncontrada := false

	for _, c := range citas {
		if c.Fecha >= fechaInicio && c.Fecha <= fechaFin {
			switch c.Estado {
			case "completada":
				res.Completadas++
			case "cancelada":
				res.Canceladas++
			case "pendiente", "confirmada":
				res.Pendientes++
			}
			res.TotalProgramadas++
		}
		if !proximaEncontrada && (c.Estado == "pendiente" || c.Estado == "confirmada") && c.Fecha >= time.Now().Format("2006-01-02") {
			res.ProximaCita = c.Fecha + " " + c.Hora
			res.ProximoDoctor = c.DoctorID
			proximaEncontrada = true
		}
	}
	return res
}

func (s *Store) agregarMedicamentos(meds []medicamentoAPI, adherencia adherenciaAPI, alertas []alertaAPI) models.ResumenMedicamentos {
	activos := 0
	for _, m := range meds {
		if m.Estado == "activo" {
			activos++
		}
	}

	alertasNoLeidas := 0
	for _, a := range alertas {
		if !a.Leida {
			alertasNoLeidas++
		}
	}

	return models.ResumenMedicamentos{
		TotalActivos:         activos,
		TomasRegistradas:     adherencia.TotalTomas,
		TomasCumplidas:       adherencia.TomasCumplidas,
		TomasNoCumplidas:     adherencia.TotalTomas - adherencia.TomasCumplidas,
		PorcentajeAdherencia: adherencia.Porcentaje,
		AlertasActivas:       alertasNoLeidas,
	}
}

func (s *Store) agregarAlimentacion(pacienteID string) models.ResumenAlimentacion {
	res := s.ObtenerResumenAlimentacion()
	porcentaje := 0.0
	if res.ComidasTotal > 0 {
		porcentaje = math.Round(float64(res.ComidasHechas)/float64(res.ComidasTotal)*100*100) / 100
	}

	ultimaComida := ""
	for i := len(res.Comidas) - 1; i >= 0; i-- {
		if res.Comidas[i].Registrada {
			ultimaComida = res.Comidas[i].TipoComida
			break
		}
	}

	saltadas := 0
	for _, c := range res.Comidas {
		if c.Saltada {
			saltadas++
		}
	}

	return models.ResumenAlimentacion{
		ComidasRegistradas: res.ComidasHechas,
		ComidasEsperadas:   res.ComidasTotal,
		PorcentajeCumplido: porcentaje,
		ComidasSaltadas:    saltadas,
		UltimaComida:       ultimaComida,
	}
}

func (s *Store) calcularEstadoGeneral(citas models.ResumenCitas, meds models.ResumenMedicamentos, alim models.ResumenAlimentacion) string {
	puntos := 0

	if meds.PorcentajeAdherencia >= 80 {
		puntos += 2
	} else if meds.PorcentajeAdherencia >= 50 {
		puntos += 1
	}

	if alim.PorcentajeCumplido >= 80 {
		puntos += 2
	} else if alim.PorcentajeCumplido >= 50 {
		puntos += 1
	}

	if meds.AlertasActivas == 0 {
		puntos += 1
	}

	if citas.Canceladas > citas.Completadas {
		puntos -= 1
	}

	switch {
	case puntos >= 4:
		return "excelente"
	case puntos >= 3:
		return "estable"
	case puntos >= 2:
		return "requiere_atencion"
	default:
		return "critico"
	}
}

func (s *Store) generarRecomendacion(citas models.ResumenCitas, meds models.ResumenMedicamentos, alim models.ResumenAlimentacion) string {
	var recs []string

	if meds.PorcentajeAdherencia < 80 {
		recs = append(recs, "Mejorar la adherencia a la medicacion")
	}
	if meds.AlertasActivas > 0 {
		recs = append(recs, "Revisar alertas de medicamentos pendientes")
	}
	if alim.ComidasSaltadas > 0 {
		recs = append(recs, "No saltar comidas, mantener horarios regulares")
	}
	if citas.Pendientes > 2 {
		recs = append(recs, "Tiene varias citas pendientes, confirmar asistencia")
	}
	if citas.Canceladas > 0 {
		recs = append(recs, "Reprogramar citas canceladas")
	}

	if len(recs) == 0 {
		return "Mantener los habitos actuales. Todo esta en orden."
	}
	return strings.Join(recs, ". ") + "."
}

// ==================== Reporte por paciente ====================

func (s *Store) GenerarReportePaciente(pacienteID string) models.ReportePaciente {
	paciente := s.FindPacienteByID(pacienteID)
	nombre := "Desconocido"
	if paciente != nil {
		nombre = paciente.Nombre
	}

	citas := s.ObtenerCitas(pacienteID)
	totalCitas := len(citas)
	citasCompletadas := 0
	var histCitas []models.ResumenCita
	for _, c := range citas {
		if c.Estado == "completada" {
			citasCompletadas++
		}
		histCitas = append(histCitas, models.ResumenCita{
			ID: c.ID, Fecha: c.Fecha, Hora: c.Hora, Doctor: c.DoctorID,
			Estado: c.Estado, Prioridad: c.Prioridad,
		})
	}

	meds := s.ObtenerMedicamentos(pacienteID)
	adherencia := s.ObtenerAdherencia(pacienteID)
	alertas := s.ObtenerAlertasMedicamentos(pacienteID)
	comidas := s.ContarComidasPaciente(pacienteID)

	var histMeds []models.ResumenMedicamento
	for _, m := range meds {
		histMeds = append(histMeds, models.ResumenMedicamento{
			Nombre: m.Nombre, Dosis: m.Dosis, Estado: m.Estado, Adherencia: adherencia.Porcentaje,
		})
	}

	alertasNoLeidas := 0
	for _, a := range alertas {
		if !a.Leida {
			alertasNoLeidas++
		}
	}

	return models.ReportePaciente{
		PacienteID:           pacienteID,
		PacienteNombre:       nombre,
		TotalCitas:           totalCitas,
		CitasCompletadas:     citasCompletadas,
		TotalMedicamentos:    len(meds),
		AdherenciaMedicacion: adherencia.Porcentaje,
		ComidasRegistradas:   comidas,
		AlertasActivas:       alertasNoLeidas,
		EstadoGeneral:        s.calcularEstadoGeneral(models.ResumenCitas{TotalProgramadas: totalCitas, Completadas: citasCompletadas}, models.ResumenMedicamentos{PorcentajeAdherencia: adherencia.Porcentaje, AlertasActivas: alertasNoLeidas}, models.ResumenAlimentacion{ComidasRegistradas: comidas, ComidasEsperadas: 3}),
		HistorialCitas:       histCitas,
		HistorialMedicamentos: histMeds,
	}
}

// ==================== Dashboard / Resumen General ====================

func (s *Store) GenerarDashboard() models.DashboardData {
	citas := s.ObtenerTodasCitas()
	hoy := time.Now().Format("2006-01-02")

	citasHoy := 0
	for _, c := range citas {
		if c.Fecha == hoy && (c.Estado == "pendiente" || c.Estado == "confirmada") {
			citasHoy++
		}
	}

	totalMedsActivos := 0
	for _, p := range s.Pacientes {
		meds := s.ObtenerMedicamentos(p.ID)
		for _, m := range meds {
			if m.Estado == "activo" {
				totalMedsActivos++
			}
		}
	}

	alertasTotales := s.ContarTodasAlertasMedicamentos()

sumaAdherencia := 0.0
pacientesConAlertas := 0
var pacienteRes []models.PacienteResumen

for _, p := range s.Pacientes {
		adh := s.ObtenerAdherencia(p.ID)
		alertas := s.ObtenerAlertasMedicamentos(p.ID)
		meds := s.ObtenerMedicamentos(p.ID)
		citasP := s.ObtenerCitas(p.ID)

		activos := 0
		for _, m := range meds {
			if m.Estado == "activo" {
				activos++
			}
		}

		proximas := 0
		for _, c := range citasP {
			if (c.Estado == "pendiente" || c.Estado == "confirmada") && c.Fecha >= hoy {
				proximas++
			}
		}

		alertasNoLeidas := 0
		for _, a := range alertas {
			if !a.Leida {
				alertasNoLeidas++
			}
		}
		if alertasNoLeidas > 0 {
			pacientesConAlertas++
		}

		estado := "estable"
		if adh.Porcentaje < 50 {
			estado = "requiere_atencion"
		} else if adh.Porcentaje >= 90 {
			estado = "excelente"
		}

		pacienteRes = append(pacienteRes, models.PacienteResumen{
			ID: p.ID, Nombre: p.Nombre, CitasProximas: proximas,
			MedicamentosActivos: activos, Adherencia: adh.Porcentaje, Estado: estado,
		})

		sumaAdherencia += adh.Porcentaje
	}

	promedioAdherencia := 0.0
	if len(s.Pacientes) > 0 {
		promedioAdherencia = math.Round(sumaAdherencia/float64(len(s.Pacientes))*100) / 100
	}

	return models.DashboardData{
		ResumenGeneral: models.ResumenGeneral{
			TotalPacientes:         len(s.Pacientes),
			TotalCitasHoy:          citasHoy,
			TotalMedicamentosActivos: totalMedsActivos,
			TotalAlertasPendientes: alertasTotales,
			PromedioAdherencia:     promedioAdherencia,
			PacientesConAlertas:    pacientesConAlertas,
		},
		Pacientes: pacienteRes,
	}
}
