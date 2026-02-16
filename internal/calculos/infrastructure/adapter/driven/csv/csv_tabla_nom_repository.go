// internal/calculos/infrastructure/adapter/driven/csv/csv_tabla_nom_repository.go
package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// impedanciaEntry holds all impedance values for a given calibre from Tabla 9.
type impedanciaEntry struct {
	ReactanciaAl    float64
	ReactanciaAcero float64
	ResCuPVC        float64
	ResCuAl         float64
	ResCuAcero      float64
	ResAlPVC        float64
	ResAlAl         float64
	ResAlAcero      float64
}

// diametroConductorEntry holds diameter values for conductors from Tabla 5.
type diametroConductorEntry struct {
	DiamTWTHW   float64
	DiamRHH_RHW float64
	DiamXHHW    float64
}

// CSVTablaNOMRepository reads NOM tables from CSV files with in-memory caching.
type CSVTablaNOMRepository struct {
	basePath             string
	tablaTierra          []valueobject.EntradaTablaTierra
	tablasAmpacidad      map[entity.TipoCanalizacion]map[valueobject.MaterialConductor]map[valueobject.Temperatura][]valueobject.EntradaTablaConductor
	tablaImpedancia      map[string]impedanciaEntry // key: calibre
	tablaConduit         []valueobject.EntradaTablaCanalizacion
	tablasCharola        map[entity.TipoCanalizacion][]valueobject.EntradaTablaCanalizacion
	estadosTemperatura   map[string]int
	factoresTemperatura  []service.EntradaTablaFactorTemperatura
	factoresAgrupamiento []service.EntradaTablaFactorAgrupamiento
	tablaDiametros       map[string]diametroConductorEntry
}

// NewCSVTablaNOMRepository creates a new repository and loads all tables into memory.
func NewCSVTablaNOMRepository(basePath string) (*CSVTablaNOMRepository, error) {
	// Verify directory exists
	info, err := os.Stat(basePath)
	if err != nil {
		return nil, fmt.Errorf("cannot access base path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("base path is not a directory: %s", basePath)
	}

	repo := &CSVTablaNOMRepository{
		basePath:        basePath,
		tablasAmpacidad: make(map[entity.TipoCanalizacion]map[valueobject.MaterialConductor]map[valueobject.Temperatura][]valueobject.EntradaTablaConductor),
	}

	// Load ground conductor table
	tablaTierra, err := repo.loadTablaTierra()
	if err != nil {
		return nil, fmt.Errorf("failed to load ground table: %w", err)
	}
	repo.tablaTierra = tablaTierra

	// Load impedance table (Tabla 9)
	tablaImpedancia, err := repo.loadTablaImpedancia()
	if err != nil {
		return nil, fmt.Errorf("failed to load impedance table: %w", err)
	}
	repo.tablaImpedancia = tablaImpedancia

	// Load conduit sizing table
	tablaConduit, err := repo.loadTablaConduit()
	if err != nil {
		return nil, fmt.Errorf("failed to load conduit sizing table: %w", err)
	}
	repo.tablaConduit = tablaConduit

	// Load cable tray sizing tables
	repo.tablasCharola = make(map[entity.TipoCanalizacion][]valueobject.EntradaTablaCanalizacion)
	repo.tablasCharola[entity.TipoCanalizacionCharolaCableEspaciado] = repo.crearTablaCharolaEspaciado()
	repo.tablasCharola[entity.TipoCanalizacionCharolaCableTriangular] = repo.crearTablaCharolaTriangular()

	// Load ampacity tables for conduit types
	for _, canalizacion := range []entity.TipoCanalizacion{
		entity.TipoCanalizacionTuberiaPVC,
		entity.TipoCanalizacionTuberiaAluminio,
		entity.TipoCanalizacionTuberiaAceroPG,
		entity.TipoCanalizacionTuberiaAceroPD,
	} {
		repo.tablasAmpacidad[canalizacion] = make(map[valueobject.MaterialConductor]map[valueobject.Temperatura][]valueobject.EntradaTablaConductor)

		for _, material := range []valueobject.MaterialConductor{
			valueobject.MaterialCobre,
			valueobject.MaterialAluminio,
		} {
			repo.tablasAmpacidad[canalizacion][material] = make(map[valueobject.Temperatura][]valueobject.EntradaTablaConductor)

			tabla, err := repo.loadTablaAmpacidad("310-15-b-16.csv", material)
			if err != nil {
				return nil, fmt.Errorf("failed to load ampacity table for %s %s: %w", canalizacion, material, err)
			}

			// Extract by temperature
			for _, temp := range []valueobject.Temperatura{valueobject.Temp60, valueobject.Temp75, valueobject.Temp90} {
				repo.tablasAmpacidad[canalizacion][material][temp] = extractByTemperature(tabla, material, temp)
			}
		}
	}

	// Load ampacity tables for cable trays (charolas)
	// Charola cable espaciado -> 310-15-b-17.csv
	repo.tablasAmpacidad[entity.TipoCanalizacionCharolaCableEspaciado] = make(map[valueobject.MaterialConductor]map[valueobject.Temperatura][]valueobject.EntradaTablaConductor)
	for _, material := range []valueobject.MaterialConductor{
		valueobject.MaterialCobre,
		valueobject.MaterialAluminio,
	} {
		repo.tablasAmpacidad[entity.TipoCanalizacionCharolaCableEspaciado][material] = make(map[valueobject.Temperatura][]valueobject.EntradaTablaConductor)

		tabla, err := repo.loadTablaAmpacidad("310-15-b-17.csv", material)
		if err != nil {
			return nil, fmt.Errorf("failed to load ampacity table for charola espaciado %s: %w", material, err)
		}

		// Extract by temperature
		for _, temp := range []valueobject.Temperatura{valueobject.Temp60, valueobject.Temp75, valueobject.Temp90} {
			repo.tablasAmpacidad[entity.TipoCanalizacionCharolaCableEspaciado][material][temp] = extractByTemperature(tabla, material, temp)
		}
	}

	// Charola cable triangular -> 310-15-b-20.csv
	repo.tablasAmpacidad[entity.TipoCanalizacionCharolaCableTriangular] = make(map[valueobject.MaterialConductor]map[valueobject.Temperatura][]valueobject.EntradaTablaConductor)
	for _, material := range []valueobject.MaterialConductor{
		valueobject.MaterialCobre,
		valueobject.MaterialAluminio,
	} {
		repo.tablasAmpacidad[entity.TipoCanalizacionCharolaCableTriangular][material] = make(map[valueobject.Temperatura][]valueobject.EntradaTablaConductor)

		tabla, err := repo.loadTablaAmpacidad("310-15-b-20.csv", material)
		if err != nil {
			return nil, fmt.Errorf("failed to load ampacity table for charola triangular %s: %w", material, err)
		}

		// Extract by temperature - charola triangular no tiene 60C, solo 75C y 90C
		for _, temp := range []valueobject.Temperatura{valueobject.Temp75, valueobject.Temp90} {
			repo.tablasAmpacidad[entity.TipoCanalizacionCharolaCableTriangular][material][temp] = extractByTemperature(tabla, material, temp)
		}
	}

	// Load estados_temperatura.csv
	estadosTemp, err := repo.loadEstadosTemperatura()
	if err != nil {
		return nil, fmt.Errorf("failed to load estados_temperatura: %w", err)
	}
	repo.estadosTemperatura = estadosTemp

	// Load factores_temperatura (310-15-b-2-a.csv)
	factoresTemp, err := repo.loadFactoresTemperatura()
	if err != nil {
		return nil, fmt.Errorf("failed to load factores_temperatura: %w", err)
	}
	repo.factoresTemperatura = factoresTemp

	// Load factores_agrupamiento (310-15-b-3-a.csv)
	factoresAgr, err := repo.loadFactoresAgrupamiento()
	if err != nil {
		return nil, fmt.Errorf("failed to load factores_agrupamiento: %w", err)
	}
	repo.factoresAgrupamiento = factoresAgr

	// Load tabla diametros (tabla-5-dimensiones-aislamiento.csv)
	tablaDiam, err := repo.loadTablaDiametros()
	if err != nil {
		return nil, fmt.Errorf("failed to load tabla diametros: %w", err)
	}
	repo.tablaDiametros = tablaDiam

	return repo, nil
}

// ObtenerTablaTierra returns the ground conductor table (250-122).
func (r *CSVTablaNOMRepository) ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error) {
	return r.tablaTierra, nil
}

// ObtenerTablaAmpacidad returns ampacity table entries for the given conduit type, material, and temperature.
func (r *CSVTablaNOMRepository) ObtenerTablaAmpacidad(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
) ([]valueobject.EntradaTablaConductor, error) {
	byMaterial, ok := r.tablasAmpacidad[canalizacion]
	if !ok {
		return nil, fmt.Errorf("no ampacity table for conduit type: %s", canalizacion)
	}

	byTemp, ok := byMaterial[material]
	if !ok {
		return nil, fmt.Errorf("no ampacity table for material: %s", material)
	}

	tabla, ok := byTemp[temperatura]
	if !ok {
		return nil, fmt.Errorf("no ampacity table for temperature: %d°C", temperatura)
	}

	return tabla, nil
}

// ObtenerCapacidadConductor returns the ampacity for a specific calibre.
func (r *CSVTablaNOMRepository) ObtenerCapacidadConductor(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
	calibre string,
) (float64, error) {
	tabla, err := r.ObtenerTablaAmpacidad(ctx, canalizacion, material, temperatura)
	if err != nil {
		return 0, fmt.Errorf("obtener tabla ampacidad: %w", err)
	}

	for _, entrada := range tabla {
		if entrada.Conductor.Calibre == calibre {
			return entrada.Capacidad, nil
		}
	}

	return 0, fmt.Errorf("calibre %s no encontrado en tabla", calibre)
}

// ObtenerImpedancia returns R and X values for the given calibre and conduit type.
func (r *CSVTablaNOMRepository) ObtenerImpedancia(
	ctx context.Context,
	calibre string,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
) (valueobject.ResistenciaReactancia, error) {
	entry, ok := r.tablaImpedancia[calibre]
	if !ok {
		return valueobject.ResistenciaReactancia{}, fmt.Errorf("calibre not found in impedance table: %s", calibre)
	}

	// Determine reactance based on conduit type
	var x float64
	switch canalizacion {
	case entity.TipoCanalizacionTuberiaAceroPG, entity.TipoCanalizacionTuberiaAceroPD:
		x = entry.ReactanciaAcero
	default:
		x = entry.ReactanciaAl
	}

	// Determine resistance based on material and conduit type
	var res float64
	if material == valueobject.MaterialCobre {
		switch canalizacion {
		case entity.TipoCanalizacionTuberiaAluminio:
			res = entry.ResCuAl
		case entity.TipoCanalizacionTuberiaAceroPG, entity.TipoCanalizacionTuberiaAceroPD:
			res = entry.ResCuAcero
		default:
			res = entry.ResCuPVC
		}
	} else { // Aluminio
		switch canalizacion {
		case entity.TipoCanalizacionTuberiaAluminio:
			res = entry.ResAlAl
		case entity.TipoCanalizacionTuberiaAceroPG, entity.TipoCanalizacionTuberiaAceroPD:
			res = entry.ResAlAcero
		default:
			res = entry.ResAlPVC
		}
	}

	return valueobject.NewResistenciaReactancia(res, x)
}

// ObtenerTablaCanalizacion returns conduit sizing table entries.
func (r *CSVTablaNOMRepository) ObtenerTablaCanalizacion(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
) ([]valueobject.EntradaTablaCanalizacion, error) {
	switch canalizacion {
	case entity.TipoCanalizacionTuberiaPVC,
		entity.TipoCanalizacionTuberiaAluminio,
		entity.TipoCanalizacionTuberiaAceroPG,
		entity.TipoCanalizacionTuberiaAceroPD:
		return r.tablaConduit, nil
	case entity.TipoCanalizacionCharolaCableEspaciado,
		entity.TipoCanalizacionCharolaCableTriangular:
		tabla, ok := r.tablasCharola[canalizacion]
		if !ok {
			return nil, fmt.Errorf("tabla de charola no cargada para: %s", canalizacion)
		}
		return tabla, nil
	default:
		return nil, fmt.Errorf("tipo de canalización no soportado: %s", canalizacion)
	}
}

// ObtenerTemperaturaPorEstado returns the average temperature for a given Mexican state.
func (r *CSVTablaNOMRepository) ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error) {
	temp, ok := r.estadosTemperatura[estado]
	if !ok {
		return 0, fmt.Errorf("estado no encontrado: %s", estado)
	}
	return temp, nil
}

// ObtenerFactorTemperatura returns the temperature correction factor based on ambient temperature and conductor temperature.
func (r *CSVTablaNOMRepository) ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error) {
	return service.CalcularFactorTemperatura(tempAmbiente, tempConductor, r.factoresTemperatura)
}

// ObtenerFactorAgrupamiento returns the grouping factor based on the number of conductors.
func (r *CSVTablaNOMRepository) ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error) {
	return service.CalcularFactorAgrupamiento(cantidadConductores, r.factoresAgrupamiento)
}

// ObtenerDiametroConductor returns the diameter in mm for a given calibre, material, and insulation type.
func (r *CSVTablaNOMRepository) ObtenerDiametroConductor(ctx context.Context, calibre string, material string, conAislamiento bool) (float64, error) {
	entry, ok := r.tablaDiametros[calibre]
	if !ok {
		return 0, fmt.Errorf("calibre no encontrado en tabla de diametros: %s", calibre)
	}

	if conAislamiento {
		return entry.DiamTWTHW, nil
	}
	return entry.DiamRHH_RHW, nil
}

// ObtenerCharolaPorAncho returns the smallest tray size that fits the required width.
func (r *CSVTablaNOMRepository) ObtenerCharolaPorAncho(ctx context.Context, anchoRequeridoMM float64) (valueobject.EntradaTablaCanalizacion, error) {
	tabla, ok := r.tablasCharola[entity.TipoCanalizacionCharolaCableEspaciado]
	if !ok {
		return valueobject.EntradaTablaCanalizacion{}, fmt.Errorf("tabla de charola no cargada")
	}

	for _, entrada := range tabla {
		anchoMM := parseAnchoCharola(entrada.Tamano)
		if anchoMM >= anchoRequeridoMM {
			return entrada, nil
		}
	}

	return valueobject.EntradaTablaCanalizacion{}, fmt.Errorf("no se encontró charola para ancho requerido: %.2f mm", anchoRequeridoMM)
}

func parseAnchoCharola(tamano string) float64 {
	var ancho float64
	if _, err := fmt.Sscanf(tamano, "%fmm", &ancho); err == nil {
		return ancho
	}
	return 0
}

func (r *CSVTablaNOMRepository) loadTablaTierra() ([]valueobject.EntradaTablaTierra, error) {
	filePath := filepath.Join(r.basePath, "250-122.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open 250-122.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read 250-122.csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("250-122.csv is empty or missing header")
	}

	// Validate header
	header := records[0]
	expectedHeader := []string{"itm_hasta", "cu_calibre", "cu_seccion_mm2", "al_calibre", "al_seccion_mm2"}
	for i, col := range expectedHeader {
		if i >= len(header) || header[i] != col {
			return nil, fmt.Errorf("250-122.csv: invalid header at position %d, expected %q got %q", i, col, header[i])
		}
	}

	var result []valueobject.EntradaTablaTierra
	for i, record := range records[1:] {
		if len(record) < 3 {
			continue
		}

		itm, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("250-122.csv line %d: invalid ITM value: %w", i+2, err)
		}

		cuSeccion, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("250-122.csv line %d: invalid cu_seccion_mm2: %w", i+2, err)
		}

		entrada := valueobject.EntradaTablaTierra{
			ITMHasta: itm,
			ConductorCu: valueobject.ConductorParams{
				Calibre:    record[1],
				Material:   valueobject.MaterialCobre,
				SeccionMM2: cuSeccion,
			},
			ConductorAl: nil,
		}

		// Parse Al columns if present and non-empty
		if len(record) >= 5 && record[3] != "" && record[4] != "" {
			alSeccion, err := strconv.ParseFloat(record[4], 64)
			if err != nil {
				return nil, fmt.Errorf("250-122.csv line %d: invalid al_seccion_mm2: %w", i+2, err)
			}
			alParams := valueobject.ConductorParams{
				Calibre:    record[3],
				Material:   valueobject.MaterialAluminio,
				SeccionMM2: alSeccion,
			}
			entrada.ConductorAl = &alParams
		}

		result = append(result, entrada)
	}

	return result, nil
}

// crearTablaCharolaEspaciado crea la tabla de dimensiones para charola cable espaciado.
// Los tamaños son anchos estándar de charola en mm.
func (r *CSVTablaNOMRepository) crearTablaCharolaEspaciado() []valueobject.EntradaTablaCanalizacion {
	// Tamaños estándar de charola (ancho en mm) y área aproximada
	// Asumiendo charola de 50mm de alto estándar
	return []valueobject.EntradaTablaCanalizacion{
		{Tamano: "50mm", AreaInteriorMM2: 2500},   // 50mm x 50mm
		{Tamano: "100mm", AreaInteriorMM2: 5000},  // 100mm x 50mm
		{Tamano: "150mm", AreaInteriorMM2: 7500},  // 150mm x 50mm
		{Tamano: "200mm", AreaInteriorMM2: 10000}, // 200mm x 50mm
		{Tamano: "300mm", AreaInteriorMM2: 15000}, // 300mm x 50mm
		{Tamano: "450mm", AreaInteriorMM2: 22500}, // 450mm x 50mm
		{Tamano: "600mm", AreaInteriorMM2: 30000}, // 600mm x 50mm
	}
}

// crearTablaCharolaTriangular crea la tabla de dimensiones para charola cable triangular.
// Los tamaños son anchos estándar de charola en mm.
func (r *CSVTablaNOMRepository) crearTablaCharolaTriangular() []valueobject.EntradaTablaCanalizacion {
	// Mismos tamaños pero el área efectiva es diferente por la forma triangular
	// Se usa 40% del área rectangular para el cálculo de llenado
	return []valueobject.EntradaTablaCanalizacion{
		{Tamano: "50mm", AreaInteriorMM2: 2500},
		{Tamano: "100mm", AreaInteriorMM2: 5000},
		{Tamano: "150mm", AreaInteriorMM2: 7500},
		{Tamano: "200mm", AreaInteriorMM2: 10000},
		{Tamano: "300mm", AreaInteriorMM2: 15000},
		{Tamano: "450mm", AreaInteriorMM2: 22500},
		{Tamano: "600mm", AreaInteriorMM2: 30000},
	}
}

// rawAmpacidadEntry holds raw data from CSV before temperature extraction.
type rawAmpacidadEntry struct {
	Capacidad60 float64
	Capacidad75 float64
	Capacidad90 float64
	Conductor   valueobject.ConductorParams
}

func (r *CSVTablaNOMRepository) loadTablaAmpacidad(filename string, material valueobject.MaterialConductor) ([]rawAmpacidadEntry, error) {
	filePath := filepath.Join(r.basePath, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", filename, err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("%s is empty or missing header", filename)
	}

	// Determine column indices based on material
	materialPrefix := "cu"
	if material == valueobject.MaterialAluminio {
		materialPrefix = "al"
	}

	// Find column indices
	header := records[0]
	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[col] = i
	}

	seccionIdx, ok := colIdx["seccion_mm2"]
	if !ok {
		return nil, fmt.Errorf("%s: missing column seccion_mm2", filename)
	}
	calibreIdx, ok := colIdx["calibre"]
	if !ok {
		return nil, fmt.Errorf("%s: missing column calibre", filename)
	}
	col60 := materialPrefix + "_60c"
	col75 := materialPrefix + "_75c"
	col90 := materialPrefix + "_90c"

	idx60, has60 := colIdx[col60]
	idx75, has75 := colIdx[col75]
	idx90, has90 := colIdx[col90]

	var result []rawAmpacidadEntry
	for i, record := range records[1:] {
		if len(record) < len(header) {
			continue // Skip incomplete rows
		}

		seccion, err := strconv.ParseFloat(record[seccionIdx], 64)
		if err != nil {
			return nil, fmt.Errorf("%s line %d: invalid seccion_mm2: %w", filename, i+2, err)
		}

		entry := rawAmpacidadEntry{
			Conductor: valueobject.ConductorParams{
				Calibre:    record[calibreIdx],
				SeccionMM2: seccion,
			},
		}

		if has60 && record[idx60] != "" {
			entry.Capacidad60, _ = strconv.ParseFloat(record[idx60], 64)
		}
		if has75 && record[idx75] != "" {
			entry.Capacidad75, _ = strconv.ParseFloat(record[idx75], 64)
		}
		if has90 && record[idx90] != "" {
			entry.Capacidad90, _ = strconv.ParseFloat(record[idx90], 64)
		}

		result = append(result, entry)
	}

	return result, nil
}

func extractByTemperature(entries []rawAmpacidadEntry, material valueobject.MaterialConductor, temp valueobject.Temperatura) []valueobject.EntradaTablaConductor {
	var result []valueobject.EntradaTablaConductor

	for _, e := range entries {
		var capacidad float64
		switch temp {
		case valueobject.Temp60:
			capacidad = e.Capacidad60
		case valueobject.Temp75:
			capacidad = e.Capacidad75
		case valueobject.Temp90:
			capacidad = e.Capacidad90
		}

		// Skip entries without capacity for this temperature
		if capacidad <= 0 {
			continue
		}

		// Set material
		params := e.Conductor
		if material == valueobject.MaterialCobre {
			params.Material = valueobject.MaterialCobre
		} else {
			params.Material = valueobject.MaterialAluminio
		}

		result = append(result, valueobject.EntradaTablaConductor{
			Capacidad: capacidad,
			Conductor: params,
		})
	}

	return result
}

func (r *CSVTablaNOMRepository) loadTablaImpedancia() (map[string]impedanciaEntry, error) {
	filePath := filepath.Join(r.basePath, "tabla-9-resistencia-reactancia.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open tabla-9-resistencia-reactancia.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read tabla-9-resistencia-reactancia.csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("tabla-9-resistencia-reactancia.csv is empty or missing header")
	}

	// Find column indices
	header := records[0]
	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[col] = i
	}

	requiredCols := []string{
		"calibre", "reactancia_al", "reactancia_acero",
		"res_cu_pvc", "res_cu_al", "res_cu_acero",
		"res_al_pvc", "res_al_al", "res_al_acero",
	}

	indices := make(map[string]int)
	for _, col := range requiredCols {
		idx, ok := colIdx[col]
		if !ok {
			return nil, fmt.Errorf("tabla-9-resistencia-reactancia.csv: missing column %s", col)
		}
		indices[col] = idx
	}

	result := make(map[string]impedanciaEntry)
	for _, record := range records[1:] {
		if len(record) < len(header) {
			continue // Skip incomplete rows
		}

		calibre := record[indices["calibre"]]

		entry := impedanciaEntry{}

		// Parse all fields
		if v, err := strconv.ParseFloat(record[indices["reactancia_al"]], 64); err == nil {
			entry.ReactanciaAl = v
		}
		if v, err := strconv.ParseFloat(record[indices["reactancia_acero"]], 64); err == nil {
			entry.ReactanciaAcero = v
		}
		if v, err := strconv.ParseFloat(record[indices["res_cu_pvc"]], 64); err == nil {
			entry.ResCuPVC = v
		}
		if v, err := strconv.ParseFloat(record[indices["res_cu_al"]], 64); err == nil {
			entry.ResCuAl = v
		}
		if v, err := strconv.ParseFloat(record[indices["res_cu_acero"]], 64); err == nil {
			entry.ResCuAcero = v
		}
		if v, err := strconv.ParseFloat(record[indices["res_al_pvc"]], 64); err == nil {
			entry.ResAlPVC = v
		}
		if v, err := strconv.ParseFloat(record[indices["res_al_al"]], 64); err == nil {
			entry.ResAlAl = v
		}
		if v, err := strconv.ParseFloat(record[indices["res_al_acero"]], 64); err == nil {
			entry.ResAlAcero = v
		}

		result[calibre] = entry
	}

	return result, nil
}

func (r *CSVTablaNOMRepository) loadTablaConduit() ([]valueobject.EntradaTablaCanalizacion, error) {
	filePath := filepath.Join(r.basePath, "tabla-conduit-dimensiones.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open tabla-conduit-dimensiones.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read tabla-conduit-dimensiones.csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("tabla-conduit-dimensiones.csv is empty or missing header")
	}

	// Find column indices
	header := records[0]
	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[col] = i
	}

	tamanoIdx, ok := colIdx["tamano"]
	if !ok {
		return nil, fmt.Errorf("tabla-conduit-dimensiones.csv: missing column tamano")
	}
	areaIdx, ok := colIdx["area_interior_mm2"]
	if !ok {
		return nil, fmt.Errorf("tabla-conduit-dimensiones.csv: missing column area_interior_mm2")
	}

	var result []valueobject.EntradaTablaCanalizacion
	for i, record := range records[1:] {
		if len(record) < len(header) {
			continue // Skip incomplete rows
		}

		area, err := strconv.ParseFloat(record[areaIdx], 64)
		if err != nil {
			return nil, fmt.Errorf("tabla-conduit-dimensiones.csv line %d: invalid area_interior_mm2: %w", i+2, err)
		}

		result = append(result, valueobject.EntradaTablaCanalizacion{
			Tamano:          record[tamanoIdx],
			AreaInteriorMM2: area,
		})
	}

	return result, nil
}

func (r *CSVTablaNOMRepository) loadEstadosTemperatura() (map[string]int, error) {
	filePath := filepath.Join(r.basePath, "estados_temperatura.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open estados_temperatura.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read estados_temperatura.csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("estados_temperatura.csv is empty or missing header")
	}

	result := make(map[string]int)
	for i, record := range records[1:] {
		if len(record) < 2 {
			continue
		}

		tempF, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("estados_temperatura.csv line %d: invalid temperatura: %w", i+2, err)
		}

		result[record[0]] = int(math.Round(tempF))
	}

	return result, nil
}

func (r *CSVTablaNOMRepository) loadFactoresTemperatura() ([]service.EntradaTablaFactorTemperatura, error) {
	filePath := filepath.Join(r.basePath, "310-15-b-2-a.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open 310-15-b-2-a.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read 310-15-b-2-a.csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("310-15-b-2-a.csv is empty or missing header")
	}

	var result []service.EntradaTablaFactorTemperatura
	for _, record := range records[1:] {
		if len(record) < 4 {
			continue
		}

		f60, _ := strconv.ParseFloat(record[1], 64)
		f75, _ := strconv.ParseFloat(record[2], 64)
		f90, _ := strconv.ParseFloat(record[3], 64)

		result = append(result, service.EntradaTablaFactorTemperatura{
			RangoTempC: record[0],
			Factor60C:  f60,
			Factor75C:  f75,
			Factor90C:  f90,
		})
	}

	return result, nil
}

func (r *CSVTablaNOMRepository) loadFactoresAgrupamiento() ([]service.EntradaTablaFactorAgrupamiento, error) {
	filePath := filepath.Join(r.basePath, "310-15-b-3-a.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open 310-15-b-3-a.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read 310-15-b-3-a.csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("310-15-b-3-a.csv is empty or missing header")
	}

	var result []service.EntradaTablaFactorAgrupamiento
	for _, record := range records[1:] {
		if len(record) < 2 {
			continue
		}

		factor, _ := strconv.ParseFloat(record[1], 64)

		min, max := parseCantidadConductores(record[0])

		result = append(result, service.EntradaTablaFactorAgrupamiento{
			CantidadMin: min,
			CantidadMax: max,
			Factor:      factor,
		})
	}

	return result, nil
}

func (r *CSVTablaNOMRepository) loadTablaDiametros() (map[string]diametroConductorEntry, error) {
	filePath := filepath.Join(r.basePath, "tabla-5-dimensiones-aislamiento.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open tabla-5-dimensiones-aislamiento.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read tabla-5-dimensiones-aislamiento.csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("tabla-5-dimensiones-aislamiento.csv is empty or missing header")
	}

	header := records[0]
	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[col] = i
	}

	calibreIdx, ok := colIdx["calibre"]
	if !ok {
		return nil, fmt.Errorf("tabla-5-dimensiones-aislamiento.csv: missing column calibre")
	}
	diamTWTHWIdx, ok := colIdx["diam_tw_thw"]
	if !ok {
		return nil, fmt.Errorf("tabla-5-dimensiones-aislamiento.csv: missing column diam_tw_thw")
	}
	diamRHH_RHWIdx, ok := colIdx["diam_rhh_rhw"]
	if !ok {
		return nil, fmt.Errorf("tabla-5-dimensiones-aislamiento.csv: missing column diam_rhh_rhw")
	}
	diamXHHWIdx, ok := colIdx["diam_xhhw"]
	if !ok {
		return nil, fmt.Errorf("tabla-5-dimensiones-aislamiento.csv: missing column diam_xhhw")
	}

	result := make(map[string]diametroConductorEntry)
	for _, record := range records[1:] {
		if len(record) < len(header) {
			continue
		}

		entry := diametroConductorEntry{}
		if v, err := strconv.ParseFloat(record[diamTWTHWIdx], 64); err == nil {
			entry.DiamTWTHW = v
		}
		if v, err := strconv.ParseFloat(record[diamRHH_RHWIdx], 64); err == nil {
			entry.DiamRHH_RHW = v
		}
		if v, err := strconv.ParseFloat(record[diamXHHWIdx], 64); err == nil {
			entry.DiamXHHW = v
		}

		result[record[calibreIdx]] = entry
	}

	return result, nil
}

func parseCantidadConductores(s string) (min, max int) {
	// Rango con "+" al final: "41+"
	if _, err := fmt.Sscanf(s, "%d+", &min); err == nil {
		return min, -1
	}
	// Rango con guión: "5-6", "7-9", "10-20", etc.
	if _, err := fmt.Sscanf(s, "%d-%d", &min, &max); err == nil {
		return
	}
	// Entero simple: "1", "2", "3", "4"
	if _, err := fmt.Sscanf(s, "%d", &min); err == nil {
		return min, min
	}
	return 0, 0
}
