package handlers

import (
	"encoding/json"
	"net/http"

	"cuidabien/alimentacion/modelo"
	"cuidabien/alimentacion/service"
)

// Handlers traduce peticiones HTTP hacia el service y las respuestas
// del service de vuelta a JSON. No contiene lógica de negocio propia.
type Handlers struct {
	Service *service.Service
}

func New(s *service.Service) *Handlers {
	return &Handlers{Service: s}
}

// HealthHandler responde el endpoint obligatorio usado por el gateway
// y por docker-compose para saber si el servicio está vivo.
func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /api/alimentacion -> lista las comidas registradas hoy
// POST /api/alimentacion -> registra una comida { "tipo_comida": "almuerzo", "descripcion": "sopa y pollo" }
func (h *Handlers) AlimentacionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.Service.ComidasDeHoy())

	case http.MethodPost:
		var entrada modelo.EntradaRegistroComida
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
			http.Error(w, "Body inválido, se espera {\"tipo_comida\": \"almuerzo\", \"descripcion\": \"...\"}", http.StatusBadRequest)
			return
		}

		nuevo, ok := h.Service.RegistrarComida(entrada)
		if !ok {
			http.Error(w, "Body inválido, se espera {\"tipo_comida\": \"almuerzo\", \"descripcion\": \"...\"}", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(nuevo)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /api/alimentacion/resumen -> estado de las comidas del día
func (h *Handlers) ResumenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Service.Resumen())
}

// POST /api/alimentacion/reset -> limpia los registros de hoy (útil para demos)
func (h *Handlers) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	h.Service.Reiniciar()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "reiniciado"})
}

// GET /api/alimentacion/historial?dias=7 -> registros de los últimos N días (por defecto 7)
func (h *Handlers) HistorialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	dias := service.ParsearDias(r.URL.Query().Get("dias"))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Service.Historial(dias))
}

// GET /api/alimentacion/hidratacion -> lista los registros de hidratación de hoy
// POST /api/alimentacion/hidratacion -> registra hidratación { "cantidad": "1 vaso" }
func (h *Handlers) HidratacionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.Service.HidratacionDeHoy())

	case http.MethodPost:
		var entrada modelo.EntradaRegistroHidratacion
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
			http.Error(w, "Body inválido, se espera {\"cantidad\": \"1 vaso\"}", http.StatusBadRequest)
			return
		}

		nuevo := h.Service.RegistrarHidratacion(entrada)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(nuevo)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /api/alimentacion/restricciones -> lista las restricciones/alergias registradas
// POST /api/alimentacion/restricciones -> agrega una restricción { "descripcion": "sin sal" }
func (h *Handlers) RestriccionesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.Service.Restricciones())

	case http.MethodPost:
		var entrada modelo.EntradaRestriccion
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
			http.Error(w, "Body inválido, se espera {\"descripcion\": \"sin sal\"}", http.StatusBadRequest)
			return
		}

		nuevo, ok := h.Service.RegistrarRestriccion(entrada)
		if !ok {
			http.Error(w, "Body inválido, se espera {\"descripcion\": \"sin sal\"}", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(nuevo)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
