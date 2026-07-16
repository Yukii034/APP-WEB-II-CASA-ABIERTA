package model

// Cuidador representa a la persona responsable del cuidado de uno o
// varios adultos mayores. Este servicio NO guarda información clínica;
// eso pertenece a otros servicios (informacion-salud, medicamentos, etc.).
type Cuidador struct {
	ID                   string   `json:"id"`
	Nombre               string   `json:"nombre"`
	Telefono             string   `json:"telefono"`
	Email                string   `json:"email"`
	Relacion             string   `json:"relacion"`
	HorarioDisponible    string   `json:"horario_disponible"`
	Pacientes            []string `json:"pacientes"`
	NivelResponsabilidad string   `json:"nivel_responsabilidad"`
	Activo               bool     `json:"activo"`
}
