package service

import (
	"context"
	"time"

	"auto-tracking/internal/domain/model"
)

type gpsRepository interface {
	Insert(ctx context.Context, p model.GPSPoint) error
	FindByTripID(ctx context.Context, tripID string) ([]model.GPSPoint, error)
	LastByTripID(ctx context.Context, tripID string) (*model.GPSPoint, error)
}

type tripRepository interface {
	Create(ctx context.Context, trip model.Trip) (string, error)
	GetByID(ctx context.Context, id string) (*model.Trip, error)
	GetActiveByVehicleID(ctx context.Context, vehicleID string) (*model.Trip, error)
	EndTrip(ctx context.Context, tripID string, endTime time.Time) error
	UpdateStats(ctx context.Context, tripID string, distKM, speed float64) error
	SetAvgSpeed(ctx context.Context, tripID string, avgSpeed float64) error
	List(ctx context.Context, vehicleID string, from, to *time.Time, limit, offset int64) ([]model.Trip, error)
	Count(ctx context.Context, vehicleID string, from, to *time.Time) (int64, error)
	AggregateStats(ctx context.Context, from, to time.Time) (*model.TripStats, error)
}
