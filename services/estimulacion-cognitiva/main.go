package main

import (
	"net/http"
	"os"

	"cuidabien/estimulacion-cognitiva/logger"
	"cuidabien/estimulacion-cognitiva/router"
	"cuidabien/estimulacion-cognitiva/storage"
)

func main() {
	log := logger.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	store := storage.New()
	if os.Getenv("SEED_DATOS") != "false" {
		store.Seed()
		log.Info("Datos de ejemplo precargados (SEED_DATOS=false para desactivar)")
	}
	handler := router.New(store, log)
	log.Info("Servicio de estimulación cognitiva escuchando en el puerto %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Error("el servidor se detuvo: %v", err)
		os.Exit(1)
	}
}
