package storage

import (
	"strconv"
	"sync"
	"time"

	"cuidabien/alimentacion/modelo"
)

// Store guarda todo el estado del servicio en memoria, protegido con
// mutex para uso concurrente. Los datos se pierden si el contenedor
// se reinicia (ver docs/arquitectura.md - limitaciones).
//
// Es la única capa que conoce cómo se guardan los datos hoy; si mañana
// se cambia a una base de datos, solo este archivo cambiaría.
type Store struct {
	muComidas         sync.Mutex
	Comidas           []modelo.RegistroComida
	siguienteIDComida int

	muHidratacion          sync.Mutex
	Hidratacion            []modelo.RegistroHidratacion
	siguienteIDHidratacion int

	muRestricciones        sync.Mutex
	Restricciones          []modelo.Restriccion
	siguienteIDRestriccion int
}

// NewStore crea el almacenamiento en memoria con datos semilla de ejemplo.
func NewStore() *Store {
	store := &Store{
		siguienteIDComida:      1,
		siguienteIDHidratacion: 1,
		siguienteIDRestriccion: 1,
	}
	store.SembrarDatos()
	return store
}

func (s *Store) GuardarComida(reg modelo.RegistroComida) {
	s.muComidas.Lock()
	defer s.muComidas.Unlock()
	s.Comidas = append(s.Comidas, reg)
}

func (s *Store) ListarComidas() []modelo.RegistroComida {
	s.muComidas.Lock()
	defer s.muComidas.Unlock()

	resultado := make([]modelo.RegistroComida, len(s.Comidas))
	copy(resultado, s.Comidas)
	return resultado
}

func (s *Store) SiguienteIDComida() string {
	s.muComidas.Lock()
	defer s.muComidas.Unlock()

	id := strconv.Itoa(s.siguienteIDComida)
	s.siguienteIDComida++
	return id
}

func (s *Store) GuardarHidratacion(reg modelo.RegistroHidratacion) {
	s.muHidratacion.Lock()
	defer s.muHidratacion.Unlock()
	s.Hidratacion = append(s.Hidratacion, reg)
}

func (s *Store) ListarHidratacion() []modelo.RegistroHidratacion {
	s.muHidratacion.Lock()
	defer s.muHidratacion.Unlock()

	resultado := make([]modelo.RegistroHidratacion, len(s.Hidratacion))
	copy(resultado, s.Hidratacion)
	return resultado
}

func (s *Store) SiguienteIDHidratacion() string {
	s.muHidratacion.Lock()
	defer s.muHidratacion.Unlock()

	id := strconv.Itoa(s.siguienteIDHidratacion)
	s.siguienteIDHidratacion++
	return id
}

func (s *Store) GuardarRestriccion(r modelo.Restriccion) {
	s.muRestricciones.Lock()
	defer s.muRestricciones.Unlock()
	s.Restricciones = append(s.Restricciones, r)
}

func (s *Store) ListarRestricciones() []modelo.Restriccion {
	s.muRestricciones.Lock()
	defer s.muRestricciones.Unlock()

	resultado := make([]modelo.Restriccion, len(s.Restricciones))
	copy(resultado, s.Restricciones)
	return resultado
}

func (s *Store) SiguienteIDRestriccion() string {
	s.muRestricciones.Lock()
	defer s.muRestricciones.Unlock()

	id := strconv.Itoa(s.siguienteIDRestriccion)
	s.siguienteIDRestriccion++
	return id
}

// SembrarDatos carga datos iniciales de ejemplo para demostración.
func (s *Store) SembrarDatos() {
	hoy := time.Now()

	comidas := []modelo.RegistroComida{
		{ID: "1", TipoComida: "desayuno", Descripcion: "Avena con frutas y cafe", Hora: time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 7, 30, 0, 0, hoy.Location())},
		{ID: "2", TipoComida: "almuerzo", Descripcion: "Arroz, pollo guisado y ensalada", Hora: time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 12, 15, 0, 0, hoy.Location())},
		{ID: "3", TipoComida: "cena", Descripcion: "Sopa de verduras y pan tostado", Hora: time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 18, 45, 0, 0, hoy.Location())},
	}
	for _, c := range comidas {
		s.GuardarComida(c)
	}
	s.siguienteIDComida = 4

	hidratacion := []modelo.RegistroHidratacion{
		{ID: "1", Hora: time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 8, 0, 0, 0, hoy.Location()), Cantidad: "1 vaso de agua"},
		{ID: "2", Hora: time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 10, 30, 0, 0, hoy.Location()), Cantidad: "1 vaso de agua"},
		{ID: "3", Hora: time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 13, 0, 0, 0, hoy.Location()), Cantidad: "2 vasos de agua"},
	}
	for _, h := range hidratacion {
		s.GuardarHidratacion(h)
	}
	s.siguienteIDHidratacion = 4

	restricciones := []modelo.Restriccion{
		{ID: "1", Descripcion: "Sin sal - presion alta"},
		{ID: "2", Descripcion: "Diabetico - evitar azucar"},
		{ID: "3", Descripcion: "Alergia a lacteos"},
	}
	for _, r := range restricciones {
		s.GuardarRestriccion(r)
	}
	s.siguienteIDRestriccion = 4
}

// Reiniciar limpia comidas e hidratación (útil para demos).
func (s *Store) Reiniciar() {
	s.muComidas.Lock()
	s.Comidas = nil
	s.muComidas.Unlock()

	s.muHidratacion.Lock()
	s.Hidratacion = nil
	s.muHidratacion.Unlock()
}
