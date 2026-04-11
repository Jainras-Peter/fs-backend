package services

import (
	"context"
	"errors"
	"fs-backend/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService interface {
	AddShipper(ctx context.Context, doc repository.ShipperDocument) (primitive.ObjectID, error)
	GetShipperList(ctx context.Context) ([]repository.ShipperDocument, error)
	UpdateShipper(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	DeleteShipper(ctx context.Context, id primitive.ObjectID) error
	SyncBooking(ctx context.Context, mblNumber, mode string, shipmentIDs []string, carrierName, estimatedDeparture, estimatedArrival string) error
	GetStatusDetails(ctx context.Context) ([]repository.BookingDocument, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error
}

type bookingService struct {
	shipperRepo  repository.ShipperRepository
	bookingRepo  repository.BookingRepository
	shipmentRepo repository.ShipmentRepository
}

func NewBookingService(shipperRepo repository.ShipperRepository, bookingRepo repository.BookingRepository, shipmentRepo repository.ShipmentRepository) BookingService {
	return &bookingService{
		shipperRepo:  shipperRepo,
		bookingRepo:  bookingRepo,
		shipmentRepo: shipmentRepo,
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

func (s *bookingService) SyncBooking(ctx context.Context, mblNumber, mode string, shipmentIDs []string, carrierName, estimatedDeparture, estimatedArrival string) error {
	if len(shipmentIDs) == 0 {
		return errors.New("No shipments selected for booking")
	}

	if mode != "FCL" && mode != "LCL" {
		return errors.New("Invalid mode. Please select FCL or LCL.")
	}

	if mode == "FCL" && len(shipmentIDs) > 1 {
		return errors.New("FCL booking can only have one shipment")
	}

	if mode == "LCL" && len(shipmentIDs) > 5 {
		return errors.New("LCL booking can have at most 5 shipments")
	}

	// Fetch shipments to validate them
	shipments, err := s.shipmentRepo.FindByShipmentIDs(ctx, shipmentIDs)
	if err != nil || len(shipments) == 0 {
		return errors.New("Shipment not found")
	}

	for _, shipment := range shipments {
		if shipment.Mode == "" {
			return errors.New("All shipments must have a mode set before syncing")
		}
		if shipment.Mode != mode {
			return errors.New("Selected shipment modes do not match the chosen booking mode")
		}
	}

	booking, err := s.bookingRepo.FindByMBLNumber(ctx, mblNumber)

	if err == nil && booking != nil {
		if booking.Mode == "" {
			return errors.New("MBL exists but mode is not set. Cannot proceed. Please delete and recreate the booking.")
		}

		if booking.Mode == "FCL" {
			return errors.New("MBL synced with FCL - cannot add more shipments")
		}

		if booking.Mode == "LCL" && mode == "FCL" {
			return errors.New("MBL is synced with LCL - cannot add FCL shipment")
		}

		for _, shipmentID := range shipmentIDs {
			if err := s.bookingRepo.AddShipmentToBooking(ctx, mblNumber, shipmentID); err != nil {
				return err
			}
		}

		return nil
	}

	newBooking := &repository.BookingDocument{
		MBLNumber:          mblNumber,
		ShipmentIDs:        shipmentIDs,
		Mode:               mode,
		CarrierName:        carrierName,
		EstimatedDeparture: estimatedDeparture,
		EstimatedArrival:   estimatedArrival,
		Status:             "Booked",
	}
	return s.bookingRepo.CreateBooking(ctx, newBooking)
}

func (s *bookingService) GetStatusDetails(ctx context.Context) ([]repository.BookingDocument, error) {
	return s.bookingRepo.GetAllBookings(ctx)
}

func (s *bookingService) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	return s.bookingRepo.UpdateBookingStatus(ctx, id, status)
}
