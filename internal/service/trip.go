package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"car-tracking/internal/domain/model"
)

type tripGPSRepo interface {
	FindByTripID(ctx context.Context, tripID string) ([]model.GPSPoint, error)
}

type tripTripRepo interface {
	Create(ctx context.Context, trip model.Trip) (string, error)
	GetByID(ctx context.Context, id string) (*model.Trip, error)
	GetActiveByVehicleID(ctx context.Context, vehicleID string) (*model.Trip, error)
	GetStaleTrips(ctx context.Context, olderThan time.Time) ([]model.Trip, error)
	EndTrip(ctx context.Context, tripID string, endTime time.Time) error
	SetAvgSpeed(ctx context.Context, tripID string, avgSpeed float64) error
	List(ctx context.Context, vehicleID string, from, to *time.Time, limit, offset int64) ([]model.Trip, error)
	Count(ctx context.Context, vehicleID string, from, to *time.Time) (int64, error)
}

type TripService struct {
	gpsRepo  tripGPSRepo
	tripRepo tripTripRepo
}

func NewTripService(tripRepo tripTripRepo, gpsRepo tripGPSRepo) *TripService {
	return &TripService{
		tripRepo: tripRepo,
		gpsRepo:  gpsRepo,
	}
}

// StartTrip closes any active trip for the vehicle and creates a new one.
func (s *TripService) StartTrip(ctx context.Context, vehicleID string) (string, error) {
	existing, err := s.tripRepo.GetActiveByVehicleID(ctx, vehicleID)
	if err != nil {
		return "", fmt.Errorf("trip service start: check active: %w", err)
	}
	if existing != nil {
		log.Printf("auto-closing stale trip %s for vehicle %s", existing.ID, vehicleID)
		if err := s.finishTrip(ctx, existing.ID); err != nil {
			return "", fmt.Errorf("trip service start: auto-close: %w", err)
		}
	}

	tripID := uuid.New().String()
	now := time.Now().UTC()

	trip := model.Trip{
		ID:        tripID,
		VehicleID: vehicleID,
		StartTime: now,
		Status:    model.TripStatusActive,
		CreatedAt: now,
	}

	id, err := s.tripRepo.Create(ctx, trip)
	if err != nil {
		return "", fmt.Errorf("trip service start: create: %w", err)
	}

	return id, nil
}

// EndTrip finalizes the active trip for a vehicle.
func (s *TripService) EndTrip(ctx context.Context, vehicleID string) error {
	trip, err := s.tripRepo.GetActiveByVehicleID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("trip service end: find active: %w", err)
	}
	if trip == nil {
		return fmt.Errorf("trip service end: no active trip for vehicle %s", vehicleID)
	}

	return s.finishTrip(ctx, trip.ID)
}

func (s *TripService) finishTrip(ctx context.Context, tripID string) error {
	now := time.Now().UTC()

	if err := s.tripRepo.EndTrip(ctx, tripID, now); err != nil {
		return fmt.Errorf("finish trip: %w", err)
	}

	points, err := s.gpsRepo.FindByTripID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("finish trip: get points: %w", err)
	}

	if len(points) > 0 {
		var totalSpeed float64
		for _, p := range points {
			totalSpeed += float64(p.Speed)
		}
		avgSpeed := totalSpeed / float64(len(points))

		if err := s.tripRepo.SetAvgSpeed(ctx, tripID, avgSpeed); err != nil {
			return fmt.Errorf("finish trip: set avg speed: %w", err)
		}
	}

	return nil
}

// ListTrips returns paginated completed trips with optional date filtering.
func (s *TripService) ListTrips(ctx context.Context, from, to *time.Time, limit, offset int64) ([]model.Trip, int64, error) {
	trips, err := s.tripRepo.List(ctx, "", from, to, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("trip service list: %w", err)
	}

	total, err := s.tripRepo.Count(ctx, "", from, to)
	if err != nil {
		return nil, 0, fmt.Errorf("trip service count: %w", err)
	}

	return trips, total, nil
}

// GetTrip returns a single trip by ID.
func (s *TripService) GetTrip(ctx context.Context, id string) (*model.Trip, error) {
	trip, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("trip service get: %w", err)
	}

	return trip, nil
}

// GetTripPoints returns all GPS points for a trip.
func (s *TripService) GetTripPoints(ctx context.Context, tripID string) ([]model.GPSPoint, error) {
	points, err := s.gpsRepo.FindByTripID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("trip service get points: %w", err)
	}

	return points, nil
}

// CloseStaleTrips finds active trips with no activity and closes them.
func (s *TripService) CloseStaleTrips(ctx context.Context, timeout time.Duration) {
	threshold := time.Now().UTC().Add(-timeout)

	stale, err := s.tripRepo.GetStaleTrips(ctx, threshold)
	if err != nil {
		log.Printf("stale trips check: %v", err)
		return
	}

	for _, trip := range stale {
		log.Printf("auto-closing stale trip %s (started %s)", trip.ID, trip.StartTime.Format(time.RFC3339))

		if err := s.finishTrip(ctx, trip.ID); err != nil {
			log.Printf("auto-close trip %s failed: %v", trip.ID, err)
		}
	}
}

// RunStaleTripsWorker periodically closes stale trips.
func (s *TripService) RunStaleTripsWorker(ctx context.Context, interval, timeout time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("stale trips worker started (check every %s, timeout %s)", interval, timeout)

	for {
		select {
		case <-ctx.Done():
			log.Println("stale trips worker stopped")
			return
		case <-ticker.C:
			s.CloseStaleTrips(ctx, timeout)
		}
	}
}
