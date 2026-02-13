# {Project Name}

{One-line description of what this project does.}

> **ORCHESTRATOR FILE — this is a navigation index, not a technical document.**
> Its only job is to point to the right AGENTS.md or skill as fast as possible.
> Every line that isn't a pointer is wasted context. Target ~150 lines (warn at 130).
> Rules, examples, conventions, commands → belong in a component AGENTS.md or a skill.

## How to Use This Guide

- Start here for cross-project norms. {Project Name} is a monorepo with several components.
- Each component has an `AGENTS.md` file with specific guidelines (e.g., `{component-1}/AGENTS.md`, `{component-2}/AGENTS.md`).
- Component docs override this file when guidance conflicts.

---

## Available Skills

Use these skills for detailed patterns on-demand:

### Generic Skills (Any Project)

| Skill | Description | URL |
|-------|-------------|-----|
| `{skill-name}` | {One-line description} | [SKILL.md](.agents/skills/{skill-name}/SKILL.md) |

### Project-Specific Skills

| Skill | Description | URL |
|-------|-------------|-----|
| `{project}-{scope}` | {One-line description} | [SKILL.md](.agents/skills/{project}-{scope}/SKILL.md) |

### Auto-invoke Skills

When performing these actions, ALWAYS invoke the corresponding skill FIRST:

| Action | Skill |
|--------|-------|
| {Action description} | `{skill-name}` |
| Creating new skills | `skill-creator` |
| Creating/auditing AGENTS.md | `agents-md-manager` |

---

## Project Overview

{One paragraph describing what the project does and who it's for.}

| Component | Location | Tech Stack |
|-----------|----------|------------|
| {Component 1} | `{path}/` | {Stack} |
| {Component 2} | `{path}/` | {Stack} |

---

## Development

### {Component 1} — {Tech Stack}

```bash
# Setup
{setup command}

# Tests
{test command}

# Lint
{lint command}

# Run dev server
{dev command}
```

### {Component 2} — {Tech Stack}

```bash
# Setup
{setup command}

# Tests
{test command}

# Lint
{lint command}

# Run dev server
{dev command}
```

---

## Commit & Pull Request Guidelines

Follow conventional-commit style: `<type>[scope]: <description>`

**Types:** `feat`, `fix`, `docs`, `chore`, `perf`, `refactor`, `style`, `test`

Before creating a PR:
1. Complete checklist in `.github/pull_request_template.md`
2. Run all relevant tests and linters
3. Link screenshots for UI changes

---

## Phases

1. **Phase 1 (current):** {Description of current phase scope}
2. **Phase 2:** {Description of next phase}
3. **Phase 3:** {Description of future phase}

**IMPORTANT:** Do not get ahead. Only implement what is needed for the current phase.

---

## Global Conventions

- **Business names in {language}** — domain terms use the project's domain language
- **Code in idiomatic {language}** — packages, internal variables, functions
- **Errors:** wrap with context, never swallow silently
- **Tests:** table-driven with subtests, one assertion per test case
- **No panic** in library/domain code — return errors instead
- **YAGNI** — current phase only: no speculative features
