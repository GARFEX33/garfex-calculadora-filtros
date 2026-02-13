# Domain Layer

Capa de negocio pura. Sin dependencias externas (sin Gin, pgx, CSV).
Recibe datos ya interpretados — no conoce archivos ni BD.

## Guias por Subdirectorio

| Subdirectorio | AGENTS.md contiene |
|---------------|--------------------|
| `entity/` | Entidades, TipoEquipo, TipoCanalizacion, MemoriaCalculo, formulas In |
| `valueobject/` | Corriente, Tension, Conductor — inmutables y validados |
| `service/` | 6 servicios de calculo, caida de tension, reglas NOM |

## Auto-invocacion

| Accion | Referencia |
|--------|-----------|
| Crear/modificar entidad o tipo | `entity/AGENTS.md` |
| Trabajar con value objects | `valueobject/AGENTS.md` |
| Crear/modificar servicio de calculo | `service/AGENTS.md` |
| Aplicar patrones Go idiomaticos | skill `golang-patterns` |
