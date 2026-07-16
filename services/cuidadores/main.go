package main

import (
	"log"
	"net/http"
	"os"

	"cuidabien/cuidadores/internal/handler"
	"cuidabien/cuidadores/internal/model"
	"cuidabien/cuidadores/internal/repository"
	"cuidabien/cuidadores/internal/router"
	"cuidabien/cuidadores/internal/service"
)

func main() {
	repo := repository.NuevaMemoriaRepository()
	svc := service.NuevoCuidadorService(repo)
	seedCuidadores(svc)
	h := handler.NuevoCuidadorHandler(svc)
	mux := router.Nuevo(h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // puerto interno por defecto dentro del contenedor
	}

	log.Printf("Servicio de cuidadores corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func seedCuidadores(svc *service.CuidadorService) {
	semillas := []model.Cuidador{
		{
			Nombre:               "Ana Torres",
			Telefono:             "555-1001",
			Email:                "ana.torres@cuidabien.local",
			Relacion:             "Cuidadora principal",
			HorarioDisponible:    "Lunes a viernes 08:00-16:00",
			Pacientes:            []string{"P001", "P003"},
			NivelResponsabilidad: "alta",
		},
		{
			Nombre:               "Luis Mendoza",
			Telefono:             "555-1002",
			Email:                "luis.mendoza@cuidabien.local",
			Relacion:             "Enfermero",
			HorarioDisponible:    "Noches 18:00-06:00",
			Pacientes:            []string{"P002"},
			NivelResponsabilidad: "media",
		},
		{
			Nombre:               "Elena Garcia",
			Telefono:             "555-1003",
			Email:                "elena.garcia@cuidabien.local",
			Relacion:             "Familiar",
			HorarioDisponible:    "Fines de semana",
			Pacientes:            []string{"P001"},
			NivelResponsabilidad: "apoyo",
		},
	}

	for _, c := range semillas {
		_, _ = svc.Crear(c)
	}
}
