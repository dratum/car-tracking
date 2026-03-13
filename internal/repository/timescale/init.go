package timescale

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const enableExtension = `CREATE EXTENSION IF NOT EXISTS timescaledb;`

const createGPSPointsTable = `
CREATE TABLE IF NOT EXISTS gps_points (
    time        TIMESTAMPTZ      NOT NULL,
    trip_id     UUID             NOT NULL,
    lat         DOUBLE PRECISION NOT NULL,
    lon         DOUBLE PRECISION NOT NULL,
    speed       REAL,
    heading     REAL,
    satellites  SMALLINT
);
`

const createHypertable = `
SELECT create_hypertable('gps_points', 'time', if_not_exists => true);
`

const createTripTimeIndex = `
CREATE INDEX IF NOT EXISTS idx_gps_points_trip_id
ON gps_points (trip_id, time DESC);
`

func InitSchema(ctx context.Context, pool *pgxpool.Pool) error {
	statements := []string{
		enableExtension,
		createGPSPointsTable,
		createHypertable,
		createTripTimeIndex,
	}
	for _, stmt := range statements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("timescale init: %w", err)
		}
	}
	return nil
}
