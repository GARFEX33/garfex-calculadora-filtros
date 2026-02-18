---
name: infrastructure-agent
description: Especialista únicamente en la capa de infraestructura de calculos. Adapters HTTP y CSV.
model: opencode/minimax-m2.5-free
---

# Calculos — Infrastructure Layer

Implementa los ports definidos en `application/port/`. Tecnologías: CSV (encoding/csv), HTTP (Gin).

> **Workflow:** Ver [`AGENTS.md` raíz](../../../AGENTS.md) → "Sistema de Agentes Especializados"

## Estructura

```
internal/calculos/infrastructure/
├── adapter/
│   ├── driven/
│   │   └── csv/              # CSVTablaNOMRepository
│   └── driver/
│       └── http/
│           ├── formatters/   # NombreTablaAmpacidad, GenerarObservaciones
│           ├── middleware/   # CORS, RequestLogger
│           └── handler.go    # CalculoHandler
└── router.go                 # Configuración de rutas Gin
```

## Dependencias permitidas

- `internal/shared/kernel/valueobject`
- `internal/calculos/domain/entity`
- `internal/calculos/application/port` (interfaces a implementar)
- `internal/calculos/application/usecase` (para llamar desde handlers)
- Gin, encoding/csv

## Dependencias prohibidas

- `internal/calculos/domain/service` — usar solo entity y valueobject
- Lógica de negocio

## Cómo modificar esta capa

> Ver flujo completo en [`AGENTS.md` raíz](../../../AGENTS.md)

## Adapters

### Driven (implementan ports)

- **CSVTablaNOMRepository** — lee tablas NOM desde CSV
- **CSVSeleccionarTemperatura** — temperaturas por estado

### Driver (HTTP)

- **CalculoHandler** — endpoints REST
  - `POST /api/v1/calculos/memoria` — memoria de cálculo completa
  - `POST /api/v1/calculos/amperaje` — calcular amperaje nominal sin memoria completa
  - `POST /api/v1/calculos/corriente-ajustada` — calcular corriente ajustada con factores NOM

### Formatters

- **NombreTablaAmpacidad** — nombres descriptivos de tablas
- **GenerarObservaciones** — observaciones del cálculo

## Mapeo de Errores HTTP

| Error domain/application | HTTP status |
|--------------------------|-------------|
| ErrModoInvalido | 400 |
| ErrTipoCanalizacionInvalido | 400 |
| ErrSistemaElectricoInvalido | 400 |
| ErrTipoEquipoInvalido | 400 |
| Validación | 400 |
| ErrConductorNoEncontrado | 422 |
| ErrCanalizacionNoDisponible | 422 |
| CALCULO_NO_POSIBLE | 422 |
| Error interno | 500 |

## Reglas de Oro

1. **Implementar exactamente el port** — no agregar métodos
2. **Sin lógica de negocio** — solo traducción de datos
3. **Handlers solo coordinan** — bind → use case → response
4. **Inyección de dependencias** — constructor, no globals
5. **Context.Context** — primer parámetro en I/O

## Referencias

- Agente: `infrastructure-agent`
- Skill: `.agents/skills/orchestrating-agents/SKILL.md`

## QA Checklist

- [ ] `go test ./internal/calculos/infrastructure/...` pasa
- [ ] Repositorios implementan ports exactamente
- [ ] Sin estado global
- [ ] Sin lógica de negocio
