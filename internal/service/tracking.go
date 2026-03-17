package service

import (
	"context"
	"fmt"

	"auto-tracking/internal/domain/geo"
	"auto-tracking/internal/domain/model"
)

type trackingGPSRepo interface {
	Insert(ctx context.Context, p model.GPSPoint) error
	LastByTripID(ctx context.Context, tripID string) (*model.GPSPoint, error)
}

type trackingTripRepo interface {
	UpdateStats(ctx context.Context, tripID string, distKM, speed float64) error
}

type TrackingService struct {
	gpsRepo  trackingGPSRepo
	tripRepo trackingTripRepo
}

func NewTrackingService(gpsRepo trackingGPSRepo, tripRepo trackingTripRepo) *TrackingService {
	return &TrackingService{
		gpsRepo:  gpsRepo,
		tripRepo: tripRepo,
	}
}

// SavePoint stores a GPS point and incrementally updates trip distance and max speed.
func (s *TrackingService) SavePoint(ctx context.Context, point model.GPSPoint) error {
	lastPoint, err := s.gpsRepo.LastByTripID(ctx, point.TripID)
	if err != nil {
		return fmt.Errorf("tracking: get last point: %w", err)
	}

	if err := s.gpsRepo.Insert(ctx, point); err != nil {
		return fmt.Errorf("tracking: insert point: %w", err)
	}

	if lastPoint != nil {
		distKM := geo.HaversineDistance(lastPoint.Lat, lastPoint.Lon, point.Lat, point.Lon)
		if err := s.tripRepo.UpdateStats(ctx, point.TripID, distKM, float64(point.Speed)); err != nil {
			return fmt.Errorf("tracking: update stats: %w", err)
		}
	}

	return nil
}
