# Calculos — Infrastructure Layer

Implementa los ports definidos en `application/port/`. Tecnologías: CSV (encoding/csv), HTTP (Gin).

## Trabajar en esta Capa

Esta capa es responsabilidad del **`infrastructure-agent`**. El agente ejecuta su ciclo completo:

```
brainstorming-infrastructure → writing-plans-infrastructure → executing-plans-infrastructure
```

**NO modificar directamente** — usar el sistema de orquestación.

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

### Para nueva feature

```bash
# Primero: domain-agent y application-agent completan sus capas
# Luego:
orchestrate-agents --agent infrastructure --feature nueva-feature
```

### Para cambios en calculos existente

```bash
# Orquestador:
# "infrastructure-agent: agregar handler para exportar resultados a CSV"
```

## Adapters

### Driven (implementan ports)

- **CSVTablaNOMRepository** — lee tablas NOM desde CSV
- **CSVSeleccionarTemperatura** — temperaturas por estado

### Driver (HTTP)

- **CalculoHandler** — endpoints REST
  - `POST /api/v1/calculos/memoria`

### Formatters

- **NombreTablaAmpacidad** — nombres descriptivos de tablas
- **GenerarObservaciones** — observaciones del cálculo

## Mapeo de Errores HTTP

| Error domain/application | HTTP status |
|--------------------------|-------------|
| ErrModoInvalido | 400 |
| ErrCanzacionNoSoportada | 400 |
| Validación | 400 |
| ErrConductorNoEncontrado | 422 |
| ErrCanalizacionNoDisponible | 422 |
| Error interno | 500 |

## Reglas de Oro

1. **Implementar exactamente el port** — no agregar métodos
2. **Sin lógica de negocio** — solo traducción de datos
3. **Handlers solo coordinan** — bind → use case → response
4. **Inyección de dependencias** — constructor, no globals
5. **Context.Context** — primer parámetro en I/O

## Referencias

- Agente: `.opencode/agents/infrastructure-agent.md`
- Comando: `.opencode/commands/orchestrate-agents.md`
- Skills: `brainstorming-infrastructure`, `writing-plans-infrastructure`, `executing-plans-infrastructure`

## QA Checklist

- [ ] `go test ./internal/calculos/infrastructure/...` pasa
- [ ] Repositorios implementan ports exactamente
- [ ] Sin estado global
- [ ] Sin lógica de negocio
