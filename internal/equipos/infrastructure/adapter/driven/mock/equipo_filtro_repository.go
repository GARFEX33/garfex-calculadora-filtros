// internal/equipos/infrastructure/adapter/driven/mock/equipo_filtro_repository.go
package mock

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/garfex/calculadora-filtros/internal/equipos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/equipos/application/port"
	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
	"github.com/google/uuid"
)

// MockEquipoFiltroRepository implements port.EquipoFiltroRepository using in-memory storage.
type MockEquipoFiltroRepository struct {
	mu      sync.RWMutex
	equipos map[uuid.UUID]*entity.EquipoFiltro
}

// NewMockEquipoFiltroRepository creates a new in-memory repository preloaded with
// the exact same records that exist in the production PostgreSQL database.
// IDs, claves, voltajes, amperajes, ITM, bornes, conexion y tipo_voltaje son
// idénticos a los de la DB para garantizar paridad total mock ↔ producción.
func NewMockEquipoFiltroRepository() *MockEquipoFiltroRepository {
	repo := &MockEquipoFiltroRepository{
		equipos: make(map[uuid.UUID]*entity.EquipoFiltro),
	}

	// Datos idénticos a producción (exportados de PostgreSQL 2026-02-09)
	seed, _ := time.Parse(time.RFC3339, "2026-02-09T23:30:25Z")

	delta := conexionPtr(entity.ConexionDelta)
	ff := tipoVoltajePtr(entity.TipoVoltajeFaseFase)

	equipos := []*entity.EquipoFiltro{
		// ── Filtros Activos (tipo A) — 480 V ──────────────────────────────────
		{
			ID:          uuid.MustParse("15cab32e-44ea-409a-9353-a8ca8e4b70ac"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE48D400MMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     480,
			Amperaje:    400,
			ITM:         600,
			Bornes:      intPtr(2),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("26be9a7d-fb4d-41e9-a16f-039410b6fd49"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE48Y050CMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     480,
			Amperaje:    50,
			ITM:         70,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("2ac53560-6142-4ec2-912b-d528a8d0dbc8"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE48Y300MMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     480,
			Amperaje:    300,
			ITM:         400,
			Bornes:      intPtr(2),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("87e6bec1-d8b2-452b-b00b-81aa89007de4"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE48Y200MMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     480,
			Amperaje:    200,
			ITM:         250,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("9fb8d082-d85a-41ee-9722-2e4f476a5f8c"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE48Y150MMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     480,
			Amperaje:    150,
			ITM:         200,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("d7a3ecd5-2792-484c-bb55-b3e0e68d32e3"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE48Y100CMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     480,
			Amperaje:    100,
			ITM:         125,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("d7cf3eeb-80ec-4637-8d3f-744e592a4b83"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE48Y075CMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     480,
			Amperaje:    75,
			ITM:         100,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("f0c86e4b-4810-4e62-86fd-4ffdfcfcb73b"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE48Y500MMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     480,
			Amperaje:    500,
			ITM:         800,
			Bornes:      intPtr(3),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		// ── Filtros Activos (tipo A) — 240 V ──────────────────────────────────
		{
			ID:          uuid.MustParse("26e64670-08bb-4350-a86f-e20549ccecda"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE24Y075CMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     240,
			Amperaje:    75,
			ITM:         100,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("9b94de98-2491-449f-9c61-3a0c3f064e21"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE24Y150MMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     240,
			Amperaje:    150,
			ITM:         200,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("e4409515-d719-4e7b-848c-477f0808d147"),
			CreatedAt:   seed,
			Clave:       strPtr("ACTISINE24Y050CMB"),
			Tipo:        entity.TipoFiltroA,
			Voltaje:     240,
			Amperaje:    50,
			ITM:         70,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		// ── Filtros de Rechazo (tipo KVAR) — 480 V ────────────────────────────
		{
			ID:          uuid.MustParse("317a688e-834f-4db3-b21c-8af6cd553023"),
			CreatedAt:   seed,
			Clave:       strPtr("FRST-40200-04IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    200,
			ITM:         400,
			Bornes:      intPtr(2),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("4ec70d76-47f3-4253-84a9-bba17976c8b8"),
			CreatedAt:   seed,
			Clave:       strPtr("BAMTS8-4200-05I"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    200,
			ITM:         350,
			Bornes:      intPtr(2),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("58c4090d-785c-4cf9-92c4-da766de1498e"),
			CreatedAt:   seed,
			Clave:       strPtr("FRS-40500-10IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    500,
			ITM:         800,
			Bornes:      intPtr(3),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("5a3737db-2b62-4b99-acc3-0b0cd16f3515"),
			CreatedAt:   seed,
			Clave:       strPtr("FRS-40250-06IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    250,
			ITM:         500,
			Bornes:      intPtr(2),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("5cfa6175-3f39-4fe6-ba0c-f59ecd7b5fc2"),
			CreatedAt:   seed,
			Clave:       strPtr("FRS-40100-04IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    100,
			ITM:         175,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("a1566f65-cc2d-470c-89d9-df45c09fcb5c"),
			CreatedAt:   seed,
			Clave:       strPtr("FRST-40125-03IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    125,
			ITM:         250,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("aa8cad63-789e-4123-b4ad-c9ce736954e6"),
			CreatedAt:   seed,
			Clave:       strPtr("FRS-40275-06IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    275,
			ITM:         500,
			Bornes:      intPtr(2),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("ace98516-60ea-4a66-86cd-cce7f813c894"),
			CreatedAt:   seed,
			Clave:       strPtr("FRST-40400-05IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    400,
			ITM:         700,
			Bornes:      intPtr(3),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("b9f34b2f-e989-4732-b246-fe915a267fb3"),
			CreatedAt:   seed,
			Clave:       strPtr("FRST-40150-04IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    150,
			ITM:         250,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("c4da3652-9354-4634-94a2-521ca43bacb3"),
			CreatedAt:   seed,
			Clave:       strPtr("FRS-40700-14IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    700,
			ITM:         1200,
			Bornes:      intPtr(4),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("d2386b85-4991-4a7a-8ae8-0d4488eaa4ca"),
			CreatedAt:   seed,
			Clave:       strPtr("FRS-40350-08IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    350,
			ITM:         600,
			Bornes:      intPtr(2),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("dee0d5c4-4485-452d-82e9-32722c101ce6"),
			CreatedAt:   seed,
			Clave:       strPtr("BAMTS8-4100-04I"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    100,
			ITM:         175,
			Bornes:      intPtr(1),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		{
			ID:          uuid.MustParse("f984a2a8-a5ae-434a-b7e3-3bf54f5ecb4d"),
			CreatedAt:   seed,
			Clave:       strPtr("FRS-40300-07IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     480,
			Amperaje:    300,
			ITM:         500,
			Bornes:      intPtr(2),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
		// ── Filtros de Rechazo (tipo KVAR) — 240 V ────────────────────────────
		{
			ID:          uuid.MustParse("ec4ce941-aefd-40dc-b988-a7aab53a6e2c"),
			CreatedAt:   seed,
			Clave:       strPtr("FRS-20200-08IP7"),
			Tipo:        entity.TipoFiltroKVAR,
			Voltaje:     240,
			Amperaje:    200,
			ITM:         700,
			Bornes:      intPtr(3),
			Conexion:    delta,
			TipoVoltaje: ff,
		},
	}

	for _, e := range equipos {
		repo.equipos[e.ID] = e
	}

	return repo
}

// Compile-time check: MockEquipoFiltroRepository must implement port.EquipoFiltroRepository.
var _ port.EquipoFiltroRepository = (*MockEquipoFiltroRepository)(nil)

// Crear persists a new equipo and returns it with the generated ID and CreatedAt.
func (r *MockEquipoFiltroRepository) Crear(ctx context.Context, equipo *entity.EquipoFiltro) (*entity.EquipoFiltro, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate clave
	for _, e := range r.equipos {
		if e.Clave != nil && equipo.Clave != nil && *e.Clave == *equipo.Clave {
			return nil, fmt.Errorf("%w: clave ya existe", dto.ErrClaveYaExiste)
		}
	}

	// Generate ID and CreatedAt
	equipo.ID = uuid.New()
	equipo.CreatedAt = time.Now().Truncate(time.Second)

	r.equipos[equipo.ID] = equipo
	return equipo, nil
}

// ObtenerPorID finds an equipo by its UUID. Returns ErrEquipoNoEncontrado if missing.
func (r *MockEquipoFiltroRepository) ObtenerPorID(ctx context.Context, id uuid.UUID) (*entity.EquipoFiltro, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	equipo, exists := r.equipos[id]
	if !exists {
		return nil, fmt.Errorf("%w: id %s", dto.ErrEquipoNoEncontrado, id)
	}
	return equipo, nil
}

// Listar returns a paginated page of equipos matching the optional filters.
func (r *MockEquipoFiltroRepository) Listar(ctx context.Context, filtros port.FiltrosListado) ([]*entity.EquipoFiltro, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.EquipoFiltro
	for _, e := range r.equipos {
		if !matchesFilters(e, filtros) {
			continue
		}
		result = append(result, e)
	}

	// Apply pagination (already filtered)
	start := filtros.Offset
	end := start + filtros.Limit
	if start > len(result) {
		return []*entity.EquipoFiltro{}, nil
	}
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}

// Contar returns the total count of equipos matching the filters (ignoring pagination).
func (r *MockEquipoFiltroRepository) Contar(ctx context.Context, filtros port.FiltrosListado) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, e := range r.equipos {
		if matchesFilters(e, filtros) {
			count++
		}
	}
	return count, nil
}

// Actualizar updates an existing equipo and returns the updated record.
// Returns ErrEquipoNoEncontrado if the ID does not exist.
func (r *MockEquipoFiltroRepository) Actualizar(ctx context.Context, equipo *entity.EquipoFiltro) (*entity.EquipoFiltro, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.equipos[equipo.ID]
	if !exists {
		return nil, fmt.Errorf("%w: id %s", dto.ErrEquipoNoEncontrado, equipo.ID)
	}

	// Check for duplicate clave (excluding current ID)
	if equipo.Clave != nil {
		for id, e := range r.equipos {
			if id != equipo.ID && e.Clave != nil && *e.Clave == *equipo.Clave {
				return nil, fmt.Errorf("%w: clave ya existe", dto.ErrClaveYaExiste)
			}
		}
	}

	// Preserve CreatedAt from existing record
	equipo.CreatedAt = existing.CreatedAt
	r.equipos[equipo.ID] = equipo
	return equipo, nil
}

// Eliminar deletes an equipo by UUID. Idempotent — no error if not found.
func (r *MockEquipoFiltroRepository) Eliminar(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.equipos, id)
	return nil
}

// matchesFilters checks if an equipo matches the given filters.
func matchesFilters(e *entity.EquipoFiltro, filtros port.FiltrosListado) bool {
	// Filter by Tipo
	if filtros.Tipo != nil && e.Tipo != *filtros.Tipo {
		return false
	}

	// Filter by Voltaje
	if filtros.Voltaje != nil && e.Voltaje != *filtros.Voltaje {
		return false
	}

	// Filter by Buscar (case-insensitive ILIKE on clave)
	if filtros.Buscar != nil && *filtros.Buscar != "" {
		if e.Clave == nil {
			return false
		}
		if !strings.Contains(strings.ToLower(*e.Clave), strings.ToLower(*filtros.Buscar)) {
			return false
		}
	}

	return true
}

// Helper functions for pointer creation

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func conexionPtr(c entity.Conexion) *entity.Conexion {
	return &c
}

func tipoVoltajePtr(tv entity.TipoVoltaje) *entity.TipoVoltaje {
	return &tv
}
