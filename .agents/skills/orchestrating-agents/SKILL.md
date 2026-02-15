---
name: orchestrating-agents
description: Orquestación de agentes especializados por capa en arquitectura hexagonal. Flujo completo desde brainstorming hasta commit final, delegando a domain-agent, application-agent e infrastructure-agent.
---

# Skill: Orchestrating Agents

Orquestación de agentes especializados por capa para proyectos con arquitectura hexagonal + vertical slices.

## Cuándo usar esta skill

- Cuando se necesita implementar una feature que cruza múltiples capas (domain, application, infrastructure)
- Cuando el trabajo es demasiado grande para un solo agente
- Cuando se quiere mantener el aislamiento entre capas durante el desarrollo

## Flujo completo

```
Usuario pide feature/cambio
         │
         ▼
┌─────────────────────────────────────┐
│         COORDINADOR                 │
│  1. Invocar skill `brainstorming`   │
│  2. Crear diseño + plan             │
│  3. Crear rama de trabajo           │
│  4. Despachar agentes en orden      │
│  5. Hacer wiring en main.go         │
│  6. Commit final                    │
└─────────────────────────────────────┘
         │
    ┌────┴────┬────────────┐
    ▼         ▼            ▼
domain-   application-  infrastructure-
agent     agent         agent
```

## Paso 1: Brainstorming inicial

Invocar `brainstorming` skill. El coordinador (este chat) refina la idea con el usuario, presenta diseño por secciones, obtiene aprobación.

**Output:** `docs/plans/YYYY-MM-DD-<feature>-design.md`

## Paso 2: Writing plans

Invocar `writing-plans` skill. El coordinador crea plan detallado con tareas para cada agente.

**Output:** `docs/plans/YYYY-MM-DD-<feature>-plan.md`

## Paso 3: Crear rama de trabajo

```bash
git checkout -b feature/nombre-de-la-feature
```

## Paso 4: Despachar agentes en orden

Orden obligatorio: **domain → application → infrastructure**

### Template para despachar agente

```
Sos el {domain-agent|application-agent|infrastructure-agent} de este proyecto. 
Tu trabajo es ejecutar {pasos específicos} del plan.

## Proyecto
Repositorio: {ruta absoluta}
Rama: {nombre de rama}
Módulo Go: {github.com/usuario/proyecto}

## Contexto — qué hicieron los agentes anteriores
{resumen de lo que ya existe y está testeado}

## Tu scope
{carpetas donde puede trabajar}

**NO toques** {carpetas prohibidas de otras capas}

## Plan a ejecutar
{ruta al archivo de plan}

## Instrucciones
1. Leé el plan completo
2. Creá tus propias tareas con TodoWrite antes de empezar
3. Ejecutá cada tarea marcando in_progress → completed
4. Verificá con `go test ./...` antes de terminar cada paso
5. Si algo falla, arreglalo antes de seguir

## Al terminar
Reportá:
- Lista exacta de archivos creados/modificados
- Output de `go test ./...`
- Issues encontrados (si hay)
- Próximos pasos sugeridos
```

### Ejemplo: Despachar domain-agent

```
Sos el domain-agent. Ejecutá los Pasos 1-2 del plan.

## Proyecto
Repositorio: C:\PROGRAMACION\mi-proyecto
Rama: feature/nueva-entidad
Módulo Go: github.com/usuario/mi-proyecto

## Contexto
Empezando desde cero. No hay agentes previos.

## Tu scope
- internal/shared/kernel/valueobject/ (si aplica)
- internal/{feature}/domain/entity/
- internal/{feature}/domain/service/

**NO toques**
- internal/{feature}/application/
- internal/{feature}/infrastructure/
- cmd/api/main.go

## Plan
 docs/plans/2026-02-15-mi-feature-plan.md

## Instrucciones
... (completar)
```

### Ejemplo: Despachar application-agent

```
Sos el application-agent. Ejecutá el Paso 3 del plan.

## Proyecto
...

## Contexto — qué hizo domain-agent
Ya están creados y testeados:
- internal/shared/kernel/valueobject/
- internal/{feature}/domain/entity/
- internal/{feature}/domain/service/

Los imports correctos que debés usar:
- Value objects: github.com/usuario/proyecto/internal/shared/kernel/valueobject
- Entities: github.com/usuario/proyecto/internal/{feature}/domain/entity

## Tu scope
- internal/{feature}/application/port/
- internal/{feature}/application/usecase/
- internal/{feature}/application/dto/

**NO toques**
- internal/{feature}/domain/
- internal/{feature}/infrastructure/
- cmd/api/main.go

## Plan
...
```

### Ejemplo: Despachar infrastructure-agent

```
Sos el infrastructure-agent. Ejecutá el Paso 4 del plan.

## Proyecto
...

## Contexto — qué hicieron los agentes anteriores
Ya están creados y testeados:
- Domain completo
- Application completo (ports, use cases, DTOs)

Los ports que debés implementar están en:
- internal/{feature}/application/port/

## Tu scope
- internal/{feature}/infrastructure/adapter/driver/
- internal/{feature}/infrastructure/adapter/driven/

**NO toques**
- internal/{feature}/domain/
- internal/{feature}/application/
- cmd/api/main.go (excepto si te lo pide específicamente)

## Plan
...
```

## Paso 5: Coordinador hace tareas finales

Después de que todos los agentes terminen, el coordinador:

1. **Actualiza `cmd/api/main.go`** — wiring de dependencias
2. **Crea placeholders** para otras features si aplica
3. **Elimina carpetas viejas** si fue refactorización
4. **Actualiza AGENTS.md** raíz y de cada capa
5. **Verifica todo:** `go test ./... && go build ./... && go vet ./...`

## Paso 6: Commit

```bash
git add -A
git commit -m "feat: implement {feature} with vertical slices

- Create domain/ with entities and services
- Create application/ with ports, use cases, DTOs
- Create infrastructure/ with adapters
- Update main.go wiring
- All tests passing"
```

## Reglas críticas

1. **Nunca en main/master** — siempre crear rama primero
2. **Esperar al agente anterior** — no despachar en paralelo
3. **Un agente a la vez** — domain termina → application empieza
4. **Verificación obligatoria** — cada agente debe reportar tests verdes
5. **No tocar fuera del scope** — cada agente respeta sus límites

## Cómo evitar duplicación de código — RESPONSABILIDAD DEL ORQUESTADOR

⚠️ **CRÍTICO:** Los agentes especializados (domain, application, infrastructure) **NO SE CONOCEN ENTRE SÍ** por diseño. Cada agente solo ve su propia capa. El **ORQUESTADOR es el único** que tiene visión global y debe:

1. **Investigar** — Buscar lo que ya existe en todas las capas
2. **Decidir** — Estrategia: ¿crear nuevo o extender existente?
3. **Comunicar** — Informar claramente al subagente qué hacer

### Principio de aislamiento

```
┌─────────────────────────────────────────────┐
│              ORQUESTADOR                    │
│  • Investiga código existente               │
│  • Toma decisiones arquitectónicas          │
│  • Comunica estrategia a subagentes         │
└──────────┬──────────────┬───────────────────┘
           │              │
     ┌─────┘              └─────┐
     ▼                          ▼
┌──────────────┐          ┌──────────────┐
│ domain-agent │          │ application- │
│              │          │    agent     │
│ NO sabe que  │          │              │
│ existe app   │          │ NO sabe qué  │
│              │          │ creó domain  │
└──────────────┘          └──────────────┘
```

### Flujo del Orquestador — 3 Pasos

#### Paso 1: Investigar (ANTES de despachar cualquier agente)

```bash
# 1. Listar servicios de dominio existentes
ls internal/{feature}/domain/service/*.go 2>/dev/null || echo "No hay servicios"

# 2. Buscar TODOs sin implementar en use cases
rg "TODO|FIXME|XXX" internal/{feature}/application/usecase --type go

# 3. Buscar métodos que calculen/processen algo similar
rg -i "func.*[Cc]alcular|func.*[Pp]rocesar" internal/{feature} --type go

# 4. Buscar por conceptos del negocio
rg -i "potencia|corriente|amperaje|tension" internal/{feature}/domain --type go
```

#### Paso 2: Decidir (El orquestador toma la decisión)

| Situación | Decisión del orquestador |
|-----------|-------------------------|
| Ya existe servicio similar en domain | Extender servicio existente, no crear nuevo |
| Use case tiene TODO que encaja | Implementar TODO, no crear nuevo use case |
| No existe nada similar | Proceder a crear nuevo |

#### Paso 3: Comunicar (Instrucciones claras al subagente)

**Ejemplo de mal prompt (agente no sabe qué existe):**
```
❌ "Creá un servicio para calcular amperaje"
```

**Ejemplo de buen prompt (orquestador investigó, decidió y comunicó):**
```
✅ "Revisé el código y encontré que CalcularCorrienteUseCase tiene 
    calcularManualPotencia() sin implementar (línea 80). 
    
    Tu tarea: Implementar ese método usando el servicio de dominio 
    CalcularAmperajeNominalCircuito que ya existe. NO crees un use case nuevo."
```

### Regla de oro: Un concepto = Un lugar

| Concepto | ¿Dónde debe vivir? | Responsable |
|----------|-------------------|-------------|
| Cálculo matemático/fórmula | `domain/service/` | domain-agent |
| Orquestación de pasos | `application/usecase/` | application-agent |
| Mapeo HTTP/JSON | `infrastructure/adapter/driver/http/` | infrastructure-agent |

### Caso de estudio: Error real de duplicación

**❌ ERROR:** Orquestador despachó domain-agent para crear `CalcularAmperajeNominalCircuito`, pero NO verificó que `CalcularCorrienteUseCase.calcularManualPotencia()` tenía un TODO sin implementar.

**Resultado:** Duplicación de lógica de cálculo.

**✅ Solución correcta:**
1. Orquestador debería haber visto el TODO en el use case existente
2. Orquestador debería haber instruido al application-agent: "Implementá calcularManualPotencia usando el servicio de dominio que ya existe"
3. NO crear nuevo servicio de dominio si ya hay uno que puede usarse

**Template de contexto para orquestador:**

```
Sos el application-agent.

## Contexto completo (responsabilidad del orquestador)

Lo que existe en domain (hecho por domain-agent previo):
- Servicio: CalcularAmperajeNominalCircuito en domain/service/
- Recibe: potenciaWatts, tension, tipoCarga, sistemaElectrico, factorPotencia
- Retorna: Corriente

Lo que existe en application (estado actual):
- Use case: CalcularCorrienteUseCase
- Método pendiente: calcularManualPotencia() con TODO sin implementar

## Tu tarea

NO crees un nuevo use case. En lugar de eso:
1. Implementá el método calcularManualPotencia() existente
2. Usá el servicio de dominio CalcularAmperajeNominalCircuito
3. Mapeá los parámetros del DTO a los parámetros del servicio
```

### Checklist del orquestador antes de cada fase

**Antes de despachar domain-agent:**
- [ ] ¿Ya existe un servicio en domain/service/ que haga algo similar?
- [ ] Si SÍ → instruir al agente que extienda, no cree nuevo
- [ ] Si NO → proceder con domain-agent

**Antes de despachar application-agent:**
- [ ] ¿Hay TODOs en use cases existentes que encajen con lo que necesitamos?
- [ ] ¿Podemos usar servicios de dominio ya existentes?
- [ ] En el prompt al agente, incluir lista de servicios de dominio disponibles

**Antes de despachar infrastructure-agent:**
- [ ] ¿Ya existe un handler similar al que necesitamos?
- [ ] ¿Podemos extender un handler existente en lugar de crear nuevo?

## Ejemplo completo

Ver referencia en: `docs/examples/orchestrating-agents-example.md`

## Troubleshooting

### El agente reporta que no puede leer el plan
- Verificar que la ruta al archivo de plan sea absoluta
- Confirmar que el archivo existe con `cat` o similar

### Tests fallan después de mover archivos
- Revisar que todos los imports fueron actualizados
- Verificar que no hay imports circulares
- Correr `go mod tidy` si es necesario

### Un agente quiere tocar código de otra capa
- Detener inmediatamente
- Recordarle su scope definido en el prompt
- Si es necesario, volver al agente anterior para que complete algo

## Ver también

- `brainstorming` skill — paso 1 del flujo
- `writing-plans` skill — paso 2 del flujo
- `domain-agent`, `application-agent`, `infrastructure-agent` — prompts en `.opencode/agents/`
