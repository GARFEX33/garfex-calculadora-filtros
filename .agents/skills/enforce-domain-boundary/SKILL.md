---
name: enforce-domain-boundary
description: Esta skill garantiza que el agente de dominio solo genere entidades, value objects y lógica de negocio. Antes de generar cualquier output, revisa que no viole las reglas del dominio.
---

# Enforce Domain Boundary

## Propósito

Esta skill actúa como auditor de dominio. Su función es garantizar que todo el código generado por un agente de dominio:

- Solo contenga **entidades**, **value objects** y **reglas de negocio**.
- No importe ni utilice código de **Application** o **Infrastructure**.
- Mantenga la coherencia con la **arquitectura Vertical Slicing + DDD + Hexagonal**.

## Reglas de la skill

1. **Ubicación de archivos**
   - Solo se permite código dentro de rutas que comiencen con `internal/domain/`.
   - Archivos fuera de `internal/domain/` deben bloquearse.

2. **Imports prohibidos**
   - Bloquear cualquier import que contenga:
     - `internal/application`
     - `internal/infrastructure`

3. **Validación de estructuras**
   - Solo se permiten structs relacionados con dominio:
     - Entidades (`Entity`)
     - Value Objects (`VO`)
     - Agregados (`Aggregate`)
   - Cualquier otro struct o tipo fuera de estos debe ser bloqueado.

4. **Reglas adicionales opcionales**
   - Evitar funciones que dependan de Application o Infrastructure.
   - Validar que no existan métodos de persistencia directa.
   - Opcional: verificar que los agregados respeten invariantes de negocio.

## Flujo de ejecución recomendado

1. El agente genera código de dominio.
2. Antes de entregar output, llama a `enforce-domain-boundary`.
3. Si el código cumple todas las reglas, se aprueba.
4. Si no cumple, se bloquea y se devuelve un error con detalle de la violación.

## Integración con agentes

- Puede ser invocado por cualquier modelo generativo que produzca código Go (Claude Code, OpenCode, GPT, etc.).
- Debe ejecutarse **siempre antes de entregar código de dominio**.
- Funciona como una **skill de auditoría**, no modifica código, solo valida.

## Extensiones posibles

- Validación de dependencias entre módulos verticales.
- Listas blancas de imports permitidos (`fmt`, `errors`, `time`, etc.).
- Generación de reportes de auditoría para cada commit o PR.
