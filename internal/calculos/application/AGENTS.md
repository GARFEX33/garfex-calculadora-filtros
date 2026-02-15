# Application Layer — calculos

Orquesta domain services. Define contratos (ports), no implementaciones.

> **Skills Reference**:
>
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — interfaces pequeñas, error wrapping, convenciones de ports

---

## Auto-invoke

| Accion                         | Skill             |
| ------------------------------ | ----------------- |
| Definir nuevo port (interface) | `golang-patterns` |
| Crear o modificar use case     | `golang-patterns` |
| Agregar o modificar DTOs       | `golang-patterns` |

---

## Estructura

```
internal/calculos/application/
├── port/           # Interfaces que infrastructure implementa
├── usecase/        # Orquestadores (ej. CalcularMemoriaUseCase)
│   └── helpers/    # Funciones auxiliares de use cases
└── dto/            # Entrada/salida de la API
```

---

## Ports (interfaces)

- **EquipoRepository** — buscar equipos en BD (PostgreSQL)
- **TablaNOMRepository** — leer tablas CSV de ampacidad, tierra, impedancia
- **SeleccionarTemperaturaPort** — selección de temperatura según reglas NOM

Reglas:

- Interfaces pequeñas y enfocadas (1–3 métodos)
- Definidas en `application/port/`
- Implementadas en `infrastructure/`
- Sin lógica de negocio
- Sin importar infrastructure en application

---

## Flujo del UseCase (orden obligatorio)

1. Corriente Nominal
2. Ajuste de Corriente
3. Seleccionar TipoCanalizacion
4. Resolver tabla ampacidad + temperatura → `SeleccionarConductorAlimentacion`
5. Conductor de Tierra (tabla 250-122)
6. Dimensionar Canalización (40% fill)
7. Resolver R y X + FP → `CalcularCaidaTension`

---

## DTOs

- Structs planos sin métodos de negocio
- Validación de input en el use case
- Nunca exponer entidades de domain directamente
- Mapping domain ↔ DTO siempre explícito

### EquipoInput

Incluye:

- Modo (LISTADO / MANUAL_AMPERAJE / MANUAL_POTENCIA)
- Datos del equipo
- Parámetros de instalación
- TipoCanalizacion
- TemperaturaOverride
- Material (Cu/Al, default Cu)

### MemoriaOutput

Resultado completo de todos los pasos.

---

## Campo Material (Cu/Al)

- Tipo: `valueobject.MaterialConductor`
- JSON: "CU" / "AL"
- Default: Cobre si vacío

Afecta:

- Selección de conductor de alimentación
- Selección de conductor de tierra
- Cálculo de caída de tensión

---

## Conductor de Tierra - Regla Material

1. Buscar fila donde `ITM <= ITMHasta`
2. Si material == Al y existe ConductorAl → usar Al
3. Si material == Al y NO existe ConductorAl → fallback a Cu
4. Si material == Cu → usar Cu

---

# CRITICAL RULES

## Dependency Direction

- Application puede depender de domain
- Domain nunca depende de application
- Application nunca depende de infrastructure

---

# Use Case Rules

## Responsibility

A use case **only orchestrates**.
It must not contain business rules.

---

## Allowed

- Call repository ports
- Build domain input structures
- Call domain services
- Map result to DTO
- Handle errors

---

## Forbidden

- Business calculations
- Mathematical formulas
- Electrical rules
- Complex conditional logic
- Large switch statements
- Infrastructure dependencies

---

## Structural Pattern

```
func (uc *XUseCase) Execute(...) (dto.Result, error) {
    data, err := uc.repo.Method(...)
    if err != nil { return dto.Result{}, err }

    result, err := service.DomainLogic(...)
    if err != nil { return dto.Result{}, err }

    return dto.Result{ ... }, nil
}
```

---

## Golden Rule

If business rules change, the use case should remain untouched.

---

## Refactor Trigger

Refactor to `domain/service` when:

- A formula appears
- More than 2 business conditionals are added
- Logic is reused in another use case
- The use case grows beyond orchestration

---

# Use Case Quality Constraints

- Use case < ~80 lineas
- No duplicación de lógica entre use cases
- Error wrapping con `%w`
- Naming cumple convenciones definidas

---

# Naming Conventions

| Entidad        | Patrón                 |
| -------------- | ---------------------- |
| Port           | `PascalCaseRepository` |
| Use case       | `PascalCaseUseCase`   |
| DTO entrada    | `PascalCaseInput`      |
| DTO salida     | `PascalCaseOutput`     |
| Error sentinel | `ErrPascalCase`        |

---

# QA CHECKLIST

## Tests

- [ ] `go test ./internal/calculos/application/...` pasa

---

## Arquitectura

- [ ] Use case solo orquesta
- [ ] Lógica delegada a domain/service
- [ ] Sin reglas eléctricas directas
- [ ] Flujo respeta el orden obligatorio

---

## Dependencias

- [ ] Sin imports de infrastructure
- [ ] Solo depende de port, dto y domain
- [ ] Domain no depende de application
- [ ] Ningún tipo de infrastructure filtrado

---

## DTO Boundary

- [ ] No se retornan entidades de domain directamente
- [ ] Mapping domain ↔ DTO explícito
- [ ] DTOs sin métodos de negocio

---

## Ports

- (1–3 [ ] Interfaces pequeñas métodos)
- [ ] Sin lógica en interfaces
- [ ] Métodos reciben context.Context si hacen I/O

---

## Consistencia

- [ ] Errores de flujo usan sentinels
- [ ] Error wrapping con `%w`
- [ ] Naming correcto
- [ ] No hay duplicación de lógica

---

## Imports

Los imports correctos en esta capa son:

```go
// Value objects
"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"

// Entities
"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"

// Services
"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"

// Application
"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
"github.com/garfex/calculadora-filtros/internal/calculos/application/port"
```
