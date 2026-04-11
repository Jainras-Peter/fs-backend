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
	SyncBooking(ctx context.Context, mblNumber, shipmentID, carrierName, estimatedDeparture, estimatedArrival string) error
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

func (s *bookingService) SyncBooking(ctx context.Context, mblNumber, shipmentID, carrierName, estimatedDeparture, estimatedArrival string) error {
	// Fetch the shipment to get its mode
	shipments, err := s.shipmentRepo.FindByShipmentIDs(ctx, []string{shipmentID})
	if err != nil || len(shipments) == 0 {
		return errors.New("Shipment not found")
	}
	shipmentMode := shipments[0].Mode // FCL or LCL

	// Validate that shipment has a mode set
	if shipmentMode == "" {
		return errors.New("Shipment mode is not set. Please set the mode (FCL/LCL) before syncing.")
	}

	booking, err := s.bookingRepo.FindByMBLNumber(ctx, mblNumber)

	if err == nil && booking != nil {
		// Booking exists - validate mode compatibility

		// If booking has no mode set (shouldn't happen but handle it), update it now
		if booking.Mode == "" {
			// This booking was created without a mode, set it now from the first shipment
			return errors.New("MBL exists but mode is not set. Cannot proceed. Please delete and recreate the booking.")
		}

		// Rule 1: If existing booking is FCL, can't add more shipments
		if booking.Mode == "FCL" {
			return errors.New("MBL synced with FCL - cannot add more shipments")
		}

		// Rule 2: If existing booking is LCL and new shipment is FCL, reject
		if booking.Mode == "LCL" && shipmentMode == "FCL" {
			return errors.New("MBL is synced with LCL - cannot add FCL shipment")
		}

		// All validations passed, append shipment to existing booking
		return s.bookingRepo.AddShipmentToBooking(ctx, mblNumber, shipmentID)
	}

	// If booking doesn't exist, create it with new details
	newBooking := &repository.BookingDocument{
		MBLNumber:          mblNumber,
		ShipmentIDs:        []string{shipmentID},
		Mode:               shipmentMode,
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
