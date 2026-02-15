# Infrastructure Layer - calculos

Implementa los ports definidos en `application/port/`.
Tecnologias: CSV (encoding/csv), HTTP (Gin).

> **Skills Reference**:
>
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — error handling, interfaces, convenciones de repositorios
> - [`golang-pro`](.agents/skills/golang-pro/SKILL.md) — connection pooling, concurrencia en queries

### Auto-invoke

| Accion                            | Skill             |
| --------------------------------- | ----------------- |
| Crear o modificar repositorio     | `golang-patterns` |
| Configurar cliente HTTP o BD       | `golang-pro`      |
| Implementar nuevo CSV reader      | `golang-patterns` |

## Estructura

```
infrastructure/
├── adapter/
│   ├── driven/
│   │   └── csv/           # CSVTablaNOMRepository, SeleccionarTemperaturaRepository
│   └── driver/
│       └── http/
│           ├── formatters/  # NombreTablaAmpacidad, GenerarObservaciones
│           └── middleware/  # CorsMiddleware, RequestLogger
├── router.go              # Configuración de rutas Gin
└── AGENTS.md
```

## Driven Adapters (CSV)

### CSVTablaNOMRepository

Lee tablas NOM desde CSV files con in-memory caching.

#### Mapeo canalizacion → tabla ampacidad

| TipoCanalizacion                             | Archivo CSV     |
| -------------------------------------------- | --------------- |
| TUBERIA_PVC / ALUMINIO / ACERO_PG / ACERO_PD | 310-15-b-16.csv |
| CHAROLA_CABLE_ESPACIADO                      | 310-15-b-17.csv |
| CHAROLA_CABLE_TRIANGULAR                     | 310-15-b-20.csv |

#### Mapeo canalizacion → columna R (Tabla 9)

| TipoCanalizacion                                     | Columna resistencia    |
| ---------------------------------------------------- | ---------------------- |
| TUBERIA_PVC / CHAROLA_ESPACIADO / CHAROLA_TRIANGULAR | `res_{material}_pvc`   |
| TUBERIA_ALUMINIO                                     | `res_{material}_al`    |
| TUBERIA_ACERO_PG / ACERO_PD                          | `res_{material}_acero` |

#### Tabla 250-122 (Conductor de Tierra)

Formato CSV con columnas Cu + Al:
```csv
itm_hasta,cu_calibre,cu_seccion_mm2,al_calibre,al_seccion_mm2
```

### SeleccionarTemperaturaRepository

Implementa `SeleccionarTemperaturaPort` delegando al servicio de dominio.

## Driver Adapters (HTTP)

### CalculoHandler

Maneja los endpoints de cálculo:
- `POST /api/v1/calculos/memoria` → MemoriaOutput

#### Mapeo errores domain → HTTP

| Error domain | HTTP status |
|--------------|-------------|
| ErrModoInvalido | 400 |
| ErrCanalizacionNoSoportada | 400 |
| ErrEquipoInputInvalido | 400 |
| Validacion de input | 400 |
| ErrConductorNoEncontrado | 422 |
| ErrCanalizacionNoDisponible | 422 |
| Error interno | 500 |

### Formatters

- `NombreTablaAmpacidad` - Genera nombre descriptivo de tabla NOM
- `GenerarObservaciones` - Genera observaciones sobre el cálculo

### Middleware

- `CorsMiddleware` - CORS para desarrollo
- `RequestLogger` - Logging de peticiones

## Variables de Entorno

`DATA_PATH` - Path a los archivos CSV de tablas NOM

---

## CRITICAL RULES

### General
- ALWAYS: Implementar exactamente el port definido en `application/port/`
- ALWAYS: `context.Context` como primer parámetro en todas las operaciones
- ALWAYS: Inyección de dependencias via constructor — sin globals ni singletons
- NEVER: Importar `domain/service` — solo `entity` y `valueobject`
- NEVER: Lógica de negocio en adapters — solo traducción datos <-> domain

### CSV Repository
- ALWAYS: Validar columnas requeridas al cargar CSV — fallar rápido si falta columna
- NEVER: Escribir datos NOM hardcodeados en Go — siempre leer del CSV

### HTTP Handler
- ALWAYS: Handler = bind input → llamar use case → traducir resultado a HTTP. Nada más.
- ALWAYS: Mapear errores de domain/application a HTTP status
- NEVER: Lógica de negocio en handlers

---

## NAMING CONVENTIONS

| Entidad                | Patrón                         | Ejemplo                         |
| ---------------------- | ------------------------------ | ------------------------------- |
| Repositorio CSV        | `CSVPascalCaseRepository`      | `CSVTablaNOMRepository`         |
| Handler struct        | `PascalCaseHandler`            | `CalculoHandler`                |
| Archivo handler       | `snake_case_handler.go`        | `calculo_handler.go`            |
| Middleware            | `PascalCaseMiddleware`         | `CorsMiddleware`                |
| Formatter             | `PascalCaseFormatter`          | `NombreTablaAmpacidad`          |

---

## QA CHECKLIST

- [ ] `go test ./internal/calculos/infrastructure/...` pasa
- [ ] Nuevo repositorio implementa el port completo
- [ ] Sin estado global mutable
- [ ] Tests de handler pasan
- [ ] Sin lógica de negocio en adapters
