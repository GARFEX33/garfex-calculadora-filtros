# Presentation Layer

Adapta HTTP <-> Application. Solo Gin handlers y middleware.

> **Skills Reference**:
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — error handling, convenciones de handlers
> - [`api-design-principles`](.agents/skills/api-design-principles/SKILL.md) — diseño REST, convenciones de endpoints

### Auto-invoke

| Accion | Skill |
|--------|-------|
| Crear o modificar handler | `golang-patterns` |
| Diseñar nuevo endpoint | `api-design-principles` |
| Definir formato de respuesta o error | `api-design-principles` |

## Estructura

- `handler/` — Gin handlers para cada endpoint
- `middleware/` — CORS, logging, recovery
- `router.go` — Setup de rutas Gin

## API Endpoints

### Fase 1 (implementados)

```
GET  /health                                          -> 200 {"status": "ok"}
POST /api/v1/calculos/memoria                         -> MemoriaOutput
```

#### POST /api/v1/calculos/memoria

**Request Body** (campos obligatorios marcados con *):

```json
{
  "modo": "MANUAL_AMPERAJE",           // * LISTADO | MANUAL_AMPERAJE | MANUAL_POTENCIA
  "amperaje_nominal": 100,             // MANUAL_AMPERAJE
  "potencia_nominal": 50,              // MANUAL_POTENCIA (KVAR o KVA)
  "tension": 220,                      // * Voltaje
  "factor_potencia": 0.9,              // Opcional, default 1.0
  "itm": 200,                          // * Capacidad del interruptor (A)
  "tipo_canalizacion": "TUBERIA_PVC",  // * Tipo de canalización
  "hilos_por_fase": 1,                 // Opcional, default 1
  "longitud_circuito": 50,             // * Metros
  "material": "Cu",                    // Opcional: "Cu" | "Al", default "Cu"
  "estado": "Sonora",                  // * Estado de México (para temperatura ambiente)
  "sistema_electrico": "DELTA"         // * DELTA | ESTRELLA | BIFASICO | MONOFASICO
}
```

**Campo `material`:**
- Opcional, default: `"Cu"` (cobre)
- Valores aceptados: `"Cu"`, `"cu"`, `"Al"`, `"al"` (case-insensitive)
- Afecta selección de conductor de alimentación, tierra e impedancia

### Fase 2 (pendientes)

```
GET  /api/v1/equipos?tipo=&min_capacidad=&max_capacidad= -> []Equipo
GET  /api/v1/equipos/{clave}                          -> Equipo
```

## Formato de Errores (consistente en todos los endpoints)

```json
{"success": false, "error": "descripcion", "code": "EQUIPO_NO_ENCONTRADO", "details": "..."}
```

El campo `success` es `false` para errores, `true` para respuestas exitosas (donde el campo `data` contiene el resultado).

## Mapeo Domain -> HTTP

| Error domain | HTTP status |
|---|---|
| ErrEquipoNoEncontrado | 404 |
| ErrModoInvalido | 400 |
| ErrCanalizacionNoSoportada | 400 |
| Validacion de input | 400 |
| Error de calculo (datos insuficientes) | 422 |
| ErrConductorNoEncontrado | 422 |
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

---

## CRITICAL RULES

### Handlers
- ALWAYS: Handler = bind input → llamar use case → traducir resultado a HTTP. Nada mas.
- ALWAYS: Mapear errores de domain/application a HTTP status usando la tabla de arriba
- NEVER: Logica de negocio en handlers
- NEVER: Importar `domain/service` directamente — solo `application/usecase`

### Errores HTTP
- ALWAYS: Formato consistente `{"error": "...", "code": "...", "details": "..."}`
- ALWAYS: Solo los handlers conocen HTTP status codes
- NEVER: Retornar stack traces o errores internos en responses de produccion

### Versionado
- ALWAYS: Nuevos endpoints en `/api/v1/` durante Fase 1 y 2
- NEVER: Cambios breaking en v1 sin crear v2

---

## NAMING CONVENTIONS

| Entidad | Patron | Ejemplo |
|---------|--------|---------|
| Handler struct | `PascalCaseHandler` | `CalculoHandler`, `EquipoHandler` |
| Metodo handler | `VerboPascalCase` | `CrearMemoria`, `ObtenerEquipo` |
| Archivo handler | `snake_case_handler.go` | `calculo_handler.go` |
| Middleware | `snake_case.go` | `cors.go`, `logging.go` |

---

## QA CHECKLIST

- [ ] `go test ./internal/presentation/...` pasa
- [ ] `go vet ./internal/presentation/...` pasa
- [ ] Todos los errores retornan formato `{"error", "code", "details"}`
- [ ] Nuevo endpoint documentado en la tabla de API Endpoints
- [ ] Sin logica de negocio en handlers
- [ ] Mapeo de errores domain → HTTP status cubierto
