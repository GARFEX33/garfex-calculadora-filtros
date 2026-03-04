# PRD: Generacion de PDF — Memoria de Calculo Electrica

**Autor:** Equipo Garfex
**Fecha:** 2026-03-02
**Estado:** Draft
**Version:** 1.3

---

## 1. Executive Summary

Garfex Calculadora de Filtros es una aplicacion web (SvelteKit + Go) que calcula memorias de calculo electrico segun normativa NOM-001-SEDE-2012 de Mexico. Actualmente, el usuario puede visualizar los resultados completos en pantalla a traves del componente `MemoriaTecnica` — que muestra 7 secciones tecnicas con formulas, tablas, diagramas SVG de arreglo de cables y conclusion de cumplimiento normativo.

Sin embargo, **no existe forma de exportar o imprimir estos resultados como un documento formal**. Los ingenieros electricos necesitan entregar memorias de calculo firmadas y selladas como parte de proyectos de instalaciones electricas. Hoy deben copiar/pegar manualmente o hacer capturas de pantalla, lo cual es ineficiente, propenso a errores y poco profesional.

Esta funcionalidad agrega un flujo de **dos pasos** para generar el PDF:

1. **Pagina intermedia de configuracion** (`/calculos/resultado/pdf`) — donde el usuario selecciona la presentacion de empresa (Summa, Garfex, Siemens), ingresa datos del proyecto (nombre, direccion), el responsable del reporte (para firma), y puede editar el nombre del equipo/carga.
2. **Generacion del PDF** — el backend Go recibe los datos de calculo + los datos de presentacion, renderiza un template HTML con branding de la empresa seleccionada y produce un PDF profesional usando `wkhtmltopdf`, listo para impresion y firma.

El resultado es un documento con identidad visual de la empresa, datos del proyecto completos y espacio de firma personalizado — sin edicion manual posterior.

---

## 2. Problem Statement

### Problema actual

Los ingenieros electricos que usan Garfex para calcular memorias tecnicas **no pueden generar un documento formal** a partir de los resultados. Esto genera:

1. **Perdida de tiempo**: Copiar datos manualmente a Word/Excel para crear el documento entregable.
2. **Riesgo de error**: Transcripcion manual de valores numericos (calibres, corrientes, porcentajes de caida) introduce errores humanos.
3. **Falta de profesionalismo**: Capturas de pantalla o documentos improvisados no cumplen con la presentacion esperada en proyectos de ingenieria.
4. **Barrera de adopcion**: Usuarios potenciales eligen herramientas que SI generan documentos formales, aunque calculen peor.

### Evidencia

- Los usuarios actuales reportan que el mayor pain point despues de calcular es "sacar el documento".
- En el mercado mexicano de ingenieria electrica, la memoria de calculo impresa y firmada es un **entregable obligatorio** en proyectos de instalaciones.
- Competidores directos (hojas de calculo custom, software de pago) ofrecen exportacion a PDF como feature basica.

---

## 3. Goals & Objectives

### Objetivo principal

Permitir a los usuarios configurar la presentacion del documento (empresa, proyecto, responsable) y generar un PDF profesional con branding de empresa, listo para impresion y firma.

### Objetivos especificos

| # | Objetivo | Medible como |
|---|----------|-------------|
| O1 | Generar PDF que replique fielmente las 7 secciones de MemoriaTecnica | PDF contiene todas las secciones con datos correctos |
| O2 | Branding de empresa seleccionable | PDF muestra logo, nombre y datos de la empresa elegida (Summa, Garfex, Siemens) |
| O3 | Datos de proyecto editables antes de generar | Usuario puede ingresar nombre de proyecto, direccion, responsable |
| O4 | Nombre de equipo/carga editable | Usuario puede modificar como aparece el nombre del equipo en el PDF sin cambiar el calculo |
| O5 | Documento profesional de ingenieria | Header con logo de empresa, pie con datos empresa, numeracion, formato tabular |
| O6 | Generacion rapida | PDF generado en <5 segundos para calculo estandar |
| O7 | Sin dependencia de entorno grafico | Funciona en servidor Linux sin X11 (Docker headless) |

### Non-goals (lo que NO buscamos resolver ahora)

- Edicion del PDF despues de generado.
- Firma digital electronica (PKI/e.firma SAT).
- Generacion batch de multiples memorias.
- Almacenamiento persistente de PDFs generados (no se guardan en DB).
- Crear/editar empresas desde la UI (las presentaciones son estaticas).
- Guardar preferencias de presentacion del usuario entre sesiones.

---

## 4. User Personas

### Persona 1: Ingeniero Electrico de Campo — "Carlos"

- **Rol**: Ingeniero electrico proyectista en empresa de instalaciones.
- **Contexto**: Necesita entregar memorias de calculo firmadas como parte de proyectos de instalaciones en industrias, plantas, edificios comerciales.
- **Frustracion**: Pierde 30-60 minutos por memoria transcribiendo datos de pantalla a un formato Word.
- **Necesidad**: Un PDF profesional que pueda imprimir, firmar y adjuntar al expediente del proyecto.
- **Frecuencia de uso**: 5-15 memorias por semana.

### Persona 2: Supervisor / Director de Ingenieria — "Ing. Martinez"

- **Rol**: Responsable tecnico que revisa y firma las memorias de calculo.
- **Contexto**: Recibe memorias de su equipo, las revisa y las aprueba con su sello.
- **Frustracion**: Documentos inconsistentes en formato y contenido entre ingenieros del equipo.
- **Necesidad**: Formato estandarizado y profesional que refleje la seriedad de la empresa.
- **Frecuencia de uso**: Revisa 20-50 memorias por semana.

### Persona 3: Ingeniero Freelance — "Ana"

- **Rol**: Consultora independiente que hace proyectos electricos para multiples clientes.
- **Contexto**: Necesita entregar documentacion profesional rapido para facturar y avanzar al siguiente proyecto.
- **Frustracion**: No tiene plantillas estandarizadas; cada documento se ve diferente.
- **Necesidad**: PDF con aspecto profesional que no requiera retoques manuales.
- **Frecuencia de uso**: 3-8 memorias por semana.

---

## 5. User Stories & Requirements

### Epic: Generacion de PDF de Memoria de Calculo

---

### Sub-Epic A: Formulario de Configuracion Pre-PDF

#### US-01: Navegar al formulario de configuracion de PDF

```
As an ingeniero electrico,
I want to press a "Generar PDF" button on the results page,
So that I'm taken to a configuration page where I can set up the document presentation.

Acceptance Criteria:
- [ ] A "Generar PDF" button is visible on /calculos/resultado in the page header
- [ ] Clicking the button navigates to /calculos/resultado/pdf
- [ ] The MemoriaOutput data is preserved during navigation (via state or query params)
- [ ] The configuration page shows a form with all editable fields
- [ ] A "Volver a Resultados" link/button allows going back without losing data
```

#### US-02: Seleccionar presentacion de empresa

```
As an ingeniero electrico,
I want to select which company branding to use for the PDF,
So that the document has the correct logo, header and footer for my client/employer.

Acceptance Criteria:
- [ ] A dropdown/radio group shows available company presentations
- [ ] Available options (static, predefined): Summa, Garfex, Siemens
- [ ] Each option shows a visual preview or label indicating the company
- [ ] Selecting a company preloads: logo, company name, address, phone, email
- [ ] The company data is displayed as read-only preview below the selector
- [ ] Default selection: first option in the list (Garfex)
- [ ] The selected company determines:
      - Logo in PDF header
      - Company name and contact info in PDF footer
      - Optional: accent colors in the PDF styling
```

**Company presentation data structure:**

| Campo | Tipo | Ejemplo (Garfex) |
|-------|------|------------------|
| id | string | "garfex" |
| nombre | string | "Garfex Ingenieria Electrica" |
| direccion | string | "Av. Insurgentes Sur 1234, CDMX" |
| telefono | string | "+52 55 1234 5678" |
| email | string | "contacto@garfex.com" |
| logo_path | string | "logos/garfex.png" |

#### US-03: Ingresar datos del proyecto

```
As an ingeniero electrico,
I want to enter the project name and address,
So that the PDF header identifies the specific project this calculation belongs to.

Acceptance Criteria:
- [ ] Text input for "Nombre del Proyecto" — required, max 200 chars
- [ ] Text input for "Direccion del Proyecto" — optional, max 300 chars
- [ ] Both fields are free text (no validation beyond length)
- [ ] Fields are empty by default (no pre-filled values)
- [ ] The project name appears prominently in the PDF header section
- [ ] The project address appears below the project name in the PDF header
```

#### US-04: Ingresar responsable del reporte

```
As an ingeniero electrico,
I want to enter the name of the person responsible for the report,
So that the signature block at the bottom of the PDF has the correct name.

Acceptance Criteria:
- [ ] Text input for "Responsable del Reporte" — required, max 150 chars
- [ ] The name appears in the signature block at the end of the PDF
- [ ] The signature block shows: line for firma, nombre del responsable, space for cedula profesional, fecha
```

#### US-05: Editar nombre del equipo o carga

```
As an ingeniero electrico,
I want to modify how the equipment/load name appears in the PDF,
So that I can use a descriptive project-specific name instead of the catalog code.

Acceptance Criteria:
- [ ] Text input for "Nombre del Equipo / Carga" pre-filled with the current equipo.clave or tipo_equipo
- [ ] User can edit the value freely (max 200 chars)
- [ ] The edited name replaces the equipment identifier in the PDF header/encabezado section
- [ ] The original calculation data (calibres, corrientes, etc.) is NOT affected by this edit
- [ ] If the user clears the field, the original value is used as fallback
```

#### US-06: Validar y enviar formulario de configuracion

```
As a user,
I want the form to validate my inputs before generating the PDF,
So that I don't get errors after waiting for the PDF to generate.

Acceptance Criteria:
- [ ] "Generar PDF" submit button at the bottom of the form
- [ ] Client-side validation: nombre_proyecto and responsable are required
- [ ] Validation errors shown inline below each field (red text, border)
- [ ] Submit button is disabled if required fields are empty
- [ ] On submit: loading state on the button ("Generando PDF...")
- [ ] On success: browser downloads PDF, button returns to normal
- [ ] On error: toast/alert shows error, button re-enables
- [ ] Timeout after 30 seconds: shows "La generacion tardo demasiado"
```

---

### Sub-Epic B: Contenido del PDF

#### US-07: PDF muestra branding de empresa seleccionada

```
As an ingeniero electrico,
I want the PDF to display the selected company's logo and information,
So that the document looks like an official company deliverable.

Acceptance Criteria:
- [ ] PDF header on every page shows: company logo (left), document title (center), date (right)
- [ ] PDF footer on every page shows: company name, address, phone, email
- [ ] Footer also shows: "Pagina X de Y"
- [ ] Logo renders correctly (PNG format, max 200x80px)
- [ ] If logo is not available, company name is shown in bold text instead
```

#### US-08: PDF muestra datos del proyecto en encabezado

```
As an ingeniero electrico,
I want the PDF to show the project data I entered in the first page section,
So that the document is clearly identified for the specific project.

Acceptance Criteria:
- [ ] First section of PDF shows: "Proyecto: <nombre_proyecto>"
- [ ] Below that: "Direccion: <direccion_proyecto>" (if provided)
- [ ] Below that: equipment/load data (clave, tipo_equipo, voltaje, amperaje)
- [ ] The nombre_equipo_override is used instead of equipo.clave if provided
- [ ] Shows instalacion data: estado, sistema electrico, tipo canalizacion, material
- [ ] Shows date of generation
```

#### US-09: PDF contiene seccion de corriente nominal

```
As an ingeniero electrico,
I want the PDF to show the current calculation section with formulas,
So that the reviewer can verify the nominal current derivation.

Acceptance Criteria:
- [ ] Shows corriente nominal calculation with formula used
- [ ] Shows factores de ajuste: temperatura, agrupamiento, total
- [ ] Shows corriente ajustada result
- [ ] Shows temperatura ambiente and temperatura referencia
- [ ] Shows tabla de ampacidad usada
- [ ] Layout matches SeccionCorriente.svelte content
```

#### US-10: PDF contiene seccion de alimentador

```
As an ingeniero electrico,
I want the PDF to show the conductor selection section,
So that the reviewer can verify the cable sizing complies with NOM.

Acceptance Criteria:
- [ ] Shows cable_fase: calibre, material, seccion_mm2, tipo_aislamiento, capacidad
- [ ] Shows num_hilos (conductores por fase en paralelo)
- [ ] If seleccion_por_caida_tension is true, shows nota_seleccion and calibre_original_ampacidad
- [ ] Layout matches SeccionAlimentador.svelte content
```

#### US-11: PDF contiene seccion de conductor de tierra

```
As an ingeniero electrico,
I want the PDF to show the ground conductor selection,
So that the reviewer can verify compliance with NOM grounding tables.

Acceptance Criteria:
- [ ] Shows cable_tierra: calibre, material, seccion_mm2
- [ ] Shows proteccion ITM value
- [ ] Layout matches SeccionTierra.svelte content
```

#### US-12: PDF contiene seccion de canalizacion

```
As an ingeniero electrico,
I want the PDF to show the raceway/conduit sizing section with formulas,
So that the reviewer can verify the conduit or cable tray selection.

Acceptance Criteria:
- [ ] Shows resultado: tamano, area_total_mm2, area_requerida_mm2, numero_de_tubos
- [ ] Shows fill_factor percentage
- [ ] For tuberia: shows detalle_tuberia (areas de fase, neutro, tierra, distribucion por tubo)
- [ ] For charola: shows detalle_charola (diametros, espaciado, anchos)
- [ ] Shows intermediate calculation formulas
- [ ] Layout matches SeccionCanalizacion.svelte content
```

#### US-13: PDF contiene seccion de caida de tension

```
As an ingeniero electrico,
I want the PDF to show the voltage drop calculation,
So that the reviewer can verify it's within NOM limits.

Acceptance Criteria:
- [ ] Shows porcentaje de caida and caida_volts
- [ ] Shows limite_porcentaje
- [ ] Shows impedancia, resistencia, reactancia
- [ ] Shows cumple/no cumple indicator
- [ ] Layout matches SeccionCaidaTension.svelte content
```

#### US-14: PDF contiene seccion de conclusion con firma personalizada

```
As an ingeniero electrico,
I want the PDF to show a conclusion with the responsible person's name in the signature block,
So that the document is ready for the reviewer to sign without manual editing.

Acceptance Criteria:
- [ ] Shows cumple_normativa status prominently (CUMPLE / NO CUMPLE NOM-001-SEDE-2012)
- [ ] Lists all observaciones
- [ ] Shows summary table: cable fase, cable tierra, canalizacion, ITM, caida de tension
- [ ] Signature block shows:
      - Line for handwritten signature
      - "Nombre: <responsable>" (pre-filled from form)
      - "Cedula Profesional: _______________" (blank line)
      - "Fecha: _______________" (blank line)
- [ ] Layout matches SeccionConclusion.svelte content + firma additions
```

#### US-15: PDF tiene formato profesional de impresion

```
As a supervisor de ingenieria,
I want the PDF to have a professional layout with page numbers and company branding,
So that it looks like a formal engineering document from my company.

Acceptance Criteria:
- [ ] Page size: Letter (8.5" x 11")
- [ ] Margins: adequate for binding (left margin slightly wider)
- [ ] Header on every page: company logo, "Memoria de Calculo Electrica", fecha
- [ ] Footer on every page: company info (nombre, direccion, telefono, email) + "Pagina X de Y"
- [ ] Tables have visible borders and proper alignment
- [ ] Font: legible sans-serif, minimum 10pt body text
- [ ] Section headings are clearly distinguishable
```

---

### Sub-Epic C: Feedback y UX

#### US-16: Feedback visual durante generacion

```
As a user,
I want to see clear feedback while the PDF is being generated,
So that I know the system is working and I don't click again.

Acceptance Criteria:
- [ ] Submit button shows spinner/loading animation during generation
- [ ] Button text changes to "Generando PDF..."
- [ ] Button is disabled (not clickable) during generation
- [ ] On success: button returns to normal state, PDF download starts
- [ ] On error: toast/alert shows error message, button re-enables
- [ ] Timeout after 30 seconds shows error "La generacion tardo demasiado"
```

---

## 6. Success Metrics

### Framework: HEART (Google)

| Dimension | Metric | Target | How to Measure |
|-----------|--------|--------|----------------|
| **Happiness** | User satisfaction with PDF quality | >4/5 rating | In-app feedback after first PDF generation |
| **Engagement** | % of calculations that generate PDF | >60% of completed calculations | Backend logs: PDF requests / total memoria requests |
| **Adoption** | Users who generate at least 1 PDF per week | >80% of active users | Backend logs: unique users with PDF requests |
| **Retention** | Users who return to generate PDFs week-over-week | >70% weekly retention | Backend logs: returning users |
| **Task Success** | PDF generation success rate | >99% | Backend logs: successful / total PDF requests |

### Metricas tecnicas

| Metric | Target | Measurement |
|--------|--------|-------------|
| Tiempo de generacion P50 | <3 segundos | Backend latency logs |
| Tiempo de generacion P95 | <5 segundos | Backend latency logs |
| Tamano del PDF | <2 MB por memoria estandar | File size check |
| Tasa de error | <1% | Error logs / total requests |

---

## 7. Scope

### In Scope (MVP)

| # | Feature | Priority |
|---|---------|----------|
| 1 | Pagina intermedia de configuracion (`/calculos/resultado/pdf`) con formulario | Must have |
| 2 | Selector de presentacion de empresa (Summa, Garfex, Siemens) con datos estaticos | Must have |
| 3 | Campos editables: nombre proyecto, direccion, responsable del reporte | Must have |
| 4 | Campo editable: nombre del equipo/carga (override visual, no afecta calculo) | Must have |
| 5 | Endpoint backend `POST /api/v1/pdf/memoria` que recibe datos de calculo + presentacion | Must have |
| 6 | Template HTML con las 7 secciones de MemoriaTecnica | Must have |
| 7 | Header PDF con logo de empresa, titulo, fecha | Must have |
| 8 | Footer PDF con datos de empresa (nombre, direccion, telefono, email) + paginacion | Must have |
| 9 | Descarga directa del PDF al navegador | Must have |
| 10 | Feedback visual (loading, error, validacion) | Must have |
| 11 | Tablas con datos de corrientes, conductores, canalizacion, caida de tension | Must have |
| 12 | Indicador visual de cumplimiento normativo (CUMPLE / NO CUMPLE) | Must have |
| 13 | Bloque de firma con nombre del responsable pre-llenado | Must have |
| 14 | Nombre de archivo descriptivo (proyecto + equipo + fecha) | Should have |
| 15 | Assets estaticos de logos para las 3 empresas predefinidas (PNG) | Must have |

### Out of Scope (Future)

| # | Feature | Reason |
|---|---------|--------|
| 1 | Diagramas SVG de arreglo de cables en el PDF | Complejidad de renderizado SVG→PDF con wkhtmltopdf; considerar en v2 |
| 2 | Firma digital electronica (e.firma SAT, PKI) | Requiere integracion con SAT; feature independiente |
| 3 | CRUD de empresas desde la UI (crear/editar/eliminar presentaciones) | Las 3 empresas son estaticas por ahora; si se necesitan mas, se agregan en codigo |
| 4 | Generacion batch de multiples PDFs | Requiere cola de jobs; future feature |
| 5 | Almacenamiento persistente de PDFs generados | Requiere storage (S3/minio); future feature |
| 6 | Preview del PDF antes de descarga | Agrega complejidad sin valor claro para MVP |
| 7 | Exportacion a Word/Excel | Formato PDF es suficiente para MVP |
| 8 | Graficos/charts embebidos en el PDF | Los datos tabulares son suficientes para la memoria tecnica |
| 9 | Watermark de "borrador" vs "final" | Future nice-to-have |
| 10 | Internacionalizacion (ingles) | El mercado target es Mexico; NOM es regulacion mexicana |
| 11 | Guardar preferencias de empresa/responsable entre sesiones | Requiere persistencia de usuario; future feature |
| 12 | Colores/tema custom por empresa | Solo logo y datos textuales en MVP; colores en v2 |

---

## 8. Technical Considerations

### 8.1 Arquitectura general

```
Frontend (SvelteKit)                          Backend (Go)
┌──────────────────────┐                ┌──────────────────────────────┐
│ /calculos/resultado   │                │                              │
│                       │                │                              │
│ [Generar PDF] ────────┼── navigate ──▶ │                              │
│                       │                │                              │
├───────────────────────┤                │                              │
│ /calculos/resultado/  │                │ POST /api/v1/pdf/memoria     │
│          pdf          │                │                              │
│                       │                │                              │
│ ┌─ Empresa: [v]─────┐│                │                              │
│ │ Summa / Garfex /   ││                │                              │
│ │ Siemens            ││                │                              │
│ ├────────────────────┤│   POST body:   │                              │
│ │ Nombre Proyecto    ││  {memoria,     │──▶ PdfHandler                │
│ │ Direccion          ││   presentacion}│     │                        │
│ │ Responsable        ││                │     ▼                        │
│ │ Nombre Equipo      ││                │  html/template (render)      │
│ ├────────────────────┤│                │  + logo empresa              │
│ │ [Generar PDF]      ││── POST ───────▶│     │                        │
│ │                    ││                │     ▼                        │
│ │ ◄── blob download  ││                │  wkhtmltopdf (HTML → PDF)    │
│ └────────────────────┘│                │     │                        │
└───────────────────────┘                │     ▼                        │
                                         │  Response: application/pdf   │
                                         └──────────────────────────────┘
```

### 8.2 Enfoque: Frontend-to-Backend con datos completos + presentacion

**Decision clave**: El frontend envia el objeto `MemoriaOutput` completo + datos de presentacion (empresa, proyecto, responsable) en el body del POST.

**Justificacion**:
- Actualmente no existe persistencia de calculos en base de datos — los resultados viajan en query params entre paginas.
- No se necesita acceso a base de datos para generar el PDF: todos los datos ya estan calculados.
- Los datos de presentacion (empresa, proyecto) son efimeros — no se guardan, solo se usan para el PDF.
- Las definiciones de empresas son estaticas en el backend (logos embebidos o en filesystem).
- Simplifica el backend: recibe datos + clave de empresa → resuelve logo/datos → renderiza → devuelve PDF.

**Trade-off**: El body del POST es mas grande (~6-12 KB de JSON), pero esto es insignificante para una operacion one-shot.

### 8.2.1 Definicion estatica de empresas

Las presentaciones de empresa se definen como configuracion estatica en el backend:

```go
// internal/pdf/domain/empresa.go
type EmpresaPresentacion struct {
    ID        string // "summa", "garfex", "siemens"
    Nombre    string
    Direccion string
    Telefono  string
    Email     string
    LogoPath  string // ruta al archivo PNG en assets/
}

// Catalogo estatico — no requiere base de datos
var EmpresasCatalogo = map[string]EmpresaPresentacion{
    "garfex": {
        ID:        "garfex",
        Nombre:    "Garfex Ingenieria Electrica",
        Direccion: "...",
        Telefono:  "...",
        Email:     "contacto@garfex.com",
        LogoPath:  "assets/logos/garfex.png",
    },
    "summa": { ... },
    "siemens": { ... },
}
```

El frontend necesita la lista de empresas disponibles. Dos opciones:

- **Opcion A (simple)**: Hardcodear la lista en el frontend tambien. Duplica datos pero es trivial para 3 items estaticos.
- **Opcion B (limpio)**: Endpoint `GET /api/v1/pdf/empresas` que retorna la lista (sin logos, solo id+nombre). Agrega un endpoint pero evita duplicacion.

**Decision recomendada**: Opcion A para MVP. Son 3 items estaticos que rara vez cambian.

### 8.3 Backend — Nuevo modulo `pdf`

Siguiendo la arquitectura hexagonal existente, la generacion de PDF se implementa como un **nuevo modulo vertical**:

```
internal/pdf/
├── domain/
│   └── empresa.go                    # EmpresaPresentacion struct + catalogo estatico
├── application/
│   ├── dto/
│   │   └── pdf_request.go            # PdfMemoriaRequest: MemoriaOutput + PresentacionInput
│   └── usecase/
│       └── generar_memoria_pdf.go    # Orquesta: resuelve empresa → renderiza HTML → exec wkhtmltopdf → retorna bytes
├── infrastructure/
│   └── adapter/
│       ├── driver/http/
│       │   └── pdf_handler.go        # HTTP handler: parsea request, llama use case, escribe response
│       └── driven/
│           ├── template/
│           │   └── html_renderer.go  # Renderiza html/template con datos + empresa
│           └── wkhtmltopdf/
│               └── pdf_generator.go  # Ejecuta wkhtmltopdf como subproceso
├── assets/
│   └── logos/
│       ├── garfex.png                # Logo Garfex (~200x80px)
│       ├── summa.png                 # Logo Summa
│       └── siemens.png               # Logo Siemens
└── templates/
    ├── memoria.html                  # Template principal
    ├── partials/
    │   ├── header.html               # Logo empresa + titulo + fecha
    │   ├── footer.html               # Datos empresa + paginacion
    │   ├── seccion_encabezado.html   # Proyecto + equipo + instalacion
    │   ├── seccion_corriente.html
    │   ├── seccion_alimentador.html
    │   ├── seccion_tierra.html
    │   ├── seccion_canalizacion.html
    │   ├── seccion_caida_tension.html
    │   └── seccion_conclusion.html   # Cumplimiento + firma con nombre responsable
    └── styles/
        └── pdf.css                   # Estilos de impresion
```

### 8.4 Endpoint

```
POST /api/v1/pdf/memoria
Content-Type: application/json

Body:
{
  "memoria": { ... },          // MemoriaOutput completo
  "presentacion": {
    "empresa_id": "garfex",    // ID de la empresa del catalogo estatico
    "nombre_proyecto": "Planta Industrial Norte",
    "direccion_proyecto": "Av. Industrial 456, Monterrey, NL",
    "responsable": "Ing. Carlos Rodriguez",
    "nombre_equipo_override": "Filtro Activo Linea 3"  // opcional
  }
}

Response:
  200 OK
  Content-Type: application/pdf
  Content-Disposition: attachment; filename="MemoriaCalculo_<proyecto>_<equipo>_<fecha>.pdf"
  Body: <binary PDF>

  400 Bad Request — body invalido, empresa_id no existe, campos requeridos faltantes
  500 Internal Server Error — fallo en generacion
```

**DTO del request:**

```go
// internal/pdf/application/dto/pdf_request.go
type PresentacionInput struct {
    EmpresaID            string `json:"empresa_id" binding:"required"`
    NombreProyecto       string `json:"nombre_proyecto" binding:"required,max=200"`
    DireccionProyecto    string `json:"direccion_proyecto" binding:"max=300"`
    Responsable          string `json:"responsable" binding:"required,max=150"`
    NombreEquipoOverride string `json:"nombre_equipo_override" binding:"max=200"`
}

type PdfMemoriaRequest struct {
    Memoria      calculos_dto.MemoriaOutput `json:"memoria" binding:"required"`
    Presentacion PresentacionInput          `json:"presentacion" binding:"required"`
}
```

### 8.5 Template HTML + CSS puro de impresion

> **Decision tecnica**: Los templates usan HTML semantico puro + CSS custom optimizado para impresion (`pdf.css`).
> NO se usa Tailwind, Bootstrap ni ningun framework CSS. Ver seccion 9.4 para justificacion y estructura del CSS.

El template HTML usa `html/template` de Go. Recibe un struct combinado con datos de calculo + presentacion:

```go
// Struct que alimenta el template
type TemplateData struct {
    // Empresa
    Empresa     EmpresaPresentacion  // logo, nombre, direccion, tel, email
    LogoBase64  string               // logo embebido como data:image/png;base64,...

    // Proyecto (del formulario)
    NombreProyecto    string
    DireccionProyecto string
    Responsable       string
    NombreEquipo      string  // override o valor original

    // Calculo (MemoriaOutput)
    Memoria     MemoriaOutput

    // Metadata
    FechaGeneracion string  // formateado: "02 de Marzo de 2026"
}
```

**Mapeo de secciones:**

| Seccion HTML | Fuente de datos | Equivalente Svelte |
|---|---|---|
| Header (cada pagina) | Empresa.Logo, "Memoria de Calculo Electrica", FechaGeneracion | N/A (nuevo) |
| Footer (cada pagina) | Empresa.Nombre, Empresa.Direccion, Empresa.Telefono, Empresa.Email, "Pagina X de Y" | N/A (nuevo) |
| Encabezado | NombreProyecto, DireccionProyecto, NombreEquipo, Memoria.tipo_equipo, Memoria.instalacion | SeccionEncabezado.svelte |
| Corriente | Memoria.corrientes (nominal, ajustada, factores) | SeccionCorriente.svelte |
| Alimentador | Memoria.cable_fase, Memoria.corrientes.corriente_por_hilo | SeccionAlimentador.svelte |
| Tierra | Memoria.cable_tierra, Memoria.proteccion.itm | SeccionTierra.svelte |
| Canalizacion | Memoria.canalizacion (resultado, detalle_charola/tuberia) | SeccionCanalizacion.svelte |
| Caida de Tension | Memoria.caida_tension (porcentaje, cumple, impedancia) | SeccionCaidaTension.svelte |
| Conclusion + Firma | Memoria.cumple_normativa, Memoria.observaciones, Responsable | SeccionConclusion.svelte + firma |

### 8.6 wkhtmltopdf

**Por que wkhtmltopdf** (y no alternativas):
- Renderiza HTML+CSS completo (no es un parser limitado como go-pdf).
- Soporte nativo de headers/footers por pagina con HTML.
- Soporte de paginacion automatica con "Pagina X de Y".
- Probado en produccion en miles de proyectos.
- Funciona headless sin entorno grafico (con `xvfb` o builds estaticos).

**Instalacion en Docker**:

```dockerfile
# En la imagen del backend
RUN apt-get update && apt-get install -y --no-install-recommends \
    wkhtmltopdf \
    xvfb \
    && rm -rf /var/lib/apt/lists/*
```

**Alternativa evaluada**: Build estatico de wkhtmltopdf sin dependencia de xvfb (wkhtmltopdf 0.12.6+ con patched Qt).

**Ejecucion desde Go**:

```go
cmd := exec.CommandContext(ctx, "wkhtmltopdf",
    "--page-size", "Letter",
    "--margin-top", "20mm",
    "--margin-bottom", "20mm",
    "--margin-left", "25mm",
    "--margin-right", "15mm",
    "--header-html", headerPath,
    "--footer-html", footerPath,
    "--enable-local-file-access",
    htmlPath,
    outputPath,
)
```

### 8.7 Frontend — Pagina intermedia + Boton de generacion

**Archivos nuevos:**

```
frontend/web/src/routes/calculos/resultado/pdf/
├── +page.svelte        # Pagina del formulario de configuracion
└── +page.ts            # Load function: recupera MemoriaOutput del state
```

**Archivo modificado:**

`frontend/web/src/routes/calculos/resultado/+page.svelte` — agregar boton "Generar PDF" que navega a `/calculos/resultado/pdf`.

**Componentes nuevos opcionales:**

```
frontend/web/src/lib/components/pdf/
├── FormularioPdf.svelte          # Formulario completo de configuracion
└── SelectorEmpresa.svelte        # Dropdown de empresa con preview de datos
```

**Flujo de navegacion:**

```
/calculos/resultado
    │
    │  click "Generar PDF"
    │  (goto con state: memoriaOutput)
    ▼
/calculos/resultado/pdf
    │
    │  formulario: empresa, proyecto, responsable, nombre equipo
    │  click "Generar PDF"
    │  POST /api/v1/pdf/memoria con {memoria, presentacion}
    │
    ▼
  blob download → PDF descargado
```

**Datos de empresa en frontend (hardcodeados para MVP):**

```typescript
// frontend/web/src/lib/config/empresas-pdf.ts
export const EMPRESAS_PDF = [
  {
    id: 'garfex',
    nombre: 'Garfex Ingenieria Electrica',
    direccion: 'Av. Insurgentes Sur 1234, CDMX',
    telefono: '+52 55 1234 5678',
    email: 'contacto@garfex.com',
  },
  {
    id: 'summa',
    nombre: 'Summa Energetica',
    direccion: '...',
    telefono: '...',
    email: '...',
  },
  {
    id: 'siemens',
    nombre: 'Siemens Energy',
    direccion: '...',
    telefono: '...',
    email: '...',
  },
] as const;
```

**Pseudocodigo del formulario:**

```typescript
// FormularioPdf.svelte
let empresaSeleccionada = $state('garfex');
let nombreProyecto = $state('');
let direccionProyecto = $state('');
let responsable = $state('');
let nombreEquipo = $state(memoriaOutput.equipo.clave || memoriaOutput.tipo_equipo);
let loading = $state(false);
let error = $state('');

async function generarPdf() {
  // Validate required fields
  if (!nombreProyecto || !responsable) { ... }

  loading = true;
  error = '';
  try {
    const response = await fetch(`${PUBLIC_API_URL}/api/v1/pdf/memoria`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        memoria: memoriaOutput,
        presentacion: {
          empresa_id: empresaSeleccionada,
          nombre_proyecto: nombreProyecto,
          direccion_proyecto: direccionProyecto,
          responsable: responsable,
          nombre_equipo_override: nombreEquipo,
        }
      })
    });
    if (!response.ok) throw new Error('Error generando PDF');

    const blob = await response.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `MemoriaCalculo_${nombreProyecto}_${nombreEquipo}_${fecha}.pdf`;
    a.click();
    URL.revokeObjectURL(url);
  } catch (e) {
    error = e instanceof Error ? e.message : 'Error desconocido';
  } finally {
    loading = false;
  }
}
```
```

### 8.8 Stack tecnico

| Componente | Tecnologia | Version |
|---|---|---|
| Backend API | Go + Gin | 1.22+ |
| PDF Engine | wkhtmltopdf | 0.12.6+ |
| HTML Templates | html/template (stdlib) | Go stdlib |
| CSS | Custom print stylesheet | N/A |
| Frontend | SvelteKit + TypeScript | 5.x + 2.x |
| Testing backend | testify | latest |
| Linting | golangci-lint | latest |

---

## 9. Design & UX Requirements

### 9.0 Design Intent — Fundamentos de diseño

Antes de cualquier decision visual, estas son las preguntas respondidas:

#### ¿Quien es esta persona?

Un ingeniero electrico mexicano sentado en su oficina o en obra. Acaba de terminar un calculo que le tomo concentracion — corrientes, factores de ajuste, tablas NOM. El resultado esta en pantalla y cumple la norma. Ahora necesita **sacar el documento** para meterlo en el expediente del proyecto, imprimirlo, firmarlo con sello y entregarlo al cliente o a la DRO (Direccion Responsable de Obra). Tiene 5 memorias mas por hacer hoy.

#### ¿Que debe lograr?

**Configurar y descargar un PDF profesional en menos de 60 segundos.** No quiere pensar. Quiere seleccionar empresa, escribir el nombre del proyecto, poner su nombre, y bajar el PDF. Si el formulario le pide pensar demasiado, le estorba.

#### ¿Como debe sentirse?

**Como llenar la caratula de un expediente tecnico.** No es divertido, no es creativo — es eficiente, claro y confiable. El formulario debe sentirse como un paso natural del flujo, no como una barrera entre el calculo y el documento. Densidad media: ni tan espacioso que parezca que hay poco (pierde tiempo scrolleando), ni tan denso que abrume (ya viene de una pantalla densa de resultados).

**Palabras clave**: Eficiente. Claro. Confiable. Sin sorpresas. Paso rapido hacia el resultado final.

**Analogia**: La portada de un folder de proyecto — llenas los campos, cierras el folder, listo.

---

#### Design intent para el formulario (frontend — Svelte + Tailwind)

```
Intent:    Ingeniero post-calculo, quiere sacar documento rapido. Paso intermedio, no destino.
Palette:   Design tokens existentes de la app (bg-card, border, primary). Sin colores nuevos.
Depth:     Borders-only, consistente con el resto de la app (cards con border-border).
Surfaces:  bg-background para pagina, bg-card para secciones del formulario. Una sola elevacion.
Typography: Misma de la app (font-sans del @theme). Labels en muted-foreground, valores en foreground.
Spacing:   Base-4 de Tailwind, consistente con el resto. gap-6 entre secciones, gap-3 dentro.
```

**Principio rector**: Este formulario NO es un feature nuevo — es una **extension del flujo de resultados**. Debe sentirse como la misma app, no como una pantalla diferente. Mismos tokens, misma densidad, misma temperatura visual.

---

#### Design intent para el PDF (backend — HTML puro + CSS impresion)

```
Intent:    Supervisor/DRO que revisa memorias impresas. Documento formal de ingenieria.
Palette:   Escala de grises + verde/rojo normativo. Sin colores decorativos.
Depth:     Lineas y bordes de tabla. Sin sombras (no existen en papel).
Surfaces:  Papel blanco. Headers de tabla en gris claro. Sin fondos de color en secciones.
Typography: Helvetica/Arial — la fuente de documentos tecnicos por excelencia.
Spacing:   Base en puntos tipograficos (pt). Generoso entre secciones, compacto dentro de tablas.
```

**Principio rector**: Este PDF NO debe verse "diseñado" — debe verse **tecnico, normativo y confiable**. Un ingeniero que lo recibe debe pensar "memoria de calculo profesional", no "app moderna". La familiaridad con el formato de documentos tecnicos es una virtud, no un defecto.

---

### 9.1 Pagina de configuracion (`/calculos/resultado/pdf`)

**Layout del formulario:**

```
┌────────────────────────────────────────────────────────┐
│  ← Volver a Resultados                                 │
│                                                        │
│  Configurar Documento PDF                              │
│  ─────────────────────────                             │
│                                                        │
│  Presentacion de Empresa                               │
│  ┌──────────────────────────────────────────────────┐  │
│  │  ○ Garfex Ingenieria Electrica                   │  │
│  │  ○ Summa Energetica                              │  │
│  │  ○ Siemens Energy                                │  │
│  ├──────────────────────────────────────────────────┤  │
│  │  Garfex Ingenieria Electrica                     │  │
│  │  Av. Insurgentes Sur 1234, CDMX                  │  │
│  │  +52 55 1234 5678 | contacto@garfex.com          │  │
│  └──────────────────────────────────────────────────┘  │
│                                                        │
│  Datos del Proyecto                                    │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Nombre del Proyecto *   [________________________] │
│  │  Direccion del Proyecto  [________________________] │
│  └──────────────────────────────────────────────────┘  │
│                                                        │
│  Datos del Reporte                                     │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Responsable del Reporte * [______________________] │
│  │  Nombre del Equipo/Carga   [_FILTRO_ACTIVO_L3____] │
│  └──────────────────────────────────────────────────┘  │
│                                                        │
│  ┌──────────────────────────────────────────────────┐  │
│  │            [ Generar PDF ]                       │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
```

**Comportamiento:**

| Elemento | Comportamiento |
|----------|---------------|
| Radio group empresa | Seleccion unica. Al cambiar, actualiza preview de datos debajo. Default: primera opcion. |
| Preview empresa | Read-only. Muestra nombre, direccion, telefono, email de la empresa seleccionada. |
| Nombre Proyecto * | Input text, required. Validacion inline si esta vacio al submit. |
| Direccion Proyecto | Input text, optional. |
| Responsable * | Input text, required. Validacion inline si esta vacio al submit. |
| Nombre Equipo | Input text, pre-filled con `equipo.clave` o `tipo_equipo`. Editable. Si se vacia, usa valor original. |
| Boton "Generar PDF" | Disabled si campos requeridos vacios. Loading state durante generacion. |
| "Volver a Resultados" | Link/boton en header, navega back sin perder datos de calculo. |

**Responsividad**:
- Desktop: formulario centrado, max-w-2xl
- Mobile: formulario full-width, padding lateral
- Radio group: stack vertical en todas las resoluciones

### 9.2 Boton "Generar PDF" en pagina de resultados

**Ubicacion**: Header de `/calculos/resultado`, alineado a la derecha, junto al boton "Nuevo Calculo".

| Estado | Apariencia |
|--------|-----------|
| Default | Icono de documento + texto "Generar PDF". Color primary. |
| Hover | Fondo oscurecido, cursor pointer |

**Responsividad**:
- Desktop: texto completo "Generar PDF"
- Mobile: solo icono (sin texto)

**Accion**: Navega a `/calculos/resultado/pdf` (NO genera PDF directamente).

### 9.3 Documento PDF

**Layout general:**

```
┌──────────────────────────────────┐
│ HEADER (cada pagina):            │
│ [Logo Empresa]  Memoria de       │
│                 Calculo Electrica│
│                          Fecha   │
├──────────────────────────────────┤
│                                  │
│  Proyecto: Planta Industrial N.  │
│  Direccion: Av. Industrial 456   │
│                                  │
│  1. Datos Generales              │
│  ┌────────────────────────────┐  │
│  │ Equipo: Filtro Activo L3   │  │
│  │ Tipo | Voltaje | Sistema   │  │
│  │ Estado | FP                │  │
│  └────────────────────────────┘  │
│                                  │
│  2. Calculo de Corriente         │
│  Formula + tabla de factores     │
│                                  │
│  3. Dimensionamiento Alimentador │
│  Cable seleccionado + tabla      │
│                                  │
│  4. Conductor de Tierra          │
│  Tabla NOM + resultado           │
│                                  │
│  5. Canalizacion                 │
│  Formulas + resultado            │
│                                  │
│  6. Caida de Tension             │
│  Formula + cumple/no cumple      │
│                                  │
│  7. Conclusion                   │
│  ┌────────────────────────────┐  │
│  │ CUMPLE NOM-001-SEDE-2012   │  │
│  └────────────────────────────┘  │
│  Observaciones                   │
│  Resumen de resultados           │
│                                  │
│  ┌────────────────────────────┐  │
│  │ Firma: __________________ │  │
│  │ Nombre: Ing. C. Rodriguez │  │
│  │ Cedula: _________________ │  │
│  │ Fecha:  _________________ │  │
│  └────────────────────────────┘  │
│                                  │
├──────────────────────────────────┤
│ FOOTER (cada pagina):            │
│ Garfex Ing. Electrica            │
│ Av. Insurgentes Sur 1234, CDMX   │
│ +52 55 1234 5678                 │
│ contacto@garfex.com              │
│              Pagina X de Y       │
└──────────────────────────────────┘
```

**Tipografia** (tokens: `--fuente-tecnica`, `--fuente-datos`):
- Titulo principal: 16pt bold — Helvetica/Arial (`--fuente-tecnica`)
- Titulos de seccion: 13pt bold con linea inferior (`--linea-seccion`)
- Cuerpo: 10pt regular
- Tablas: 9pt
- Valores numericos: Courier (`--fuente-datos`) — corrientes, calibres, porcentajes
- Footer: 7-8pt en `--tinta-tenue`

**Colores** (tokens: `--tinta`, `--papel`, `--sello-*`):
- Tinta y papel: negro tinta (`--tinta: #1a1a1a`) sobre blanco papel (`--papel: #ffffff`)
- Texto secundario: gris grafito (`--tinta-suave: #4a4a4a`)
- Headers de tabla: papel gris (`--papel-gris: #f2f0ed`)
- Dictamen normativo: verde sello (`--sello-aprobado: #1b5e20`) / rojo sello (`--sello-rechazado: #b71c1c`)
- Sin colores decorativos. Sin colores de la app. Solo tinta, papel y sellos.

**Tablas** (el corazon de la memoria):
- Bordes en `--linea` (#cccccc), 0.5pt solid, `border-collapse: collapse`
- Header de tabla con fondo `--papel-gris` y texto bold
- Filas alternas con fondo `--papel-gris` para legibilidad en impresion
- Alineacion numerica a la derecha con `--fuente-datos` (monospace)
- Valores con unidades explicitas (A, V, mm², %, Ω)
- Padding uniforme (4pt 8pt)

### 9.4 Identidad visual del PDF — Exploracion de dominio

Antes de definir colores y tipografia, exploramos el mundo donde vive este documento.

#### Dominio: Ingenieria electrica mexicana

| Concepto | Vocabulario visual |
|----------|-------------------|
| **Memoria de calculo** | Documento tecnico formal, numerado, con secciones y tablas. Se imprime, se firma, se archiva en folder. |
| **NOM-001-SEDE-2012** | Normativa oficial. Sellos, referencias a articulos, lenguaje prescriptivo ("debe cumplir", "no menor a"). |
| **Planos electricos** | Lineas negras sobre fondo blanco. Diagramas unifilares. Cuadros de datos. Nada de color decorativo. |
| **Oficina de ingenieria** | Escritorios con planos enrollados, folders beige, sellos de goma, impresoras laser (blanco y negro). |
| **DRO (Director Responsable de Obra)** | Autoridad que revisa y sella. Busca estructura, datos claros, cumplimiento. No quiere "diseño" — quiere informacion. |

#### Color world — colores que existen naturalmente en este dominio

| Color | Donde existe | Rol en el PDF |
|-------|-------------|--------------|
| **Negro tinta** (#1a1a1a) | Texto impreso, lineas de planos, firmas | Texto principal, bordes de tabla, lineas |
| **Gris grafito** (#4a4a4a) | Lapiz de ingeniero, anotaciones, metadata | Texto secundario, subtitulos, etiquetas |
| **Gris papel** (#f2f0ed) | Papel bond envejecido, folders de expediente | Fondo de headers de tabla, separadores |
| **Blanco papel** (#ffffff) | Hoja carta nueva, superficie de trabajo | Fondo del documento |
| **Verde normativo** (#1b5e20) | Sello de "APROBADO", palomita de cumplimiento, luz verde | Indicador CUMPLE |
| **Rojo normativo** (#b71c1c) | Sello de "RECHAZADO", marca de error, señalizacion de peligro electrico | Indicador NO CUMPLE |
| **Azul acero** (#37474f) | Color de tableros electricos, ductos metalicos, herramienta | Acento sutil en headers de seccion (opcional) |

#### Signature — elemento unico de este producto

**La franja de cumplimiento normativo.** En la conclusion, un bloque prominente que dice "CUMPLE CON NOM-001-SEDE-2012" en verde o "REQUIERE REVISION" en rojo — similar a un sello oficial. Esto NO es un badge generico de "success/error": es el dictamen tecnico del documento. Es lo primero que busca quien revisa la memoria. Si ves ese bloque, sabes que es una memoria de calculo de Garfex.

#### Defaults que rechazamos

| Default generico | Por que lo rechazamos | Alternativa |
|-----------------|----------------------|-------------|
| Colores primarios de la app (primary/accent) en el PDF | El PDF no es la app — es un documento tecnico. Los colores de UI no pertenecen al papel | Escala de grises + solo verde/rojo normativo |
| Tipografia moderna sans-serif thin (Inter, Geist) | Demasiado "tech/startup". Los documentos de ingenieria usan fuentes con mas peso y tradicion | Helvetica/Arial en pesos regulares y bold. Courier para valores numericos |
| Cards con border-radius y sombras | No existen en papel impreso. Se ven artificiales en PDF | Tablas con bordes rectos, secciones separadas con lineas horizontales |
| Badges pill-shaped para cumplimiento | Parecen UI de app, no un dictamen tecnico | Bloque rectangular con fondo solido, tipo sello oficial |
| Spacing generoso entre elementos (tipo landing page) | Desperdicia papel. Los ingenieros estan acostumbrados a documentos densos | Compacto dentro de secciones, generoso ENTRE secciones |

#### Decision explicita: Familiar > Creativo

Este PDF debe verse como **las mejores memorias de calculo que los ingenieros ya conocen** — no como algo nuevo. La familiaridad con el formato de documentos tecnicos genera confianza inmediata. Un supervisor que recibe este PDF debe poder navegarlo sin instrucciones porque la estructura es la que espera: encabezado → datos → calculos → resultado → firma.

**"Sameness" en documentos tecnicos de ingenieria es una VIRTUD, no un defecto.** La consistencia con el formato establecido comunica seriedad, normatividad y profesionalismo.

---

### 9.5 Estrategia de CSS: HTML semantico puro + CSS de impresion

**Decision**: Los templates del PDF usan **HTML semantico puro + CSS custom optimizado para impresion**. NO se usa Tailwind, Bootstrap ni ningun framework CSS.

**Justificacion**:

| Criterio | Tailwind / Framework | CSS puro para impresion |
|----------|---------------------|------------------------|
| Compatibilidad wkhtmltopdf | Parcial — custom properties y features modernas pueden fallar | Total — CSS2/3 basico es 100% soportado |
| Unidades | `rem`, `px` (pantalla) | `pt`, `mm`, `cm` (impresion) |
| Layout | Flexbox/Grid (soporte parcial en Qt WebKit) | `float`, `table`, flexbox basico (seguro) |
| Page breaks | No tiene utilities para `page-break-*` | `page-break-before`, `page-break-after`, `page-break-inside: avoid` |
| Peso | ~300KB+ de CSS purged | ~5KB de CSS dedicado |
| Mantenimiento | Requiere build pipeline | Archivo CSS plano, cero dependencias |
| Predecibilidad | Varia por version de Tailwind y motor | Identico en cualquier renderizado |

**Estructura del stylesheet `pdf.css`:**

> **Tokens semanticos del dominio**: Los nombres de variables evocan el mundo de la ingenieria electrica y documentos tecnicos — no un UI framework generico. Alguien leyendo solo los tokens deberia poder adivinar que esto es un documento tecnico de ingenieria, no una app web.

```css
/* ═══════════════════════════════════════════════════════════
   TOKENS — Vocabulario visual de ingenieria electrica
   ═══════════════════════════════════════════════════════════ */
:root {
  /* --- Tinta y papel --- */
  --tinta:              #1a1a1a;  /* Negro tinta — texto principal, lineas, firmas */
  --tinta-suave:        #4a4a4a;  /* Gris grafito — subtitulos, etiquetas, metadata */
  --tinta-tenue:        #808080;  /* Gris lapiz — texto terciario, notas al pie */
  --papel:              #ffffff;  /* Hoja carta nueva — fondo del documento */
  --papel-gris:         #f2f0ed;  /* Papel bond envejecido — headers de tabla, separadores */

  /* --- Lineas y estructura --- */
  --linea:              #cccccc;  /* Borde de tabla, lineas divisorias */
  --linea-seccion:      #999999;  /* Separador entre secciones principales */

  /* --- Dictamen normativo --- */
  --sello-aprobado:     #1b5e20;  /* Verde normativo — CUMPLE NOM */
  --sello-rechazado:    #b71c1c;  /* Rojo normativo — NO CUMPLE / requiere revision */

  /* --- Tipografia --- */
  --fuente-tecnica:     "Helvetica Neue", Helvetica, Arial, sans-serif;
  --fuente-datos:       "Courier New", Courier, monospace;  /* Valores numericos, calibres */
}

/* ═══════════════════════════════════════════════════════════
   PAGINA — Tamano carta, margenes para encuadernacion
   ═══════════════════════════════════════════════════════════ */
@page {
  size: Letter;
  margin: 20mm 15mm 20mm 25mm; /* top right bottom left — izquierda mayor para encuadernar */
}

/* ═══════════════════════════════════════════════════════════
   TIPOGRAFIA — Puntos tipograficos (pt) para impresion
   ═══════════════════════════════════════════════════════════ */
body {
  font-family: var(--fuente-tecnica);
  font-size: 10pt;
  line-height: 1.4;
  color: var(--tinta);
}

h1 {
  font-size: 16pt;
  font-weight: bold;
  color: var(--tinta);
}

h2 {
  font-size: 13pt;
  font-weight: bold;
  border-bottom: 1pt solid var(--linea-seccion);
  padding-bottom: 4pt;
  margin-top: 16pt;
  color: var(--tinta);
}

h3 {
  font-size: 11pt;
  font-weight: bold;
  color: var(--tinta-suave);
}

/* Texto secundario (etiquetas, metadata) */
.label, dt {
  color: var(--tinta-suave);
  font-size: 9pt;
}

/* Valores numericos (corrientes, calibres, porcentajes) */
.dato, .numeric {
  font-family: var(--fuente-datos);
  color: var(--tinta);
}

/* ═══════════════════════════════════════════════════════════
   TABLAS — El corazon de una memoria de calculo
   ═══════════════════════════════════════════════════════════ */
table {
  width: 100%;
  border-collapse: collapse;
  font-size: 9pt;
  margin: 8pt 0;
}

th {
  background: var(--papel-gris);
  font-weight: bold;
  text-align: left;
  padding: 4pt 8pt;
  border: 0.5pt solid var(--linea);
  color: var(--tinta);
}

td {
  padding: 4pt 8pt;
  border: 0.5pt solid var(--linea);
  vertical-align: top;
  color: var(--tinta);
}

/* Alineacion numerica — valores a la derecha, monospace */
td.numeric, th.numeric {
  text-align: right;
  font-family: var(--fuente-datos);
}

/* Fila alterna para tablas largas (legibilidad en impresion) */
tbody tr:nth-child(even) {
  background: var(--papel-gris);
}

/* ═══════════════════════════════════════════════════════════
   PAGE BREAKS — Control de paginacion
   ═══════════════════════════════════════════════════════════ */
.seccion {
  page-break-inside: avoid;
}

.pagina-nueva {
  page-break-before: always;
}

/* Evitar que tablas se corten a mitad de fila */
tr {
  page-break-inside: avoid;
}

/* ═══════════════════════════════════════════════════════════
   DICTAMEN NORMATIVO — El sello del documento
   Signature element: bloque rectangular tipo sello oficial
   ═══════════════════════════════════════════════════════════ */
.dictamen {
  text-align: center;
  padding: 8pt 16pt;
  font-weight: bold;
  font-size: 13pt;
  letter-spacing: 0.5pt;
  margin: 16pt 0;
  border: 2pt solid;
}

.dictamen-cumple {
  background: var(--sello-aprobado);
  border-color: var(--sello-aprobado);
  color: white;
}

.dictamen-no-cumple {
  background: var(--sello-rechazado);
  border-color: var(--sello-rechazado);
  color: white;
}

/* ═══════════════════════════════════════════════════════════
   FIRMA — Bloque de responsable al pie del documento
   ═══════════════════════════════════════════════════════════ */
.firma-block {
  margin-top: 48pt;
  width: 55%;
  margin-left: auto;
  margin-right: auto;
  text-align: center;
}

.firma-linea {
  border-top: 1pt solid var(--tinta);
  margin-top: 64pt;
  padding-top: 4pt;
  font-size: 10pt;
}

.firma-cargo {
  font-size: 9pt;
  color: var(--tinta-suave);
  margin-top: 2pt;
}

/* ═══════════════════════════════════════════════════════════
   HEADER/FOOTER DE PAGINA (wkhtmltopdf templates separados)
   ═══════════════════════════════════════════════════════════ */
.page-header {
  font-size: 8pt;
  color: var(--tinta-suave);
  border-bottom: 0.5pt solid var(--linea);
  padding-bottom: 4pt;
}

.page-footer {
  font-size: 7pt;
  color: var(--tinta-tenue);
  border-top: 0.5pt solid var(--linea);
  padding-top: 4pt;
}
```

**Test de tokens**: Lee los nombres de variables en voz alta — `--tinta`, `--papel`, `--papel-gris`, `--linea-seccion`, `--sello-aprobado`, `--fuente-tecnica`, `--fuente-datos`. ¿Suenan como un documento tecnico de ingenieria? Si. ¿Podrian pertenecer a cualquier app web? No. Eso es diseño con intencion.

**Reglas de oro para los templates HTML:**

1. **Solo HTML semantico**: `<table>`, `<thead>`, `<tbody>`, `<h1>`-`<h3>`, `<p>`, `<dl>`, `<dt>`, `<dd>`. Sin `<div>` soup.
2. **Unidades en `pt`** para tipografia, `mm` para margenes de pagina — nunca `rem` ni `px`.
3. **`float` para layouts de 2 columnas** — wkhtmltopdf lo soporta al 100%. Flexbox solo para cosas simples.
4. **`page-break-inside: avoid`** en secciones y filas de tabla — evita cortes feos.
5. **Logo embebido como base64** en el template — evita problemas de rutas relativas con wkhtmltopdf.
6. **Sin JavaScript** en los templates — wkhtmltopdf puede ejecutar JS pero agrega latencia y riesgo.
7. **Sin `@import`** — todo el CSS inline o en un solo archivo referenciado con `--enable-local-file-access`.

---

## 10. Timeline & Milestones

### Fase 1 — Backend: Estructura + Empresas + Endpoint (3-5 dias)

| Tarea | Estimacion |
|-------|-----------|
| Crear estructura `internal/pdf/` (domain, application, infrastructure) | 0.5 dia |
| Implementar `empresa.go` con catalogo estatico de 3 empresas | 0.25 dia |
| Implementar `pdf_request.go` (DTO con MemoriaOutput + PresentacionInput) | 0.25 dia |
| Recopilar y agregar assets: 3 logos PNG en `assets/logos/` | 0.25 dia |
| Implementar templates HTML (7 secciones + header con logo + footer con empresa) | 2 dias |
| Implementar CSS de impresion | 0.5 dia |
| Implementar `html_renderer.go` (renderizar template con datos + empresa) | 0.5 dia |
| Implementar `pdf_generator.go` (wrapper de wkhtmltopdf) | 0.5 dia |
| Implementar `generar_memoria_pdf.go` (use case: resuelve empresa → renderiza → genera) | 0.5 dia |
| Implementar `pdf_handler.go` (HTTP handler: parsea request compuesto) | 0.5 dia |
| Registrar ruta en router | 0.25 dia |

### Fase 2 — Docker + wkhtmltopdf (1 dia)

| Tarea | Estimacion |
|-------|-----------|
| Crear/actualizar Dockerfile con wkhtmltopdf | 0.5 dia |
| Verificar funcionamiento headless (sin X11) | 0.5 dia |

### Fase 3 — Frontend: Pagina configuracion + Formulario (2-3 dias)

| Tarea | Estimacion |
|-------|-----------|
| Agregar boton "Generar PDF" en resultado/+page.svelte (navega a /pdf) | 0.25 dia |
| Crear `empresas-pdf.ts` con datos estaticos de las 3 empresas | 0.25 dia |
| Crear pagina `/calculos/resultado/pdf/+page.svelte` con formulario | 1 dia |
| Implementar SelectorEmpresa con radio group + preview datos | 0.5 dia |
| Campos: nombre proyecto, direccion, responsable, nombre equipo | 0.5 dia |
| Validacion client-side (required, max length) | 0.25 dia |
| Implementar fetch + blob download + loading/error states | 0.5 dia |
| Responsividad mobile del formulario | 0.25 dia |

### Fase 4 — Testing (2-3 dias)

| Tarea | Estimacion |
|-------|-----------|
| Tests unitarios: html_renderer (template renderiza con empresa correcta) | 0.5 dia |
| Tests unitarios: pdf_generator (mock de wkhtmltopdf) | 0.5 dia |
| Tests unitarios: pdf_handler (validacion de PresentacionInput, empresa invalida) | 0.5 dia |
| Tests integracion: endpoint genera PDF con branding correcto | 0.5 dia |
| Tests frontend: formulario, validacion, loading, error states | 0.5 dia |
| QA manual: revisar PDF con cada empresa (3 logos diferentes) | 0.5 dia |

### Fase 5 — Polish + Documentacion (1 dia)

| Tarea | Estimacion |
|-------|-----------|
| Ajustar estilos CSS del PDF y posicion de logos | 0.5 dia |
| Actualizar Swagger/OpenAPI con nuevo endpoint (incluir PresentacionInput schema) | 0.25 dia |
| Actualizar AGENTS.md con nueva feature | 0.25 dia |

### Total estimado: 9-13 dias

---

## 11. Risks & Mitigation

| # | Riesgo | Probabilidad | Impacto | Mitigacion |
|---|--------|-------------|---------|-----------|
| R1 | wkhtmltopdf no disponible o falla en Docker | Media | Alto | Usar build estatico de wkhtmltopdf sin dependencia de xvfb. Tener fallback con chromedp/headless Chrome como alternativa. |
| R2 | Renderizado HTML→PDF no replica fielmente las secciones | Media | Medio | Crear tests visuales (snapshot tests). Disenar CSS print-specific desde el inicio. |
| R3 | Tiempo de generacion >5s para memorias complejas | Baja | Medio | Timeout configurable. Medir P95 en staging. Optimizar templates si necesario. |
| R4 | SVG de diagramas no renderiza en wkhtmltopdf | Alta | Bajo | Excluido del MVP (out of scope). En v2, convertir SVG→PNG server-side antes de renderizar. |
| R5 | Tamano de PDF excesivo (>5MB) | Baja | Bajo | Comprimir imagenes. Optimizar CSS. Limitar resolucion. |
| R6 | Memory leak por procesos wkhtmltopdf sin terminar | Baja | Alto | Context con timeout. Matar proceso si excede limite. Rate limiting en endpoint. |
| R7 | Body de MemoriaOutput manipulado maliciosamente | Media | Medio | Validar estructura del JSON en el handler. Sanitizar valores antes de inyectar en HTML template (html/template ya escapa por defecto). |
| R8 | Concurrencia: multiples PDFs generandose simultaneamente saturan CPU/memoria | Baja | Alto | Semaforo/pool de workers limitado (ej: max 3 generaciones concurrentes). |

---

## 12. Dependencies & Assumptions

### Dependencies

| # | Dependencia | Tipo | Estado |
|---|------------|------|--------|
| D1 | wkhtmltopdf binario disponible en el entorno de ejecucion | Externa | Pendiente — requiere Dockerfile |
| D2 | MemoriaOutput del backend es estable y completa | Interna | Estable — no se esperan cambios de estructura |
| D3 | Templates HTML deben reflejar las 7 secciones de MemoriaTecnica.svelte | Interna | Las secciones ya estan definidas y estables |
| D4 | Frontend ya tiene la data de MemoriaOutput disponible en la pagina de resultados | Interna | Disponible — via page data |
| D5 | Gin router permite agregar nuevas rutas sin conflictos | Interna | OK — patron ya establecido |

### Assumptions

| # | Asuncion | Riesgo si es falsa |
|---|---------|-------------------|
| A1 | Los usuarios tienen acceso a impresora o visor de PDF | Bajo — todos los navegadores modernos abren PDF |
| A2 | El formato Letter (carta) es el estandar en Mexico para memorias de calculo | Bajo — es el estandar de facto; A4 como alternativa futura |
| A3 | No se necesita autenticacion para el endpoint de PDF (mismo contexto que el calculo) | Medio — si se agrega auth en el futuro, el endpoint debera protegerse |
| A4 | wkhtmltopdf 0.12.6+ soporta CSS3 suficiente para tablas y layout profesional | Bajo — ampliamente probado |
| A5 | El volumen de generacion de PDFs no requiere cola de jobs (sync es suficiente) | Medio — si el uso crece, considerar async con job queue |

---

## 13. Open Questions

| # | Pregunta | Responsable | Fecha limite | Estado |
|---|---------|-------------|-------------|--------|
| Q1 | ~~Debemos incluir el logo en el header del PDF?~~ | — | — | **Resuelta**: Si, el logo viene de la empresa seleccionada (Summa/Garfex/Siemens) |
| Q2 | El PDF debe tener tamanio Letter o A4? O configurable? | Product | Antes de Fase 1 | Abierta |
| Q3 | ~~Se necesita campo editable para nombre del proyecto?~~ | — | — | **Resuelta**: Si, formulario pre-PDF con nombre proyecto, direccion, responsable |
| Q4 | ~~Que datos debe llevar el bloque de firma?~~ | — | — | **Resuelta**: Firma (linea), Nombre (pre-llenado), Cedula Profesional (vacio), Fecha (vacio) |
| Q5 | Se quiere rate limiting en el endpoint de PDF? Cuantas generaciones por minuto por IP? | Engineering | Antes de Fase 2 | Abierta |
| Q6 | Alternativa a wkhtmltopdf: Evaluar chromedp (headless Chrome) vs wkhtmltopdf? chromedp tiene mejor soporte CSS pero mas pesado. | Engineering | Antes de Fase 2 | Abierta |
| Q7 | Los diagramas SVG de arreglo de cables (charola/tuberia) deberian incluirse en v2 del PDF? | Product | Post-MVP | Abierta |
| Q8 | Se necesita versionado del PDF (v1, v2) si el usuario recalcula con diferentes parametros? | Product | Post-MVP | Abierta |
| Q9 | Datos exactos de las 3 empresas (Summa, Garfex, Siemens): direccion, telefono, email. Necesitamos los valores reales para la configuracion estatica. | Product | Antes de Fase 1 | **Abierta — CRITICA** |
| Q10 | Assets de logos: se necesitan los 3 logos en formato PNG con fondo transparente, max 200x80px. Donde estan o quien los provee? | Design | Antes de Fase 1 | **Abierta — CRITICA** |
| Q11 | El campo "Nombre del Equipo/Carga" debe tener algun prefijo o formato? Ej: "FA-" para filtro activo? | Product | Antes de Fase 3 | Abierta |
| Q12 | Si en el futuro se necesitan mas empresas, como se agregan? Deploy con codigo actualizado es aceptable? | Engineering | Post-MVP | Abierta |

---

## Appendix A: Data Flow Detail

### Request flow completo

```
1. Usuario completa calculo en pagina principal (/)
2. Frontend hace POST /api/v1/calculos/memoria → recibe MemoriaOutput
3. Frontend redirige a /calculos/resultado?data=<encoded MemoriaOutput>
4. Pagina de resultados muestra MemoriaTecnica con las 7 secciones
5. Usuario presiona "Generar PDF"
6. Frontend navega a /calculos/resultado/pdf (pasando MemoriaOutput via state)
7. Pagina de configuracion muestra formulario:
   a. Selector de empresa (Summa / Garfex / Siemens)
   b. Nombre del proyecto (requerido)
   c. Direccion del proyecto (opcional)
   d. Responsable del reporte (requerido)
   e. Nombre del equipo/carga (pre-llenado, editable)
8. Usuario llena formulario y presiona "Generar PDF"
9. Frontend hace POST /api/v1/pdf/memoria con { memoria, presentacion }
10. Backend:
    a. Parsea y valida request (MemoriaOutput + PresentacionInput)
    b. Resuelve empresa del catalogo estatico por empresa_id
    c. Carga logo de empresa desde assets/logos/
    d. Construye TemplateData (empresa + proyecto + calculo)
    e. Renderiza HTML con html/template (header con logo + 7 secciones + firma + footer con empresa)
    f. Escribe HTML temporal a disco
    g. Ejecuta wkhtmltopdf: HTML → PDF
    h. Lee PDF resultante
    i. Limpia archivos temporales
    j. Retorna PDF como application/pdf
11. Frontend recibe blob, crea URL, dispara descarga
12. Usuario obtiene PDF descargado con branding de empresa
```

### Estructura del body del POST

```json
{
  "memoria": {
    "equipo": { ... },
    "tipo_equipo": "FILTRO_ACTIVO",
    "factor_potencia": 0.90,
    "estado": "Nuevo Leon",
    "instalacion": { ... },
    "corrientes": { ... },
    "cable_fase": { ... },
    "cable_neutro": null,
    "cable_tierra": { ... },
    "canalizacion": { ... },
    "proteccion": { ... },
    "caida_tension": { ... },
    "cumple_normativa": true,
    "observaciones": [ ... ],
    "pasos": [ ... ]
  },
  "presentacion": {
    "empresa_id": "garfex",
    "nombre_proyecto": "Planta Industrial Norte",
    "direccion_proyecto": "Av. Industrial 456, Monterrey, NL",
    "responsable": "Ing. Carlos Rodriguez",
    "nombre_equipo_override": "Filtro Activo Linea 3"
  }
}
```

---

## Appendix B: Template HTML Reference

Mapeo entre secciones del frontend Svelte, templates HTML del PDF y fuentes de datos:

| # | Componente Svelte | Template HTML | Datos principales |
|---|---|---|---|
| — | N/A (nuevo) | header.html | Empresa.Logo, "Memoria de Calculo Electrica", FechaGeneracion |
| — | N/A (nuevo) | footer.html | Empresa.Nombre, Empresa.Direccion, Empresa.Telefono, Empresa.Email, "Pagina X de Y" |
| 1 | SeccionEncabezado.svelte | seccion_encabezado.html | NombreProyecto, DireccionProyecto, NombreEquipo, tipo_equipo, instalacion |
| 2 | SeccionCorriente.svelte | seccion_corriente.html | corrientes.* |
| 3 | SeccionAlimentador.svelte | seccion_alimentador.html | cable_fase, corrientes.corriente_por_hilo |
| 4 | SeccionTierra.svelte | seccion_tierra.html | cable_tierra, proteccion.itm |
| 5 | SeccionCanalizacion.svelte | seccion_canalizacion.html | canalizacion.* |
| 6 | SeccionCaidaTension.svelte | seccion_caida_tension.html | caida_tension.* |
| 7 | SeccionConclusion.svelte | seccion_conclusion.html | cumple_normativa, observaciones, Responsable (firma) |

---

## Appendix C: Catalogo de Empresas (Presentaciones)

Empresas predefinidas para MVP. Datos estaticos, no editables por el usuario.

| Campo | Garfex | Summa | Siemens |
|-------|--------|-------|---------|
| **id** | `garfex` | `summa` | `siemens` |
| **nombre** | Garfex Ingenieria Electrica | Summa Energetica | Siemens Energy |
| **direccion** | *(por definir — ver Q9)* | *(por definir)* | *(por definir)* |
| **telefono** | *(por definir)* | *(por definir)* | *(por definir)* |
| **email** | *(por definir)* | *(por definir)* | *(por definir)* |
| **logo** | `assets/logos/garfex.png` | `assets/logos/summa.png` | `assets/logos/siemens.png` |

> **ACCION REQUERIDA**: Completar los datos reales de cada empresa antes de iniciar Fase 1. Ver Q9 y Q10 en Open Questions.

**Requisitos de logos:**
- Formato: PNG con fondo transparente
- Tamano maximo: 200x80 px (se escala automaticamente)
- Resolucion: minimo 150 DPI para impresion

**Para agregar una nueva empresa en el futuro:**
1. Agregar entrada en `empresa.go` (catalogo Go)
2. Agregar entrada en `empresas-pdf.ts` (catalogo frontend)
3. Agregar logo en `assets/logos/<id>.png`
4. Deploy

---

## Appendix D: Validacion del PRD

### Changelog

| Version | Fecha | Cambios |
|---------|-------|---------|
| 1.0 | 2026-03-02 | PRD inicial: boton → genera PDF directamente |
| 1.1 | 2026-03-02 | Pagina intermedia de configuracion, selector de empresa (Summa/Garfex/Siemens), datos de proyecto, responsable, nombre equipo editable. 16 user stories (vs 10 originales). Nuevo endpoint con PresentacionInput. |
| 1.2 | 2026-03-02 | Decision tecnica: CSS puro para impresion (no Tailwind). Seccion 9.5 con estrategia CSS, estructura de pdf.css, @page rules, reglas de oro para templates. Evaluacion Tailwind vs CSS puro en tabla comparativa. |
| 1.3 | 2026-03-02 | Design intent (seccion 9.0): personas, contexto, feeling para formulario y PDF. Identidad visual del PDF (seccion 9.4): exploracion de dominio, color world, signature (dictamen normativo), defaults rechazados. Tokens CSS renombrados de genericos (--color-text) a semanticos del dominio (--tinta, --papel, --sello-aprobado, --fuente-datos). |

### Checklist de completitud

- [x] Problema claramente articulado
- [x] Usuarios identificados con personas
- [x] Exito es medible con metricas concretas
- [x] Scope acotado con in/out claros
- [x] Requerimientos son testeables (acceptance criteria)
- [x] Timeline estimado con fases
- [x] Riesgos identificados con mitigacion
- [x] Stakeholders pueden alinearse con este documento
- [x] No hay texto placeholder pendiente
- [x] Consideraciones tecnicas validadas contra codebase existente
