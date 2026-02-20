# Estructura del Proyecto

```
internal/
  shared/
    kernel/
      valueobject/
  calculos/
    domain/
      entity/
      service/
    application/
      port/
      usecase/
        helpers/
      dto/
    infrastructure/
      adapter/
        driver/http/
        driven/csv/
  equipos/
    domain/
    application/
    infrastructure/
cmd/api/main.go
data/tablas_nom/
tests/integration/
```

## Guias por Capa

| Capa | Ubicacion | AGENTS.md |
|------|-----------|-----------|
| Shared Kernel | `internal/shared/kernel/` | [AGENTS.md](../../internal/shared/kernel/AGENTS.md) |
| Feature Calculos | `internal/calculos/` | ver subcapas abajo |
| Domain — Entity | `internal/calculos/domain/entity/` | [AGENTS.md](../../internal/calculos/domain/AGENTS.md) |
| Domain — Services | `internal/calculos/domain/service/` | [AGENTS.md](../../internal/calculos/domain/AGENTS.md) |
| Application | `internal/calculos/application/` | [AGENTS.md](../../internal/calculos/application/AGENTS.md) |
| Infrastructure | `internal/calculos/infrastructure/` | [AGENTS.md](../../internal/calculos/infrastructure/AGENTS.md) |
| Feature Equipos | `internal/equipos/` | [AGENTS.md](../../internal/equipos/AGENTS.md) |
| Datos NOM | `data/tablas_nom/` | [AGENTS.md](../../data/tablas_nom/AGENTS.md) |

## Reglas de Aislamiento Entre Features

- `calculos/` NUNCA importa `equipos/` y viceversa
- `shared/kernel/` NO importa ninguna feature
- `cmd/api/main.go` es el ÚNICO archivo que puede importar múltiples features
- Comunicación entre features: solo vía interfaces en `shared/kernel/`

## Reglas de Dependencias por Capa

> Estas reglas aplican a TODAS las features. Ver [AGENTS.md](../../AGENTS.md) por feature.

| Capa | Dependencias permitidas | Dependencias PROHIBIDAS |
|------|------------------------|------------------------|
| **Domain** | `shared/kernel/valueobject`, stdlib | `application/`, `infrastructure/`, frameworks |
| **Application** | `domain/`, `shared/kernel/valueobject`, stdlib | `infrastructure/`, frameworks |
| **Infrastructure** | `domain/`, `application/port`, `shared/kernel/valueobject`, frameworks | Lógica de negocio (va en domain) |

### Anti-duplicación

- NO definir reglas de dependencias en múltiples AGENTS.md
- Referenciar siempre `docs/reference/structure.md` para reglas de arquitectura

### Subcarpetas sin AGENTS.md

Las siguientes subcarpetas **heredan** del AGENTS.md de su padre (no necesitan AGENTS.md propio):

| Subcarpeta | Hereda de |
|------------|-----------|
| `internal/calculos/domain/entity/` | `internal/calculos/domain/AGENTS.md` |
| `internal/calculos/domain/service/` | `internal/calculos/domain/AGENTS.md` |
| `internal/calculos/application/dto/` | `internal/calculos/application/AGENTS.md` |
| `internal/calculos/application/port/` | `internal/calculos/application/AGENTS.md` |
| `internal/calculos/application/usecase/` | `internal/calculos/application/AGENTS.md` |
| `internal/calculos/infrastructure/adapter/driver/http/` | `internal/calculos/infrastructure/AGENTS.md` |
| `internal/calculos/infrastructure/adapter/driven/csv/` | `internal/calculos/infrastructure/AGENTS.md` |
