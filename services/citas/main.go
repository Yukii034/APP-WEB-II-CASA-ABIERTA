package main

import (
	"cuidabien/citas/handlers"
	"cuidabien/citas/middleware"
	"cuidabien/citas/router"
	"cuidabien/citas/storage"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	store := storage.NewStore()
	h := handlers.New(store)
	r := router.New(h)

	limiter := middleware.NewRateLimiter(100, time.Minute)
	wrapped := middleware.RateLimit(limiter)(middleware.Auth(r))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servicio de citas corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, wrapped))
}
