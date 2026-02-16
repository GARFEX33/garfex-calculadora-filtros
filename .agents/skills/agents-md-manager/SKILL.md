---
name: agents-md-manager
description: Gestionar archivos AGENTS.md jerárquicos para proyectos Go con arquitectura hexagonal. Trigger: Crear/auditar AGENTS.md, agregar capas, registrar skills, después de features.
license: Apache-2.0
metadata:
  author: garfex
  version: "2.0"
  agent: agents-md-curator
---

# AGENTS.md Manager

Gestiona la estructura jerárquica de AGENTS.md para proyectos Go con arquitectura hexagonal + vertical slices. Optimizado para eficiencia de tokens y coordinación con sistema de agentes especializados.

## Cuándo Usar

- Crear AGENTS.md para nuevo proyecto Go
- Auditar AGENTS.md existentes
- Agregar nueva capa o feature
- Después de completar una feature (post-merge)
- Después de 10+ commits desde última auditoría

## Agente Asociado

Este skill es usado por el **`agents-md-curator`** — un agente especializado que:

- Solo lee código y documentación, nunca modifica código
- Propone cambios uno a uno, esperando confirmación
- Se invoca después de merges o cambios significativos

Ver: `.opencode/agents/agents-md-curator.md`

---

## Principios Fundamentales

### Eficiencia de Tokens

| Archivo         | Propósito             | Límite        |
| --------------- | --------------------- | ------------- |
| Root AGENTS.md  | Índice de navegación  | ~150 líneas   |
| Layer AGENTS.md | Reglas de la capa     | ~150 líneas   |
| Skills          | Patrones con ejemplos | Autocontenido |

**Carga por acción:** root (~50 líneas) + 1 layer (~50 líneas) + 1 skill (~100 líneas) = ~200 líneas total

### Jerarquía de Precedencia

```
Layer AGENTS.md > Root AGENTS.md > Skills (solo para patrones)
```

Si hay conflicto no resoluble → STOP y pedir confirmación.

### Nunca Duplicar

- Si está en root → no en layer
- Si está en skill → no en AGENTS.md
- Referenciar, no copiar

---

## Modos de Operación

### `/agents-md-manager create`

Generar jerarquía completa desde cero.

**Pasos:**

1. **Escanear** estructura del proyecto
2. **Evaluar** granularidad por directorio
3. **Generar** root AGENTS.md (usar template)
4. **Generar** layer AGENTS.md por cada capa
5. **Presentar** cada archivo para confirmación

### `/agents-md-manager audit`

Revisar AGENTS.md existentes y proponer correcciones.

**Output:**

```
=== AUDITORÍA AGENTS.md ===

Root AGENTS.md: {OK|WARN|FAIL}
  Estructura: {OK|FAIL} — secciones faltantes: {lista}
  Contenido: {OK|WARN|FAIL} — {advertencias}

{layer}/AGENTS.md: {OK|WARN|FAIL}
  Estructura: {OK|FAIL}
  Sincronización: {OK|WARN} — {discrepancias con código}

Skills: {n} advertencias
Agentes: {OK|WARN} — {verificación de referencias}

=== PROPUESTAS ({n}) ===

1. {descripción}
2. {descripción}

¿Aplicar propuesta 1? [si/no]
```

---

## Checklists de Auditoría

### Root AGENTS.md

| Sección                     | Requerida | Regla                                               |
| --------------------------- | --------- | --------------------------------------------------- |
| `## Cómo Usar Esta Guía`    | Sí        | 3 bullets: empezar aquí, docs por capa, precedencia |
| `## Regla Anti-Duplicación` | Sí        | Flujo Investigar → Decidir → Comunicar              |
| `## Sistema de Agentes`     | Sí        | Tabla de agentes con skills asociados               |
| `## Guías por Capa`         | Sí        | Tabla con ubicación y AGENTS.md                     |
| `## Skills Disponibles`     | Sí        | Tablas de skills genéricos y de proyecto            |
| `## Auto-invocación`        | Sí        | Mapeo acción → agente/skill                         |
| `## Stack`                  | Sí        | Una línea                                           |
| `## Comandos`               | Sí        | Build, test, lint                                   |
| `## Documentación`          | Opcional  | Planes completados y en progreso                    |

**Límite:** ~150 líneas (warn a 130)

### Layer AGENTS.md

| Sección                      | Requerida | Regla                                              |
| ---------------------------- | --------- | -------------------------------------------------- |
| Título con agente            | Sí        | `# {Feature} — {Layer} Layer` + agente responsable |
| `## Trabajar en esta Capa`   | Sí        | Qué agente, qué skills, flujo                      |
| `## Estructura`              | Sí        | Árbol de directorios                               |
| `## Dependencias permitidas` | Sí        | Lista de imports válidos                           |
| `## Dependencias prohibidas` | Sí        | Lista de imports bloqueados                        |
| `## Reglas de Oro`           | Sí        | 5-6 puntos máximo                                  |
| `## QA Checklist`            | Sí        | Checkboxes de verificación                         |

**Límite:** ~150 líneas (warn a 120)

### Consistencia de Skills

| Check         | Regla                                           |
| ------------- | ----------------------------------------------- |
| Registrados   | Todo skill en `.agents/skills/` aparece en root |
| Existen       | Todo skill referenciado tiene SKILL.md          |
| Sin huérfanos | No hay skills sin referencia                    |

### Consistencia de Agentes

| Check            | Regla                                               |
| ---------------- | --------------------------------------------------- |
| Registrados      | Todo agente en `.opencode/agents/` aparece en root  |
| Referencias      | Cada layer AGENTS.md menciona su agente responsable |
| Skills asociados | Cada agente tiene skills listados                   |

---

## Algoritmo de Granularidad

Evaluar 3 métricas por directorio:

| Métrica                         | Umbral | Significado          |
| ------------------------------- | ------ | -------------------- |
| Archivos `.go` (sin `_test.go`) | > 8    | Demasiados conceptos |
| Líneas del AGENTS.md            | > 150  | Demasiado contenido  |
| Subdirectorios con `.go`        | >= 3   | Patrones distintos   |

**Decisión:**

```
0 umbrales excedidos → mantener actual
1 umbral excedido → WARN, revisar después
2+ umbrales excedidos → proponer split o skill extraction
```

---

## Integración con Sistema de Agentes

### Flujo de Orquestación

```
Usuario pide feature
        │
        ▼
┌─────────────────────────────────┐
│         ORQUESTADOR             │
│  1. Investigar qué existe       │
│  2. Decidir estrategia          │
│  3. Comunicar a subagentes      │
└─────────────────────────────────┘
        │
   ┌────┴────┬────────────┐
   ▼         ▼            ▼
domain-   application- infrastructure-
agent     agent        agent
        │
        ▼ (post-merge)
┌─────────────────────────────────┐
│      agents-md-curator          │
│  Auditar y actualizar AGENTS.md │
└─────────────────────────────────┘
```

### Regla Anti-Duplicación del Orquestador

El orquestador DEBE verificar antes de despachar agentes:

```bash
# Investigar
ls internal/{feature}/domain/service/*.go
rg "TODO|FIXME" internal/{feature}/application/usecase --type go
rg -i "func.*[Cc]alcular" internal/{feature} --type go

# Decidir
# ¿Extender existente o crear nuevo?

# Comunicar
# Instrucciones claras al agente: qué hacer y qué NO hacer
```

Esta regla debe estar documentada en el root AGENTS.md.

---

## Enforcement de Aislamiento (Hexagonal)

Durante auditoría, validar:

| Capa           | Puede importar                     | NO puede importar           |
| -------------- | ---------------------------------- | --------------------------- |
| domain         | shared/kernel, stdlib              | application, infrastructure |
| application    | domain, shared/kernel              | infrastructure              |
| infrastructure | application, domain, shared/kernel | —                           |

Si hay violación → FAIL en auditoría, proponer refactor.

---

## Detección de Drift

Disparar advertencia si:

- Nuevo directorio con `.go` pero sin AGENTS.md
- Directorio excede umbrales de granularidad
- Skill referenciado en commits pero no registrado
- Nuevo patrón no documentado en layer AGENTS.md
- Nuevo agente no registrado en root

Si hay drift → WARN, proponer actualización.

---

## Reglas de Ubicación

| Tipo               | Ubicación                                   | Naming            |
| ------------------ | ------------------------------------------- | ----------------- |
| Root AGENTS.md     | `AGENTS.md`                                 | Raíz del proyecto |
| Layer AGENTS.md    | `internal/{feature}/{layer}/AGENTS.md`      | Por capa          |
| Feature AGENTS.md  | `internal/{feature}/AGENTS.md`              | Por feature       |
| Data AGENTS.md     | `data/{name}/AGENTS.md`                     | Por dataset       |
| Skills genéricos   | `.agents/skills/{name}/SKILL.md`            | Sin prefijo       |
| Skills de proyecto | `.agents/skills/{project}-{scope}/SKILL.md` | Con prefijo       |
| Agentes            | `.opencode/agents/{name}.md`                | Por rol           |

---

## Comandos

```bash
/agents-md-manager create   # Generar jerarquía completa
/agents-md-manager audit    # Auditar y proponer correcciones
```

---

## Templates

- **Root AGENTS.md**: [assets/ROOT-AGENTS-TEMPLATE.md](assets/ROOT-AGENTS-TEMPLATE.md)
- **Layer AGENTS.md**: [assets/LAYER-AGENTS-TEMPLATE.md](assets/LAYER-AGENTS-TEMPLATE.md)

---

## Supuestos de Arquitectura

Optimizado para:

- **Go 1.22+**
- **Hexagonal / Clean Architecture** con vertical slices
- **Layout:** `internal/`, `cmd/`, `data/`, `.agents/skills/`, `.opencode/agents/`
- **Sistema de agentes especializados** por capa
