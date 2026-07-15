package service

import (
	"strconv"
	"time"

	"cuidabien/alimentacion/modelo"
	"cuidabien/alimentacion/storage"
)

// Service contiene la lógica de negocio del servicio de alimentación.
// Depende del storage.Store para leer/guardar datos, pero toda la
// lógica de "qué significa una comida saltada" o "cuál es el nivel de
// alerta" vive aquí, no en storage ni en handlers.
type Service struct {
	store            *storage.Store
	comidasEsperadas []modelo.ComidaEsperada
}

// New crea el service con las comidas esperadas por defecto (desayuno,
// almuerzo, cena).
func New(store *storage.Store) *Service {
	return &Service{
		store: store,
		comidasEsperadas: []modelo.ComidaEsperada{
			{Tipo: "desayuno", HoraLimite: "10:00"},
			{Tipo: "almuerzo", HoraLimite: "15:00"},
			{Tipo: "cena", HoraLimite: "21:00"},
		},
	}
}

// ConfigurarHoraLimite permite ajustar la hora límite de una comida
// esperada (usado en main.go a partir de config.Config).
func (s *Service) ConfigurarHoraLimite(tipo, horaLimite string) {
	if horaLimite == "" {
		return
	}
	for i := range s.comidasEsperadas {
		if s.comidasEsperadas[i].Tipo == tipo {
			s.comidasEsperadas[i].HoraLimite = horaLimite
		}
	}
}

// RegistrarComida valida y guarda una comida registrada hoy.
func (s *Service) RegistrarComida(entrada modelo.EntradaRegistroComida) (modelo.RegistroComida, bool) {
	if entrada.TipoComida == "" {
		return modelo.RegistroComida{}, false
	}

	nuevo := modelo.RegistroComida{
		ID:          s.store.SiguienteIDComida(),
		TipoComida:  entrada.TipoComida,
		Descripcion: entrada.Descripcion,
		Hora:        time.Now(),
	}
	s.store.GuardarComida(nuevo)
	return nuevo, true
}

// ComidasDeHoy devuelve las comidas registradas en el día actual.
func (s *Service) ComidasDeHoy() []modelo.RegistroComida {
	return filtrarPorDia(s.store.ListarComidas(), time.Now())
}

// Resumen calcula el estado de las comidas del día (hechas, saltadas,
// nivel de alerta).
func (s *Service) Resumen() modelo.Resumen {
	return calcularResumen(s.ComidasDeHoy(), s.comidasEsperadas, time.Now())
}

// Historial devuelve las comidas registradas en los últimos "dias" días.
func (s *Service) Historial(dias int) []modelo.RegistroComida {
	if dias <= 0 {
		dias = 7
	}
	desde := time.Now().AddDate(0, 0, -dias)

	resultado := []modelo.RegistroComida{}
	for _, reg := range s.store.ListarComidas() {
		if reg.Hora.After(desde) {
			resultado = append(resultado, reg)
		}
	}
	return resultado
}

// Reiniciar limpia comidas e hidratación registradas (útil para demos).
func (s *Service) Reiniciar() {
	s.store.Reiniciar()
}

// RegistrarHidratacion guarda un registro de líquidos tomados.
func (s *Service) RegistrarHidratacion(entrada modelo.EntradaRegistroHidratacion) modelo.RegistroHidratacion {
	nuevo := modelo.RegistroHidratacion{
		ID:       s.store.SiguienteIDHidratacion(),
		Hora:     time.Now(),
		Cantidad: entrada.Cantidad,
	}
	s.store.GuardarHidratacion(nuevo)
	return nuevo
}

// HidratacionDeHoy devuelve los registros de hidratación del día actual.
func (s *Service) HidratacionDeHoy() []modelo.RegistroHidratacion {
	hoy := time.Now().Format("2006-01-02")

	resultado := []modelo.RegistroHidratacion{}
	for _, reg := range s.store.ListarHidratacion() {
		if reg.Hora.Format("2006-01-02") == hoy {
			resultado = append(resultado, reg)
		}
	}
	return resultado
}

// RegistrarRestriccion valida y guarda una restricción/alergia alimentaria.
func (s *Service) RegistrarRestriccion(entrada modelo.EntradaRestriccion) (modelo.Restriccion, bool) {
	if entrada.Descripcion == "" {
		return modelo.Restriccion{}, false
	}

	nuevo := modelo.Restriccion{
		ID:          s.store.SiguienteIDRestriccion(),
		Descripcion: entrada.Descripcion,
	}
	s.store.GuardarRestriccion(nuevo)
	return nuevo, true
}

// Restricciones devuelve todas las restricciones/alergias registradas.
func (s *Service) Restricciones() []modelo.Restriccion {
	resultado := s.store.ListarRestricciones()
	if resultado == nil {
		resultado = []modelo.Restriccion{}
	}
	return resultado
}

// ParsearDias interpreta el query param "dias" (usado por el handler).
func ParsearDias(valor string) int {
	if valor == "" {
		return 0
	}
	n, err := strconv.Atoi(valor)
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

// ---- Funciones puras (sin estado), fáciles de probar con go test ----

func filtrarPorDia(regs []modelo.RegistroComida, dia time.Time) []modelo.RegistroComida {
	hoy := dia.Format("2006-01-02")
	resultado := []modelo.RegistroComida{}
	for _, reg := range regs {
		if reg.Hora.Format("2006-01-02") == hoy {
			resultado = append(resultado, reg)
		}
	}
	return resultado
}

// calcularResumen es una función pura para que sea fácil de probar con
// go test, sin necesidad de un storage.Store real.
func calcularResumen(regs []modelo.RegistroComida, esperadas []modelo.ComidaEsperada, ahora time.Time) modelo.Resumen {
	registradas := map[string]bool{}
	for _, reg := range regs {
		registradas[reg.TipoComida] = true
	}

	var comidas []modelo.EstadoComida
	comidasHechas := 0
	haySaltadas := false

	for _, esperada := range esperadas {
		hecha := registradas[esperada.Tipo]
		saltada := false

		if !hecha {
			limite, err := time.ParseInLocation("15:04", esperada.HoraLimite, ahora.Location())
			if err == nil {
				limiteHoy := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), limite.Hour(), limite.Minute(), 0, 0, ahora.Location())
				if ahora.After(limiteHoy) {
					saltada = true
					haySaltadas = true
				}
			}
		} else {
			comidasHechas++
		}

		comidas = append(comidas, modelo.EstadoComida{
			TipoComida: esperada.Tipo,
			Registrada: hecha,
			Saltada:    saltada,
			HoraLimite: esperada.HoraLimite,
		})
	}

	mensaje := ""
	if haySaltadas {
		mensaje = "Hay una o más comidas que no se registraron a tiempo hoy."
	}

	comidasSaltadas := 0
	for _, c := range comidas {
		if c.Saltada {
			comidasSaltadas++
		}
	}
	nivelAlerta := "ok"
	if comidasSaltadas == 1 {
		nivelAlerta = "atencion"
	} else if comidasSaltadas >= 2 {
		nivelAlerta = "urgente"
	}

	return modelo.Resumen{
		Comidas:       comidas,
		ComidasHechas: comidasHechas,
		ComidasTotal:  len(esperadas),
		HaySaltadas:   haySaltadas,
		Mensaje:       mensaje,
		NivelAlerta:   nivelAlerta,
	}
}
