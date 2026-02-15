package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fs-backend/models/mbl_schema"
	"fs-backend/repository"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

// DocumentConvertService defines the interface for document conversion operations
type DocumentConvertService interface {
	ConvertMBL(ctx context.Context, fileBytes []byte, filename string) (*mbl_schema.ConvertMBLResponse, error)
}

type documentConvertService struct {
	extractionBaseURL string
	mblRepo           repository.MBLRepository
	mblCacheRepo      repository.MBLCacheRepository
	bookingRepo       repository.BookingRepository
	shipmentRepo      repository.ShipmentRepository
	shipperRepo       repository.ShipperRepository
}

// NewDocumentConvertService creates a new DocumentConvertService with all dependencies
func NewDocumentConvertService(
	extractionBaseURL string,
	mblRepo repository.MBLRepository,
	mblCacheRepo repository.MBLCacheRepository,
	bookingRepo repository.BookingRepository,
	shipmentRepo repository.ShipmentRepository,
	shipperRepo repository.ShipperRepository,
) DocumentConvertService {
	return &documentConvertService{
		extractionBaseURL: extractionBaseURL,
		mblRepo:           mblRepo,
		mblCacheRepo:      mblCacheRepo,
		bookingRepo:       bookingRepo,
		shipmentRepo:      shipmentRepo,
		shipperRepo:       shipperRepo,
	}
}

// ConvertMBL orchestrates the full MBL conversion flow with deduplication:
// 1. Hash file → check MBL_Cache → if hit, skip extraction
// 2. Extract data from document via extraction server (on cache miss)
// 3. Save extraction result to MBL_Cache
// 4. Map extracted data to MBL schema
// 5. Check if MBL number already exists → skip insert if duplicate
// 6. Lookup linked shippers via Booking → Shipment → Shipper chain
// 7. Return response
func (s *documentConvertService) ConvertMBL(ctx context.Context, fileBytes []byte, filename string) (*mbl_schema.ConvertMBLResponse, error) {
	// Step 1: Compute file hash and check MBL_Cache
	fileHash := computeFileHash(fileBytes)
	var extractedData map[string]interface{}

	cached, err := s.mblCacheRepo.FindByFileHash(ctx, fileHash)
	if err == nil && cached != nil {
		// Cache HIT — skip extraction
		log.Printf("CACHE HIT: File hash %s found in MBL_Cache, skipping extraction", fileHash[:12])
		extractedData = cached.ExtractedData
	} else {
		// Cache MISS — call extraction server
		log.Printf("CACHE MISS: File hash %s not found, calling extraction server", fileHash[:12])
		extractedData, err = extractMBLFromDocument(s.extractionBaseURL, fileBytes, filename)
		if err != nil {
			return nil, err
		}
		log.Printf("MBL extraction completed for file: %s", filename)

		// Step 2: Save extraction result to MBL_Cache
		mblNumberForCache := ""
		if val, ok := extractedData["mbl_number"]; ok && val != nil {
			mblNumberForCache, _ = val.(string)
		}
		cacheDoc := &repository.MBLCacheDocument{
			FileHash:      fileHash,
			MBLNumber:     mblNumberForCache,
			ExtractedData: extractedData,
		}
		if err := s.mblCacheRepo.Insert(ctx, cacheDoc); err != nil {
			log.Printf("Warning: failed to save to MBL_Cache: %v", err)
		}
	}

	// Step 3: Map the flat extracted data to the structured MBL document
	mblDoc := mapExtractionToMBLDocument(extractedData)
	mblNumber := mblDoc.MBL.BillOfLadingNo
	log.Printf("MBL number extracted: %s", mblNumber)

	// Step 4: Check if MBL already exists in DB by mbl_number → skip insert if duplicate
	existingMBL, err := s.mblRepo.FindByMBLNumber(ctx, mblNumber)
	if err == mongo.ErrNoDocuments || existingMBL == nil {
		// Not found → insert
		if err := s.mblRepo.InsertMBL(ctx, mblDoc); err != nil {
			return nil, err
		}
		log.Printf("MBL document stored in DB: %s", mblNumber)
	} else {
		log.Printf("MBL %s already exists in DB, skipping insert", mblNumber)
	}

	// Step 5: Lookup linked shippers
	shipperList, err := lookupShippersByMBLNumber(ctx, mblNumber, s.bookingRepo, s.shipmentRepo, s.shipperRepo)
	if err != nil {
		log.Printf("Warning: shipper lookup failed for MBL %s: %v", mblNumber, err)
		shipperList = []mbl_schema.ShipperDetail{}
	}

	// Step 6: Build and return response
	return &mbl_schema.ConvertMBLResponse{
		MBLNumber:   mblNumber,
		ShipperList: shipperList,
	}, nil
}

// computeFileHash returns a SHA-256 hex digest of the file bytes
func computeFileHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
