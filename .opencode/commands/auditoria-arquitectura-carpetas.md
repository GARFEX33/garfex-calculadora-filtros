---
description: Invoca al agente auditor-arquitectura para realizar una auditoría rápida de la estructura de carpetas y archivos, verificando cumplimiento de Clean Architecture + DDD + Hexagonal + Vertical Slicing en Go. Solo analiza nombres y organización, sin leer el código.
agent: auditor-arquitectura
disable-model-invocation: false
---

# Auditoría rápida de estructura de carpetas y archivos

invoke-skill: clean-ddd-hexagonal-vertical-go-enterprise
context:
tipo_auditoria: "estructura-carpetas-archivos"
objetivo: "Revisar organización de módulos, capas y vertical slices sin inspeccionar código"
scope: "estructura-sin-contenido"
checklist: - [ ] ¿Existe separación clara entre carpetas de dominio, aplicación e infraestructura? - [ ] ¿Cada módulo o vertical slice tiene su propia carpeta independiente? - [ ] ¿Los ports y adapters están correctamente ubicados según Hexagonal Architecture? - [ ] ¿Las entidades, value objects y aggregates están en la carpeta de dominio? - [ ] ¿Los casos de uso o servicios de aplicación están en la carpeta de aplicación? - [ ] ¿Los handlers, repositorios y adaptadores externos están en infraestructura? - [ ] ¿Cada paquete tiene un archivo README o documentación mínima? - [ ] ¿Se respeta la convención de nombres de carpetas y archivos Go? - [ ] ¿No hay acoplamientos indebidos entre vertical slices o módulos? - [ ] ¿La estructura comunica claramente la intención de cada módulo (Screaming Architecture)?
output-format: "checklist-report"
