package storage

// factory.go
//
// Responsabilidad: crear y configurar la conexión a PostgreSQL mediante GORM.
// Ver SPEC.md §10.3 y §10.6.
//
// No contiene lógica de negocio ni queries de dominio: los repositorios
// concretos (paciente_postgres.go, signos_vitales_postgres.go, usuario_gorm.go)
// consumirán la conexión que aquí se construye.

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"monitoreo-signos-vitales/internal/config"
)

// NewPostgresConnection abre una conexión a PostgreSQL usando GORM,
// a partir de la configuración centralizada (internal/config).
func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := buildDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("storage: error al conectar con PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("storage: error al obtener *sql.DB: %w", err)
	}

	// Configuración del pool de conexiones.
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("storage: error al verificar conexión con PostgreSQL: %w", err)
	}

	return db, nil
}

// buildDSN construye el Data Source Name para PostgreSQL a partir
// de la configuración cargada desde variables de entorno.
func buildDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)
}
