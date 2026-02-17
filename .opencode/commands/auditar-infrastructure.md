---
description: "Audita la capa de infrastructure de una feature verificando adapters, handlers, repositorios, seguridad y Go idiomático."
agent: auditor-infrastructure
disable-model-invocation: false
---

# Auditoría de Capa de Infrastructure

Invoca al agente `auditor-infrastructure` para realizar una auditoría estricta de la capa de infrastructure.

## Qué verifica

### Arquitectura Hexagonal

- [ ] Import de application/port para implementar
- [ ] Import de application/usecase para llamar
- [ ] Import de domain/entity para mapear
- [ ] Sin imports de domain/service para lógica

### Driven Adapters (Repositorios)

- [ ] Implementa interface de application/port
- [ ] Constructor recibe dependencias (db, config)
- [ ] context.Context como primer parámetro
- [ ] Solo traduce datos, sin lógica de negocio
- [ ] Manejo correcto de errores de I/O
- [ ] Cerrar recursos (defer rows.Close())

### Driver Adapters (HTTP Handlers)

- [ ] Constructor recibe use cases inyectados
- [ ] Handler solo: bind → validate → call UC → respond
- [ ] Sin lógica de negocio
- [ ] Mapeo correcto de errores a HTTP status
- [ ] Context propagado a use cases

### Anti-Patterns a detectar

- [ ] God Handler (handler que hace todo)
- [ ] Repository con Lógica (calcula en repo)
- [ ] Adapter Leaky (expone detalles de BD)
- [ ] Missing Context (sin context.Context)

### Go Idiomático (golang-patterns + golang-pro)

- [ ] gofmt aplicado
- [ ] Error wrapping con `%w`
- [ ] Defer para cleanup (file.Close, rows.Close)
- [ ] Context propagado correctamente
- [ ] Timeouts configurados en clientes HTTP/DB
- [ ] Graceful shutdown manejado
- [ ] Sin goroutines huérfanas

### Seguridad

- [ ] Sin SQL injection (prepared statements)
- [ ] Sin path traversal en file operations
- [ ] Input sanitizado antes de logs
- [ ] Secrets no hardcodeados
- [ ] CORS configurado si es API web

## Uso

```
/auditar-infrastructure calculos
/auditar-infrastructure equipos
/auditar-infrastructure {feature}
```

## Output

Reporte estructurado con severidades:

- **CRÍTICO** — debe corregirse antes de merge
- **IMPORTANTE** — debería corregirse pronto
- **SUGERENCIA** — nice to have

## Sección especial: Seguridad

El reporte incluye una sección dedicada a seguridad:

```
SEGURIDAD
---------
- SQL injection: {OK|WARN}
- Path traversal: {OK|WARN}
- Secrets: {OK|WARN}
```

## Skills de referencia

- `clean-ddd-hexagonal-vertical-go-enterprise`
- `brainstorming-infrastructure`
- `golang-patterns`
- `golang-pro`
- `api-design-principles`
