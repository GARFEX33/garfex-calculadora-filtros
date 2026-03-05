# PRD: Modernización de Interfaz de Memoria de Cálculo (NOM-001)

## 1. Objetivo del Cambio

Migrar el diseño actual basado en tablas rígidas y estilos antiguos (serif/negro) hacia una interfaz de ingeniería moderna, limpia y de alta legibilidad, optimizada para generación de PDF mediante `wkhtmltopdf`.

## 2. Lineamientos Visuales (Design Tokens)

Todos los archivos deben heredar los estilos definidos en el nuevo `pdf.css`:

- **Tipografía:** Sans-serif (`Segoe UI`, `Helvetica`). Nada de `Times New Roman`.
- **Colores:** \* `Primary`: `#1e40af` (Azul Navy) para títulos y acentos.
- `Text-Main`: `#1e293b` (Slate 800) para valores de datos.
- `Text-Muted`: `#64748b` (Slate 500) para etiquetas/labels.

- **Contenedores:** Reemplazar tablas con bordes por la clase `.card` (border 1pt suave, padding 12pt).

## 3. Estructura de Datos (Layout)

Queda estrictamente prohibido el uso de `<table>` para disposición de datos generales. Se debe implementar el sistema de **Data Grid**:

- **Contenedor:** `<div class="data-grid">` (simulado con `display: table`).
- **Elemento:** `<div class="data-item">` con un `<span class="data-label">` y un `<span class="data-value">`.
- **Distribución:** 2 columnas por fila para datos técnicos.

## 4. Requisitos por Sección

### A. Secciones de Cálculo (Corriente, Alimentador, Tierra)

- **Fórmulas:** Deben ir dentro de un `<div class="formula-box">`.
- **Variables:** Usar formato matemático ($I_n$, $I_a$, $f_t$).
- **Resultados:** El valor final de cada cálculo debe usar la clase `.resultado-destacado` (fuente mono, negrita, alineado a la derecha).

### B. Secciones de Validación (Dictámenes)

- **Semáforo:** Los cuadros de "CUMPLE" o "NO CUMPLE" deben usar las clases:
- `.dictamen.cumple` (Verde esmeralda suave).
- `.dictamen.no-cumple` (Rojo suave).

- **Badges:** Los estados pequeños dentro de líneas de texto deben usar la clase `.badge`.

### C. Sección de Canalización

- **Diagramas SVG:** Si se incluyen diagramas de llenado de tubería, deben estar centrados y con un borde sutil (`.card`).

## 5. Especificaciones Técnicas (wkhtmltopdf)

- **Paginación:** Mantener `page-break-inside: avoid;` en todas las `.card` para evitar que un dato se corte entre dos páginas.
- **Unidades:** Usar `pt` para fuentes y `mm` para márgenes de página. No usar `px`.
