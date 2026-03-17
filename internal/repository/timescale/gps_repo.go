package timescale

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"car-tracking/internal/domain/model"
)

type GPSRepo struct {
	pool *pgxpool.Pool
}

func NewGPSRepo(pool *pgxpool.Pool) *GPSRepo {
	return &GPSRepo{pool: pool}
}

func (r *GPSRepo) Insert(ctx context.Context, p model.GPSPoint) error {
	const query = `
		insert into gps_points (
							    time
							  , trip_id
							  , lat
							  , lon
							  , speed
							  , heading
							  , satellites
							  )
						values (
							    @time
							  , @trip_id
							  , @lat
							  , @lon
							  , @speed
							  , @heading
							  , @satellites
							  )
	`

	_, err := r.pool.Exec(ctx, query, pgx.NamedArgs{
		"time":       p.Time,
		"trip_id":    p.TripID,
		"lat":        p.Lat,
		"lon":        p.Lon,
		"speed":      p.Speed,
		"heading":    p.Heading,
		"satellites": p.Satellites,
	})
	if err != nil {
		return fmt.Errorf("gps_repo insert: %w", err)
	}

	return nil
}

func (r *GPSRepo) FindByTripID(ctx context.Context, tripID string) ([]model.GPSPoint, error) {
	const query = `
		select
			  time
			, trip_id
			, lat
			, lon
			, speed
			, heading
			, satellites
		from gps_points
		where trip_id = @trip_id
		order by time asc`

	rows, err := r.pool.Query(ctx, query, pgx.NamedArgs{"trip_id": tripID})
	if err != nil {
		return nil, fmt.Errorf("gps_repo find by trip: %w", err)
	}
	defer rows.Close()

	var points []model.GPSPoint
	for rows.Next() {
		var p model.GPSPoint
		if err := rows.Scan(&p.Time, &p.TripID, &p.Lat, &p.Lon, &p.Speed, &p.Heading, &p.Satellites); err != nil {
			return nil, fmt.Errorf("gps_repo scan: %w", err)
		}
		points = append(points, p)
	}

	return points, rows.Err()
}

// LastByTripID returns the most recent GPS point for a trip, or nil if none exist.
func (r *GPSRepo) LastByTripID(ctx context.Context, tripID string) (*model.GPSPoint, error) {
	const query = `
		select
			  time
			, trip_id
			, lat
			, lon
			, speed
			, heading
			, satellites
		from gps_points
		where trip_id = @trip_id
		order by time desc
		limit 1`

	var p model.GPSPoint
	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{"trip_id": tripID}).Scan(
		&p.Time, &p.TripID, &p.Lat, &p.Lon, &p.Speed, &p.Heading, &p.Satellites,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gps_repo last by trip: %w", err)
	}

	return &p, nil
}
