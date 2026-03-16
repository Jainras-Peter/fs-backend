package services

import (
	"context"
	"fs-backend/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService interface {
	AddShipper(ctx context.Context, doc repository.ShipperDocument) (primitive.ObjectID, error)
	GetShipperList(ctx context.Context) ([]repository.ShipperDocument, error)
	UpdateShipper(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	DeleteShipper(ctx context.Context, id primitive.ObjectID) error
}

type bookingService struct {
	shipperRepo repository.ShipperRepository
}

func NewBookingService(shipperRepo repository.ShipperRepository) BookingService {
	return &bookingService{
		shipperRepo: shipperRepo,
	}
}

func (s *bookingService) AddShipper(ctx context.Context, doc repository.ShipperDocument) (primitive.ObjectID, error) {
	res, err := s.shipperRepo.CreateShipper(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (s *bookingService) GetShipperList(ctx context.Context) ([]repository.ShipperDocument, error) {
	return s.shipperRepo.FindAllShippers(ctx)
}

func (s *bookingService) UpdateShipper(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	_, err := s.shipperRepo.UpdateShipper(ctx, id, updates)
	return err
}

func (s *bookingService) DeleteShipper(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.shipperRepo.DeleteShipper(ctx, id)
	return err
}
