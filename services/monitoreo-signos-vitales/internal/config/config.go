package config

// config.go
//
// Responsabilidad: centralizar la configuración del microservicio
// obtenida desde variables de entorno (SPEC.md §17, §10.6).
//
// No contiene lógica de negocio ni de conexión a PostgreSQL.
// La conexión a la base de datos se implementa en TASK-004 (storage/factory.go),
// que consumirá los valores expuestos por este paquete.

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config agrupa toda la configuración del microservicio.
type Config struct {
	AppName    string
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
}

// Load carga las variables de entorno y construye la configuración
// centralizada del microservicio.
//
// Si existe un archivo .env en la raíz del proyecto, sus valores se
// cargan en el entorno del proceso. Su ausencia no es un error: en
// entornos como Docker, las variables ya llegan definidas externamente.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("config: no se encontró archivo .env, se usarán variables de entorno del sistema")
	}

	cfg := &Config{
		AppName:    getEnv("APP_NAME", "monitoreo-signos-vitales"),
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "monitoreo_signos_vitales"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate garantiza que la configuración obligatoria esté presente.
// No valida credenciales de conexión: eso ocurre al intentar conectar
// en storage/factory.go (TASK-004).
func (c *Config) validate() error {
	if c.DBUser == "" {
		return fmt.Errorf("config: DB_USER es obligatorio")
	}
	if c.DBName == "" {
		return fmt.Errorf("config: DB_NAME es obligatorio")
	}
	return nil
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
// si no se encuentra definida.
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return defaultValue
}
