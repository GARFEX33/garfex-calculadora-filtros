# Design: claude-md-manager Skill

**Fecha:** 2026-02-12
**Estado:** Validado
**Tipo:** Skill en `.agents/skills/claude-md-manager/SKILL.md`

---

## Resumen

Skill para crear y auditar la jerarquia de CLAUDE.md en proyectos Go con arquitectura hexagonal. Sigue el patron de Prowler: root como orquestador puro, CLAUDE.md por capa, skills genericos y de proyecto, tabla de auto-invocacion unificada.

**Principio central:** Token efficiency = precision. Cada archivo carga SOLO lo que necesita.

---

## Modos de Uso

- `/claude-md-manager create` — escanea el repo, genera la jerarquia completa
- `/claude-md-manager audit` — revisa CLAUDE.md existentes, propone correcciones

---

## Estructura que Genera el Skill

### Root CLAUDE.md (~50 lineas, orquestador puro)

```
1. Como Usar Esta Guia (3 lineas)
2. Guias por Capa (tabla: capa → CLAUDE.md)
3. Skills Disponibles
   3a. Skills Genericos (tabla: skill → descripcion → ruta)
   3b. Skills de Proyecto (tabla: skill → descripcion → ruta)
4. Auto-invocacion (tabla unificada: accion → CLAUDE.md O skill)
5. Stack (1 linea), Comandos (6 lineas), Fases (3 lineas), Convenciones (5 lineas)
```

**Regla:** CERO contenido tecnico detallado en root. Solo tablas de direcciones.

### CLAUDE.md de Capa (~30-50 lineas)

- Solo reglas de ESA capa
- NO duplica contenido del root
- Referencia skills relevantes: "Para patrones Go, usa skill `golang-patterns`"
- Dice "que hacer", no "como hacerlo con ejemplos"

### Skills (~50-200 lineas)

- Autocontenidos: patrones detallados con ejemplos de codigo
- NO repiten estructura del proyecto
- Dicen "como hacerlo con ejemplos"

### Regla anti-duplicacion

Si algo esta en root, no va en CLAUDE.md de capa.
Si algo esta en un skill, no va en CLAUDE.md.

---

## Jerarquia de Carga (Token Efficiency)

```
Paso 1: Claude Code lee root CLAUDE.md (~50 lineas)
        → Solo tablas de direcciones

Paso 2: Segun la accion, carga UN CLAUDE.md de capa (~40-60 lineas)
        → Solo reglas de esa capa

Paso 3: Si necesita un skill, carga ESE skill (~50-200 lineas)
        → Solo el patron especifico
```

Total por accion: ~140-310 lineas vs ~800+ si todo estuviera en un archivo.

---

## Estructura de Carpetas

```
proyecto/
├── CLAUDE.md                          ← Orquestador puro
├── .agents/
│   └── skills/
│       ├── golang-patterns/           ← Generico: patrones Go
│       │   └── SKILL.md
│       ├── golang-pro/                ← Generico: Go avanzado
│       │   └── SKILL.md
│       ├── api-design-principles/     ← Generico: diseno API
│       │   └── SKILL.md
│       ├── garfex-domain/             ← Proyecto: entidades, VOs, formulas
│       │   └── SKILL.md
│       ├── garfex-infrastructure/     ← Proyecto: repos, CSV, BD
│       │   └── SKILL.md
│       └── claude-md-manager/         ← Meta-skill: este mismo
│           └── SKILL.md
├── internal/
│   ├── domain/
│   │   └── CLAUDE.md                  ← Reglas de capa + ref a skills
│   ├── application/
│   │   └── CLAUDE.md
│   ├── infrastructure/
│   │   └── CLAUDE.md
│   └── presentation/
│       └── CLAUDE.md
└── data/
    └── tablas_nom/
        └── CLAUDE.md
```

### Reglas de Placement

1. **Skills genericos** → `.agents/skills/{nombre}/SKILL.md` — sin prefijo de proyecto
2. **Skills de proyecto** → `.agents/skills/{proyecto}-{capa}/SKILL.md` — con prefijo
3. **CLAUDE.md de capa** → `internal/{capa}/CLAUDE.md` — reglas + referencia a skills
4. **Root CLAUDE.md** → raiz — solo tablas de direcciones

---

## Algoritmo de Granularidad Adaptativa

### Metricas de Evaluacion

| Metrica | Umbral | Significado |
|---------|--------|-------------|
| Archivos .go en directorio (sin _test.go) | > 8 | Demasiados conceptos para un CLAUDE.md |
| Lineas del CLAUDE.md | > 80 | Excede procesamiento eficiente |
| Sub-directorios con codigo | >= 3 | Patrones distintos por sub-directorio |

### Decision

```
0 de 3 metricas → mantener CLAUDE.md actual
1 de 3 metricas → advertir, sugerir revision futura
2+ de 3 metricas → proponer split en sub-directorios
```

### Comportamiento por Modo

- **CREATE:** Escanea → evalua → genera estructura minima necesaria. Nunca crea CLAUDE.md en sub-directorio si no pasa umbral.
- **AUDIT:** Re-evalua metricas contra estado actual. Si proyecto crecio, propone split con secciones sugeridas. Espera confirmacion.

---

## Modo AUDIT — Checklists

### Root CLAUDE.md

```
[ ] Existe CLAUDE.md en raiz
[ ] Tiene seccion "Como Usar Esta Guia"
[ ] Tiene tabla "Guias por Capa" — cada directorio con CLAUDE.md listado
[ ] Tiene tabla "Skills Genericos" — cada skill sin prefijo proyecto
[ ] Tiene tabla "Skills de Proyecto" — cada skill con prefijo proyecto
[ ] Tiene tabla "Auto-invocacion" — acciones mapeadas a CLAUDE.md + skills
[ ] NO tiene contenido detallado (> 5 lineas por seccion = advertencia)
[ ] Stack, Comandos, Fases, Convenciones presentes y concisos
```

### CLAUDE.md de Capa

```
[ ] Existe el archivo
[ ] <= 80 lineas (si no, sugerir split)
[ ] NO duplica contenido del root
[ ] Referencia skills relevantes
[ ] Contenido refleja el codigo actual
```

### Consistencia Skills ↔ CLAUDE.md

```
[ ] Todo skill en .agents/skills/ aparece en root CLAUDE.md
[ ] Todo skill referenciado en auto-invocacion existe como archivo
[ ] Skills de proyecto tienen prefijo consistente
[ ] No hay skills huerfanos
```

### Output del Audit

```
=== AUDIT CLAUDE.md ===

Root CLAUDE.md: OK (8/8 checks)
internal/domain/CLAUDE.md: OK (5/5 checks)
...

Skills: 1 advertencia
  WARN: golang-pro no esta en tabla de root

Granularidad: OK

=== PROPUESTAS (N) ===
1. Descripcion del cambio propuesto
Aplicar propuesta 1? [si/no]
```

---

## Modo CREATE — Flujo

1. Escanear estructura (`internal/`, `cmd/`, `data/`, `.agents/skills/`)
2. Detectar capas existentes (por presencia de archivos .go)
3. Detectar skills existentes (por presencia de SKILL.md)
4. Evaluar granularidad (algoritmo de 3 metricas)
5. Generar root CLAUDE.md como orquestador
6. Por cada capa detectada, generar su CLAUDE.md
7. Presentar cada archivo al usuario para confirmacion antes de escribir

---

## Decisiones Clave

1. **Un solo archivo SKILL.md** — sin assets ni references (YAGNI)
2. **Go + hexagonal only** — optimizado para este stack, no generico
3. **Proponer + confirmar** — nunca escribe sin aprobacion del usuario
4. **Granularidad adaptativa** — minimalista hoy, sugiere split cuando crece
5. **Token efficiency como principio** — toda la arquitectura optimiza contexto minimo
6. **Patron Prowler** — root orquestador con skills genericos + proyecto + auto-invocacion
7. **CLAUDE.md = que hacer, Skill = como hacerlo** — separacion clara de responsabilidades
