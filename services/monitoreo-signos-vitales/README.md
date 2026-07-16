# Microservicio: Monitoreo de Signos Vitales

API estática del proyecto CuidaBien para registrar y consultar mediciones de adultos mayores. Usa almacenamiento en memoria y carga datos de demostración al iniciar; no requiere base de datos ni autenticación para la casa abierta.

## Puerto

Se inicia en el puerto **8083**. Puede cambiarse con la variable de entorno `PORT`.

## Endpoints

| Método | Ruta | Descripción |
| --- | --- | --- |
| GET | `/health` | Estado del servicio |
| POST | `/api/vitales` | Crea y evalúa un registro |
| GET | `/api/vitales/{id_adulto_mayor}` | Historial, descendente por fecha |
| GET | `/api/vitales/{id_adulto_mayor}/ultimo` | Última medición |
| GET | `/api/vitales/{id_adulto_mayor}/tendencia?parametro=frecuencia_cardiaca&dias=30` | Serie para gráficos |

Los parámetros de tendencia admitidos son `presion_sistolica`, `presion_diastolica`, `frecuencia_cardiaca`, `temperatura`, `saturacion_oxigeno` y `nivel_glucosa`.

### Ejemplo de creación

```json
{
  "id_adulto_mayor": "1",
  "registrado_por": "Cuidador Ana",
  "presion_sistolica": 122,
  "presion_diastolica": 78,
  "frecuencia_cardiaca": 72,
  "temperatura": 36.7,
  "saturacion_oxigeno": 97,
  "nivel_glucosa": 108,
  "nivel_dolor": 2,
  "observaciones": "Control de rutina"
}
```

La respuesta incluye los valores y su evaluación (`normal`, `bajo`, `alto` o `critico`) junto con el estado general. Los rangos son los del prototipo entregado; el catálogo parametrizable queda fuera del alcance de esta versión estática.

## Ejecución

```bash
cd services/monitoreo-signos-vitales
go run ./cmd/API
```

## Persistencia

Los datos se almacenan únicamente en memoria y se pierden al detener el servicio. Esto es intencional para la demostración estática de Casa Abierta.
