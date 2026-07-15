package models

// paciente.go
//
// Responsabilidad: definición de la entidad Paciente (SPEC.md §8.4, §9).
// Entidad principal del dominio: representa al adulto mayor monitoreado.

import (
	"time"

	"github.com/google/uuid"
)

// Paciente representa a un adulto mayor monitoreado dentro del sistema.
// Contiene información básica y relativamente estática (SPEC.md §8.4).
//
// Nota: el esquema sigue estrictamente los campos definidos en SPEC.md
// §8.4 y §9 (tabla `pacientes`). No se incluye borrado lógico (soft delete)
// por no estar especificado; el endpoint DELETE (SPEC.md §14) se resolverá
// como eliminación física en la capa de storage.
type Paciente struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Nombres         string    `gorm:"type:varchar(150);not null"`
	Apellidos       string    `gorm:"type:varchar(150);not null"`
	Cedula          string    `gorm:"type:varchar(20);not null;uniqueIndex"`
	FechaNacimiento time.Time `gorm:"type:date;not null"`
	Sexo            string    `gorm:"type:varchar(20);not null"`
	Telefono        string    `gorm:"type:varchar(20)"`
	Direccion       string    `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TableName fija explícitamente el nombre de la tabla en PostgreSQL,
// conforme al Modelo Relacional definido en SPEC.md §9.
func (Paciente) TableName() string {
	return "pacientes"
}
