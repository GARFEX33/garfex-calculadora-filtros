// internal/equipos/infrastructure/adapter/driven/postgres/equipo_filtro_repository.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	appdto "github.com/garfex/calculadora-filtros/internal/equipos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/equipos/application/port"
	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresEquipoFiltroRepository implements port.EquipoFiltroRepository using pgx.
type PostgresEquipoFiltroRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresEquipoFiltroRepository creates a new repository with the given pool.
func NewPostgresEquipoFiltroRepository(pool *pgxpool.Pool) *PostgresEquipoFiltroRepository {
	return &PostgresEquipoFiltroRepository{pool: pool}
}

// Compile-time check: PostgresEquipoFiltroRepository must implement port.EquipoFiltroRepository.
var _ port.EquipoFiltroRepository = (*PostgresEquipoFiltroRepository)(nil)

// Crear inserts a new equipo_filtro and returns the created record with DB-generated fields.
func (r *PostgresEquipoFiltroRepository) Crear(ctx context.Context, equipo *entity.EquipoFiltro) (*entity.EquipoFiltro, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO equipos_filtros (clave, tipo, voltaje, "qn/In", itm, bornes, conexion, tipo_voltaje)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, clave, tipo, voltaje, "qn/In", itm, bornes, conexion, tipo_voltaje
	`

	row := r.pool.QueryRow(ctx, query,
		equipo.Clave,
		mapTipoFiltroToDB(equipo.Tipo),
		equipo.Voltaje,
		equipo.Amperaje,
		equipo.ITM,
		equipo.Bornes,
		mapConexionToDB(equipo.Conexion),
		mapTipoVoltajeToDB(equipo.TipoVoltaje),
	)

	created, err := scanEquipoFiltro(row)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("%w: clave ya existe", appdto.ErrClaveYaExiste)
		}
		return nil, fmt.Errorf("insertar equipo: %w", err)
	}

	return created, nil
}

// ObtenerPorID fetches a single equipo by UUID.
func (r *PostgresEquipoFiltroRepository) ObtenerPorID(ctx context.Context, id uuid.UUID) (*entity.EquipoFiltro, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, created_at, clave, tipo, voltaje, "qn/In", itm, bornes, conexion, tipo_voltaje
		FROM equipos_filtros
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	equipo, err := scanEquipoFiltro(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: id %s", appdto.ErrEquipoNoEncontrado, id)
		}
		return nil, fmt.Errorf("obtener equipo por id: %w", err)
	}

	return equipo, nil
}

// buildWhereClause constructs the shared WHERE clause and args for Listar and Contar.
func buildWhereClause(filtros port.FiltrosListado) (string, []any, int) {
	where := " WHERE 1=1"
	args := []any{}
	argIdx := 1

	if filtros.Tipo != nil {
		where += fmt.Sprintf(" AND tipo = $%d", argIdx)
		args = append(args, mapTipoFiltroToDB(*filtros.Tipo))
		argIdx++
	}
	if filtros.Buscar != nil && *filtros.Buscar != "" {
		where += fmt.Sprintf(" AND clave ILIKE $%d", argIdx)
		args = append(args, "%"+*filtros.Buscar+"%")
		argIdx++
	}
	if filtros.Voltaje != nil {
		where += fmt.Sprintf(" AND voltaje = $%d", argIdx)
		args = append(args, *filtros.Voltaje)
		argIdx++
	}
	return where, args, argIdx
}

// Listar returns a paginated page of equipos matching the optional filters.
func (r *PostgresEquipoFiltroRepository) Listar(ctx context.Context, filtros port.FiltrosListado) ([]*entity.EquipoFiltro, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	where, args, argIdx := buildWhereClause(filtros)

	query := `SELECT id, created_at, clave, tipo, voltaje, "qn/In", itm, bornes, conexion, tipo_voltaje FROM equipos_filtros` +
		where +
		" ORDER BY created_at DESC" +
		fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)

	args = append(args, filtros.Limit, filtros.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listar equipos: %w", err)
	}
	defer rows.Close()

	equipos := make([]*entity.EquipoFiltro, 0, filtros.Limit)
	for rows.Next() {
		equipo, err := scanEquipoFiltroFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("escanear equipo: %w", err)
		}
		equipos = append(equipos, equipo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %w", err)
	}

	return equipos, nil
}

// Contar returns the total count of equipos matching the filters (ignores pagination).
func (r *PostgresEquipoFiltroRepository) Contar(ctx context.Context, filtros port.FiltrosListado) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	where, args, _ := buildWhereClause(filtros)
	query := `SELECT COUNT(*) FROM equipos_filtros` + where

	var count int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("contar equipos: %w", err)
	}
	return count, nil
}

// Actualizar updates an existing equipo and returns the updated record.
func (r *PostgresEquipoFiltroRepository) Actualizar(ctx context.Context, equipo *entity.EquipoFiltro) (*entity.EquipoFiltro, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE equipos_filtros
		SET clave = $1, tipo = $2, voltaje = $3, "qn/In" = $4, itm = $5, bornes = $6, conexion = $7, tipo_voltaje = $8
		WHERE id = $9
		RETURNING id, created_at, clave, tipo, voltaje, "qn/In", itm, bornes, conexion, tipo_voltaje
	`

	row := r.pool.QueryRow(ctx, query,
		equipo.Clave,
		mapTipoFiltroToDB(equipo.Tipo),
		equipo.Voltaje,
		equipo.Amperaje,
		equipo.ITM,
		equipo.Bornes,
		mapConexionToDB(equipo.Conexion),
		mapTipoVoltajeToDB(equipo.TipoVoltaje),
		equipo.ID,
	)

	updated, err := scanEquipoFiltro(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: id %s", appdto.ErrEquipoNoEncontrado, equipo.ID)
		}
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("%w: clave ya existe", appdto.ErrClaveYaExiste)
		}
		return nil, fmt.Errorf("actualizar equipo: %w", err)
	}

	return updated, nil
}

// Eliminar deletes an equipo by UUID. Idempotent — no error if not found.
func (r *PostgresEquipoFiltroRepository) Eliminar(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx, `DELETE FROM equipos_filtros WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("eliminar equipo: %w", err)
	}

	return nil
}

// ─── Mappers ────────────────────────────────────────────────────────────────

// mapTipoFiltroToDB converts a domain TipoFiltro to the PostgreSQL enum string.
func mapTipoFiltroToDB(t entity.TipoFiltro) string {
	return string(t)
}

// mapTipoFiltroFromDB converts a PostgreSQL enum string to a domain TipoFiltro.
func mapTipoFiltroFromDB(s string) entity.TipoFiltro {
	t, err := entity.ParseTipoFiltro(s)
	if err != nil {
		// Fallback — should not happen with a properly constrained DB
		return entity.TipoFiltroA
	}
	return t
}

// mapConexionToDB converts a nullable *Conexion to a nullable *string for PostgreSQL.
func mapConexionToDB(c *entity.Conexion) *string {
	if c == nil {
		return nil
	}
	s := string(*c)
	return &s
}

// mapConexionFromDB converts a nullable *string from PostgreSQL to a domain *Conexion.
func mapConexionFromDB(s *string) *entity.Conexion {
	if s == nil {
		return nil
	}
	c, err := entity.ParseConexion(*s)
	if err != nil {
		// Fallback — should not happen with a properly constrained DB
		return nil
	}
	return &c
}

// mapTipoVoltajeToDB converts a nullable *TipoVoltaje to a nullable *string for PostgreSQL.
func mapTipoVoltajeToDB(tv *entity.TipoVoltaje) *string {
	if tv == nil {
		return nil
	}
	s := string(*tv)
	return &s
}

// mapTipoVoltajeFromDB converts a nullable *string from PostgreSQL to a domain *TipoVoltaje.
func mapTipoVoltajeFromDB(s *string) *entity.TipoVoltaje {
	if s == nil {
		return nil
	}
	tv, err := entity.ParseTipoVoltaje(*s)
	if err != nil {
		// Fallback — should not happen with a properly constrained DB
		return nil
	}
	return &tv
}

// scanEquipoFiltro scans a single pgx.Row into a domain entity.
func scanEquipoFiltro(row pgx.Row) (*entity.EquipoFiltro, error) {
	var (
		id          uuid.UUID
		createdAt   time.Time
		clave       *string
		tipo        string
		voltaje     int
		amperaje    int
		itm         int
		bornes      *int
		conexion    *string
		tipoVoltaje *string
	)

	err := row.Scan(&id, &createdAt, &clave, &tipo, &voltaje, &amperaje, &itm, &bornes, &conexion, &tipoVoltaje)
	if err != nil {
		return nil, err
	}

	return &entity.EquipoFiltro{
		ID:          id,
		CreatedAt:   createdAt,
		Clave:       clave,
		Tipo:        mapTipoFiltroFromDB(tipo),
		Voltaje:     voltaje,
		Amperaje:    amperaje,
		ITM:         itm,
		Bornes:      bornes,
		Conexion:    mapConexionFromDB(conexion),
		TipoVoltaje: mapTipoVoltajeFromDB(tipoVoltaje),
	}, nil
}

// scanEquipoFiltroFromRows scans a pgx.Rows (multi-row scan) into a domain entity.
func scanEquipoFiltroFromRows(rows pgx.Rows) (*entity.EquipoFiltro, error) {
	var (
		id          uuid.UUID
		createdAt   time.Time
		clave       *string
		tipo        string
		voltaje     int
		amperaje    int
		itm         int
		bornes      *int
		conexion    *string
		tipoVoltaje *string
	)

	err := rows.Scan(&id, &createdAt, &clave, &tipo, &voltaje, &amperaje, &itm, &bornes, &conexion, &tipoVoltaje)
	if err != nil {
		return nil, err
	}

	return &entity.EquipoFiltro{
		ID:          id,
		CreatedAt:   createdAt,
		Clave:       clave,
		Tipo:        mapTipoFiltroFromDB(tipo),
		Voltaje:     voltaje,
		Amperaje:    amperaje,
		ITM:         itm,
		Bornes:      bornes,
		Conexion:    mapConexionFromDB(conexion),
		TipoVoltaje: mapTipoVoltajeFromDB(tipoVoltaje),
	}, nil
}

// isUniqueViolation detects PostgreSQL unique constraint violations (SQLSTATE 23505).
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
