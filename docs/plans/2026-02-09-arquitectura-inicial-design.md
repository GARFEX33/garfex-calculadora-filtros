# Diseño de Arquitectura Inicial - Garfex Calculadora Filtros

**Fecha:** 2026-02-09
**Versión:** 1.2 (Buenas prácticas Go + API design aplicadas)
**Autor:** Diseño colaborativo con usuario

---

## 1. Resumen Ejecutivo

Sistema backend en Go/Gin para calcular memorias de cálculo de instalaciones eléctricas según normativa NOM (México). Arquitectura hexagonal/clean architecture para desacoplamiento total entre frontend, backend y base de datos.

**Fase 1 (MVP):**
- 2 tipos de equipos: Filtros Activos (FA) y Filtros de Rechazo (FR/KVAR)
- 3 formas de entrada: selección desde listado (Supabase), entrada manual por amperaje, entrada manual por potencia
- 4 pasos de cálculo: conductor de alimentación, conductor de tierra, canalización, caída de tensión

---

## 1.1 Non-Goals (Fase 1)

Lo que este sistema **no hará** en la fase actual:
- No genera PDF (Fase 3)
- No maneja autenticación ni múltiples usuarios
- No es multi-tenant
- No optimiza rendimiento (sin caché, sin pooling avanzado)
- No valida normas eléctricas más allá de NOM básica
- No gestiona proyectos ni historial de cálculos
- No tiene frontend (repositorio separado, Fase 3)

---

## 2. Stack Tecnológico

- **Backend:** Go 1.22+, Gin (framework web)
- **Base de datos:** PostgreSQL (Supabase en Docker Compose, Ubuntu)
- **Driver BD:** pgx/v5 (conexión directa a PostgreSQL)
- **Tablas de referencia:** CSV (normativa NOM)
- **Testing:** testing (stdlib) + testify
- **Linting:** golangci-lint

### Topología de Red

```
Servidor Ubuntu (192.168.1.X)
  ├── Supabase (Docker Compose) → PostgreSQL :5432
  ├── Backend Go (producción)   → conecta a localhost:5432
  └── Cloudflare Tunnel         → expone API REST al exterior

Laptop Windows (misma red WiFi - desarrollo)
  └── Backend Go (dev) → conecta a 192.168.1.X:5432
```

### Variables de Entorno por Ambiente

```bash
# Desarrollo (laptop Windows)
DB_HOST=192.168.1.X
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=...
DB_NAME=postgres

# Producción (servidor Ubuntu)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=...
DB_NAME=postgres
```

---

## 3. Arquitectura Hexagonal

```
┌─────────────────────────────────────────────────────────┐
│                    PRESENTATION                          │
│  (Gin handlers, JSON, HTTP)                             │
│                                                          │
│  POST /api/v1/calculos/memoria                          │
│  GET  /api/v1/equipos                                   │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   APPLICATION                            │
│  (Use Cases, Ports/Interfaces, DTOs)                    │
│                                                          │
│  CalcularMemoriaUseCase                                 │
│  ├─ EquipoRepository (interface)                        │
│  └─ TablaNOMRepository (interface)                      │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                      DOMAIN                              │
│  (Entities, Value Objects, Services)                    │
│                                                          │
│  Entities: Equipo, FiltroActivo, FiltroRechazo         │
│  Value Objects: Corriente, Tension, Conductor          │
│  Services: Cálculos eléctricos (6 servicios)           │
└─────────────────────────────────────────────────────────┘
                     ▲
                     │
┌────────────────────┴────────────────────────────────────┐
│                 INFRASTRUCTURE                           │
│  (Repositories, Clients)                                │
│                                                          │
│  PostgresEquipoRepository (implementa EquipoRepository) │
│  CSVTablaNOMRepository (implementa TablaNOMRepository)  │
│  PostgresClient (pgx connection pool)                   │
└─────────────────────────────────────────────────────────┘
```

---

## 4. Estructura de Carpetas (Go)

```
garfex-calculadora-filtros/
├── cmd/
│   └── api/
│       └── main.go                  # Punto de entrada
│
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── equipo.go            # Struct base + interface CalculadorCorriente
│   │   │   ├── filtro_activo.go     # Embeds Equipo, implementa interfaz
│   │   │   ├── filtro_rechazo.go    # Embeds Equipo, implementa interfaz
│   │   │   └── memoria_calculo.go   # Resultado de los pasos
│   │   │
│   │   ├── valueobject/
│   │   │   ├── corriente.go         # Valor inmutable con validación
│   │   │   ├── tension.go           # Valor inmutable + método en_kv()
│   │   │   └── conductor.go         # Calibre, material, aislamiento
│   │   │
│   │   └── service/
│   │       ├── calculo_corriente_nominal.go  # Paso 1a
│   │       ├── ajuste_corriente.go           # Paso 1b
│   │       ├── calculo_conductor.go          # Paso 2
│   │       ├── calculo_tierra.go             # Paso 3
│   │       ├── calculo_canalizacion.go       # Paso 4
│   │       └── calculo_caida_tension.go      # Paso 5
│   │
│   ├── application/
│   │   ├── port/
│   │   │   ├── equipo_repository.go    # Interface (contrato)
│   │   │   └── tabla_nom_repository.go # Interface (contrato)
│   │   │
│   │   ├── usecase/
│   │   │   └── calcular_memoria.go     # Orquesta los servicios
│   │   │
│   │   └── dto/
│   │       ├── equipo_input.go         # 3 formas de entrada
│   │       └── memoria_output.go       # Resultado estructurado
│   │
│   ├── infrastructure/
│   │   ├── repository/
│   │   │   ├── postgres_equipo_repository.go
│   │   │   └── csv_tabla_nom_repository.go
│   │   │
│   │   └── client/
│   │       └── postgres_client.go      # pgx connection pool
│   │
│   └── presentation/
│       ├── handler/
│       │   ├── calculo_handler.go      # POST /calculos/memoria
│       │   └── equipo_handler.go       # GET /equipos
│       ├── middleware/
│       │   └── cors.go
│       └── router.go                   # Setup de rutas Gin
│
├── data/
│   └── tablas_nom/
│       ├── 310-15-b-16.csv            # Conductores en tubería (Fase 1)
│       └── 250-122.csv                # Conductores de tierra (Fase 1)
│
├── tests/
│   ├── domain/
│   ├── application/
│   └── infrastructure/
│
├── docs/
│   └── plans/
│
├── .env
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

---

## 5. Domain Layer

### 5.1 Entities

En Go no hay clases abstractas. Se usa una **interface** para definir el comportamiento común:

#### Interface `CalculadorCorriente`
```go
type CalculadorCorriente interface {
    CalcularCorrienteNominal() (valueobject.Corriente, error)
}
```

#### Struct `Equipo` (base embebida)
```go
type Equipo struct {
    Clave   string
    Tipo    TipoFiltro  // enum: ACTIVO, RECHAZO
    Voltaje int         // en Voltios
    ITM     int         // interruptor termomagnético de fábrica
    Bornes  int
}
```

#### `FiltroActivo` (implementa CalculadorCorriente)
```go
type FiltroActivo struct {
    Equipo              // embedded
    Amperaje int        // qn/In de la BD
}

// Retorna amperaje directamente (no calcula)
func (fa *FiltroActivo) CalcularCorrienteNominal() (valueobject.Corriente, error)
```

#### `FiltroRechazo` (implementa CalculadorCorriente)
```go
type FiltroRechazo struct {
    Equipo              // embedded
    KVAR int            // qn/In de la BD
}

// Aplica fórmula: I = KVAR / (KV × √3)
// donde KV = Voltaje / 1000
func (fr *FiltroRechazo) CalcularCorrienteNominal() (valueobject.Corriente, error)
```

#### `MemoriaCalculo`
```go
type MemoriaCalculo struct {
    Equipo                 CalculadorCorriente
    CorrienteNominal       valueobject.Corriente
    CorrienteAjustada      valueobject.Corriente
    FactoresAjuste         map[string]float64
    ConductorAlimentacion  valueobject.Conductor
    HilosPorFase           int
    ConductorTierra        valueobject.Conductor
    Canalizacion           Canalizacion
    CaidaTension           float64  // porcentaje
    CumpleNormativa        bool
}
```

### 5.2 Value Objects (inmutables)

#### `Corriente`
```go
type Corriente struct {
    Valor  float64
    Unidad string  // "A"
}
// Validación: Valor > 0
// Construcción solo via constructor: NewCorriente(valor float64) (Corriente, error)
```

#### `Tension`
```go
type Tension struct {
    Valor  int
    Unidad string  // "V"
}
// Validación: voltajes NOM válidos (127, 220, 440, 480, etc.)
// Método: EnKilovoltios() float64
```

#### `Conductor`
```go
type Conductor struct {
    Calibre          string  // ej: "12 AWG", "1/0 AWG"
    Material         string  // "Cu" o "Al"
    TipoAislamiento  string  // "THHN", "THW", etc.
}
// Validación: calibres válidos NOM
```

### 5.3 Services

#### `CalculoCorrienteNominalService`
```go
func CalcularCorrienteNominal(equipo CalculadorCorriente) (Corriente, error)
  // Delega al método de la entidad
  // FA → amperaje directo
  // FR → I = KVAR / (KV × √3)
```

#### `AjusteCorrienteService`
```go
func AjustarCorriente(cn Corriente, factores map[string]float64) (Corriente, error)
  // Aplica factores: temperatura, agrupamiento, etc.
```

#### `CalculoConductorService`
```go
func CalcularConductor(ca Corriente, hilosPorFase int, repo port.TablaNOMRepository) (Conductor, error)
  // Divide corriente si hilosPorFase > 1
  // Consulta tabla NOM
```

#### `CalculoTierraService`
```go
func CalcularConductorTierra(itm int, repo port.TablaNOMRepository) (Conductor, error)
  // Usa ITM del equipo → tabla 250-122
```

#### `CalculoCanalizacionService`
```go
func CalcularCanalizacion(conductores []Conductor, tipo string, repo port.TablaNOMRepository) (Canalizacion, error)
  // tipo: "TUBERIA" | "CHAROLA"
```

#### `CalculoCaidaTensionService`
```go
func CalcularCaidaTension(conductor Conductor, corriente Corriente, distancia float64, tension Tension) (float64, error)
  // Retorna porcentaje
  // Valida límites NOM (3% o 5%)
```

---

## 6. Application Layer

### 6.1 Ports (Interfaces en Go)

#### `EquipoRepository`
```go
type EquipoRepository interface {
    ObtenerPorClave(ctx context.Context, clave string) (*entity.Equipo, error)
    ListarPorTipo(ctx context.Context, tipo entity.TipoFiltro) ([]*entity.Equipo, error)
    ListarTodos(ctx context.Context) ([]*entity.Equipo, error)
    FiltrarPorRangoCapacidad(ctx context.Context, min, max int, tipo *entity.TipoFiltro) ([]*entity.Equipo, error)
}
```

#### `TablaNOMRepository`
```go
type TablaNOMRepository interface {
    CargarTabla(nombre string) ([][]string, error)
    ObtenerConductorPorCorriente(corriente float64, tabla string, temperatura int, material string) (entity.Conductor, error)
    ObtenerConductorTierra(itm int) (entity.Conductor, error)
    ObtenerCanalizacion(areaConductores float64, tipo string) (string, error)
}
```

### 6.2 Use Case Principal

#### `CalcularMemoriaUseCase`

**Input:** `EquipoInput`
```go
type EquipoInput struct {
    Modo        string      // "LISTADO" | "AMPERAJE" | "POTENCIA"
    ClaveEquipo *string
    Tipo        *TipoFiltro
    Voltaje     *int
    Amperaje    *int
    Potencia    *int        // KVAR
    ITM         *int
    Bornes      *int
    HilosPorFase    int
    Distancia       float64
    FactoresAjuste  map[string]float64
}
```

**Flujo:**
1. Obtener/construir entidad `Equipo` según `Modo`
2. Calcular corriente nominal
3. Ajustar corriente
4. Calcular conductor de alimentación (consulta tabla NOM)
5. Calcular conductor de tierra (consulta tabla NOM)
6. Calcular canalización (consulta tabla NOM)
7. Calcular caída de tensión
8. Construir y retornar `MemoriaOutput`

**Output:** `MemoriaOutput`
```go
type MemoriaOutput struct {
    Equipo                map[string]interface{}
    CorrienteNominal      float64
    CorrienteAjustada     float64
    FactoresAplicados     map[string]float64
    ConductorAlimentacion map[string]string
    HilosPorFase          int
    ConductorTierra       map[string]string
    Canalizacion          map[string]string
    CaidaTension          float64
    CumpleNormativa       bool
}
```

---

## 7. Infrastructure Layer

### 7.1 PostgresClient (pgx)
```go
type PostgresClient struct {
    pool *pgxpool.Pool  // privado: acceso solo via métodos
}

func NewPostgresClient(ctx context.Context) (*PostgresClient, error)
  // Lee DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME del entorno
  // Retorna pool de conexiones

func (c *PostgresClient) Pool() *pgxpool.Pool
  // Acceso controlado al pool
```

### 7.2 PostgresEquipoRepository
```go
type PostgresEquipoRepository struct {
    client *PostgresClient
}
// Implementa port.EquipoRepository
// Mapea rows de PostgreSQL a entidades del domain
// Interpreta "qn/In" según tipo:
//   ACTIVO → FiltroActivo{Amperaje: qnIn}
//   RECHAZO → FiltroRechazo{KVAR: qnIn}
```

### 7.3 CSVTablaNOMRepository

**Fase 1 - Archivos iniciales:**
```
data/tablas_nom/
  ├── 310-15-b-16.csv    # Conductores en tubería
  └── 250-122.csv        # Conductores de tierra
```

**Futuro (escalable):**
```
data/tablas_nom/
  ├── 310-15-b-20.csv    # Triangular en charola
  ├── 310-15-b-17.csv    # Cable con espacio
  ├── 310-60-69.csv      # Individual +2000V Cu
  └── ...
```

---

## 8. Presentation Layer

### 8.1 Router (Gin)
```go
// GET  /health                         → health check
// GET  /api/v1/equipos                 → listar equipos
// GET  /api/v1/equipos/:clave          → obtener equipo
// POST /api/v1/calculos/memoria        → calcular memoria
```

### 8.2 Handlers

#### `CalcuHandler`
```go
// POST /api/v1/calculos/memoria
// 1. Bind JSON → EquipoInput
// 2. Llama use case
// 3. Retorna MemoriaOutput como JSON
// Status codes:
//   200 → cálculo exitoso
//   400 → input inválido (binding error)
//   404 → equipo no encontrado (modo LISTADO)
//   422 → error de cálculo (datos insuficientes)
//   500 → error interno
```

#### `EquipoHandler`
```go
// GET /api/v1/equipos  (query params: tipo, min_capacidad, max_capacidad)
// → 200: []Equipo
// GET /api/v1/equipos/:clave
// → 200: Equipo | 404: no encontrado
```

### 8.3 Formato de Errores (consistente en todos los endpoints)
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code"`
    Details string `json:"details,omitempty"`
}
// Ejemplo: {"error": "equipo no encontrado", "code": "EQUIPO_NO_ENCONTRADO"}
```

### 8.5 Política de Versionado de API

- **Estrategia:** URL versioning (`/api/v1/`)
- **Non-breaking** (no requieren nueva versión): agregar campos opcionales en responses, nuevos endpoints, query params opcionales
- **Breaking** (requieren `/api/v2/`): eliminar campos, cambiar tipos, modificar contratos existentes
- `v1` se mantiene estable durante Fase 1 y Fase 2

### 8.6 Errores: Domain vs Presentation

| Capa | Tipo de error | Ejemplo |
|---|---|---|
| `domain/entity` | Reglas de negocio | `ErrCorrienteInvalida`, `ErrVoltajeInvalido` |
| `domain/service` | Errores de cálculo | `ErrDivisionPorCero` |
| `application/usecase` | Errores de flujo | `ErrEquipoNoEncontrado`, `ErrModoInvalido` |
| `presentation/handler` | **Traduce** errores a HTTP | `ErrEquipoNoEncontrado` → `404 JSON` |

**Regla:** Solo los handlers conocen HTTP status codes. Domain y application solo retornan `error`.

### 8.4 Inyección de Dependencias (manual en main.go)
```go
// En cmd/api/main.go:
pgClient     := client.NewPostgresClient(ctx)
equipoRepo   := repository.NewPostgresEquipoRepository(pgClient)
tablaNOMRepo := repository.NewCSVTablaNOMRepository("data/tablas_nom")
useCase      := usecase.NewCalcularMemoriaUseCase(equipoRepo, tablaNOMRepo)
calcHandler  := handler.NewCalcuHandler(useCase)
equipHandler := handler.NewEquipoHandler(equipoRepo)
router       := presentation.SetupRouter(calcHandler, equipHandler)

// Graceful shutdown con SIGINT/SIGTERM
```

---

## 9. Flujo de Datos Completo

```
1. Request HTTP  → Gin Router
2. Router        → CalcuHandler (valida JSON → EquipoInput)
3. Handler       → CalcularMemoriaUseCase.Execute(input)
4. Use Case      → Obtiene equipo (PostgresEquipoRepository → pgx → PostgreSQL)
5. Use Case      → Ejecuta servicios del domain
6. Services      → Consultan tablas NOM (CSVTablaNOMRepository → CSV)
7. Use Case      → Construye MemoriaCalculo
8. Use Case      → Retorna MemoriaOutput
9. Handler       → Serializa a JSON
10. Response HTTP → Cliente
```

---

## 10. Dependencias Go (go.mod)

```
module github.com/garfex/calculadora-filtros

go 1.22

require (
    github.com/gin-gonic/gin         v1.10.x  // Framework web
    github.com/jackc/pgx/v5          v5.x.x   // Driver PostgreSQL
    github.com/joho/godotenv         v1.x.x   // Variables de entorno
    github.com/stretchr/testify      v1.x.x   // Testing assertions
)
```

---

## 11. Principios de Diseño

1. **Domain sin dependencias externas** - Sin Gin, sin pgx, sin ninguna librería externa
2. **Domain no conoce NOM como archivos** - Solo recibe datos ya interpretados (conductores, valores numéricos). La lectura de CSV ocurre en `infrastructure/`, nunca en `domain/`
3. **Accept interfaces, return structs** - Funciones aceptan interfaces, retornan tipos concretos
4. **Interfaces definidas donde se consumen** - Los ports viven en `application/port/`, no en `infrastructure/`
5. **Interfaces implícitas de Go** - Duck typing: infrastructure implementa ports sin declararlo explícitamente
6. **Inyección de dependencias manual** - Todo se construye en `cmd/api/main.go`, sin frameworks DI
7. **Errores explícitos** - `error` como segundo valor, nunca ignorar, envolver con `fmt.Errorf("%w", err)`
8. **Errores de dominio tipados** - Sentinel errors: `ErrEquipoNoEncontrado`, `ErrCorrienteInvalida`
9. **Context en todas las operaciones I/O** - Para cancelación y timeouts
10. **No estado global mutable** - Sin variables de package-level, todo via inyección
11. **Receptores consistentes** - Usar punteros o valores, no mezclar en el mismo tipo
12. **Graceful shutdown** - Manejo de señales `SIGINT`/`SIGTERM`
13. **YAGNI** - Solo lo necesario para Fase 1

---

## 12. Buenas Prácticas Go (Obligatorias)

### Convenciones de Nomenclatura
- **Nombres de negocio en español** (`MemoriaCalculo`, `CorrienteNominal`, `FiltroRechazo`)
- **Código Go en inglés idiomático** (nombres de packages, variables internas)
- **Packages:** cortos, minúsculas, sin guiones bajos (`entity`, `usecase`, `valueobject`)
- **Constructores:** `NewXxx()` para structs con validación

### Errores de Dominio
```go
// En internal/domain/entity/errors.go
var (
    ErrEquipoNoEncontrado = errors.New("equipo no encontrado")
    ErrCorrienteInvalida  = errors.New("corriente debe ser mayor que cero")
    ErrVoltajeInvalido    = errors.New("voltaje no válido según normativa NOM")
)
```

### Testing
```bash
go test -race ./...     # Siempre con race detector
go test -cover ./...    # Verificar cobertura
```
- Table-driven tests con `t.Run` para subtests
- Tests del domain sin dependencias externas

### Tooling (obligatorio antes de commit)
```bash
go build ./...
go vet ./...
golangci-lint run
```

### Golangci-lint (`.golangci.yml`)
Linters mínimos: `errcheck`, `govet`, `staticcheck`, `gofmt`, `goimports`.

### Anti-patrones Prohibidos
- `panic` para errores de negocio
- `context.Context` almacenado en structs
- Mezclar receptores valor/puntero en el mismo tipo
- Estado global mutable a nivel de package

---

## 14. Roadmap

### Fase 1 (MVP) - Actual
- 2 tipos de equipos (FA, FR)
- 3 formas de entrada
- 6 servicios de cálculo
- 2 tablas NOM (tubería, tierra)
- API REST con Gin

### Fase 2 (Futuro)
- Más tipos de equipos (motores, contactos, alumbrado)
- Más tablas NOM
- Generación de PDF

### Fase 3 (Futuro)
- Frontend (repo separado)
- Autenticación y usuarios
- Proyectos y múltiples memorias

---

## 15. Conclusión

Arquitectura sólida en Go que permite desarrollo incremental y mantenible.

**Próximo paso:** Crear estructura de carpetas e implementar domain layer comenzando con entidades y value objects.
