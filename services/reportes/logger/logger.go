package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

type LogEntry struct {
	Nivel   string `json:"nivel"`
	Mensaje string `json:"mensaje"`
	Tipo    string `json:"tipo"`
	Path    string `json:"path,omitempty"`
	IP      string `json:"ip,omitempty"`
	Tiempo  string `json:"tiempo"`
}

func LogJSON(nivel, mensaje, tipo, path, ip string) {
	entry := LogEntry{
		Nivel:   nivel,
		Mensaje: mensaje,
		Tipo:    tipo,
		Path:    path,
		IP:      ip,
		Tiempo:  time.Now().Format(time.RFC3339),
	}
	data, _ := json.Marshal(entry)
	fmt.Println(string(data))
}
