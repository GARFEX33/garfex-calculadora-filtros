# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Garfex Calculadora Filtros** - Backend API para calcular memorias de cálculo de instalaciones eléctricas según normativa NOM (México).

### Stack Tecnológico
- **Backend:** Go 1.22+, Gin (framework web)
- **Base de datos:** PostgreSQL (Supabase en Docker Compose, Ubuntu)
- **Driver BD:** pgx/v5 (conexión directa a PostgreSQL)
- **Tablas de referencia:** CSV (normativa NOM)
- **Testing:** testing (stdlib) + testify
- **Linting:** golangci-lint

### Frontend
El frontend está en un **repositorio separado** (arquitectura desacoplada). Este repo es solo el backend.

## Arquitectura

Este proyecto usa **Arquitectura Hexagonal / Clean Architecture** con separación estricta de capas:

```
garfex-calculadora-filtros/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada, inyección de dependencias
│
├── internal/
│   ├── domain/
│   │   ├── entity/              # Equipo, FiltroActivo, FiltroRechazo, MemoriaCalculo
│   │   ├── valueobject/         # Corriente, Tension, Conductor (inmutables)
│   │   └── service/             # Lógica de cálculos eléctricos (6 servicios)
│   │
│   ├── application/
│   │   ├── port/                # Interfaces (EquipoRepository, TablaNOMRepository)
│   │   ├── usecase/             # CalcularMemoriaUseCase (orquesta servicios)
│   │   └── dto/                 # EquipoInput, MemoriaOutput
│   │
│   ├── infrastructure/
│   │   ├── repository/          # PostgresEquipoRepository, CSVTablaNOMRepository
│   │   └── client/              # PostgresClient (pgx pool)
│   │
│   └── presentation/
│       ├── handler/             # Gin handlers (HTTP ↔ Application)
│       ├── middleware/          # CORS, logging, recovery
│       └── router.go            # Setup de rutas Gin
│
├── data/
│   └── tablas_nom/              # Tablas NOM en CSV
│
└── tests/                       # Tests por capa
```

### Principios de Diseño

1. **Domain sin dependencias externas** - Sin Gin, sin pgx, sin ninguna librería externa
2. **Domain no conoce NOM como archivos** - Solo recibe datos ya interpretados (conductores, valores numéricos). La lectura de CSV ocurre en `infrastructure/`, nunca en `domain/`
3. **Application** define contratos (ports/interfaces), no implementaciones
4. **Infrastructure** implementa los ports con tecnologías específicas (pgx, CSV)
5. **Presentation** solo adapta HTTP ↔ Application (Gin handlers)
6. **Accept interfaces, return structs** - Funciones aceptan interfaces, retornan tipos concretos
7. **Interfaces definidas donde se consumen** - Los ports viven en `application/port/`
8. **Inyección de dependencias manual** - Se construye todo en `cmd/api/main.go`, sin frameworks DI
9. **YAGNI:** solo implementar lo necesario para la fase actual

## Non-Goals (Fase 1)

Lo que este sistema **no hará** en la fase actual:
- No genera PDF (Fase 3)
- No maneja múltiples usuarios ni autenticación
- No es multi-tenant
- No optimiza rendimiento (sin caché, sin pooling avanzado)
- No valida normas eléctricas más allá de NOM básica
- No gestiona proyectos ni historial de cálculos
- No tiene frontend (repositorio separado, Fase 3)

## Estructura del Proyecto

Ver documento de diseño completo: `docs/plans/2026-02-09-arquitectura-inicial-design.md`

### Tablas NOM (CSV)
Las tablas de referencia de normativa NOM están en `data/tablas_nom/`:
- `310-15-b-16.csv` - Conductores en tubería
- `250-122.csv` - Conductores de tierra
- (Más tablas se agregarán según se necesiten)

### Base de Datos (Supabase)

Tabla `equipos_filtros`:
```sql
- clave: text (unique) -- Identificador del equipo
- tipo: tipo_filtro (enum: ACTIVO, RECHAZO)
- voltaje: integer (V)
- "qn/In": integer -- KVAR para Filtros de Rechazo, Amperaje para Filtros Activos
- itm: integer -- Interruptor termomagnético de fábrica
- bornes: smallint
```

## Topología de Red

```
Servidor Ubuntu (192.168.1.X)
  ├── Supabase (Docker Compose) → PostgreSQL :5432
  ├── Backend Go (producción)   → conecta a localhost:5432
  └── Cloudflare Tunnel         → expone API REST al exterior

Laptop Windows (misma red WiFi - desarrollo)
  └── Backend Go (dev)          → conecta a 192.168.1.X:5432
```

## Comandos de Desarrollo

### Setup
```bash
# Instalar dependencias
go mod tidy

# Configurar variables de entorno
cp .env.example .env
# Editar .env con DB_HOST (IP del servidor Ubuntu en dev, localhost en prod)
```

### Desarrollo
```bash
# Levantar servidor de desarrollo
go run cmd/api/main.go

# Correr tests
go test ./...

# Linting
golangci-lint run
```

### Variables de Entorno
```bash
# Desarrollo (laptop Windows - misma red WiFi)
DB_HOST=192.168.1.X
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=...
DB_NAME=postgres

# Producción (servidor Ubuntu)
DB_HOST=localhost
DB_PORT=5432
```

### Testing
- Tests unitarios del domain (sin dependencias externas, solo stdlib de Go)
- Tests de integración de repositories (mocks/interfaces)
- Tests de API (httptest de Go + Gin)

## Flujo de Trabajo

### Desarrollo Paso a Paso
Este proyecto se desarrolla **incrementalmente**:
1. **Fase 1 (actual):** 2 tipos de equipos (FA, FR), 6 servicios de cálculo, 2 tablas NOM
2. **Fase 2:** Más tipos de equipos, más tablas NOM
3. **Fase 3:** Generación de PDF, frontend (repo separado)

**IMPORTANTE:** No adelantarse. Implementar solo lo necesario para la fase actual.

### Reglas de Negocio Clave

#### Tipos de Equipos
- **Filtro Activo (FA):** Corriente nominal = amperaje directo (no se calcula)
- **Filtro de Rechazo (FR):** Corriente nominal = `I = KVAR / (KV × √3)`

#### Pasos de Cálculo
1. **Corriente Nominal:** Calcular In según tipo de equipo
2. **Ajuste de Corriente:** Aplicar factores (temperatura, agrupamiento, etc.)
3. **Conductor de Alimentación:** Seleccionar calibre de tabla NOM, considerar hilos por fase
4. **Conductor de Tierra:** Usar ITM del equipo, tabla 250-122
5. **Canalización:** Tubería o charola según área de conductores
6. **Caída de Tensión:** Validar límites NOM (3% o 5%)

## Convenciones de Código

### Idioma y Nomenclatura
- **Nombres de negocio en español** (`MemoriaCalculo`, `CorrienteNominal`, `FiltroRechazo`)
- **Código Go en inglés idiomático** (nombres de packages, variables internas)
- **Packages:** cortos, minúsculas, sin guiones bajos (`entity`, `usecase`, `valueobject`)
- **Constructores:** `NewXxx()` para structs con validación

### Formato y Estilo
- `gofmt` y `goimports` en todo el código
- Retorno temprano: manejar errores primero, mantener el happy path sin indentar
- Receptores consistentes: usar punteros o valores, no mezclar por tipo

### Errores
- Siempre retornar `error` como segundo valor, nunca ignorar con `_`
- Envolver errores con contexto: `fmt.Errorf("calcular corriente: %w", err)`
- Errores de dominio como tipos propios: `ErrEquipoNoEncontrado`, `ErrCorrienteInvalida`
- Usar `errors.Is` y `errors.As` para verificar tipos de error

### Context
- `context.Context` como primer parámetro en todas las operaciones I/O
- Nunca almacenar context en structs

### Interfaces
- Definidas en `application/port/`, implementadas en `infrastructure/`
- Pequeñas y enfocadas (pocos métodos por interface)
- Go implementa interfaces implícitamente (duck typing)

### Testing
- Tests con `go test -race ./...` para detectar race conditions
- Table-driven tests con `t.Run` para subtests
- Tests del domain sin dependencias externas (solo stdlib)

## Variables de Entorno

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=...
DB_NAME=postgres
ENVIRONMENT=development|production
```

## API Endpoints

### Health Check
```
GET /health
Response 200: { "status": "ok" }
```

### Calcular Memoria
```
POST /api/v1/calculos/memoria
Body: EquipoInput (modo, datos del equipo, parámetros de instalación)
Response 200: MemoriaOutput (resultados de todos los pasos)
Response 400: Error de validación de input
Response 404: Equipo no encontrado (si modo=LISTADO)
Response 422: Error de cálculo (datos insuficientes)
Response 500: Error interno
```

### Listar Equipos
```
GET /api/v1/equipos?tipo=ACTIVO&min_capacidad=100&max_capacidad=300
Response 200: []Equipo
```

### Obtener Equipo
```
GET /api/v1/equipos/{clave}
Response 200: Equipo
Response 404: Equipo no encontrado
```

### Formato de Errores (consistente en todos los endpoints)
```json
{
  "error": "descripción del error",
  "code": "EQUIPO_NO_ENCONTRADO",
  "details": "información adicional opcional"
}
```

### Política de Versionado de API

- **URL versioning:** `/api/v1/` (estrategia actual)
- **Non-breaking changes** (no requieren nueva versión): agregar campos opcionales en responses, agregar nuevos endpoints, agregar query params opcionales
- **Breaking changes** (requieren `/api/v2/`): eliminar campos, cambiar tipos, modificar contratos existentes, cambiar semántica de endpoints
- La versión actual `v1` se mantendrá estable durante toda la Fase 1 y Fase 2

### Errores: Domain vs Presentation

| Capa | Tipo de error | Ejemplo |
|---|---|---|
| `domain/entity` | Errores de reglas de negocio | `ErrCorrienteInvalida`, `ErrVoltajeInvalido` |
| `domain/service` | Errores de cálculo | `ErrDivisionPorCero` |
| `application/usecase` | Errores de flujo | `ErrEquipoNoEncontrado`, `ErrModoInvalido` |
| `presentation/handler` | **Traduce** errores a HTTP | `domain.ErrEquipoNoEncontrado` → `404 JSON` |

**Regla:** Los handlers son los únicos que conocen HTTP status codes. El domain y application solo retornan `error`.

## Buenas Prácticas Go (Obligatorias)

### Graceful Shutdown
El servidor Gin debe manejar señales `SIGINT`/`SIGTERM` para cerrar conexiones limpiamente antes de terminar.

### Golangci-lint
Configurar `.golangci.yml` con al menos: `errcheck`, `govet`, `staticcheck`, `gofmt`, `goimports`.

### Comandos de Verificación
```bash
go build ./...          # Sin errores de compilación
go vet ./...            # Análisis estático básico
go test -race ./...     # Tests con race detector
golangci-lint run       # Linting completo
```

### Anti-patrones a Evitar
- No usar `panic` para manejo de errores de negocio
- No almacenar `context.Context` en structs
- No mezclar receptores valor/puntero en el mismo tipo
- No usar estado global mutable (`var db *pgxpool.Pool` a nivel de package)

## Contacto y Documentación

- Diseño de arquitectura: `docs/plans/2026-02-09-arquitectura-inicial-design.md`
- Issues y mejoras: GitHub Issues
