package modelo

import "time"

// RegistroComida representa una comida registrada.
type RegistroComida struct {
	ID          string    `json:"id"`
	TipoComida  string    `json:"tipo_comida"` // desayuno | almuerzo | cena | merienda
	Descripcion string    `json:"descripcion,omitempty"`
	Hora        time.Time `json:"hora"`
}

// EntradaRegistroComida es el body esperado en POST /api/alimentacion.
type EntradaRegistroComida struct {
	TipoComida  string `json:"tipo_comida"`
	Descripcion string `json:"descripcion"`
}

// RegistroHidratacion representa un registro de líquidos tomados.
type RegistroHidratacion struct {
	ID       string    `json:"id"`
	Hora     time.Time `json:"hora"`
	Cantidad string    `json:"cantidad,omitempty"` // ej. "1 vaso", "poco", opcional
}

// EntradaRegistroHidratacion es el body esperado en POST /api/alimentacion/hidratacion.
type EntradaRegistroHidratacion struct {
	Cantidad string `json:"cantidad"`
}

// Restriccion representa una restricción o alergia alimentaria del adulto mayor.
type Restriccion struct {
	ID          string `json:"id"`
	Descripcion string `json:"descripcion"` // ej. "sin sal", "diabético", "alergia a lácteos"
}

// EntradaRestriccion es el body esperado en POST /api/alimentacion/restricciones.
type EntradaRestriccion struct {
	Descripcion string `json:"descripcion"`
}

// EstadoComida indica si una comida del día ya se registró o si se saltó.
type EstadoComida struct {
	TipoComida string `json:"tipo_comida"`
	Registrada bool   `json:"registrada"`
	Saltada    bool   `json:"saltada"`
	HoraLimite string `json:"hora_limite"`
}

// Resumen es el estado del día: qué comidas van, cuáles faltan y si hay
// alguna que ya se saltó (pasó la hora límite y no se registró).
type Resumen struct {
	Comidas       []EstadoComida `json:"comidas"`
	ComidasHechas int            `json:"comidas_hechas"`
	ComidasTotal  int            `json:"comidas_total"`
	HaySaltadas   bool           `json:"hay_saltadas"`
	Mensaje       string         `json:"mensaje,omitempty"`
	NivelAlerta   string         `json:"nivel_alerta"` // "ok" | "atencion" | "urgente"
}

// ComidaEsperada define, para una comida del día, la hora tras la cual
// se considera saltada si no fue registrada.
type ComidaEsperada struct {
	Tipo       string
	HoraLimite string // formato "HH:MM"
}
