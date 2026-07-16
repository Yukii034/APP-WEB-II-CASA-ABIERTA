package router

import (
	"net/http"

	"cuidabien/cuidadores/internal/handler"
)

// Nuevo arma el mux con todas las rutas del servicio de cuidadores.
// Usa el enrutador nativo de Go 1.22 (soporta método + parámetros de
// ruta tipo {id}), así que no se necesitan dependencias externas.
func Nuevo(h *handler.CuidadorHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /api/cuidadores", h.Crear)
	mux.HandleFunc("GET /api/cuidadores", h.Listar)
	mux.HandleFunc("GET /api/cuidadores/paciente/{pacienteId}", h.ListarPorPaciente)
	mux.HandleFunc("GET /api/cuidadores/{id}", h.ObtenerPorID)
	mux.HandleFunc("PUT /api/cuidadores/{id}", h.Actualizar)
	mux.HandleFunc("DELETE /api/cuidadores/{id}", h.Eliminar)

	return mux
}
