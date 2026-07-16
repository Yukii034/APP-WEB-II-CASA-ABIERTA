package storage

import (
	"strconv"
	"sync"

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

// NewStore crea el almacenamiento en memoria con los contadores de ID
// listos para usarse.
func NewStore() *Store {
	return &Store{
		siguienteIDComida:      1,
		siguienteIDHidratacion: 1,
		siguienteIDRestriccion: 1,
	}
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

// Reiniciar limpia comidas e hidratación (útil para demos).
func (s *Store) Reiniciar() {
	s.muComidas.Lock()
	s.Comidas = nil
	s.muComidas.Unlock()

	s.muHidratacion.Lock()
	s.Hidratacion = nil
	s.muHidratacion.Unlock()
}
