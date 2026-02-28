# Calculos — Application Layer

Orquesta domain services. Define contratos (ports), no implementaciones.

## Estructura

```
internal/calculos/application/
├── port/           # Interfaces que infrastructure implementa
│   ├── TablaNOMRepository
│   ├── SeleccionarTemperatura
│   └── EquipoRepository
├── usecase/        # Orquestadores
│   └── helpers/    # Funciones auxiliares
└── dto/            # Entrada/salida
```

> **Nota:** Las subcarpetas `port/`, `usecase/` y `dto/` heredan las reglas de este AGENTS.md. No necesitan AGENTS.md propio.

## Dependencias permitidas

- `internal/shared/kernel/valueobject`
- `internal/calculos/domain/entity`
- `internal/calculos/domain/service`
- stdlib de Go

## Dependencias prohibidas

> Ver reglas consolidadas en [docs/reference/structure.md](../../../docs/reference/structure.md)

## Flujo de Use Cases (orden obligatorio)

1. Corriente Nominal
2. Ajuste de Corriente
3. Seleccionar TipoCanalizacion
4. Resolver tabla ampacidad + temperatura
5. Conductor de Tierra
6. Dimensionar Canalización
7. Calcular Caída de Tensión

### Firma de CalcularCaidaTensionUseCase.Execute

```go
func (uc *CalcularCaidaTensionUseCase) Execute(
    ctx              context.Context,
    calibre          string,
    material         valueobject.MaterialConductor,
    corrienteAjustada valueobject.Corriente,
    longitudCircuito float64,
    tension          valueobject.Tension,
    limiteCaida      float64,
    tipoCanalizacion entity.TipoCanalizacion,
    sistemaElectrico entity.SistemaElectrico,
    tipoVoltaje      entity.TipoVoltaje,
    hilosPorFase     int,
    factorPotencia   float64,   // cosθ ∈ (0, 1] — obligatorio
) (dto.ResultadoCaidaTension, error)
```

> El campo `Impedancia` en `dto.ResultadoCaidaTension` es `Zef = R·cosθ + X·senθ` por conductor (Ω/km).
> La caída se calcula como: `e = factor × (I/N) × L × Zef`

## Reglas de Oro — Capa Application

*Estas reglas son específicas para la capa Application de cálculos. Ver [docs/reference/structure.md](../../../docs/reference/structure.md) para reglas globales.*

### Use cases solo orquestan

```go
// BIEN: solo coordina
data, err := uc.repo.Find(...)
result, err := service.Calculate(...)
return dto.FromDomain(result), nil

// MAL: lógica de negocio aquí
if valor > 100 {  // esto va en domain
    return error
}
```

### Ports

- **Driver**: interfaces que expone application (para HTTP, gRPC, CLI)
- **Driven**: interfaces que application necesita (Repository, Clientes externos)

### DTOs

- **Structs con tipos PRIMITIVOS** (string, int, float64)
- Nunca exponer value objects ni entities de domain
- Métodos helper permitidos: `Validate()`, `ToDomain*()`
- Mapping explícito `domain ↔ DTO` dentro del use case

**Campos con unidad en `EquipoInput`:**

| Campo           | Tipo    | Unidad campo    | Valores válidos         | Default |
| --------------- | ------- | --------------- | ----------------------- | ------- |
| `Tension`       | float64 | `TensionUnidad` | V, kV                   | "V"     |
| `PotenciaNominal` | float64 | `PotenciaUnidad` | W, KW, KVA, KVAR     | "KW"    |

Los defaults se aplican en `ApplyDefaults()`. La conversión al value object ocurre en `ToDomainTension()` y `ToDomainPotencia()`.

```go
// ✅ CORRECTO — DTO con primitivos
type MiInput struct {
    Corriente        float64  // primitivo
    TipoCanalizacion string   // primitivo
    Material         string   // primitivo
    Temperatura      *int     // primitivo opcional
}

func (i MiInput) Validate() error { ... }
func (i MiInput) ToDomainMaterial() valueobject.MaterialConductor { ... }

// ❌ INCORRECTO — DTO con value objects
type MiInput struct {
    Corriente valueobject.Corriente  // NO — esto es domain bleeding
}
```

### Use Case como Puente DTO ↔ Domain

El use case es responsable de la conversión:

```go
func (uc *MiUseCase) Execute(ctx context.Context, input dto.MiInput) (dto.MiOutput, error) {
    // 1. Validar DTO
    if err := input.Validate(); err != nil { return ..., err }

    // 2. Convertir primitivos → value objects
    corriente, err := valueobject.NewCorriente(input.Corriente)
    tipoCanalizacion, err := entity.ParseTipoCanalizacion(input.TipoCanalizacion)
    material := input.ToDomainMaterial()

    // 3. Llamar servicio de dominio
    resultado, err := service.MiServicio(corriente, material, ...)

    // 4. Convertir domain → DTO output
    return dto.MiOutput{
        Calibre:    resultado.Calibre(),
        Material:   resultado.Material().String(),
        SeccionMM2: resultado.SeccionMM2(),
    }, nil
}
```

### Single Responsibility

- **Un use case = una responsabilidad**
- Si un use case hace 2+ cosas distintas → separar
- El orquestador coordina múltiples use cases

## DTOs de Charola

Los DTOs de salida para charola incluyen el campo `FactorControl`:

| DTO                        | Campo                      | Descripción                |
| -------------------------- | -------------------------- | -------------------------- |
| `CharolaTriangularOutput`  | `FactorControl float64`   | Siempre 1.0 para cables control |
| `CharolaEspaciadoOutput`   | `FactorControl float64`   | Siempre 1.0 para cables control |
| `DetalleCharola` (memoria_output.go) | `FactorControl float64 \`json:"factor_control"\`` | Incluido en respuesta API |

**Regla:** Los use cases NO recalculan valores que el domain service ya calculó — el campo viene mapeado directamente desde el servicio de dominio.
