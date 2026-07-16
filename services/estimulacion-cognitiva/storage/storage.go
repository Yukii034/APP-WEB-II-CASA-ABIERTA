package storage

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"cuidabien/estimulacion-cognitiva/models"
)

var ErrTipoRequerido = errors.New("tipo es requerido, ej. 'memoria', 'trivia', 'sopa_letras'")

const umbralDiasAlerta = 2

type Storage struct {
	mu         sync.Mutex
	ejercicios []models.Ejercicio
	siguiente  int
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Agregar(tipo string) (models.Ejercicio, error) {
	tipo = strings.TrimSpace(tipo)
	if tipo == "" {
		return models.Ejercicio{}, ErrTipoRequerido
	}

	return s.agregarConFecha(tipo, time.Now()), nil
}

func (s *Storage) agregarConFecha(tipo string, fecha time.Time) models.Ejercicio {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.siguiente++
	e := models.Ejercicio{
		ID:    strconv.Itoa(s.siguiente),
		Tipo:  tipo,
		Fecha: fecha,
	}
	s.ejercicios = append(s.ejercicios, e)
	return e
}

func (s *Storage) Listar() []models.Ejercicio {
	s.mu.Lock()
	defer s.mu.Unlock()
	copia := make([]models.Ejercicio, len(s.ejercicios))
	copy(copia, s.ejercicios)
	return copia
}

func diasEntre(desde, hasta time.Time) int {
	d := time.Date(desde.Year(), desde.Month(), desde.Day(), 0, 0, 0, 0, desde.Location())
	h := time.Date(hasta.Year(), hasta.Month(), hasta.Day(), 0, 0, 0, 0, hasta.Location())
	return int(h.Sub(d).Hours() / 24)
}

func (s *Storage) Resumen() models.Resumen {
	ejercicios := s.Listar()
	ahora := time.Now()
	if len(ejercicios) == 0 {
		return models.Resumen{
			Ejercicios: []models.Ejercicio{},
			Total:      0,
			HayAlerta:  true,
			Mensaje:    "Aún no se ha registrado ningún ejercicio de estimulación cognitiva.",
		}
	}

	ultimo := ejercicios[0]
	for _, e := range ejercicios[1:] {
		if e.Fecha.After(ultimo.Fecha) {
			ultimo = e
		}
	}

	ejerciciosHoy := 0
	for _, e := range ejercicios {
		if diasEntre(e.Fecha, ahora) == 0 {
			ejerciciosHoy++
		}
	}

	diasDesdeUltimo := diasEntre(ultimo.Fecha, ahora)
	hayAlerta := diasDesdeUltimo >= umbralDiasAlerta

	mensaje := "Actividad cognitiva al día."
	if hayAlerta {
		mensaje = "Han pasado varios días sin actividad cognitiva registrada."
	}

	ultimoCopia := ultimo
	return models.Resumen{
		Ejercicios:      ejercicios,
		Total:           len(ejercicios),
		EjerciciosHoy:   ejerciciosHoy,
		UltimoEjercicio: &ultimoCopia,
		DiasDesdeUltimo: diasDesdeUltimo,
		HayAlerta:       hayAlerta,
		Mensaje:         mensaje,
	}
}

func (s *Storage) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ejercicios = nil
	s.siguiente = 0
}

func (s *Storage) Seed() {
	ahora := time.Now()
	s.agregarConFecha("trivia", ahora.AddDate(0, 0, -4))
	s.agregarConFecha("memoria", ahora.AddDate(0, 0, -umbralDiasAlerta))
}
