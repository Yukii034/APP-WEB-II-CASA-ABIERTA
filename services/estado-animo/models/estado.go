package models

type RegistroEstadoAnimo struct {
	ID         string `json:"id"`
	Fecha      string `json:"fecha"`
	Nivel      int    `json:"nivel"`
	Emocion    string `json:"emocion"`
	Comentario string `json:"comentario"`
}

type Alerta struct {
	GenerarAlerta bool   `json:"generar_alerta"`
	Mensaje       string `json:"mensaje"`
}