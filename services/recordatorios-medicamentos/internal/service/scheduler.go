package service

import (
	"context"
	"time"
)

// IniciarVerificador revisa periódicamente la hora actual.
// Evita generar dos veces la misma alerta en un mismo minuto.
func (s *RecordatorioMedicamentoService) IniciarVerificador(
	ctx context.Context,
	intervalo time.Duration,
	ubicacion *time.Location,
) {
	if intervalo <= 0 {
		intervalo = time.Minute
	}

	if ubicacion == nil {
		ubicacion = time.Local
	}

	ticker := time.NewTicker(intervalo)
	defer ticker.Stop()

	ultimaHora := ""

	for {
		select {
		case <-ctx.Done():
			return

		case ahora := <-ticker.C:
			horaActual := ahora.
				In(ubicacion).
				Format("15:04")

			if horaActual == ultimaHora {
				continue
			}

			ultimaHora = horaActual

			_, _ = s.VerificarHora(horaActual)
		}
	}
}