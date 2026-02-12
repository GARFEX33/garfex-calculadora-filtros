---
name: claude-md-manager
description: Create and audit hierarchical CLAUDE.md files for Go hexagonal projects. Generates token-efficient orchestrator pattern with skills integration.
---

# CLAUDE.md Manager

Manage the hierarchical CLAUDE.md structure for Go projects with hexagonal/clean architecture. Follows the Prowler orchestrator pattern optimized for token efficiency.

## When to Activate

- Setting up CLAUDE.md for a new Go project
- Auditing existing CLAUDE.md files for consistency
- Adding a new layer or component to the project
- After significant code growth (new entities, services, etc.)

## Modes

### `/claude-md-manager create`

Generate the full CLAUDE.md hierarchy from scratch.

**Steps:**

1. **Scan** project structure:
   - Find directories with `.go` files under `internal/`, `cmd/`, `data/`
   - Find skills in `.agents/skills/` (read each SKILL.md frontmatter)
   - Classify skills: generic (no project prefix) vs project-specific (has prefix)

2. **Evaluate granularity** per directory (see Granularity Algorithm below)

3. **Generate root CLAUDE.md** as pure orchestrator with these sections only:
   ```
   # {Project Name}
   {One-line description}

   ## Como Usar Esta Guia
   - 3 bullet points: start here, layer docs, precedence rule

   ## Guias por Capa
   | Capa | Ubicacion | CLAUDE.md contiene |
   {one row per directory with CLAUDE.md}

   ## Skills Disponibles
   ### Skills Genericos
   | Skill | Descripcion | Ruta |
   {one row per generic skill}

   ### Skills de Proyecto
   | Skill | Descripcion | Ruta |
   {one row per project skill}

   ## Auto-invocacion
   | Accion | Referencia |
   {maps actions to CLAUDE.md files AND skills}

   ## Stack
   {one line}

   ## Comandos
   {build, test, lint commands}

   ## Fases
   {numbered list, current phase marked}

   ## Convenciones Globales
   {5-6 bullet points max}
   ```

4. **Generate layer CLAUDE.md** for each detected layer:
   - Max 80 lines
   - Only rules specific to that layer
   - Reference relevant skills: "Para patrones Go, usa skill `golang-patterns`"
   - Never duplicate root content (no stack, no conventions, no commands)

5. **Present each file** to user for confirmation before writing

### `/claude-md-manager audit`

Review existing CLAUDE.md files and propose fixes.

**Checklist — Root CLAUDE.md:**

| Check | Rule |
|-------|------|
| Exists | CLAUDE.md in project root |
| Como Usar Esta Guia | Section present |
| Guias por Capa | Every directory with CLAUDE.md is listed |
| Skills Genericos | Every skill in `.agents/skills/` without project prefix is listed |
| Skills de Proyecto | Every skill with project prefix is listed |
| Auto-invocacion | Actions mapped to CLAUDE.md files AND skills |
| No detailed content | No section exceeds 5 content lines (tables excluded) |
| Essentials present | Stack, Comandos, Fases, Convenciones exist |

**Checklist — Layer CLAUDE.md:**

| Check | Rule |
|-------|------|
| Exists | File present in layer directory |
| Line count | <= 80 lines (warn if approaching, propose split if over) |
| No duplication | Does not repeat root content (conventions, stack, commands) |
| Skill references | References relevant skills for detailed patterns |
| Reflects code | Content matches current files and structures |

**Checklist — Skills consistency:**

| Check | Rule |
|-------|------|
| All registered | Every SKILL.md in `.agents/skills/` appears in root tables |
| All exist | Every skill referenced in auto-invoke has a SKILL.md file |
| Consistent prefix | Project skills use consistent prefix (`{project}-`) |
| No orphans | No skills exist without being referenced in root |

**Checklist — Granularity:**

Run the Granularity Algorithm on each directory with a CLAUDE.md.

**Output format:**

```
=== AUDIT CLAUDE.md ===

Root CLAUDE.md: {OK|WARN|FAIL} ({n}/{total} checks)
  {list warnings/failures}

{layer}/CLAUDE.md: {OK|WARN|FAIL} ({n}/{total} checks)
  {list warnings/failures}

Skills: {n} advertencias
  {list warnings}

Granularidad: {OK|action needed}
  {list recommendations}

=== PROPUESTAS ({n}) ===
1. {description}
2. {description}

Aplicar propuesta 1? [si/no]
```

Apply proposals one at a time, waiting for user confirmation each time.

## Granularity Algorithm

Evaluate 3 metrics per directory:

| Metric | Threshold | Meaning |
|--------|-----------|---------|
| `.go` files (excluding `_test.go`) | > 8 | Too many concepts for one CLAUDE.md |
| CLAUDE.md line count | > 80 | Exceeds efficient processing |
| Sub-directories with `.go` files | >= 3 | Distinct patterns per sub-directory |

**Decision:**

```
0 of 3 thresholds exceeded → keep current CLAUDE.md
1 of 3 thresholds exceeded → warn, suggest future review
2+ of 3 thresholds exceeded → propose split into sub-directory CLAUDE.md files
```

When proposing a split:
- Suggest which sub-directories need their own CLAUDE.md
- Draft the content for each new CLAUDE.md
- Update the parent CLAUDE.md to reference the new files
- Present all changes for confirmation

## Placement Rules

| Type | Location | Naming |
|------|----------|--------|
| Generic skills | `.agents/skills/{name}/SKILL.md` | No project prefix |
| Project skills | `.agents/skills/{project}-{scope}/SKILL.md` | With project prefix |
| Layer CLAUDE.md | `internal/{layer}/CLAUDE.md` | Standard |
| Data CLAUDE.md | `data/{name}/CLAUDE.md` | Standard |
| Root CLAUDE.md | `CLAUDE.md` | Project root |

## Token Efficiency Rules

These rules are the foundation of every decision this skill makes:

1. **Root = directions only.** Zero detailed technical content. Only tables pointing to CLAUDE.md files and skills.
2. **Layer CLAUDE.md = rules only.** What to do, not how. Max 80 lines. References skills for detailed patterns.
3. **Skills = patterns with examples.** How to do it. Self-contained. Does not repeat project structure.
4. **Never duplicate.** If info is in root, not in layer. If in a skill, not in CLAUDE.md.
5. **Load path per action:** root (~50 lines) + 1 layer CLAUDE.md (~50 lines) + 1 skill if needed (~100 lines) = ~200 lines total context.

## Architecture Assumptions

This skill is optimized for:
- **Go 1.22+** projects
- **Hexagonal / Clean Architecture** with layers: domain, application, infrastructure, presentation
- **Standard layout:** `internal/`, `cmd/`, `data/`, `.agents/skills/`
- **Single developer or small team** workflow
