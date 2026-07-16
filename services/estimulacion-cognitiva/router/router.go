package router

import (
	"net/http"

	"cuidabien/estimulacion-cognitiva/handlers"
	"cuidabien/estimulacion-cognitiva/logger"
	"cuidabien/estimulacion-cognitiva/middleware"
	"cuidabien/estimulacion-cognitiva/storage"
)

func New(store *storage.Storage, log *logger.Logger) http.Handler {
	h := handlers.New(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/ejercicios", h.Listar)
	mux.HandleFunc("POST /api/ejercicios", h.Agregar)
	mux.HandleFunc("GET /api/ejercicios/resumen", h.Resumen)
	mux.HandleFunc("POST /api/ejercicios/reset", h.Reset)

	return middleware.Recover(log)(middleware.Logging(log)(mux))
}
