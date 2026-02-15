---
name: finishing-a-development-branch
description: Usar cuando la implementación esté completa, todos los tests pasen y se necesite decidir cómo integrar el trabajo - guía la finalización presentando opciones estructuradas para merge, PR o descarte
---

# Finalizar una Rama de Desarrollo

## Descripción General

Guiar la finalización del trabajo presentando opciones claras y ejecutando el flujo elegido.

**Principio central:** Verificar tests → Presentar opciones → Ejecutar elección → Limpieza.

**Anunciar al inicio:** "Estoy usando la skill finishing-a-development-branch para completar este trabajo."

## El Proceso

### Paso 1: Verificar Tests

**Antes de presentar opciones, verificar que los tests pasen:**

```bash
# Ejecutar la suite de tests del proyecto
npm test / cargo test / pytest / go test ./...
```

**Si los tests fallan:**

```
Tests fallando (<N> fallos). Deben corregirse antes de completar:

[Mostrar fallos]

No se puede proceder con merge o PR hasta que los tests pasen.
```

Detenerse. No continuar al Paso 2.

**Si los tests pasan:** Continuar al Paso 2.

### Paso 2: Determinar Rama Base

```bash
# Intentar ramas base comunes
git merge-base HEAD main 2>/dev/null || git merge-base HEAD master 2>/dev/null
```

O preguntar: "Esta rama se creó desde main — ¿es correcto?"

### Paso 3: Presentar Opciones

Presentar exactamente estas 4 opciones:

```
Implementación completa. ¿Qué deseas hacer?

1. Hacer merge local a <rama-base>
2. Hacer push y crear un Pull Request
3. Mantener la rama tal como está (lo gestiono después)
4. Descartar este trabajo

¿Qué opción?
```

**No agregar explicación adicional.**

### Paso 4: Ejecutar la Elección

#### Opción 1: Merge Local

```bash
# Cambiar a la rama base
git checkout <rama-base>

# Traer últimos cambios
git pull

# Hacer merge de la rama feature
git merge <rama-feature>

# Verificar tests en el resultado del merge
<comando de tests>

# Si los tests pasan
git branch -d <rama-feature>
```

#### Opción 2: Push y Crear PR

```bash
# Subir rama
git push -u origin <rama-feature>

# Crear PR
gh pr create --title "<título>" --body "$(cat <<'EOF'
## Resumen
<2-3 puntos de lo que cambió>

## Plan de Pruebas
- [ ] <pasos de verificación>
EOF
)"
```

#### Opción 3: Mantener Tal Como Está

Reportar: "Manteniendo la rama <nombre>."

#### Opción 4: Descartar

**Confirmar primero:**

```
Esto eliminará permanentemente:
- Rama <nombre>
- Todos los commits: <lista-de-commits>

Escribe 'discard' para confirmar.
```

Esperar confirmación exacta.

Si se confirma:

```bash
git checkout <rama-base>
git branch -D <rama-feature>
```

## Referencia Rápida

| Opción         | Merge | Push | Eliminar Rama |
| -------------- | ----- | ---- | ------------- |
| 1. Merge local | ✓     | -    | ✓             |
| 2. Crear PR    | -     | ✓    | -             |
| 3. Mantener    | -     | -    | -             |
| 4. Descartar   | -     | -    | ✓ (forzado)   |

## Errores Comunes

**Omitir verificación de tests**

- **Problema:** Merge de código roto o PR fallando
- **Solución:** Siempre verificar tests antes de ofrecer opciones

**Preguntas abiertas**

- **Problema:** "¿Qué hago ahora?" → ambiguo
- **Solución:** Presentar exactamente 4 opciones estructuradas

**No confirmar descarte**

- **Problema:** Eliminación accidental de trabajo
- **Solución:** Requerir confirmación escrita "discard"

## Señales de Alerta

**Nunca:**

- Continuar con tests fallando
- Hacer merge sin verificar tests en el resultado
- Eliminar trabajo sin confirmación
- Forzar push sin solicitud explícita

**Siempre:**

- Verificar tests antes de ofrecer opciones
- Presentar exactamente 4 opciones
- Obtener confirmación escrita para Opción 4

## Integración

**Invocado por:**

- **subagent-driven-development** (Paso 7) — Después de completar todas las tareas
- **executing-plans** (Paso 5) — Después de completar todos los lotes

**Prohibido:**
No se permite el uso de worktrees ni mecanismos relacionados bajo ninguna circunstancia.
