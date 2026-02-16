---
name: agents-md-manager
description: Create and audit hierarchical AGENTS.md files for Go hexagonal projects. Trigger: When creating/auditing AGENTS.md files, adding new layers, or registering new skills.
license: Apache-2.0
metadata:
  author: garfex
  version: "1.0"
---

# AGENTS.md Manager

Manage the hierarchical AGENTS.md structure for Go projects with hexagonal/clean architecture. Follows the Prowler orchestrator pattern optimized for token efficiency.

## When to Use

- Setting up AGENTS.md for a new Go project
- Auditing existing AGENTS.md files for consistency
- Adding a new layer or component to the project
- After significant code growth (new entities, services, etc.)

## Modes

### `/claude-md-manager create`

Generate the full AGENTS.md hierarchy from scratch.

**Steps:**

1. **Scan** project structure:
   - Find directories with `.go` files under `internal/`, `cmd/`, `data/`
   - Find skills in `.agents/skills/` (read each SKILL.md frontmatter)
   - Classify skills: generic (no project prefix) vs project-specific (has prefix)

2. **Evaluate granularity** per directory (see Granularity Algorithm below)

3. **Generate root AGENTS.md** as pure orchestrator with these sections only:

   ```
   # {Project Name}
   {One-line description}

   ## Como Usar Esta Guia
   - 3 bullet points: start here, layer docs, precedence rule
   ```

## Guias por Capa

| Capa | Ubicacion | AGENTS.md contiene |
{one row per directory with AGENTS.md}

## Skills Disponibles

### Skills Genericos

| Skill | Descripcion | Ruta |
{one row per generic skill}

### Skills de Proyecto

| Skill | Descripcion | Ruta |
{one row per project skill}

    ## Auto-invocacion
    | Accion | Referencia |
    {maps actions to AGENTS.md files AND skills}

## Stack

{one line}

## Comandos

{build, test, lint commands}

## Fases

{numbered list, current phase marked}

## Convenciones Globales

{5-6 bullet points max}

```

4. **Generate layer AGENTS.md** for each detected layer:
- Max ~300 lines (warn at 250, propose skill extraction if over 300)
- Only rules specific to that layer
- Reference relevant skills: "Para patrones Go, usa skill `golang-patterns`"
- Never duplicate root content (no stack, no conventions, no commands)

5. **Present each file** to user for confirmation before writing

### `/claude-md-manager audit`

Review existing AGENTS.md files and propose fixes.

**Checklist — Root AGENTS.md:**

| Check | Rule |
|-------|------|
| Exists | AGENTS.md in project root |
| Navigation index | No section exceeds 5 content lines (tables excluded) — every line that isn't a pointer is wasted context |
| Line count | Target ~150 lines (warn at 130, tables excluded) |
| How to Use | Section present with 3 bullet points |
| Available Skills — Generic | Every skill in `.agents/skills/` without project prefix is listed |
| Available Skills — Project | Every skill with project prefix is listed |
| Available Skills — Auto-invoke | Table maps actions to skills |
| Project Overview | Section present with component table |
| Development | At least one component with commands block |
| Commit & PR Guidelines | Section present |
| Phases | Current phase marked, YAGNI note present |
| Global Conventions | Section present, max 6 bullet points |

**Checklist — Root AGENTS.md — Template Structure:**

| Section | Required | Rule |
|---------|----------|------|
| `## How to Use This Guide` | Yes | 3 bullets: start here, component docs, override rule |
| `## Available Skills` | Yes | 3 subsections: Generic, Project-Specific, Auto-invoke |
| `## Project Overview` | Yes | One paragraph + component table (name, location, stack) |
| `## Development` | Yes | One block per component with setup/test/lint/run commands |
| `## Commit & Pull Request Guidelines` | Yes | Conventional-commit types listed |
| `## Phases` | Yes | Numbered list, current phase marked, YAGNI note |
| `## Global Conventions` | Yes | Max 6 bullets, no detailed technical rules |

**Checklist — Layer AGENTS.md:**

| Check | Rule |
|-------|------|
| Exists | AGENTS.md present in layer directory |
| Line count | <= 300 lines (warn at 250, propose skill extraction if over 300) |
| No duplication | Does not repeat root AGENTS.md content (conventions, stack, commands) |
| Skill references | Skills Reference block present at top |
| Auto-invoke | Layer-specific auto-invoke table present |
| Reflects code | Content matches current files and structures |

**Checklist — Layer AGENTS.md — Template Structure:**

| Section | Required | Rule |
|---------|----------|------|
| Skills Reference block | Yes | Blockquote at top listing relevant skills with links |
| `### Auto-invoke Skills` | Yes | Table mapping layer-specific actions to skills |
| `## CRITICAL RULES` | Yes | At least one concept block with ALWAYS/NEVER rules |
| `## DECISION TREES` | Yes | At least one decision tree |
| `## TECH STACK` | Yes | One-liner with versions |
| `## PROJECT STRUCTURE` | Yes | Directory tree or file format spec (for data/ layers) |
| `## COMMANDS` | Yes | Grouped by category (dev, test, lint) |
| `## QA CHECKLIST` | Yes | Checkboxes covering test, lint, layer-specific checks |
| `## NAMING CONVENTIONS` | Yes | Table with entity, pattern, example columns |

**Checklist — Skills consistency:**

| Check | Rule |
|-------|------|
| All registered | Every SKILL.md in `.agents/skills/` appears in root tables |
| All exist | Every skill referenced in auto-invoke has a SKILL.md file |
| Consistent prefix | Project skills use consistent prefix (`{project}-`) |
| No orphans | No skills exist without being referenced in root |

**Checklist — Granularity:**

Run the Granularity Algorithm on each directory with an AGENTS.md.

**Output format:**

```

=== AUDIT AGENTS.md ===

Root AGENTS.md: {OK|WARN|FAIL} ({n}/{total} checks)
Structure: {OK|FAIL} — missing sections: {list or "none"}
Content: {OK|WARN|FAIL} — {list warnings/failures}

{layer}/AGENTS.md: {OK|WARN|FAIL} ({n}/{total} checks)
Structure: {OK|FAIL} — missing sections: {list or "none"}
Content: {OK|WARN|FAIL} — {list warnings/failures}

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
| `.go` files (excluding `_test.go`) | > 8 | Too many concepts for one AGENTS.md |
| Root AGENTS.md line count | > 150 | Orchestrator has too much content — move to layer or skill |
| Layer AGENTS.md line count | > 300 | Too much content — extract to skill or sub-directory AGENTS.md |
| Sub-directories with `.go` files | >= 3 | Distinct patterns per sub-directory |

**Decision:**

```

0 of 3 thresholds exceeded → keep current AGENTS.md
1 of 3 thresholds exceeded → warn, suggest future review
2+ of 3 thresholds exceeded → propose split into sub-directory AGENTS.md files

````

When proposing a split:
- Suggest which sub-directories need their own AGENTS.md
- Draft the content for each new AGENTS.md
- Update the parent AGENTS.md to reference the new files
- Present all changes for confirmation

## Placement Rules

| Type | Location | Naming |
|------|----------|--------|
| Generic skills | `.agents/skills/{name}/SKILL.md` | No project prefix |
| Project skills | `.agents/skills/{project}-{scope}/SKILL.md` | With project prefix |
| Layer AGENTS.md | `internal/{layer}/AGENTS.md` | Standard |
| Data AGENTS.md | `data/{name}/AGENTS.md` | Standard |
| Root AGENTS.md | `AGENTS.md` | Project root |

## Token Efficiency Rules

These rules are the foundation of every decision this skill makes:

1. **Root = navigation index only.** Its only job is to point to the right AGENTS.md or skill as fast as possible. Every line that isn't a pointer is wasted context. Target ~150 lines.
2. **Layer AGENTS.md = rules only.** What to do, not how. Max ~300 lines. References skills for detailed patterns.
3. **Skills = patterns with examples.** How to do it. Self-contained. Does not repeat project structure.
4. **Never duplicate.** If info is in root, not in layer. If in a skill, not in AGENTS.md.
5. **Load path per action:** root (~50 lines) + 1 layer AGENTS.md (~50 lines) + 1 skill if needed (~100 lines) = ~200 lines total context.

## Architecture Assumptions

This skill is optimized for:
- **Go 1.22+** projects
- **Hexagonal / Clean Architecture** with layers: domain, application, infrastructure, presentation
- **Standard layout:** `internal/`, `cmd/`, `data/`, `.agents/skills/`
- **Single developer or small team** workflow

---

## Commands

```bash
/claude-md-manager create   # Generate full AGENTS.md hierarchy from scratch
/claude-md-manager audit    # Review existing AGENTS.md files and propose fixes
````

---

## Resources

- **Root AGENTS.md template**: See [assets/ROOT-AGENTS-TEMPLATE.md](assets/ROOT-AGENTS-TEMPLATE.md) for the monorepo root structure
- **Layer AGENTS.md template**: See [assets/LAYER-AGENTS-TEMPLATE.md](assets/LAYER-AGENTS-TEMPLATE.md) for component/layer structure (also covers `data/` layers)
- **Skill creator spec**: See [../skill-creator/SKILL.md](../skill-creator/SKILL.md) for skill authoring conventions
- **Skill template**: See [../skill-creator/assets/SKILL-TEMPLATE.md](../skill-creator/assets/SKILL-TEMPLATE.md) for the canonical SKILL.md template
