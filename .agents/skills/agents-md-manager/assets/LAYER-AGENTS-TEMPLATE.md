# {Component/Layer Name} - AI Agent Ruleset

> **LAYER FILE — target under ~300 lines. Warn at 250, extract to a skill or sub-directory AGENTS.md if over 300.**
> This file contains rules, decision trees, naming conventions and structure for this layer.
> No deep code examples here — those belong in a skill. No root content (stack, global conventions, commands) — those stay in root AGENTS.md.

> **Skills Reference**: For detailed patterns, use these skills:
> - [`{primary-skill}`](../.agents/skills/{primary-skill}/SKILL.md) - {What it covers}
> - [`{secondary-skill}`](../.agents/skills/{secondary-skill}/SKILL.md) - {What it covers}
> - [`{generic-skill}`](../.agents/skills/{generic-skill}/SKILL.md) - {What it covers}

### Auto-invoke Skills

When performing these actions, ALWAYS invoke the corresponding skill FIRST:

| Action | Skill |
|--------|-------|
| {Action description} | `{skill-name}` |
| {Action description} | `{skill-name}` |
| Creating a git commit | `{commit-skill}` |

---

## CRITICAL RULES - NON-NEGOTIABLE

### {Concept 1 — e.g., Models, Entities, Handlers}
- ALWAYS: {Rule that must always be followed}
- ALWAYS: {Rule that must always be followed}
- NEVER: {Anti-pattern that must never appear}
- NEVER: {Anti-pattern that must never appear}

### {Concept 2 — e.g., Services, Use Cases, Repositories}
- ALWAYS: {Rule that must always be followed}
- ALWAYS: {Rule that must always be followed}
- NEVER: {Anti-pattern that must never appear}

### {Concept 3 — e.g., Error Handling, Validation, Security}
- ALWAYS: {Rule that must always be followed}
- NEVER: {Anti-pattern that must never appear}

---

## DECISION TREES

### {Decision 1 — e.g., Which serializer?, Which layer owns this logic?}
```
{Condition A} → {Action A}
{Condition B} → {Action B}
Otherwise     → {Default action}
```

### {Decision 2 — e.g., Sync vs Async, Repository vs Service}
```
{Condition A} → {Action A}
{Condition B} → {Action B}
```

---

## TECH STACK

{Technology 1} {version} | {Technology 2} {version} | {Technology 3} {version}

---

## PROJECT STRUCTURE

```
{component-root}/
├── {dir-1}/               # {What it contains}
│   ├── {subdir}/          # {What it contains}
│   └── {file}.{ext}       # {What it contains}
├── {dir-2}/               # {What it contains}
└── {dir-3}/               # {What it contains}
```

> **Note for data/ layers**: Describe file formats, naming rules, and validation constraints here instead of or in addition to the directory tree. Example:
> - Files are CSV, UTF-8, comma-separated, header row required
> - Naming: `{tabla}_{nom_article}.csv` (e.g., `conductores_310_16.csv`)
> - Required columns: `{col1}`, `{col2}`, `{col3}`
> - Numeric values use `.` as decimal separator, no thousand separators

---

## COMMANDS

```bash
# {Category 1 — e.g., Development}
{command}   # {description}
{command}   # {description}

# {Category 2 — e.g., Testing}
{command}   # {description}
{command}   # {description}

# {Category 3 — e.g., Linting / Quality}
{command}   # {description}
```

---

## QA CHECKLIST

- [ ] {Test command} passes
- [ ] {Lint command} passes
- [ ] {Layer-specific check — e.g., migrations created if models changed}
- [ ] {Layer-specific check — e.g., new endpoints documented}
- [ ] Tests cover success and error cases

---

## NAMING CONVENTIONS

| Entity | Pattern | Example |
|--------|---------|---------|
| {Entity type} | `{Pattern}` | `{ConcreteExample}` |
| {Entity type} | `{Pattern}` | `{ConcreteExample}` |
| {Entity type} | `{Pattern}` | `{ConcreteExample}` |

---

## {COMPONENT-SPECIFIC SECTION — optional, remove if not needed}

> Add here anything that doesn't fit the sections above but is critical for this layer.
> Examples: API response format, CSV column spec, message queue contracts, auth flows.

```{json|go|python|yaml}
{example}
```
