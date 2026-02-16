---
description: "Audita la capa de application de una feature verificando use cases, ports, DTOs, orquestación y Go idiomático."
agent: auditor-application
disable-model-invocation: true
---

# Auditoría de Capa de Application

Invoca al agente `auditor-application` para realizar una auditoría estricta de la capa de application.

## Qué verifica

### Arquitectura Hexagonal
- [ ] Sin imports de infrastructure/
- [ ] Sin imports de frameworks externos
- [ ] Imports de domain/ correctos
- [ ] Imports de shared/kernel/ si necesario

### Ports (Driver)
- [ ] Son interfaces, no structs
- [ ] Métodos reciben/retornan DTOs o primitivos
- [ ] No exponen tipos de domain directamente

### Ports (Driven)
- [ ] Son interfaces, no structs
- [ ] context.Context como primer parámetro
- [ ] Sin detalles de implementación (SQL, HTTP)

### Use Cases
- [ ] Constructor `New*()` recibe interfaces
- [ ] Solo orquestación, no lógica de negocio
- [ ] Tamaño < 100 líneas (idealmente < 80)
- [ ] Error wrapping con contexto

### DTOs
- [ ] Structs planos sin métodos de negocio
- [ ] Funciones de mapping `FromDomain()` / `ToDomain()`
- [ ] No exponen detalles internos de domain

### Anti-Patterns a detectar
- [ ] Anemic Use Case (solo pasa datos)
- [ ] Fat Use Case (> 150 líneas)
- [ ] Leaky Abstraction (expone SQL, etc.)
- [ ] Domain Bleeding (retorna entidad al exterior)

### Go Idiomático (golang-patterns + golang-pro)
- [ ] gofmt aplicado
- [ ] Error wrapping con `%w`
- [ ] Context.Context como primer parámetro
- [ ] Interfaces pequeñas (1-3 métodos ideal)
- [ ] Constructores retornan struct concreto
- [ ] Sin estado global mutable
- [ ] Table-driven tests con subtests

## Uso

```
/auditar-application calculos
/auditar-application equipos
/auditar-application {feature}
```

## Output

Reporte estructurado con severidades:
- **CRÍTICO** — debe corregirse antes de merge
- **IMPORTANTE** — debería corregirse pronto
- **SUGERENCIA** — nice to have

## Skills de referencia

- `clean-ddd-hexagonal-vertical-go-enterprise`
- `brainstorming-application`
- `golang-patterns`
- `golang-pro`
