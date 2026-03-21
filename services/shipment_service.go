package services

import (
	"context"
	"errors"
	"fs-backend/repository"
)

type ShipmentService interface {
	GetAllShipments(ctx context.Context) ([]ShipmentWithStatusDTO, error)
	InsertShipment(ctx context.Context, doc *repository.ShipmentDocument) (string, error)
	UpdateShipment(ctx context.Context, id string, doc *repository.ShipmentDocument) error
	DeleteShipment(ctx context.Context, id string) error
}

type shipmentService struct {
	shipmentRepo repository.ShipmentRepository
	bookingRepo  repository.BookingRepository
	shipperRepo  repository.ShipperRepository
}

func NewShipmentService(shipmentRepo repository.ShipmentRepository, bookingRepo repository.BookingRepository, shipperRepo repository.ShipperRepository) ShipmentService {
	return &shipmentService{
		shipmentRepo: shipmentRepo,
		bookingRepo:  bookingRepo,
		shipperRepo:  shipperRepo,
	}
}

// ShipmentWithStatusDTO extends ShipmentDocument for the frontend
type ShipmentWithStatusDTO struct {
	repository.ShipmentDocument `bson:",inline"`
	MBLNumber                   string `json:"mbl_number"`
	Status                      string `json:"status"`
}

func (s *shipmentService) GetAllShipments(ctx context.Context) ([]ShipmentWithStatusDTO, error) {
	shipments, err := s.shipmentRepo.GetAllShipments(ctx)
	if err != nil {
		return nil, err
	}
	
	var dtos []ShipmentWithStatusDTO
	for _, shipment := range shipments {
		dto := ShipmentWithStatusDTO{
			ShipmentDocument: shipment,
			MBLNumber:        "-",
			Status:           "Yet to sync",
		}
		
		booking, err := s.bookingRepo.FindByShipmentID(ctx, shipment.ShipmentID)
		if err == nil && booking != nil {
			dto.MBLNumber = booking.MBLNumber
			dto.Status = "MBL number Sycned"
		}
		dtos = append(dtos, dto)
	}
	
	return dtos, nil
}

func (s *shipmentService) InsertShipment(ctx context.Context, doc *repository.ShipmentDocument) (string, error) {
	// Validate shipper ID
	if doc.ShipperID == "" {
		return "", errors.New("shipper ID is required")
	}
	shippers, err := s.shipperRepo.FindByShipperIDs(ctx, []string{doc.ShipperID})
	if err != nil {
		return "", err
	}
	if len(shippers) == 0 {
		return "", errors.New("Invalid shipper id")
	}

	// Auto ID logic
	newID, err := s.shipmentRepo.GetNextShipmentID(ctx)
	if err != nil {
		return "", err
	}
	doc.ShipmentID = newID
	
	err = s.shipmentRepo.InsertShipment(ctx, doc)
	if err != nil {
		return "", err
	}
	return newID, nil
}

func (s *shipmentService) UpdateShipment(ctx context.Context, id string, doc *repository.ShipmentDocument) error {
	return s.shipmentRepo.UpdateShipment(ctx, id, doc)
}

func (s *shipmentService) DeleteShipment(ctx context.Context, id string) error {
	err := s.shipmentRepo.DeleteShipment(ctx, id)
	if err != nil {
		return err
	}
	
	// Ensure it is also removed from any Booking that might contain it
	return s.bookingRepo.RemoveShipmentFromBooking(ctx, id)
}
