package services

import (
	"context"
	"fmt"
	"fs-backend/models/hbl_schema"
	"fs-backend/repository"
	"log"
)

// DocumentPreviewService defines the interface for document preview operations
type DocumentPreviewService interface {
	PreviewHBL(ctx context.Context, req hbl_schema.PreviewHBLRequest) (*hbl_schema.PreviewHBLResponse, error)
}

type documentPreviewService struct {
	mblRepo      repository.MBLRepository
	hblRepo      repository.HBLRepository
	shipmentRepo repository.ShipmentRepository
	shipperRepo  repository.ShipperRepository
}

// NewDocumentPreviewService creates a new DocumentPreviewService with all dependencies
func NewDocumentPreviewService(
	mblRepo repository.MBLRepository,
	hblRepo repository.HBLRepository,
	shipmentRepo repository.ShipmentRepository,
	shipperRepo repository.ShipperRepository,
) DocumentPreviewService {
	return &documentPreviewService{
		mblRepo:      mblRepo,
		hblRepo:      hblRepo,
		shipmentRepo: shipmentRepo,
		shipperRepo:  shipperRepo,
	}
}

// PreviewHBL generates multiple HBLs from an MBL and a list of shipper IDs.
// Flow:
// 1. Fetch MBL from DB
// 2. Fetch shipments by shipper IDs
// 3. Fetch shipper details
// 4. For each shipper: map MBL + shipment + shipper → HBL, generate HBL number, store in DB
// 5. Return all generated HBLs
func (s *documentPreviewService) PreviewHBL(ctx context.Context, req hbl_schema.PreviewHBLRequest) (*hbl_schema.PreviewHBLResponse, error) {
	// Step 1: Fetch MBL from DB
	mblDoc, err := s.mblRepo.FindByMBLNumber(ctx, req.MBLNumber)
	if err != nil {
		return nil, fmt.Errorf("MBL not found for number %s: %w", req.MBLNumber, err)
	}
	log.Printf("Fetched MBL: %s", req.MBLNumber)

	// Step 2: Fetch shipments for the given shipper IDs
	shipments, err := s.shipmentRepo.FindByShipperIDs(ctx, req.ShipperList)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shipments: %w", err)
	}

	// Build shipper_id → shipment map
	shipmentByShipperID := make(map[string]repository.ShipmentDocument)
	for _, shipment := range shipments {
		shipmentByShipperID[shipment.ShipperID] = shipment
	}

	// Step 3: Fetch shipper details
	shippers, err := s.shipperRepo.FindByShipperIDs(ctx, req.ShipperList)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shippers: %w", err)
	}

	// Build shipper_id → shipper map
	shipperByID := make(map[string]repository.ShipperDocument)
	for _, shipper := range shippers {
		shipperByID[shipper.ShipperID] = shipper
	}

	// Step 4: Generate HBLs — one per shipper
	var hblList []hbl_schema.HBLData
	hblIndex := 1

	for _, shipperID := range req.ShipperList {
		shipment, shipmentFound := shipmentByShipperID[shipperID]
		shipper, shipperFound := shipperByID[shipperID]

		if !shipmentFound {
			log.Printf("Warning: no shipment found for shipper_id %s, skipping", shipperID)
			continue
		}
		if !shipperFound {
			log.Printf("Warning: no shipper details found for shipper_id %s, skipping", shipperID)
			continue
		}

		// Generate HBL number
		hblNumber := generateHBLNumber(req.MBLNumber, hblIndex)

		// Map MBL + shipment + shipper → HBL
		hblData := mapMBLToHBL(mblDoc.MBL, shipment, shipper, hblNumber, mblDoc.Mode)

		// Store HBL in DB
		hblDoc := &hbl_schema.HBLDocument{
			ShipmentID: shipment.ShipmentID,
			HBLNumber:  hblNumber,
			HBL:        hblData,
		}
		if err := s.hblRepo.InsertHBL(ctx, hblDoc); err != nil {
			log.Printf("Warning: failed to store HBL %s: %v", hblNumber, err)
		} else {
			log.Printf("HBL stored in DB: %s (shipment: %s, shipper: %s)", hblNumber, shipment.ShipmentID, shipperID)
		}

		hblList = append(hblList, hblData)
		hblIndex++
	}

	// Step 5: Return response
	return &hbl_schema.PreviewHBLResponse{
		MBLNumber:  req.MBLNumber,
		TotalCount: len(hblList),
		HBLList:    hblList,
	}, nil
}
