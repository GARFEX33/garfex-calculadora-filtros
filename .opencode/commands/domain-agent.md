---
description: "Invoca al agente de Dominio para crear entidades y lógica de negocio según la arquitectura hexagonal."
disable-model-invocation: true
---

# Flujo del comando

1. Recibe la descripción del requerimiento de dominio.
2. Llama al agente `domain-agent`.
3. Asegura que la skill `enforce-domain-boundary` se ejecute antes de generar output.
4. Recibe el output y lo devuelve al Gestor o al siguiente paso.
