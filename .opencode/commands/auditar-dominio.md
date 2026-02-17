---
description: "Audita la capa de dominio de una feature verificando DDD, pureza, value objects, entidades, servicios y Go idiomático."
agent: auditor-domain
disable-model-invocation: false
---

# Auditoría de Capa de Dominio

Invoca al agente `auditor-domain` para realizar una auditoría estricta de la capa de dominio.

## Qué verifica

### Arquitectura DDD

- [ ] Sin imports de application/ o infrastructure/
- [ ] Sin imports de frameworks externos (Gin, pgx, csv)
- [ ] Sin I/O (no archivos, no HTTP, no DB)
- [ ] Sin context.Context en domain (excepto interfaces)
- [ ] Sin tags JSON en structs de dominio
- [ ] Sin panic() — solo errores

### Value Objects

- [ ] Constructor `New*()` con validación
- [ ] Retorna `(T, error)` no solo `T`
- [ ] Campos no exportados (minúscula)
- [ ] Inmutables — sin setters ni mutación

### Entidades

- [ ] Tiene identidad (ID)
- [ ] Constructor valida invariantes
- [ ] Métodos de comportamiento, no solo data

### Servicios de Dominio

- [ ] Sin estado (stateless)
- [ ] Funciones puras (sin I/O)
- [ ] Tests unitarios sin mocks de I/O

### Go Idiomático (golang-patterns + golang-pro)

- [ ] gofmt aplicado
- [ ] Error wrapping con `%w`
- [ ] Return early pattern
- [ ] Sin naked returns
- [ ] GoDoc en funciones exportadas
- [ ] Zero value útil
- [ ] Accept interfaces, return structs

## Uso

```
/auditar-dominio calculos
/auditar-dominio equipos
/auditar-dominio {feature}
```

## Output

Reporte estructurado con severidades:

- **CRÍTICO** — debe corregirse antes de merge
- **IMPORTANTE** — debería corregirse pronto
- **SUGERENCIA** — nice to have

## Skills de referencia

- `clean-ddd-hexagonal-vertical-go-enterprise`
- `enforce-domain-boundary`
- `golang-patterns`
- `golang-pro`
