# HTTP Handlers

Handlers HTTP para la API REST.

## Handlers

| Handler | Endpoint | Descripción |
|---------|----------|-------------|
| `CalculoHandler` | POST /api/v1/calculos/memoria | Memoria de cálculo |
| `CalculoHandler` | POST /api/v1/calculos/amperaje | Cálculo rápido de amperaje |

## Subpaquetes

| Paquete | Descripción |
|---------|-------------|
| `formatters/` | Formateo de respuestas |
| `middleware/` | Middleware HTTP |

## Endpoints

```bash
# Memoria de cálculo
POST /api/v1/calculos/memoria

# Amperaje rápido
POST /api/v1/calculos/amperaje
```

## Errores

Retorna errores con código HTTP apropiado:
- 400: Bad Request
- 500: Internal Server Error
