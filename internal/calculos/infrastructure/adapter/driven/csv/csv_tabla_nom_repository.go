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
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// Local types to avoid importing domain/service (infrastructure should not contain business logic)
type factorTemperaturaEntry struct {
	rangoTempC string
	factor60C  float64
	factor75C  float64
	factor90C  float64
}

type factorAgrupamientoEntry struct {
	cantidadMin int
	cantidadMax int
	factor      float64
}

// impedanciaEntry holds all impedance values for a given calibre from Tabla 9.
type impedanciaEntry struct {
	SeccionMM2      float64
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
	AreaTWTHW   float64
}

// conductorDesnudoEntry holds values for bare conductors from Tabla 8.
type conductorDesnudoEntry struct {
	SeccionMM2          float64
	AreaConductorTierra float64
	DiametroMM          float64
	NumeroHilos         int
}

// tuboOcupacionEntry holds occupation table entries for conduits (40% fill).
type tuboOcupacionEntry struct {
	Tamano             string
	AreaOcupacionMM2   float64
	AreaInteriorMM2    float64
	DesignacionMetrica string
}

// CSVTablaNOMRepository reads NOM tables from CSV files with in-memory caching.
type CSVTablaNOMRepository struct {
	basePath               string
	tablaTierra            []valueobject.EntradaTablaTierra
	tablasAmpacidad        map[entity.TipoCanalizacion]map[valueobject.MaterialConductor]map[valueobject.Temperatura][]valueobject.EntradaTablaConductor
	tablaImpedancia        map[string]impedanciaEntry // key: calibre
	tablaConduit           []valueobject.EntradaTablaCanalizacion
	tablasCharola          map[entity.TipoCanalizacion][]valueobject.EntradaTablaCanalizacion
	estadosTemperatura     map[string]int
	factoresTemperatura    []factorTemperaturaEntry
	factoresAgrupamiento   []factorAgrupamientoEntry
	tablaDiametros         map[string]diametroConductorEntry
	tablaConductorDesnudo  map[string]conductorDesnudoEntry // Tabla 8 - conductores desnudos
	tablasOcupacionTuberia map[entity.TipoCanalizacion][]valueobject.EntradaTablaOcupacion
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
	tablaEspaciado, err := repo.crearTablaCharolaEspaciado()
	if err != nil {
		return nil, fmt.Errorf("failed to load charola espaciado table: %w", err)
	}
	tablaTriangular, err := repo.crearTablaCharolaTriangular()
	if err != nil {
		return nil, fmt.Errorf("failed to load charola triangular table: %w", err)
	}
	repo.tablasCharola[entity.TipoCanalizacionCharolaCableEspaciado] = tablaEspaciado
	repo.tablasCharola[entity.TipoCanalizacionCharolaCableTriangular] = tablaTriangular

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

	// Load tabla conductors desenudos (tabla-8-conductor-desnudo.csv) - para conductor de tierra
	tablaDesnudo, err := repo.loadTablaConductorDesnudo()
	if err != nil {
		return nil, fmt.Errorf("failed to load tabla conductor desnudo: %w", err)
	}
	repo.tablaConductorDesnudo = tablaDesnudo

	// Load conduit occupation tables (40% fill)
	tablasOcupacion, err := repo.loadTablasOcupacionTuberia()
	if err != nil {
		return nil, fmt.Errorf("failed to load tablas ocupacion tuberia: %w", err)
	}
	repo.tablasOcupacionTuberia = tablasOcupacion

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
// Simple table lookup - no business logic here.
func (r *CSVTablaNOMRepository) ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error) {
	if tempAmbiente < -10 {
		return 0, fmt.Errorf("temperatura ambiente inválida: %d°C", tempAmbiente)
	}

	for _, entrada := range r.factoresTemperatura {
		if rangoContiene(entrada.rangoTempC, tempAmbiente) {
			switch tempConductor {
			case valueobject.Temp60:
				return entrada.factor60C, nil
			case valueobject.Temp75:
				return entrada.factor75C, nil
			case valueobject.Temp90:
				return entrada.factor90C, nil
			default:
				return 0, fmt.Errorf("temperatura de conductor no soportada: %v", tempConductor)
			}
		}
	}
	return 0, fmt.Errorf("no se encontró factor para temperatura ambiente %d°C", tempAmbiente)
}

// ObtenerFactorAgrupamiento returns the grouping factor based on the number of conductors.
// Simple table lookup - no business logic here.
func (r *CSVTablaNOMRepository) ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error) {
	if cantidadConductores <= 0 {
		return 0, fmt.Errorf("cantidad de conductores debe ser mayor que cero: %d", cantidadConductores)
	}

	for _, entrada := range r.factoresAgrupamiento {
		if entrada.cantidadMax == -1 {
			if cantidadConductores >= entrada.cantidadMin {
				return entrada.factor, nil
			}
		} else {
			if cantidadConductores >= entrada.cantidadMin && cantidadConductores <= entrada.cantidadMax {
				return entrada.factor, nil
			}
		}
	}
	// Default fallback per NOM
	return 0.30, nil
}

// rangoContiene checks if a temperature range contains the given temperature.
func rangoContiene(rango string, temp int) bool {
	var min, max int
	if _, err := fmt.Sscanf(rango, "%d-%d", &min, &max); err == nil {
		return temp >= min && temp <= max
	}
	if _, err := fmt.Sscanf(rango, "%d+", &min); err == nil {
		return temp >= min
	}
	return false
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

// ObtenerAreaConductor returns the area with insulation (area_tw_thw) for a given calibre.
func (r *CSVTablaNOMRepository) ObtenerAreaConductor(ctx context.Context, calibre string) (float64, error) {
	entry, ok := r.tablaDiametros[calibre]
	if !ok {
		return 0, fmt.Errorf("calibre no encontrado en tabla de áreas: %s", calibre)
	}

	if entry.AreaTWTHW <= 0 {
		return 0, fmt.Errorf("área no disponible para calibre: %s", calibre)
	}

	return entry.AreaTWTHW, nil
}

// ObtenerAreaConductorDesnudo returns the area for bare conductor (Tabla 8) - used for ground conductors.
func (r *CSVTablaNOMRepository) ObtenerAreaConductorDesnudo(ctx context.Context, calibre string) (float64, error) {
	entry, ok := r.tablaConductorDesnudo[calibre]
	if !ok {
		return 0, fmt.Errorf("calibre no encontrado en tabla de conductor desnudo: %s", calibre)
	}

	if entry.AreaConductorTierra <= 0 {
		return 0, fmt.Errorf("área no disponible para calibre desnudo: %s", calibre)
	}

	return entry.AreaConductorTierra, nil
}

// ObtenerTablaOcupacionTuberia returns the conduit occupancy table for 40% fill.
func (r *CSVTablaNOMRepository) ObtenerTablaOcupacionTuberia(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]valueobject.EntradaTablaOcupacion, error) {
	tabla, ok := r.tablasOcupacionTuberia[canalizacion]
	if !ok {
		return nil, fmt.Errorf("tabla de ocupación no disponible para tipo de canalización: %s", canalizacion)
	}

	return tabla, nil
}

// ObtenerTablaCharola returns the complete charola sizing table for the given type.
func (r *CSVTablaNOMRepository) ObtenerTablaCharola(ctx context.Context, tipo entity.TipoCanalizacion) ([]valueobject.EntradaTablaCanalizacion, error) {
	switch tipo {
	case entity.TipoCanalizacionCharolaCableEspaciado:
		tabla, ok := r.tablasCharola[entity.TipoCanalizacionCharolaCableEspaciado]
		if !ok {
			return nil, fmt.Errorf("tabla de charola espaciado no cargada")
		}
		return tabla, nil
	case entity.TipoCanalizacionCharolaCableTriangular:
		tabla, ok := r.tablasCharola[entity.TipoCanalizacionCharolaCableTriangular]
		if !ok {
			return nil, fmt.Errorf("tabla de charola triangular no cargada")
		}
		return tabla, nil
	default:
		return nil, fmt.Errorf("tipo de canalización no válido para charola: %s", tipo)
	}
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
// Lee del archivo CSV: charola_dimensiones.csv
func (r *CSVTablaNOMRepository) crearTablaCharolaEspaciado() ([]valueobject.EntradaTablaCanalizacion, error) {
	// Usar la tabla del archivo CSV - el ancho ya está en mm
	return r.loadTablaCharolaDimensiones()
}

// crearTablaCharolaTriangular crea la tabla de dimensiones para charola cable triangular.
// Lee del archivo CSV: charola_dimensiones.csv
func (r *CSVTablaNOMRepository) crearTablaCharolaTriangular() ([]valueobject.EntradaTablaCanalizacion, error) {
	// Usar la tabla del archivo CSV - el ancho ya está en mm
	return r.loadTablaCharolaDimensiones()
}

// loadTablaCharolaDimensiones carga las dimensiones de charolas desde el archivo CSV.
func (r *CSVTablaNOMRepository) loadTablaCharolaDimensiones() ([]valueobject.EntradaTablaCanalizacion, error) {
	filePath := filepath.Join(r.basePath, "charola_dimensiones.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open charola_dimensiones.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read charola_dimensiones.csv: %w", err)
	}

	var result []valueobject.EntradaTablaCanalizacion
	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 2 {
			continue
		}
		tamanoPulgadas := record[0]
		anchoMM, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("charola_dimensiones.csv line %d: invalid ancho_mm: %w", i+1, err)
		}
		result = append(result, valueobject.EntradaTablaCanalizacion{
			Tamano:          tamanoPulgadas,
			AreaInteriorMM2: anchoMM, // En este CSV, el valor es el ancho directo en mm
		})
	}

	return result, nil
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
		"calibre", "seccion_mm2", "reactancia_al", "reactancia_acero",
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
		if v, err := strconv.ParseFloat(record[indices["seccion_mm2"]], 64); err == nil {
			entry.SeccionMM2 = v
		}
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

func (r *CSVTablaNOMRepository) loadFactoresTemperatura() ([]factorTemperaturaEntry, error) {
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

	var result []factorTemperaturaEntry
	for _, record := range records[1:] {
		if len(record) < 4 {
			continue
		}

		f60, _ := strconv.ParseFloat(record[1], 64)
		f75, _ := strconv.ParseFloat(record[2], 64)
		f90, _ := strconv.ParseFloat(record[3], 64)

		result = append(result, factorTemperaturaEntry{
			rangoTempC: record[0],
			factor60C:  f60,
			factor75C:  f75,
			factor90C:  f90,
		})
	}

	return result, nil
}

func (r *CSVTablaNOMRepository) loadFactoresAgrupamiento() ([]factorAgrupamientoEntry, error) {
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

	var result []factorAgrupamientoEntry
	for _, record := range records[1:] {
		if len(record) < 2 {
			continue
		}

		factor, _ := strconv.ParseFloat(record[1], 64)

		min, max := parseCantidadConductores(record[0])

		result = append(result, factorAgrupamientoEntry{
			cantidadMin: min,
			cantidadMax: max,
			factor:      factor,
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

		// Parse area_tw_thw column
		areaTWTHWIdx, ok := colIdx["area_tw_thw"]
		if ok && areaTWTHWIdx < len(record) {
			if v, err := strconv.ParseFloat(record[areaTWTHWIdx], 64); err == nil {
				entry.AreaTWTHW = v
			}
		}

		result[record[calibreIdx]] = entry
	}

	return result, nil
}

// loadTablaConductorDesnudo loads bare conductor table (Tabla 8) for ground conductor area.
func (r *CSVTablaNOMRepository) loadTablaConductorDesnudo() (map[string]conductorDesnudoEntry, error) {
	filePath := filepath.Join(r.basePath, "tabla-8-conductor-desnudo.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open tabla-8-conductor-desnudo.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read tabla-8-conductor-desnudo.csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("tabla-8-conductor-desnudo.csv is empty or missing header")
	}

	header := records[0]
	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[col] = i
	}

	calibreIdx, ok := colIdx["calibre"]
	if !ok {
		return nil, fmt.Errorf("tabla-8-conductor-desnudo.csv: missing column calibre")
	}
	areaTierraIdx, ok := colIdx["area_conductor_tierra"]
	if !ok {
		return nil, fmt.Errorf("tabla-8-conductor-desnudo.csv: missing column area_conductor_tierra")
	}

	result := make(map[string]conductorDesnudoEntry)
	for _, record := range records[1:] {
		if len(record) < len(header) {
			continue
		}

		entry := conductorDesnudoEntry{}
		if v, err := strconv.ParseFloat(record[areaTierraIdx], 64); err == nil {
			entry.AreaConductorTierra = v
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

// loadTablasOcupacionTuberia loads the conduit occupation tables for 40% fill.
func (r *CSVTablaNOMRepository) loadTablasOcupacionTuberia() (map[entity.TipoCanalizacion][]valueobject.EntradaTablaOcupacion, error) {
	result := make(map[entity.TipoCanalizacion][]valueobject.EntradaTablaOcupacion)

	// Define mapping from canalizacion to CSV file
	files := map[entity.TipoCanalizacion]string{
		entity.TipoCanalizacionTuberiaPVC:     "tubo-ocupacion-pvc-40.csv",
		entity.TipoCanalizacionTuberiaAceroPG: "tubo-ocupacion-acero-pg-40.csv",
		entity.TipoCanalizacionTuberiaAceroPD: "tubo-ocupacion-acero-pd-40.csv",
	}

	for canalizacion, filename := range files {
		tabla, err := r.loadTablaOcupacionTuberia(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to load %s: %w", filename, err)
		}
		result[canalizacion] = tabla
	}

	return result, nil
}

// loadTablaOcupacionTuberia loads a single conduit occupation table.
func (r *CSVTablaNOMRepository) loadTablaOcupacionTuberia(filename string) ([]valueobject.EntradaTablaOcupacion, error) {
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

	// Find column indices
	header := records[0]
	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[col] = i
	}

	tamanoIdx, ok := colIdx["tamano"]
	if !ok {
		return nil, fmt.Errorf("%s: missing column tamano", filename)
	}
	areaOcupIdx, ok := colIdx["area_ocupacion_mm2"]
	if !ok {
		return nil, fmt.Errorf("%s: missing column area_ocupacion_mm2", filename)
	}
	designIdx, ok := colIdx["designacion_metrica"]
	if !ok {
		return nil, fmt.Errorf("%s: missing column designacion_metrica", filename)
	}

	var result []valueobject.EntradaTablaOcupacion
	for _, record := range records[1:] {
		if len(record) < len(header) {
			continue
		}

		areaOcup, err := strconv.ParseFloat(record[areaOcupIdx], 64)
		if err != nil {
			return nil, fmt.Errorf("%s line: invalid area_ocupacion_mm2: %w", filename, err)
		}

		// Area interior is calculated from area_ocupacion / 0.40 (since area_ocupacion is 40% fill)
		areaInterior := areaOcup / 0.40

		result = append(result, valueobject.EntradaTablaOcupacion{
			Tamano:             record[tamanoIdx],
			AreaOcupacionMM2:   areaOcup,
			AreaInteriorMM2:    areaInterior,
			DesignacionMetrica: record[designIdx],
		})
	}

	return result, nil
}

// ObtenerSeccionConductor returns the cross-sectional area in mm² for a given calibre from Tabla 9.
func (r *CSVTablaNOMRepository) ObtenerSeccionConductor(ctx context.Context, calibre string) (float64, error) {
	entry, ok := r.tablaImpedancia[calibre]
	if !ok {
		return 0, fmt.Errorf("calibre no encontrado en tabla de impedancia: %s", calibre)
	}

	if entry.SeccionMM2 <= 0 {
		return 0, fmt.Errorf("sección no disponible para calibre: %s", calibre)
	}

	return entry.SeccionMM2, nil
}
