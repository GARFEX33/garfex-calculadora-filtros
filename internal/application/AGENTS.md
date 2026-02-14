Aquí lo tienes completo, limpio y reorganizado:

---

# Application Layer

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

- `port/` — Interfaces que infrastructure implementa
- `usecase/` — Orquestadores (ej. `CalcularMemoriaUseCase`)
- `dto/` — Entrada/salida de la API

---

## Ports (interfaces)

- **EquipoRepository** — buscar equipos en BD (PostgreSQL)
- **TablaNOMRepository** — leer tablas CSV de ampacidad, tierra, impedancia

Reglas:

- Interfaces pequeñas y enfocadas (1–3 metodos)
- Definidas en `application/port/`
- Implementadas en `infrastructure/`
- Sin logica de negocio
- Sin importar infrastructure en application

---

## Flujo del UseCase (orden obligatorio)

1. Corriente Nominal
2. Ajuste de Corriente
3. Seleccionar TipoCanalizacion
4. Resolver tabla ampacidad + temperatura → `SeleccionarConductorAlimentacion`
5. Conductor de Tierra (tabla 250-122)
6. Dimensionar Canalizacion (40% fill)
7. Resolver R y X + FP → `CalcularCaidaTension`

---

## DTOs

- Structs planos sin metodos de negocio
- Validacion de input en el use case
- Nunca exponer entidades de domain directamente
- Mapping domain ↔ DTO siempre explicito

### EquipoInput

Incluye:

- Modo (LISTADO / MANUAL_AMPERAJE / MANUAL_POTENCIA)
- Datos del equipo
- Parametros de instalacion
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

- Seleccion de conductor de alimentacion
- Seleccion de conductor de tierra
- Calculo de caida de tension

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
- No duplicacion de logica entre use cases
- Error wrapping con `%w`
- Naming cumple convenciones definidas

---

# Naming Conventions

| Entidad        | Patron                 |
| -------------- | ---------------------- |
| Port           | `PascalCaseRepository` |
| Use case       | `PascalCaseUseCase`    |
| DTO entrada    | `PascalCaseInput`      |
| DTO salida     | `PascalCaseOutput`     |
| Error sentinel | `ErrPascalCase`        |

---

# QA CHECKLIST

## Tests

- [ ] `go test ./internal/application/...` pasa

---

## Arquitectura

- [ ] Use case solo orquesta
- [ ] Logica delegada a domain/service
- [ ] Sin reglas electricas directas
- [ ] Flujo respeta el orden obligatorio

---

## Dependencias

- [ ] Sin imports de infrastructure
- [ ] Solo depende de port, dto y domain
- [ ] Domain no depende de application
- [ ] Ningun tipo de infrastructure filtrado

---

## DTO Boundary

- [ ] No se retornan entidades de domain directamente
- [ ] Mapping domain ↔ DTO explicito
- [ ] DTOs sin metodos de negocio

---

## Ports

- [ ] Interfaces pequeñas (1–3 metodos)
- [ ] Sin logica en interfaces
- [ ] Metodos reciben context.Context si hacen I/O

---

## Consistencia

- [ ] Errores de flujo usan sentinels
- [ ] Error wrapping con `%w`
- [ ] Naming correcto
- [ ] No hay duplicacion de logica

---

Con esto tienes una guía clara, coherente y fuerte para cualquier agente o refactor futuro.
