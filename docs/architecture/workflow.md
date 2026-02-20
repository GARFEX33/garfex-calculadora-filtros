# Workflow de Desarrollo

Para cualquier feature o bugfix, seguir este flujo de skills en orden:

| Paso | Skill                    | Trigger                     | Que hace                                                                                                                          |
| ---- | ------------------------ | --------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| 1    | `brainstorming`          | Usuario pide feature/cambio | Refina ideas con preguntas, explora alternativas, presenta diseño por secciones para validar. Guarda documento de diseño.         |
| 2    | `writing-plans`          | Diseño aprobado             | Divide el trabajo en tareas pequeñas (2-5 min cada una). Cada tarea tiene: rutas exactas, código completo, pasos de verificación. |
| 3    | `executing-plans`        | Plan listo                  | Despacha subagente fresco por tarea con revisión de dos etapas (spec + calidad)                                                   |

**IMPORTANTE:** No saltear pasos. Si el usuario dice "agregá X", empezar con `brainstorming`, NO con código.

---

## 🔄 Workflow Completo: Desde Idea hasta Merge

### Fase 1: Diseño (Orquestador)
```
Usuario pide feature
    │
    ▼
brainstorming → writing-plans → Crear rama
```

### Fase 2: Implementación (Agentes especializados en orden)
```
domain-agent → application-agent → infrastructure-agent
    │                │                    │
    ▼                ▼                    ▼
 tests green    tests green         tests green
```

### Fase 3: Integración (Orquestador)
```
Wiring en main.go → go test ./... → ✅ Todo pasa
```

### Fase 4: Documentación PRE-merge (OBLIGATORIO)
```
Auditar AGENTS.md con agents-md-manager
    │
    ▼
¿Hay drift? ──Si──→ Aplicar correcciones → Commit
    │
   No
    │
    ▼
Merge feature a main
```

**⚠️ Importante:** Los cambios a AGENTS.md son parte del mismo PR/feature. NUNCA mergear sin sincronizar la documentación.

---

## Regla Anti-Duplicación (OBLIGATORIO) — RESPONSABILIDAD DEL ORQUESTADOR

⚠️ **Los agentes especializados NO se conocen entre sí.** El orquestador es el único con visión global de todas las capas y debe:

1. **Investigar** — Buscar lo que ya existe
2. **Decidir** — Extender vs crear nuevo
3. **Comunicar** — Instrucciones claras al subagente

### Flujo del Orquestador (antes de despachar agentes)

**Paso 1: Investigar**
```bash
ls internal/{feature}/domain/service/*.go 2>/dev/null
rg "TODO|FIXME|XXX" internal/{feature}/application/usecase --type go
rg -i "func.*[Cc]alcular" internal/{feature} --type go
```

**Paso 2: Decidir**
| Situación | Decisión |
|-----------|----------|
| Existe servicio similar | Extender, no crear nuevo |
| Use case tiene TODO | Implementar TODO primero |
| Nada similar | Crear nuevo |

**Paso 3: Comunicar (en el prompt al agente)**

❌ Mal: "Creá un servicio para calcular amperaje"

✅ Bien: "Implementá el método calcularManualPotencia() que tiene un TODO en 
          CalcularCorrienteUseCase. Usá el servicio CalcularAmperajeNominalCircuito 
          que ya existe en domain/service/. NO crees un use case nuevo."

### Checklist (orquestador)
- [ ] ¿Investigué qué ya existe en domain/ y application/?
- [ ] ¿Tomé la decisión de extender vs crear?
- [ ] ¿Comuniqué claramente al agente qué hacer y qué NO hacer?
- [ ] ¿Verifiqué si el cambio requiere actualizar AGENTS.md? (nuevo endpoint, nueva regla, nuevo agent, nuevo skill)

**Error real:** Orquestador despachó domain-agent para crear servicio nuevo sin verificar que el use case existente tenía un TODO sin implementar. Resultado: duplicación.

---

## Actualizacion de Documentacion

⚠️ **REGLA OBLIGATORIA:** Al terminar cada tarea, ANTES de hacer commit:
1. Ejecutar `git status` para ver archivos modificados
2. Si hay cambios en código (domain/application/infrastructure), verificar si corresponde actualizar:
   - AGENTS.md de la capa afectada
   - AGENTS.md raíz (si hay nuevos skills o agentes)
3. Actualizar AGENTS.md si es necesario
4. Luego hacer commit (incluyendo cambios de AGENTS.md)

** Esta regla es parte de la definition of done. NO hacer commit sin verificar AGENTS.md.**
