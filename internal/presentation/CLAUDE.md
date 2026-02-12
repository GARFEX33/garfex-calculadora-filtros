# Presentation Layer

Adapta HTTP <-> Application. Solo Gin handlers y middleware.

## Estructura

- `handler/` — Gin handlers para cada endpoint
- `middleware/` — CORS, logging, recovery
- `router.go` — Setup de rutas Gin

## API Endpoints

```
GET  /health                                          -> 200 {"status": "ok"}
POST /api/v1/calculos/memoria                         -> MemoriaOutput
GET  /api/v1/equipos?tipo=&min_capacidad=&max_capacidad= -> []Equipo
GET  /api/v1/equipos/{clave}                          -> Equipo
```

## Formato de Errores (consistente en todos los endpoints)

```json
{"error": "descripcion", "code": "EQUIPO_NO_ENCONTRADO", "details": "..."}
```

## Mapeo Domain -> HTTP

| Error domain | HTTP status |
|---|---|
| ErrEquipoNoEncontrado | 404 |
| Validacion de input | 400 |
| Error de calculo (datos insuficientes) | 422 |
| Error interno | 500 |

**Regla:** Solo los handlers conocen HTTP status codes. El domain y application solo retornan `error`.

## Versionado de API

- URL versioning: `/api/v1/`
- **Non-breaking** (no requieren nueva version): agregar campos opcionales, nuevos endpoints, query params opcionales
- **Breaking** (requieren `/api/v2/`): eliminar campos, cambiar tipos, modificar contratos
- v1 se mantiene estable durante Fase 1 y Fase 2

## Convenciones

- Handlers son delegadores: llaman use case, traducen resultado a HTTP
- Graceful shutdown: manejar SIGINT/SIGTERM para cerrar conexiones limpiamente
- No logica de negocio en handlers
- No importar domain/service directamente — solo application/usecase
