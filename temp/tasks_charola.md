# Tareas Infrastructure Agent - Endpoint Charola

## Tarea 1: Implementar ObtenerTablaCharola en CSVTablaNOMRepository
- **Archivo**: internal/calculos/infrastructure/adapter/driven/csv/csv_tabla_nom_repository.go
- **Acción**: Agregar el método ObtenerTablaCharola que retorne la tabla según el tipo de canalización

## Tarea 2: Crear CharolaHandler HTTP
- **Archivo**: internal/calculos/infrastructure/adapter/driver/http/charola_handler.go
- **Acción**: Crear handler con endpoints POST /charola/espaciado y POST /charola/triangular

## Tarea 3: Actualizar Router
- **Archivo**: internal/calculos/infrastructure/router.go
- **Acción**: Agregar rutas de charola y parámetros de use cases

## Tarea 4: Verificar compilación
- **Comando**: go build ./...
