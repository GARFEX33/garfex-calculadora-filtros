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
Las tablas de referencia de normativa NOM están en `data/tablas_nom/`.

**Tablas de ampacidad (selección de conductor de alimentación — dependen del tipo de canalización):**
- `310-15-b-16.csv` - Tubería conduit (`TUBERIA_PVC/ALUMINIO/ACERO_PG/ACERO_PD`) — 14 AWG a 2000 MCM, Cu/Al 60/75/90°C
- `310-15-b-17.csv` - Charola cable espaciado (`CHAROLA_CABLE_ESPACIADO`) — 14 AWG a 2000 MCM, Cu/Al 60/75/90°C
- `310-15-b-20.csv` - Charola triangular (`CHAROLA_CABLE_TRIANGULAR`) — 8 AWG a 1000 MCM, Cu/Al 75/90°C (sin columna 60°C)

**Tabla de conductor de tierra:**
- `250-122.csv` - Conductores de tierra (independiente del tipo de canalización)

**Tablas de referencia para caída de tensión (método impedancia):**
- `tabla-9-resistencia-reactancia.csv` - Resistencia AC por tipo conduit (PVC/Al/Acero) y reactancia — 14 AWG a 1000 MCM
- `tabla-5-dimensiones-aislamiento.csv` - Diámetro exterior con aislamiento (THW, RHH, XHHW) — para cálculo DMG
- `tabla-8-conductor-desnudo.csv` - Diámetro desnudo + número de hilos — para cálculo RMG

**Formato CSV ampacidad:** `seccion_mm2,calibre,cu_60c,cu_75c,cu_90c,al_60c,al_75c,al_90c` — celdas vacías donde no aplica.

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

### Actualización de Documentación (OBLIGATORIO al terminar cada tarea)

Al completar cualquier tarea de implementación, **antes del commit final**, verificar si hubo cambios que divergen del plan original y actualizar:

| Documento | Cuándo actualizar |
|-----------|------------------|
| `docs/plans/2026-02-10-domain-layer.md` | Cuando la implementación diverge del plan (API diferente, campos extra, decisiones tomadas) |
| `docs/plans/2026-02-09-arquitectura-inicial-design.md` | Cuando cambia el diseño de una entidad, VO o servicio |
| `CLAUDE.md` | Cuando cambia una convención de código o regla de negocio del proyecto |
| `.claude/memory/MEMORY.md` | Cuando hay algo que debe persistir entre sesiones (decisiones, patrones, paths) |

**Regla:** Si la implementación real diverge del plan → marcar la tarea como `✅ COMPLETADO` en el plan y documentar qué cambió y por qué, incluyendo el impacto en tareas futuras.

### Desarrollo Paso a Paso
Este proyecto se desarrolla **incrementalmente**:
1. **Fase 1 (actual):** 4 tipos de equipos (FA, FR, Transformador, Carga), 6 servicios de cálculo, 7 tablas NOM (3 ampacidad + 1 tierra + 3 referencia impedancia)
2. **Fase 2:** Más tipos de equipos, más tablas NOM
3. **Fase 3:** Generación de PDF, frontend (repo separado)

**IMPORTANTE:** No adelantarse. Implementar solo lo necesario para la fase actual.

### Reglas de Negocio Clave

#### Tipos de Equipos
- **Filtro Activo (FA):** Corriente nominal = AmperajeNominal directo (no se calcula)
- **Filtro de Rechazo (FR):** Corriente nominal = `I = KVAR / (KV × √3)`
- **Transformador:** Corriente nominal = `I = KVA / (KV × √3)`
- **Carga:** Corriente nominal = `I = KW / (KV × factor × FP)` donde factor depende de fases (1/2/3)

#### Pasos de Cálculo
1. **Corriente Nominal:** Calcular In según tipo de equipo
2. **Ajuste de Corriente:** Aplicar factores (temperatura, agrupamiento, etc.)
3. **Selección de Canalización:** Elegir `TipoCanalizacion` (6 tipos) — **determina qué tabla NOM usar para ampacidad y qué columna R para impedancia**
4. **Conductor de Alimentación:** Cargar tabla NOM según canalización → auto-seleccionar columna de temperatura (≤100A→60°C, >100A→75°C, override para 90°C) → seleccionar calibre
5. **Conductor de Tierra:** Usar ITM del equipo, tabla 250-122 (independiente de la canalización)
6. **Canalización:** Calcular dimensiones de tubería/charola según área de conductores (40% fill para tubería)
7. **Caída de Tensión (método impedancia NOM):** Z=√(R²+X²), VD=√3×I×Z×L_km — R de Tabla 9, X calculada geométricamente con DMG/RMG (Tablas 5 y 8)

#### Selección de Temperatura (NOM)
- Circuitos ≤ 100 A o calibres 14–1 AWG → columna **60°C**
- Circuitos > 100 A o calibres > 1 AWG → columna **75°C**
- **90°C:** solo con `temperatura_override: 90` explícito (muy raro — requiere todos los equipos certificados 90°C)
- Tabla 310-15(b)(20) Charola triangular **no tiene columna 60°C** → fallback automático a 75°C

#### TipoCanalizacion (enum de dominio — 6 valores)
```go
TUBERIA_PVC              → tabla 310-15-b-16.csv, R: columna pvc
TUBERIA_ALUMINIO         → tabla 310-15-b-16.csv, R: columna al
TUBERIA_ACERO_PG         → tabla 310-15-b-16.csv, R: columna acero
TUBERIA_ACERO_PD         → tabla 310-15-b-16.csv, R: columna acero
CHAROLA_CABLE_ESPACIADO  → tabla 310-15-b-17.csv, R: columna pvc
CHAROLA_CABLE_TRIANGULAR → tabla 310-15-b-20.csv, R: columna pvc
```

#### Caída de Tensión — Método Impedancia
- **R** (resistencia AC): Tabla 9 → columna según material conductor + tipo canalización
- **X** (reactancia inductiva): calculada geométricamente = `(1/n) × 2π×60 × 2×10⁻⁷ × ln(DMG/RMG) × 1000`
- **RMG**: `(diametro_desnudo/2) × factor_hilos` — datos de Tabla 8
- **DMG**: `diametro_exterior_thw × factor_canalizacion` — datos de Tabla 5, factor: tubería/triangular=1.0, espaciado=2.0
- Diseño completo: `docs/plans/2026-02-11-caida-tension-impedancia-design.md`

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
