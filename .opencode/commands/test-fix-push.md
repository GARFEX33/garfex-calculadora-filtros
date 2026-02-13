---
description: Ejecuta tests, aplica fixes, hace commit y push real a GitHub
agent: commit-agent # tu agente Git que ya tiene commit-work
model: anthropic/claude-haiku-4-5-20251001
temperature: 0.1
---

# Flujo completo optimizado

1. **Detecta archivos modificados**
   - `git diff --name-only HEAD`
   - Lista todos los archivos cambiados desde el último commit.

2. **Filtra archivos de código y documentación**
   - Archivos de código: `.py`, `.go`, `.js`, `.ts`, `.dart`, etc. → **para tests**
   - Archivos de documentación o configuración: `README.md`, `AGENTS.md`, `.gitignore`, assets → **no se testean**.

3. **Ejecuta tests solo de archivos de código modificados**
   - Python: `pytest path/al/archivo_modificado.py`
   - Go: `go test ./ruta/paquete`
   - Node.js: `jest archivoModificado.test.js`
   - Si no hay cambios en archivos de código, se salta la ejecución de tests.

3.1 **Revisa los `AGENTS.md` y valida si hay que modificar algo respecto al codigo realizado**

4. **Detecta cambios en `AGENTS.md` o archivos de agentes**
   - Si hay cambios en agentes o en `AGENTS.md`, se ejecuta la **auditoría de agentes** usando la skill `agents-md-manager`.

5. **Identifica tests que fallan y errores en auditoría**
   - Marca los tests que necesitan corrección.
   - Revisa alertas o errores detectados en la auditoría de agentes.

6. **Aplica correcciones en el código y/o agentes**
   - Corrige los tests fallidos.
   - Corrige errores de auditoría si se detectaron.

7. **Usa `commit-work` para Git**
   - **Staging:** `git add` de los archivos modificados.
   - **Commit:** Mensajes claros describiendo los fixes aplicados, indicando si fueron **tests** o **auditoría de agentes**.

8. **Revisa el estado del repositorio**
   - `git status` → archivos modificados y staged
   - `git diff` → cambios no commiteados
   - `git tree` → estructura de commits y ramas

9. **Haz push a GitHub**
   - Envía los commits al repositorio remoto (`git push`).

10. **Genera un reporte final**
    - Incluye:
      - Tests que fallaban y ahora pasan
      - Resultados de auditoría de AGENTS.md
      - Archivos modificados
      - Commits creados
      - Estado final del repo (`git tree` y `git diff` si aplica)
