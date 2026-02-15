# Ejemplo: Orquestar Agentes (Flujo Completo)

Este ejemplo muestra cómo el orquestador delega TODO el ciclo de trabajo a los agentes, quienes piensan, planifican e implementan autónomamente.

## Escenario

**Usuario:** "Quiero agregar una feature `proyectos` para guardar memorias de cálculo agrupadas por proyecto. Cada proyecto tiene nombre, cliente, y puede tener múltiples memorias."

## Actores

- **Orquestador** (este chat) — coordina, pregunta al usuario, aprueba diseños
- **Domain-Agent** — piensa, planifica e implementa el dominio
- **Application-Agent** — piensa, planifica e implementa application
- **Infrastructure-Agent** — piensa, planifica e implementa infraestructura

---

## FASE 1: Orquestador Inicia

### Paso 1: Crear rama

```bash
git checkout -b feature/proyectos
```

### Paso 2: Despachar Domain-Agent

**Orquestador invoca:**
```bash
orchestrate-agents --agent domain --feature proyectos
```

**Prompt enviado:**

```
Sos el domain-agent.

Características deseadas por el usuario:
"Quiero agregar una feature 'proyectos' para guardar memorias de cálculo 
agrupadas por proyecto. Cada proyecto tiene nombre, cliente, y puede tener 
múltiples memorias."

Feature: proyectos
Scope: internal/proyectos/domain/ + shared/kernel/ (si necesitás VOs nuevos)

Contexto disponible:
- shared/kernel/ ya existe con Corriente, Tension, MaterialConductor, etc.
- No hay dependencias previas para esta feature
- La feature 'calculos' ya existe con MemoriaCalculo

Tu trabajo completo:
1. SKILL: brainstorming-dominio
   - Explorá: ¿qué entidades, VOs, agregados necesitás?
   - Pensá en invariants y reglas de negocio
   - Diseñá el modelo de dominio
   - PRESENTÁ el diseño al orquestador para aprobación

2. SKILL: writing-plans-dominio
   - Creá tus propias tareas con TodoWrite
   - Planificá la implementación paso a paso

3. SKILL: executing-plans-dominio
   - Implementá cada tarea
   - Verificá con go test después de cada una

Reglas:
- Domain nunca depende de Application ni Infrastructure
- Si tenés dudas sobre requisitos, preguntá al orquestador
- Reportá progreso después de cada fase (brainstorming, writing, executing)
```

---

## FASE 2: Domain-Agent Trabaja (Autónomo)

### Brainstorming-Dominio

**Domain-Agent invoca skill:** `brainstorming-dominio`

**Domain-Agent → Orquestador:**
```
✅ Brainstorming-dominio completado

DISEÑO PROPUESTO:

Entidades:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Proyecto (Aggregate Root)
├── ID: IDProyecto (VO)
├── Nombre: string (validado: no vacío, max 100 chars)
├── Cliente: string (validado: no vacío)
├── FechaCreacion: time.Time
├── Estado: EstadoProyecto (VO: enum ACTIVO, ARCHIVADO, CANCELADO)
└── Memorias: []MemoriaAsociada

Value Objects:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
- IDProyecto: UUID con validación
- EstadoProyecto: enum con métodos de transición
- MemoriaAsociada: referencia a MemoriaCalculo (solo ID + FechaAgregada)

Domain Services:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
- (ninguno necesario por ahora, lógica en entidad)

Repository Interface:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
- ProyectoRepository en domain/repository/
- Métodos: Save, FindByID, List, Update, Delete

Reglas de Negocio Identificadas:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Nombre no puede estar vacío
2. Cliente no puede estar vacío
3. FechaCreacion se setea automáticamente al crear
4. Estado inicial siempre ACTIVO
5. Solo proyectos ACTIVOS pueden recibir nuevas memorias

PREGUNTAS AL ORQUESTADOR:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. ¿Una memoria puede estar en múltiples proyectos? (Sugiero: NO)
2. ¿Necesitamos histórico de cambios de estado? (Sugiero: NO por ahora)
3. ¿El cliente es solo string o entidad separada? (Sugiero: string simple)

¿Aprobás este diseño? ¿Tenés cambios?
```

**Orquestador → Domain-Agent:**
```
Respuestas:
1. NO, una memoria solo en un proyecto ✓
2. NO, sin histórico por ahora ✓
3. Cliente como string simple ✓

Diseño APROBADO ✅

Continuá con writing-plans-dominio.
```

### Writing-Plans-Dominio

**Domain-Agent invoca skill:** `writing-plans-dominio`

**Domain-Agent → Orquestador:**
```
✅ Writing-plans-dominio completado

PLAN DE TAREAS:

Phase 1: Value Objects
- [ ] T1: Crear VO IDProyecto (UUID + validación)
- [ ] T2: Crear VO EstadoProyecto (enum + métodos)
- [ ] T3: Tests de VOs

Phase 2: Entidad Proyecto
- [ ] T4: Crear struct Proyecto
- [ ] T5: Implementar constructor NewProyecto
- [ ] T6: Implementar método AgregarMemoria
- [ ] T7: Implementar método Archivar
- [ ] T8: Tests de entidad

Phase 3: Repository
- [ ] T9: Crear interface ProyectoRepository
- [ ] T10: Definir errores de dominio

Verificación: go test ./internal/proyectos/domain/...

¿Aprobás este plan de tareas?
```

**Orquestador → Domain-Agent:**
```
✅ Plan APROBADO

Ejecutá las tareas con executing-plans-dominio.
```

### Executing-Plans-Dominio

**Domain-Agent invoca skill:** `executing-plans-dominio`

**Domain-Agent → Orquestador (updates en progreso):**
```
Progreso de ejecución:

Phase 1: Value Objects
✅ T1: IDProyecto creado (uuid validation, tests pass)
✅ T2: EstadoProyecto creado (enum con CanAgregarMemoria())
✅ T3: Tests VOs - 6/6 pass

Phase 2: Entidad Proyecto (en progreso)
✅ T4: Struct Proyecto creado
✅ T5: Constructor NewProyecto con validaciones
✅ T6: AgregarMemoria con regla de estado
⏳ T7: Método Archivar (implementando)
⏳ T8: Tests de entidad (pendiente)

Tests actuales: 10/10 pass ✅
```

**Domain-Agent → Orquestador (final):**
```
✅ Executing-plans-dominio COMPLETADO

ARCHIVOS CREADOS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
internal/proyectos/domain/
├── entity/
│   ├── proyecto.go              (entidad + métodos)
│   ├── proyecto_test.go         (tests de comportamiento)
│   ├── id_proyecto.go           (VO ID)
│   ├── estado_proyecto.go       (VO Estado)
│   └── vo_test.go               (tests de VOs)
├── repository/
│   └── proyecto_repository.go   (interface)
└── errors.go                    (ErrProyectoNoEncontrado, etc.)

TEST RESULTS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
go test ./internal/proyectos/domain/...
✅ PASS: 18/18 tests
✅ Coverage: 96%
✅ Race detector: pass

DECISIONES TOMADAS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
- Agregué validación de nombre único (no mostrada en diseño inicial)
- Agregué método CanAgregarMemoria() para encapsular regla
- Error específico ErrProyectoArchivado para intentos de agregar

Sin issues bloqueantes.

🎯 LISTO PARA APPLICATION-AGENT
```

---

## FASE 3: Despachar Application-Agent

**Orquestador invoca:**
```bash
orchestrate-agents --agent application --feature proyectos
```

**Prompt enviado:**

```
Sos el application-agent.

CONTEXTO (Domain-Agent completó):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Entidades:
- Proyecto: ID, Nombre, Cliente, FechaCreacion, Estado, Memorias[]
- VOs: IDProyecto, EstadoProyecto, MemoriaAsociada

Repository Interface:
- ProyectoRepository: Save, FindByID, List, Update, Delete

Reglas de Negocio ya implementadas:
- Validaciones en constructor
- Solo ACTIVOS pueden recibir memorias
- Método Archivar() cambia estado

Características deseadas por usuario:
"Necesito poder crear proyectos y agregar memorias de cálculo a ellos"

Feature: proyectos
Scope: internal/proyectos/application/

Tu trabajo completo:
1. SKILL: brainstorming-application
   - Analizá el dominio
   - Diseñá casos de uso necesarios
   - Definí ports (driver y driven)
   - Diseñá DTOs
   - PRESENTÁ al orquestador

2. SKILL: writing-plans-application
   - Creá tus tareas

3. SKILL: executing-plans-application
   - Implementá ports, use cases, DTOs
   - Verificá tests

Reglas:
- Use cases solo orquestan (sin lógica de negocio)
- No expongas entidades de domain directamente
- Usa DTOs para entrada/salida
```

### Brainstorming-Application

**Application-Agent → Orquestador:**
```
✅ Brainstorming-application completado

CASOS DE USO IDENTIFICADOS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. CrearProyectoUseCase
   Input: nombre, cliente
   Output: proyecto creado con ID

2. AgregarMemoriaAProyectoUseCase
   Input: proyectoID, memoriaID
   Output: confirmación
   Errores: proyecto no existe, proyecto archivado

3. ObtenerProyectoUseCase
   Input: proyectoID
   Output: proyecto con sus memorias

4. ListarProyectosUseCase
   Input: filtros (opcional)
   Output: lista de proyectos

5. ArchivarProyectoUseCase
   Input: proyectoID
   Output: confirmación

PORTS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Driver (entry points):
- (ninguno explicito, los handlers HTTP llaman use cases)

Driven (dependencias):
- ProyectoRepository (ya existe en domain)

DTOs:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Input:
- CrearProyectoInput
- AgregarMemoriaInput
- ProyectoIDInput

Output:
- ProyectoOutput (sin lógica, campos planos)
- ProyectoListOutput

¿Aprobás estos casos de uso?
```

**Orquestador → Application-Agent:**
```
✅ Casos de uso APROBADOS

Agregué uno más:
6. CambiarNombreProyectoUseCase (para edición)

Continuá con writing-plans.
```

### Writing + Executing Plans

**Application-Agent ejecuta** (resumen):**

```
✅ Todos los pasos completados

Archivos creados:
- port/
  ├── proyecto_repository.go (re-export de domain)
  └── errors.go
- usecase/
  ├── crear_proyecto.go
  ├── crear_proyecto_test.go
  ├── agregar_memoria.go
  ├── agregar_memoria_test.go
  ├── obtener_proyecto.go
  ├── listar_proyectos.go
  ├── archivar_proyecto.go
  └── cambiar_nombre.go
- dto/
  ├── proyecto_input.go
  └── proyecto_output.go

Tests: 24/24 pass ✅

🎯 LISTO PARA INFRASTRUCTURE-AGENT
```

---

## FASE 4: Despachar Infrastructure-Agent

**Orquestador invoca:**
```bash
orchestrate-agents --agent infrastructure --feature proyectos
```

**Contexto enviado:**
```
CONTEXTO (Application-Agent completó):
- Ports: ProyectoRepository interface
- Use Cases: 6 casos de uso implementados
- DTOs: Inputs y outputs definidos

Características deseadas:
"Implementar con PostgreSQL para producción y CSV para tests.
HTTP REST handlers para la API."

Tu trabajo completo (brainstorm → write → execute)
```

### Brainstorming-Infrastructure

**Infrastructure-Agent → Orquestador:**
```
✅ Brainstorming-infrastructure completado

ADAPTERS PROPUESTOS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Driven (implementan ProyectoRepository):
1. PostgresProyectoRepository
   - Tabla: proyectos (id, nombre, cliente, fecha_creacion, estado)
   - Tabla: proyecto_memorias (proyecto_id, memoria_id, fecha_agregada)
   - pgx/v5 para conexión

2. InMemoryProyectoRepository (para tests)
   - map[string]Proyecto
   - thread-safe con sync.RWMutex

Driver (HTTP Handlers):
1. ProyectoHandler
   - POST /api/v1/proyectos → CrearProyectoUseCase
   - GET /api/v1/proyectos/{id} → ObtenerProyectoUseCase
   - GET /api/v1/proyectos → ListarProyectosUseCase
   - POST /api/v1/proyectos/{id}/memorias → AgregarMemoriaUseCase
   - PATCH /api/v1/proyectos/{id}/archivar → ArchivarProyectoUseCase
   - PATCH /api/v1/proyectos/{id}/nombre → CambiarNombreUseCase

Config:
- ProyectoConfig con DB_HOST, DB_PORT, DB_NAME, etc.

¿Aprobás este diseño?
```

**Orquestador → Infrastructure-Agent:**
```
✅ Diseño APROBADO

Notas:
- Usar transacciones para operaciones complejas
- Manejar errores de DB con wrapping

Continuá.
```

### Ejecución (resumen)

```
✅ Infrastructure completo

Archivos:
- adapter/driven/postgres/
  ├── proyecto_repository.go
  └── proyecto_repository_test.go (integración)
- adapter/driven/memory/
  └── proyecto_repository.go
- adapter/driver/http/
  ├── proyecto_handler.go
  └── proyecto_handler_test.go
- config.go

Tests: 18/18 pass ✅
(Incluye tests de integración con testcontainers)

🎯 LISTO. Orquestador debe hacer wiring.
```

---

## FASE 5: Orquestador Finaliza

### Wiring en main.go

```go
// cmd/api/main.go

import (
    // ... otros imports
    proyectosapp "github.com/garfex/calculadora-filtros/internal/proyectos/application/usecase"
    proyectosinfra "github.com/garfex/calculadora-filtros/internal/proyectos/infrastructure/adapter/driven/postgres"
    proyectoshttp "github.com/garfex/calculadora-filtros/internal/proyectos/infrastructure/adapter/driver/http"
)

func main() {
    // ... repos existentes
    
    // Proyectos
    proyectoRepo := proyectosinfra.NewPostgresProyectoRepository(db)
    crearProyectoUC := proyectosapp.NewCrearProyectoUseCase(proyectoRepo)
    agregarMemoriaUC := proyectosapp.NewAgregarMemoriaUseCase(proyectoRepo)
    // ... otros use cases
    
    // Handlers
    proyectoHandler := proyectoshttp.NewProyectoHandler(
        crearProyectoUC,
        agregarMemoriaUC,
        // ...
    )
    
    // Router
    router := gin.New()
    proyectoHandler.RegisterRoutes(router)
    // ...
}
```

### Verificación Final

```bash
go test ./...
# ✅ PASS: todos los tests de todas las capas

go build ./...
# ✅ Sin errores

git add -A
git commit -m "feat: add proyectos feature with full vertical slices

Domain:
- Proyecto aggregate with ID, Nombre, Cliente, Estado
- Value objects: IDProyecto, EstadoProyecto
- Repository interface

Application:
- 6 use cases: Crear, AgregarMemoria, Obtener, Listar, Archivar, CambiarNombre
- DTOs for all operations
- Ports clearly defined

Infrastructure:
- PostgresProyectoRepository with migrations
- InMemory repository for testing
- HTTP REST handlers with gin
- Full test coverage including integration tests

All tests passing: 60/60 ✅"
```

---

## Timeline del Proceso

| Fase | Actor | Duración | Output |
|------|-------|----------|--------|
| 1 | Orquestador | 5 min | Rama creada, domain-agent despachado |
| 2a | Domain-Agent | 15 min | Diseño aprobado |
| 2b | Domain-Agent | 10 min | Plan de tareas aprobado |
| 2c | Domain-Agent | 30 min | Domain implementado y testeado ✅ |
| 3a | Application-Agent | 10 min | Casos de uso diseñados |
| 3b | Application-Agent | 45 min | Application implementada ✅ |
| 4a | Infrastructure-Agent | 10 min | Adapters diseñados |
| 4b | Infrastructure-Agent | 60 min | Infrastructure implementada ✅ |
| 5 | Orquestador | 15 min | Wiring, tests, commit |

**Total: ~3 horas** para feature completa con 3 capas, testeada.

---

## Lecciones Aprendidas

1. **Los agentes son autónomos** — el orquestador solo aprueba/ajusta
2. **Brainstorming es crucial** — evita retrabajo posterior
3. **Tests en cada capa** — asegura calidad antes de pasar al siguiente agente
4. **Comunicación clara** — el orquestador debe responder rápido a preguntas
5. **Scope bien definido** — cada agente sabe exactamente qué no tocar
