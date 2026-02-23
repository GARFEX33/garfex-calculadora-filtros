// internal/equipos/infrastructure/adapter/driven/postgres/pool.go
package postgres

import (
	"context"
	"fmt"
	"time"

	sharedpostgres "github.com/garfex/calculadora-filtros/internal/shared/infrastructure/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a pgxpool.Pool from the given DBConfig.
// Uses sensible defaults: max 10 connections, 5s connect timeout.
func NewPool(cfg sharedpostgres.DBConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parsear config de pool: %w", err)
	}

	poolCfg.MaxConns = 10
	poolCfg.MinConns = 1
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("crear pool de conexiones: %w", err)
	}

	// Verify connectivity
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("no se pudo conectar a PostgreSQL (%s:%s): %w", cfg.Host, cfg.Port, err)
	}

	return pool, nil
}
