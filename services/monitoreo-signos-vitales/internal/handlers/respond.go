package handlers

import (
	"encoding/json"
	"net/http"
)

func responderJSON(w http.ResponseWriter, estado int, valor any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(estado)
	_ = json.NewEncoder(w).Encode(valor)
}
func responderError(w http.ResponseWriter, estado int, mensaje string) {
	responderJSON(w, estado, map[string]string{"error": mensaje})
}
