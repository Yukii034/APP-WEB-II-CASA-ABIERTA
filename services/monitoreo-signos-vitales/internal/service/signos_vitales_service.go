package service

import (
	"fmt"
	"strings"
	"time"

	"monitoreo-signos-vitales/internal/models"
	"monitoreo-signos-vitales/internal/storage"
)

// SignosVitalesService implements validation, classification and queries.
type SignosVitalesService struct {
	repo storage.SignosVitalesRepository
}

func NuevoSignosVitalesService(repo storage.SignosVitalesRepository) *SignosVitalesService {
	return &SignosVitalesService{repo: repo}
}

func (s *SignosVitalesService) Crear(entrada models.EntradaSignosVitales) (models.SignosVitales, error) {
	if err := validarEntrada(entrada); err != nil {
		return models.SignosVitales{}, err
	}
	registro := models.SignosVitales{ID: s.repo.SiguienteID(), IDAdultoMayor: entrada.IDAdultoMayor, RegistradoPor: entrada.RegistradoPor, FechaRegistro: time.Now(), PresionSistolica: entrada.PresionSistolica, PresionDiastolica: entrada.PresionDiastolica, FrecuenciaCardiaca: entrada.FrecuenciaCardiaca, Temperatura: entrada.Temperatura, SaturacionOxigeno: entrada.SaturacionOxigeno, NivelGlucosa: entrada.NivelGlucosa, Peso: entrada.Peso, Altura: entrada.Altura, NivelDolor: entrada.NivelDolor, Observaciones: entrada.Observaciones}
	registro.Evaluacion = Evaluar(registro)
	s.repo.Guardar(registro)
	return registro, nil
}

func (s *SignosVitalesService) Historial(idAdultoMayor string) []models.SignosVitales {
	registros := s.repo.PorAdultoMayor(idAdultoMayor)
	for i := range registros {
		registros[i].Evaluacion = Evaluar(registros[i])
	}
	return registros
}

func (s *SignosVitalesService) Ultimo(idAdultoMayor string) (models.SignosVitales, bool) {
	historial := s.Historial(idAdultoMayor)
	if len(historial) == 0 {
		return models.SignosVitales{}, false
	}
	return historial[0], true
}

func (s *SignosVitalesService) Tendencia(idAdultoMayor, parametro string, dias int) ([]models.PuntoTendencia, error) {
	if dias <= 0 {
		dias = 30
	}
	if dias > 365 {
		return nil, fmt.Errorf("dias no puede ser mayor que 365")
	}
	parametro = strings.ToLower(parametro)
	if !parametroValido(parametro) {
		return nil, fmt.Errorf("parametro no válido")
	}
	limite := time.Now().AddDate(0, 0, -dias)
	puntos := make([]models.PuntoTendencia, 0)
	for _, registro := range s.Historial(idAdultoMayor) {
		if registro.FechaRegistro.Before(limite) {
			continue
		}
		valor, estado, disponible := valorParametro(registro, parametro)
		if disponible {
			puntos = append(puntos, models.PuntoTendencia{Fecha: registro.FechaRegistro, Valor: valor, Estado: estado})
		}
	}
	return puntos, nil
}

// Evaluar classifies values using the demo reference ranges from the specification.
// In production those ranges should be supplied by a configurable clinical catalog.
func Evaluar(registro models.SignosVitales) models.EvaluacionRegistro {
	valores := []models.EvaluacionValor{
		{"presion_sistolica", clasificar(float64(registro.PresionSistolica), 100, 135, 90, 160)},
		{"presion_diastolica", clasificar(float64(registro.PresionDiastolica), 60, 88, 50, 100)},
		{"frecuencia_cardiaca", clasificar(float64(registro.FrecuenciaCardiaca), 55, 95, 45, 120)},
	}
	if registro.Temperatura != nil {
		valores = append(valores, models.EvaluacionValor{"temperatura", clasificar(*registro.Temperatura, 36, 37.4, 35, 38.3)})
	}
	if registro.SaturacionOxigeno != nil {
		valores = append(valores, models.EvaluacionValor{"saturacion_oxigeno", clasificar(float64(*registro.SaturacionOxigeno), 93, 100, 88, 101)})
	}
	if registro.NivelGlucosa != nil {
		valores = append(valores, models.EvaluacionValor{"nivel_glucosa", clasificar(*registro.NivelGlucosa, 70, 140, 55, 200)})
	}
	general := models.EstadoNormal
	for _, valor := range valores {
		if valor.Estado == models.EstadoCritico {
			general = models.EstadoCritico
			break
		}
		if valor.Estado == models.EstadoBajo || valor.Estado == models.EstadoAlto {
			general = models.EstadoAdvertencia
		}
	}
	return models.EvaluacionRegistro{EstadoGeneral: general, Valores: valores}
}

func validarEntrada(e models.EntradaSignosVitales) error {
	if strings.TrimSpace(e.IDAdultoMayor) == "" {
		return fmt.Errorf("id_adulto_mayor es obligatorio")
	}
	if e.PresionSistolica <= 0 || e.PresionDiastolica <= 0 || e.FrecuenciaCardiaca <= 0 {
		return fmt.Errorf("presion_sistolica, presion_diastolica y frecuencia_cardiaca deben ser mayores que cero")
	}
	if e.NivelDolor != nil && (*e.NivelDolor < 0 || *e.NivelDolor > 10) {
		return fmt.Errorf("nivel_dolor debe estar entre 0 y 10")
	}
	return nil
}

func clasificar(valor, minimo, maximo, criticoBajo, criticoAlto float64) models.Estado {
	if valor <= criticoBajo || valor >= criticoAlto {
		return models.EstadoCritico
	}
	if valor < minimo {
		return models.EstadoBajo
	}
	if valor > maximo {
		return models.EstadoAlto
	}
	return models.EstadoNormal
}

func parametroValido(p string) bool { _, _, ok := valorParametro(models.SignosVitales{}, p); return ok }
func valorParametro(r models.SignosVitales, p string) (float64, models.Estado, bool) {
	for _, e := range Evaluar(r).Valores {
		if e.Parametro != p {
			continue
		}
		switch p {
		case "presion_sistolica":
			return float64(r.PresionSistolica), e.Estado, true
		case "presion_diastolica":
			return float64(r.PresionDiastolica), e.Estado, true
		case "frecuencia_cardiaca":
			return float64(r.FrecuenciaCardiaca), e.Estado, true
		case "temperatura":
			if r.Temperatura != nil {
				return *r.Temperatura, e.Estado, true
			}
		case "saturacion_oxigeno":
			if r.SaturacionOxigeno != nil {
				return float64(*r.SaturacionOxigeno), e.Estado, true
			}
		case "nivel_glucosa":
			if r.NivelGlucosa != nil {
				return *r.NivelGlucosa, e.Estado, true
			}
		}
	}
	return 0, "", p == "presion_sistolica" || p == "presion_diastolica" || p == "frecuencia_cardiaca" || p == "temperatura" || p == "saturacion_oxigeno" || p == "nivel_glucosa"
}
