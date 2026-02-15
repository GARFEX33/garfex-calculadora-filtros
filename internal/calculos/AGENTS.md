# Feature: calculos

Memoria de cálculo eléctrico según normativa NOM (México).

Esta feature implementa el cálculo completo de una instalación eléctrica:
corriente nominal → ajuste por temperatura/agrupamiento → selección de conductor →
conductor de tierra → dimensionamiento de canalización → caída de tensión (IEEE-141).

## Estructura

```
internal/calculos/
├── domain/          ← entidades y servicios de cálculo puro
│   ├── entity/      ← Proyecto, TipoCanalizacion, SistemaElectrico, etc.
│   └── service/     ← 7 servicios de cálculo NOM + IEEE-141
├── application/     ← ports, use cases, DTOs
│   ├── port/        ← TablaNOMRepository, EquipoRepository
│   ├── usecase/     ← OrquestadorMemoriaCalculo y micro use cases
│   └── dto/         ← EquipoInput, MemoriaOutput
└── infrastructure/  ← adapters HTTP (driver) y CSV (driven)
    └── adapter/
        ├── driver/http/     ← CalculoHandler, formatters
        └── driven/csv/      ← CSVTablaNOMRepository
```

## Reglas de Aislamiento

- **NO importa** `equipos/` ni ninguna otra feature
- **Solo importa** `shared/kernel/` para value objects compartidos
- **`cmd/api/main.go`** es el único que instancia y conecta esta feature

## Cómo modificar esta feature

**NUNCA modificar directamente.** Usar el sistema de orquestación.

### Para cambios en calculos:

```bash
# Según la capa que necesites modificar:

# Capa domain:
orchestrate-agents --agent domain --feature calculos

# Capa application:
orchestrate-agents --agent application --feature calculos

# Capa infrastructure:
orchestrate-agents --agent infrastructure --feature calculos
```

### Flujo típico para nueva funcionalidad:

```
1. domain-agent implementa entidades/servicios
2. application-agent implementa use cases
3. infrastructure-agent implementa adapters
4. Orquestador actualiza main.go
```

## Agente por Capa

| Capa | Agente | AGENTS.md | Flujo |
|------|--------|-----------|-------|
| domain/ | `domain-agent` | [`domain/AGENTS.md`](domain/AGENTS.md) | `brainstorming-dominio → writing-plans-dominio → executing-plans-dominio` |
| application/ | `application-agent` | [`application/AGENTS.md`](application/AGENTS.md) | `brainstorming-application → writing-plans-application → executing-plans-application` |
| infrastructure/ | `infrastructure-agent` | [`infrastructure/AGENTS.md`](infrastructure/AGENTS.md) | `brainstorming-infrastructure → writing-plans-infrastructure → executing-plans-infrastructure` |

## Referencias

- Comando: `.opencode/commands/orchestrate-agents.md`
- Ejemplo: `.opencode/commands/orchestrate-agents-example.md`
- Sistema: `AGENTS.md` (raíz)
