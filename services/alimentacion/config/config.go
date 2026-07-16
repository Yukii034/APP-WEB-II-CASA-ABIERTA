package config

import "os"

// Config agrupa la configuración del servicio leída del entorno.
type Config struct {
	Puerto        string
	DesayunoHasta string
	AlmuerzoHasta string
	CenaHasta     string
}

// Load lee la configuración desde variables de entorno, aplicando
// valores por defecto cuando no están definidas.
func Load() Config {
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	return Config{
		Puerto:        puerto,
		DesayunoHasta: os.Getenv("DESAYUNO_HASTA"),
		AlmuerzoHasta: os.Getenv("ALMUERZO_HASTA"),
		CenaHasta:     os.Getenv("CENA_HASTA"),
	}
}
