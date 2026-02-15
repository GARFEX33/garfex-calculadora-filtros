# Feature: equipos

Catálogo de equipos Garfex — PLACEHOLDER FUTURO.

Esta feature estará a cargo del catálogo de equipos eléctricos (filtros activos,
filtros de rechazo, transformadores, cargas) con su búsqueda y persistencia en PostgreSQL.

## Estado actual

**Estructura vacía.** Solo existe como placeholder para mantener los límites de la arquitectura.

## Cuando se implemente

Seguir el flujo completo de agentes:

1. `domain-agent` — entidades de equipo, value objects propios
2. `application-agent` — ports (EquipoRepository), use cases de búsqueda, DTOs
3. `infrastructure-agent` — PostgresEquipoRepository, handlers HTTP GET /equipos

## Reglas de aislamiento

- Esta feature NO importa `calculos/` ni ninguna otra feature
- Solo importa `shared/kernel/` si necesita value objects eléctricos compartidos
- `cmd/api/main.go` es el único que conecta esta feature con las demás
