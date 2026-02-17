---
description: "Audita las 3 capas de una feature (domain, application, infrastructure) ejecutando los 3 auditores en secuencia."
disable-model-invocation: false
---

# Auditoría Completa de Feature

Ejecuta una auditoría completa de las 3 capas de una feature, invocando los auditores en orden:

```
auditor-domain → auditor-application → auditor-infrastructure
```

## Flujo

1. **Auditar Domain** — DDD, pureza, value objects, entidades, servicios
2. **Auditar Application** — Use cases, ports, DTOs, orquestación
3. **Auditar Infrastructure** — Adapters, handlers, repos, seguridad

Cada auditor genera un reporte independiente con severidades.

## Uso

```
/auditar-feature calculos
/auditar-feature equipos
/auditar-feature {nombre-feature}
```

## Qué verifican los auditores

### auditor-domain

- Pureza del dominio (sin imports externos)
- Value Objects inmutables con validación
- Entidades con identidad e invariantes
- Servicios de dominio stateless
- Go idiomático

### auditor-application

- Use cases que solo orquestan
- Ports como interfaces
- DTOs con mapping explícito
- Sin lógica de negocio en use cases
- Go idiomático

### auditor-infrastructure

- Adapters que implementan ports exactamente
- Handlers que solo coordinan
- Repositorios sin lógica de negocio
- Seguridad (SQL injection, path traversal)
- Go idiomático

## Output

Reporte consolidado con:

```
=== AUDITORÍA COMPLETA: {feature} ===

DOMAIN LAYER
------------
✅ Passed: {n}
⚠️ Warnings: {n}
❌ Failed: {n}

APPLICATION LAYER
-----------------
✅ Passed: {n}
⚠️ Warnings: {n}
❌ Failed: {n}

INFRASTRUCTURE LAYER
--------------------
✅ Passed: {n}
⚠️ Warnings: {n}
❌ Failed: {n}

RESUMEN TOTAL
-------------
Total críticos: {n}
Total importantes: {n}
Total sugerencias: {n}

¿Listo para merge? {SÍ|NO - corregir {n} críticos}
```

## Cuándo usar

- **Antes de merge** a main
- **Como PR review** automatizado
- **Después de refactorizaciones** grandes
- **Antes de deploy** a producción

## Skills invocados

- `clean-ddd-hexagonal-vertical-go-enterprise`
- `enforce-domain-boundary`
- `golang-patterns`
- `golang-pro`
- `api-design-principles`
