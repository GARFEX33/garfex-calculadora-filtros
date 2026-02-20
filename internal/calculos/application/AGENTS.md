---
name: application-agent
description: Especialista únicamente en la capa de aplicación de calculos. Ports, use cases y DTOs.
model: opencode/minimax-m2.5-free
---

# Calculos — Application Layer

Orquesta domain services. Define contratos (ports), no implementaciones.

> **Workflow:** Ver [`AGENTS.md` raíz](../../../AGENTS.md) → "Sistema de Agentes Especializados"

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

## Dependencias permitidas

- `internal/shared/kernel/valueobject`
- `internal/calculos/domain/entity`
- `internal/calculos/domain/service`
- stdlib de Go

## Dependencias prohibidas

- `internal/calculos/infrastructure/` — **nunca**
- Frameworks (Gin, pgx, etc.)

## Cómo modificar esta capa

> Ver flujo completo en [`AGENTS.md` raíz](../../../AGENTS.md)

## Flujo de Use Cases (orden obligatorio)

1. Corriente Nominal
2. Ajuste de Corriente
3. Seleccionar TipoCanalizacion
4. Resolver tabla ampacidad + temperatura
5. Conductor de Tierra
6. Dimensionar Canalización
7. Calcular Caída de Tensión

## Reglas de Application

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

## Referencias

- Agente: `application-agent`
- Skill: `.agents/skills/orchestrating-agents/SKILL.md`

## QA Checklist

- [ ] `go test ./internal/calculos/application/...` pasa
- [ ] Use case < ~80 líneas
- [ ] Sin lógica de negocio (fórmulas, validaciones complejas)
- [ ] Sin imports de infrastructure
- [ ] Error wrapping con `%w`
- [ ] DTOs usan SOLO primitivos (string, int, float64)
- [ ] Use case convierte DTO → value objects antes de llamar a domain
- [ ] Use case convierte resultado domain → DTO antes de retornar
- [ ] Un use case = una responsabilidad (no combinar funcionalidades)
