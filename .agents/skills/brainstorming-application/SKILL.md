---
name: brainstorming-application
description: "DEBES usar esto solo para diseñar la capa de Application según la arquitectura Hexagonal. No puede crear ni pensar en dominio; si necesita elementos de dominio que no existen, deberá solicitarlo a domain-agent y esperar la información."
---

# Convertir Diseño de Domain en Diseño de Application Layer

## Descripción General

Este agente se enfoca exclusivamente en la **capa de Application**:

- Define **use cases**, **ports** (driver/driven) y **DTOs**.
- Coordina la lógica de dominio ya existente.
- Nunca crea ni modifica entidades, agregados, value objects o reglas de negocio.
- Nunca piensa en implementación de infraestructura o detalles técnicos externos.
- Si necesita información de dominio que no existe, **detiene su proceso y solicita al domain-agent** la creación o aclaración antes de continuar.

<HARD-GATE>
NO invocar ninguna skill de dominio directamente, NO diseñar entidades ni agregados, NO crear reglas de negocio.
Si se necesita un elemento de dominio inexistente → invocar domain-agent y esperar la respuesta antes de continuar.
</HARD-GATE>

<REPOSITORY-CONSTRAINT>
No usar `.worktree`, `git worktree` ni estrategias basadas en múltiples worktrees.
</REPOSITORY-CONSTRAINT>

## Anti-Patrón: "Application se inventa dominio"

Nunca asumir ni crear lógica de dominio. Todo elemento de dominio necesario debe provenir de Domain o ser solicitado explícitamente a domain-agent.

## Checklist de Brainstorming de Application

1. **Usar solo lo disponible en Domain** — Entidades, agregados, value objects y errores existentes.
2. **Hacer preguntas aclaratorias** — Sobre flujos de negocio, inputs/outputs y coordinación de lógica existente.
3. **Si falta información de Domain** — Invocar `domain-agent` y **detenerse hasta recibir la información**.
4. **Proponer 2-3 enfoques de diseño de Application** — Incluyendo use cases, driver/driven ports y DTOs.
5. **Presentar diseño de Application** — Confirmar boundaries, inputs/outputs y dependencias.
6. **Escribir documento de diseño de Application** — Guardarlo en `docs/plans/YYYY-MM-DD-<tema>-application-design.md` y hacer commit.
7. **Transición a planificación de implementación de Application** — Invocar `writing-plans-application` solo después de aprobación.

## Flujo del Proceso

```dot
digraph brainstorming_application {
    "Usar elementos existentes de Domain" [shape=box];
    "Hacer preguntas aclaratorias" [shape=box];
    "¿Falta información de Domain?" [shape=diamond];
    "Solicitar información a domain-agent" [shape=box];
    "Proponer 2-3 enfoques de Application" [shape=box];
    "Presentar diseño de Application" [shape=box];
    "¿Usuario aprueba diseño de Application?" [shape=diamond];
    "Escribir documento de diseño de Application" [shape=box];
    "Invocar writing-plans-application" [shape=doublecircle];

    "Usar elementos existentes de Domain" -> "Hacer preguntas aclaratorias";
    "Hacer preguntas aclaratorias" -> "¿Falta información de Domain?";
    "¿Falta información de Domain?" -> "Solicitar información a domain-agent" [label="sí"];
    "¿Falta información de Domain?" -> "Proponer 2-3 enfoques de Application" [label="no"];
    "Solicitar información a domain-agent" -> "Proponer 2-3 enfoques de Application";
    "Proponer 2-3 enfoques de Application" -> "Presentar diseño de Application";
    "Presentar diseño de Application" -> "¿Usuario aprueba diseño de Application?";
    "¿Usuario aprueba diseño de Application?" -> "Presentar diseño de Application" [label="no, revisar"];
    "¿Usuario aprueba diseño de Application?" -> "Escribir documento de diseño de Application" [label="sí"];
    "Escribir documento de diseño de Application" -> "Invocar writing-plans-application";
}
```

## Principios Clave

- **Nunca crear ni modificar dominio** — Solo usar lo que existe o solicitarlo al domain-agent.
- **Una pregunta a la vez** — No abrumar con múltiples preguntas.
- **Opciones múltiples cuando sea posible** — Facilita comparar enfoques de orchestration.
- **YAGNI en Application** — Solo diseñar lo que coordina la lógica de dominio y entradas/salidas.
- **Explorar alternativas de diseño** — Presentar 2-3 enfoques antes de decidir.
- **Validación incremental** — Presentar diseño por secciones, obtener aprobación antes de continuar.
- **Documentar todo** — Guardar diseño aprobado en documento formal.
- **Transición clara a implementación** — Solo invocar `writing-plans-application` tras aprobación.

## Proceso Detallado

**Entendiendo la capa de Application:**

- Revisar lo que ya existe en Domain.
- Hacer preguntas sobre flujos de negocio, orchestration de agregados, inputs/outputs y transacciones.
- Identificar los **ports**: driver ports (inputs) y driven ports (dependencias externas del dominio).
- Si falta información de dominio → invocar `domain-agent` y esperar la respuesta antes de continuar.

**Explorando enfoques:**

- Presentar 2-3 alternativas de diseño de use cases y DTOs.
- Explicar ventajas, desventajas y recomendar la mejor opción.

**Presentando el diseño de Application:**

- Explicar use cases, handler y flujo de información.
- Explicar ports y DTOs.
- Preguntar al usuario si el diseño tiene sentido antes de continuar.

**Después del brainstorming:**

- Guardar el diseño aprobado en `docs/plans/YYYY-MM-DD-<tema>-application-design.md`.
- Solo después de aprobación, invocar `writing-plans-application` para planificar implementación.
