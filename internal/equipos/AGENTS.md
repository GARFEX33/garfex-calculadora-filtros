# Feature: equipos

Catálogo de equipos Garfex — PLACEHOLDER FUTURO.

Esta feature estará a cargo del catálogo de equipos eléctricos (filtros activos,
filtros de rechazo, transformadores, cargas) con su búsqueda y persistencia en PostgreSQL.

## Estado actual

**Estructura vacía.** Solo existe como placeholder para mantener los límites de la arquitectura.

## Implementar esta feature

Cuando se decida implementar, usar el sistema de orquestación:

```bash
# Paso 1: Domain
orchestrate-agents --agent domain --feature equipos

# Paso 2: Application  
orchestrate-agents --agent application --feature equipos

# Paso 3: Infrastructure
orchestrate-agents --agent infrastructure --feature equipos

# Paso 4: Orquestador actualiza main.go
```

## Estructura esperada

```
internal/equipos/
├── domain/
│   ├── entity/      ← Equipo, FiltroActivo, FiltroRechazo, etc.
│   └── service/     ← Búsqueda, Validaciones
├── application/
│   ├── port/        ← EquipoRepository
│   ├── usecase/     ← BuscarEquipo, ListarEquipos
│   └── dto/         ← EquipoInput, EquipoOutput
└── infrastructure/
    └── adapter/
        ├── driver/http/      ← EquipoHandler
        └── driven/postgres/  ← PostgresEquipoRepository
```

## Reglas de Aislamiento

- **NO importa** `calculos/` ni ninguna otra feature
- **Solo importa** `shared/kernel/` si necesita VOs eléctricos compartidos
- **`cmd/api/main.go`** es el único que conecta esta feature

## Referencias

- Comando: `.opencode/commands/orchestrate-agents.md`
- Agente domain: `.opencode/agents/domain-agent.md`
- Agente application: `.opencode/agents/application-agent.md`
- Agente infrastructure: `.opencode/agents/infrastructure-agent.md`
