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
