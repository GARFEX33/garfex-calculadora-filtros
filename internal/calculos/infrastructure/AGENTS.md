# Calculos — Infrastructure Layer

Implementa los ports definidos en `application/port/`. Tecnologías: CSV (encoding/csv), HTTP (Gin).

> **Workflow:** Ver [docs/architecture/agents.md](../../../docs/architecture/agents.md)

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

> **Nota:** Las subcarpetas `adapter/driver/http/` y `adapter/driven/csv/` heredan las reglas de este AGENTS.md. No necesitan AGENTS.md propio.

## Dependencias permitidas

- `internal/shared/kernel/valueobject`
- `internal/calculos/domain/entity`
- `internal/calculos/application/port` (interfaces a implementar)
- `internal/calculos/application/usecase` (para llamar desde handlers)
- Gin, encoding/csv

## Dependencias prohibidas

> Ver reglas consolidadas en [docs/reference/structure.md](../../../docs/reference/structure.md)

- `internal/calculos/domain/service` — usar solo entity y valueobject
- Lógica de negocio

## Cómo modificar esta capa

> Ver flujo completo en [docs/architecture/workflow.md](../../../docs/architecture/workflow.md)

## Adapters

### Driven (implementan ports)

- **CSVTablaNOMRepository** — lee tablas NOM desde CSV
- **CSVSeleccionarTemperatura** — temperaturas por estado

### Driver (HTTP)

- **CalculoHandler** — endpoints REST
  - `POST /api/v1/calculos/amperaje` — calcular amperaje nominal
  - `POST /api/v1/calculos/corriente-ajustada` — calcular corriente ajustada con factores NOM
  - `POST /api/v1/calculos/conductor-alimentacion` — seleccionar conductor de alimentación
  - `POST /api/v1/calculos/conductor-tierra` — seleccionar conductor de tierra
  - `POST /api/v1/calculos/tuberia` — dimensionar tubería
  - `POST /api/v1/calculos/charola/espaciado` — calcular espaciado en charola
  - `POST /api/v1/calculos/charola/triangular` — calcular configuración triangular
  - `POST /api/v1/calculos/caida-tension` — calcular caída de tensión

### Formatters

- **NombreTablaAmpacidad** — nombres descriptivos de tablas
- **GenerarObservaciones** — observaciones del cálculo

## Mapeo de Errores HTTP

| Error domain/application    | HTTP status |
| --------------------------- | ----------- |
| ErrModoInvalido             | 400         |
| ErrTipoCanalizacionInvalido | 400         |
| ErrSistemaElectricoInvalido | 400         |
| ErrTipoEquipoInvalido       | 400         |
| Validación                  | 400         |
| ErrConductorNoEncontrado    | 422         |
| ErrCanalizacionNoDisponible | 422         |
| CALCULO_NO_POSIBLE          | 422         |
| Error interno               | 500         |

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
