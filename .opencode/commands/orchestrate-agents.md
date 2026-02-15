# Orchestrate Agents Command

## Description

Orquesta agentes especializados que ejecutan todo su ciclo de trabajo autónomamente: **brainstorming → writing-plans → executing**. El orquestador da características de alto nivel, los agentes piensan, planifican e implementan.

## Concepto Clave

```
ORQUESTADOR (coordinador)          AGENTES ESPECIALIZADOS
─────────────────────────          ─────────────────────
"Necesito crear Proyectos"    →    domain-agent:
                                   - Piensa (brainstorming-dominio)
                                   - Planifica (writing-plans-dominio)
                                   - Implementa (executing-plans-dominio)
                                   → "Listo, domain completo"

"Continuar con Application"   →    application-agent:
                                   - Analiza domain hecho
                                   - Piensa casos de uso
                                   - Planifica e implementa
                                   → "Listo, application completo"

"Continuar con Infra"         →    infrastructure-agent:
                                   - Implementa adapters
                                   → "Listo, infra completa"

ORQUESTADOR hace wiring en main.go y commit
```

## Usage

```bash
# Despachar domain-agent para crear dominio desde cero
orchestrate-agents --agent domain --feature proyectos

# Despachar application-agent (requiere domain ya hecho)
orchestrate-agents --agent application --feature proyectos

# Despachar infrastructure-agent (requiere application ya hecho)
orchestrate-agents --agent infrastructure --feature proyectos
```

## Parameters

| Parameter | Description | Required |
|-----------|-------------|----------|
| `--agent` | Tipo: `domain`, `application`, `infrastructure` | Yes |
| `--feature` | Nombre de la feature: `calculos`, `proyectos`, etc. | Yes |
| `--context` | Contexto adicional (objetivos, restricciones) | No |

## Flujo Completo

### Paso 1: Orquestador inicia feature

```bash
# Usuario solicita nueva feature
> Quiero agregar soporte para proyectos que agrupen memorias

# Orquestador (este chat):
1. Crear rama: git checkout -b feature/proyectos
2. Despachar domain-agent
```

### Paso 2: Domain-Agent (autónomo)

**Prompt enviado por orquestador:**

```
Sos el domain-agent.

Características deseadas por el usuario:
"Necesito una entidad Proyecto que agrupe memorias de cálculo. 
Cada proyecto tiene nombre, cliente, fecha de creación."

Feature: proyectos
Scope: internal/proyectos/domain/ + shared/kernel/ (si necesitás VOs nuevos)

Contexto disponible:
- shared/kernel ya existe con Corriente, Tension, etc.
- No hay dependencias previas para esta feature

Tu trabajo (flujo completo):
1. SKILL: brainstorming-dominio
   - Explorá requisitos
   - Identificá entidades, VOs, agregados
   - Presentá diseño al orquestador para aprobación

2. SKILL: writing-plans-dominio  
   - Creá tus propias tareas con TodoWrite
   - Planificá implementación

3. SKILL: executing-plans-dominio
   - Ejecutá cada tarea
   - Verificá con go test

Reglas:
- Domain nunca depende de Application ni Infrastructure
- Si tenés dudas sobre requisitos, preguntá al orquestador
- Reportá progreso después de cada fase
```

**Domain-agent ejecuta:**
- Invoca `brainstorming-dominio` → presenta diseño
- Espera aprobación del orquestador
- Invoca `writing-plans-dominio` → crea tareas
- Invoca `executing-plans-dominio` → implementa
- Reporta: "✅ Domain completo. Archivos: [...]. Tests: pass"

### Paso 3: Application-Agent (autónomo)

**Prompt enviado:**

```
Sos el application-agent.

Contexto (hecho por domain-agent):
- Entidad Proyecto: ID, Nombre, Cliente, FechaCreacion, Memorias[]
- VO: IDProyecto
- Repository interface en domain/repository/

Características deseadas:
"Necesito poder crear proyectos y agregar memorias de cálculo a ellos"

Feature: proyectos
Scope: internal/proyectos/application/

Tu trabajo (flujo completo):
1. SKILL: brainstorming-application
   - Analizá el dominio
   - Diseñá casos de uso
   - Definí ports y DTOs
   - Presentá al orquestador

2. SKILL: writing-plans-application
   - Creá tus tareas

3. SKILL: executing-plans-application
   - Implementá ports, use cases, DTOs
   - Verificá tests

Reglas:
- Use cases solo orquestan, sin lógica de negocio
- No expongas entidades de domain directamente
```

**Application-agent ejecuta:**
- Invoca `brainstorming-application` → presenta diseño
- Espera aprobación
- Invoca `writing-plans-application` → crea tareas
- Invoca `executing-plans-application` → implementa
- Reporta: "✅ Application completo. Ports: [...]. Tests: pass"

### Paso 4: Infrastructure-Agent (autónomo)

**Prompt enviado:**

```
Sos el infrastructure-agent.

Contexto (hecho por application-agent):
- Port: ProyectoRepository (Save, FindByID, List)
- Use Cases: CrearProyectoUseCase, AgregarMemoriaUseCase
- DTOs: CrearProyectoInput, ProyectoOutput

Características deseadas:
"Implementar con PostgreSQL. También CSV para tests. HTTP REST handlers."

Feature: proyectos
Scope: internal/proyectos/infrastructure/

Tu trabajo (flujo completo):
1. SKILL: brainstorming-infrastructure
   - Diseñá adapters (Postgres, CSV, HTTP)
   - Presentá al orquestador

2. SKILL: writing-plans-infrastructure
   - Creá tus tareas

3. SKILL: executing-plans-infrastructure
   - Implementá adapters
   - Verificá tests

Reglas:
- Implementá exactamente los ports de application
- Sin lógica de negocio
- Manejá errores técnicos
```

**Infrastructure-agent ejecuta:**
- Invoca `brainstorming-infrastructure` → presenta diseño
- Espera aprobación
- Invoca `writing-plans-infrastructure` → crea tareas
- Invoca `executing-plans-infrastructure` → implementa
- Reporta: "✅ Infrastructure completo. Adapters: [...]. Tests: pass"

### Paso 5: Orquestador finaliza

```bash
# Coordinador:
1. Actualizar cmd/api/main.go (wiring)
2. go test ./...
3. go build ./...
4. git add -A
5. git commit -m "feat: add proyectos feature"
```

## Interacción Agente ↔ Orquestador

### Durante Brainstorming

**Agente → Orquestador:**
```
✅ Brainstorming-dominio completado

Diseño propuesto:
- Entidad: Proyecto (ID, Nombre, Cliente, FechaCreacion, Estado)
- VO: IDProyecto (UUID), EstadoProyecto (enum: ACTIVO, ARCHIVADO)
- Repository: ProyectoRepository

Preguntas:
1. ¿El estado es necesario o siempre activo?
2. ¿Un proyecto puede tener múltiples clientes?

¿Aprobás este diseño?
```

**Orquestador → Agente:**
```
Respuestas:
1. Sí, agregar estado
2. No, un solo cliente por proyecto

✅ Diseño aprobado. Continuá con writing-plans.
```

### Durante Writing Plans

**Agente → Orquestador:**
```
✅ Writing-plans-dominio completado

Tareas creadas:
1. [ ] Crear VO IDProyecto con validación UUID
2. [ ] Crear VO EstadoProyecto enum
3. [ ] Crear entidad Proyecto con métodos
4. [ ] Crear ProyectoRepository interface
5. [ ] Tests de entidad
6. [ ] Tests de VOs

¿Aprobás este plan?
```

**Orquestador → Agente:**
```
✅ Plan aprobado. Ejecutá las tareas.
```

### Durante Execution

**Agente → Orquestador (updates):**
```
Progreso:
- [x] Tarea 1: IDProyecto completado
- [x] Tarea 2: EstadoProyecto completado
- [x] Tarea 3: Entidad Proyecto completado
- [ ] Tarea 4: Repository interface en progreso

Tests pasando hasta ahora: 8/8 ✅
```

**Agente → Orquestador (final):**
```
✅ Executing-plans-dominio completado

Archivos creados:
- internal/proyectos/domain/entity/proyecto.go
- internal/proyectos/domain/entity/id_proyecto.go
- internal/proyectos/domain/entity/estado_proyecto.go
- internal/proyectos/domain/repository/proyecto_repository.go
- internal/proyectos/domain/entity/proyecto_test.go
- internal/proyectos/domain/entity/vo_test.go

Tests: ✅ 12/12 pass
Coverage: 94%

Sin issues encontrados.
Listo para application-agent.
```

## Checklist del Orquestador

### ANTES de despachar el primer agente (Prevención de duplicación)

```bash
# 1. Buscar si ya existe funcionalidad similar
rg -i "concepto|calcular|procesar" internal/{feature} --type go

# 2. Buscar TODOs sin implementar que puedan reutilizarse
rg "TODO|FIXME|XXX" internal/{feature} --type go

# 3. Revisar servicios de dominio existentes
ls internal/{feature}/domain/service/*.go
```

**Preguntas clave:**
- [ ] ¿Ya existe un servicio que haga este cálculo/proceso?
- [ ] ¿Hay algún método con TODO que debería implementarse primero?
- [ ] ¿Estoy pidiendo crear duplicación en lugar de reutilizar?

### Antes de despachar cada agente:

- [ ] Rama creada y activa (`git checkout -b feature/x`)
- [ ] Agente anterior completó su trabajo (excepto domain-agent)
- [ ] Contexto claro: características deseadas
- [ ] Scope definido: qué carpetas puede tocar
- [ ] **Verificado que no hay duplicación potencial** ⚠️

### Durante la orquestación:

- [ ] Revisar diseño propuesto por agente (brainstorming)
- [ ] Aprobar o solicitar cambios
- [ ] Revisar plan de tareas (writing-plans)
- [ ] Esperar reporte de ejecución (executing-plans)
- [ ] Verificar tests reportados por agente

### Después de todos los agentes:

- [ ] Hacer wiring en `cmd/api/main.go`
- [ ] Correr `go test ./...` completo
- [ ] Correr `go build ./...`
- [ ] Hacer commit

## Reglas de Oro

1. **Un agente a la vez** — no paralelizar
2. **Esperar aprobación** — después de brainstorming
3. **Verificar tests** — cada agente reporta su estado
4. **No micromanejar** — los agentes crean sus propias tareas
5. **Preguntar es válido** — agentes pueden consultar al orquestador

## Troubleshooting

### Agente se queda trabado
- Preguntar: "¿Necesitás ayuda? ¿Hay algún blocker?"
- Si no responde, cancelar y reiniciar

### Diseño propuesto no es correcto
- Rechazar amablemente: "No, necesito que ajustes X porque..."
- Pedir reconsideración con constraints adicionales

### Tests fallan después de agente
- Pedir al agente que arregle: "Tests fallan, revisá el output"
- No arreglar manualmente (rompe el flujo)

### Agente quiere tocar otra capa
- Recordar scope: "Recordá que solo podés tocar domain/"
- Si insiste, detener y reportar

## Ver También

- `.opencode/agents/domain-agent.md` — configuración del agente de dominio
- `.opencode/agents/application-agent.md` — configuración del agente de aplicación
- `.opencode/agents/infrastructure-agent.md` — configuración del agente de infraestructura
- `.opencode/commands/orchestrate-agents-example.md` — ejemplo completo paso a paso
- `AGENTS.md` raíz — reglas del proyecto
