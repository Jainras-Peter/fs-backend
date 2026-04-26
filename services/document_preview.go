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
	UpdateHBL(ctx context.Context, hblNumber string, data hbl_schema.HBLData) error
}

type documentPreviewService struct {
	mblRepo      repository.MBLRepository
	hblRepo      repository.HBLRepository
	shipmentRepo repository.ShipmentRepository
	shipperRepo  repository.ShipperRepository
	mblCacheRepo repository.MBLCacheRepository
}

// NewDocumentPreviewService creates a new DocumentPreviewService with all dependencies
func NewDocumentPreviewService(
	mblRepo repository.MBLRepository,
	hblRepo repository.HBLRepository,
	shipmentRepo repository.ShipmentRepository,
	shipperRepo repository.ShipperRepository,
	mblCacheRepo repository.MBLCacheRepository,
) DocumentPreviewService {
	return &documentPreviewService{
		mblRepo:      mblRepo,
		hblRepo:      hblRepo,
		shipmentRepo: shipmentRepo,
		shipperRepo:  shipperRepo,
		mblCacheRepo: mblCacheRepo,
	}
}

// PreviewHBL generates multiple HBLs from an MBL and a list of shipment IDs.
// Flow:
// 1. Fetch MBL from DB
// 2. Fetch shipments by shipment IDs
// 3. Extract shipper IDs from shipments and fetch shipper details
// 4. For each shipment: map MBL + shipment + shipper → HBL, generate HBL number, store in DB
// 5. Return all generated HBLs
func (s *documentPreviewService) PreviewHBL(ctx context.Context, req hbl_schema.PreviewHBLRequest) (*hbl_schema.PreviewHBLResponse, error) {
	// Step 1: Fetch MBL from DB
	mblDoc, err := s.mblRepo.FindByMBLNumber(ctx, req.MBLNumber)
	if err != nil {
		return nil, fmt.Errorf("MBL not found for number %s: %w", req.MBLNumber, err)
	}
	log.Printf("Fetched MBL: %s", req.MBLNumber)

	// Fetch MBL Cache to get raw extracted data for accurate scores
	mblCacheDoc, err := s.mblCacheRepo.FindByMBLNumber(ctx, req.MBLNumber)
	var validationScore, accuracyScore float64
	if err == nil && mblCacheDoc != nil {
		validationScore, accuracyScore = CalculateScores(mblCacheDoc.ExtractedData)
	} else {
		log.Printf("Warning: MBL_Cache not found for %s, scores will be 0", req.MBLNumber)
	}

	// Step 2: Fetch shipments for the given shipment IDs
	shipments, err := s.shipmentRepo.FindByShipmentIDs(ctx, req.ShipmentList)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shipments: %w", err)
	}

	// Build shipment_id → shipment map, and collect unique shipper IDs
	shipmentByID := make(map[string]repository.ShipmentDocument)
	shipperIDMap := make(map[string]bool)
	var shipperIDs []string

	for _, shipment := range shipments {
		shipmentByID[shipment.ShipmentID] = shipment
		if !shipperIDMap[shipment.ShipperID] {
			shipperIDMap[shipment.ShipperID] = true
			shipperIDs = append(shipperIDs, shipment.ShipperID)
		}
	}

	// Step 3: Fetch shipper details
	shippers, err := s.shipperRepo.FindByShipperIDs(ctx, shipperIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shippers: %w", err)
	}

	// Build shipper_id → shipper map
	shipperByID := make(map[string]repository.ShipperDocument)
	for _, shipper := range shippers {
		shipperByID[shipper.ShipperID] = shipper
	}

	// Step 4: Generate HBLs — one per shipment
	var hblList []hbl_schema.HBLData
	
	totalCount, err := s.hblRepo.CountTotal(ctx)
	if err != nil {
		totalCount = 0
	}
	hblIndex := int(totalCount) + 1

	for _, shipmentID := range req.ShipmentList {
		shipment, shipmentFound := shipmentByID[shipmentID]
		if !shipmentFound {
			log.Printf("Warning: no shipment found for shipment_id %s, skipping", shipmentID)
			continue
		}

		shipper, shipperFound := shipperByID[shipment.ShipperID]
		if !shipperFound {
			log.Printf("Warning: no shipper details found for shipper_id %s (shipment %s), skipping", shipment.ShipperID, shipmentID)
			continue
		}

		// Generate HBL number based on global HBL collection count
		hblNumber := generateHBLNumber(req.MBLNumber, hblIndex)

		// Map MBL + shipment + shipper → HBL
		hblData := mapMBLToHBL(mblDoc.MBL, shipment, shipper, hblNumber, mblDoc.Mode, validationScore, accuracyScore)

		// Store HBL in DB
		hblDoc := &hbl_schema.HBLDocument{
			ShipmentID: shipment.ShipmentID,
			HBLNumber:  hblNumber,
			HBL:        hblData,
		}
		if err := s.hblRepo.InsertHBL(ctx, hblDoc); err != nil {
			log.Printf("Warning: failed to store HBL %s: %v", hblNumber, err)
		} else {
			log.Printf("HBL stored in DB: %s (shipment: %s, shipper: %s)", hblNumber, shipmentID, shipment.ShipperID)
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

// UpdateHBL updates an existing HBL document
func (s *documentPreviewService) UpdateHBL(ctx context.Context, hblNumber string, data hbl_schema.HBLData) error {
	return s.hblRepo.UpdateHBL(ctx, hblNumber, data)
}
