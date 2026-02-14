# Diseño: Ports y CSV Repository (Infrastructure)

## Resumen

Diseño de la capa Application (ports/interfaces) y primera implementación de Infrastructure (CSV reader para tablas NOM). PostgreSQL se deja para la siguiente iteración.

## Contexto

El Domain Layer está completo con 6 servicios de cálculo. Los servicios reciben datos pre-resueltos como slices (`[]EntradaTablaConductor`, etc.), no hacen I/O. Los ports definen los contratos que Infrastructure debe implementar para proveer esos datos.

## Decisiones de Diseño

### 1. Ports en Application

Ubicación: `internal/application/port/`

Dos interfaces separadas, pequeñas y enfocadas:

#### TablaNOMRepository

```go
type TablaNOMRepository interface {
    ObtenerTablaAmpacidad(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura) ([]service.EntradaTablaConductor, error)
    ObtenerTablaTierra(ctx context.Context) ([]service.EntradaTablaTierra, error)
    ObtenerImpedancia(ctx context.Context, calibre string, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor) (ResistenciaReactancia, error)
    ObtenerTablaCanalizacion(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]service.EntradaTablaCanalizacion, error)
}
```

**ResistenciaReactancia**: Struct simple con `R float64` y `X float64`.

#### EquipoRepository (para iteración PostgreSQL)

```go
type EquipoRepository interface {
    BuscarPorClave(ctx context.Context, clave string) (entity.Equipo, error)
}
```

### 2. CSV Repository en Infrastructure

Ubicación: `internal/infrastructure/repository/`

**Estrategia de caché en memoria**: Las tablas CSV se cargan UNA VEZ en el constructor y se mantienen en memoria.

- **Ventaja**: Sin I/O de disco en cada request
- **Memoria**: ~50-100KB totales (tablas pequeñas)
- **Simplicidad**: Lookup directo en slices

#### Constructor

```go
func NewCSVTablaNOMRepository(basePath string) (*CSVTablaNOMRepository, error)
```

#### Estructura de archivos

```
internal/infrastructure/repository/
├── csv_tabla_nom_repository.go      # Implementación
├── csv_tabla_nom_repository_test.go # Tests
└── testdata/                        # Copias de CSVs para tests
    ├── 310-15-b-16.csv
    ├── 310-15-b-17.csv
    ├── 310-15-b-20.csv
    ├── 250-122.csv              # A crear
    └── tabla-9-resistencia-reactancia.csv
```

### 3. Mapeos Implementados

#### Tabla Ampacidad por TipoCanalizacion

| TipoCanalizacion | Archivo CSV |
|------------------|-------------|
| TUBERIA_PVC / ALUMINIO / ACERO_PG / ACERO_PD | 310-15-b-16.csv |
| CHAROLA_CABLE_ESPACIADO | 310-15-b-17.csv |
| CHAROLA_CABLE_TRIANGULAR | 310-15-b-20.csv |

#### Columna de Resistencia (Tabla 9)

| TipoCanalizacion | Columna R |
|------------------|-----------|
| TUBERIA_PVC, CHAROLA_* | res_{material}_pvc |
| TUBERIA_ALUMINIO | res_{material}_al |
| TUBERIA_ACERO_PG, ACERO_PD | res_{material}_acero |

#### Columna de Reactancia (Tabla 9)

| TipoCanalizacion | Columna X |
|------------------|-----------|
| PVC, ALUMINIO, CHAROLA_* | reactancia_al |
| ACERO_PG, ACERO_PD | reactancia_acero |

### 4. Tipos Nuevos en Domain

Para que los ports sean type-safe, agregamos enums en `valueobject/`:

```go
// MaterialConductor enum
type MaterialConductor int
const (
    MaterialCobre MaterialConductor = iota
    MaterialAluminio
)

// Temperatura enum  
type Temperatura int
const (
    Temp60 Temperatura = 60
    Temp75 Temperatura = 75
    Temp90 Temperatura = 90
)
```

### 5. DTOs en Application

Ubicación: `internal/application/dto/`

```go
// EquipoInput - entrada del use case
type EquipoInput struct {
    Modo string // LISTADO, MANUAL_AMPERAJE, MANUAL_POTENCIA
    // ... campos según modo
}

// MemoriaOutput - salida completa
type MemoriaOutput struct {
    // Resultados de los 7 pasos
}
```

### 6. Actualización AGENTS.md

Se agregaron errores al mapeo HTTP en `presentation/AGENTS.md`:
- `ErrModoInvalido` → 400
- `ErrCanalizacionNoSoportada` → 400
- `ErrConductorNoEncontrado` → 422

## Out of Scope (próxima iteración)

- PostgreSQL repository
- Use case completo (orquestación)
- Handlers HTTP
- Tabla 250-122.csv (creada en esta iteración)

## QA Checklist

- [ ] Tests de CSV repository con testdata/
- [ ] Validación de cabeceras CSV (fail-fast)
- [ ] Manejo de errores con contexto (línea, columna)
- [ ] Sin estado global (caché en struct, no var package-level)
- [ ] `go test ./internal/infrastructure/...` pasa

## Referencias

- Domain services: `internal/domain/service/`
- CSV tablas: `data/tablas_nom/`
- Application AGENTS.md: `internal/application/AGENTS.md`
- Infrastructure AGENTS.md: `internal/infrastructure/AGENTS.md`
