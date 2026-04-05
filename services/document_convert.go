package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"fs-backend/models/mbl_schema"
	"fs-backend/repository"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrUnsupportedExtractionModel = errors.New("unsupported extraction model")

// DocumentConvertService defines the interface for document conversion operations
type DocumentConvertService interface {
	ConvertMBL(ctx context.Context, fileBytes []byte, filename string, model string) (*mbl_schema.ConvertMBLResponse, error)
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
func (s *documentConvertService) ConvertMBL(ctx context.Context, fileBytes []byte, filename string, model string) (*mbl_schema.ConvertMBLResponse, error) {
	extractionEngine, err := normalizeExtractionEngine(model)
	if err != nil {
		return nil, err
	}

	// Step 1: Compute file hash and check MBL_Cache
	fileHash := computeFileHash(fileBytes)
	var extractedData map[string]interface{}

	cached, err := s.mblCacheRepo.FindByFileHashAndEngine(ctx, fileHash, extractionEngine)
	if err == nil && cached != nil {
		// Cache HIT — skip extraction
		log.Printf("CACHE HIT: File hash %s found in MBL_Cache for engine %s, skipping extraction", fileHash[:12], extractionEngine)
		extractedData = cached.ExtractedData
	} else {
		// Cache MISS — call extraction server
		log.Printf("CACHE MISS: File hash %s not found for engine %s, calling extraction server", fileHash[:12], extractionEngine)
		extractedData, err = extractMBLFromDocument(s.extractionBaseURL, fileBytes, filename, extractionEngine)
		if err != nil {
			return nil, err
		}
		log.Printf("MBL extraction completed for file: %s using engine: %s", filename, extractionEngine)

		// Validation: Ensure essential MBL details are present before caching
		mblNumberCheck := getStr(extractedData, "mbl_number", "")
		carrierNameCheck := getStr(extractedData, "carrier_name", "")
		shipperNameCheck := getStr(extractedData, "shipper_name", "")
		vesselNameCheck := getStr(extractedData, "vessel_name", "")

		validCount := 0
		if mblNumberCheck != "" {
			validCount++
		}
		if carrierNameCheck != "" {
			validCount++
		}
		if shipperNameCheck != "" {
			validCount++
		}
		if vesselNameCheck != "" {
			validCount++
		}

		if mblNumberCheck == "" || validCount < 2 {
			log.Printf("Validation failed: MBL details not found for file %s", filename)
			return nil, errors.New("MBL details not found")
		}

		// Step 2: Save extraction result to MBL_Cache
		mblNumberForCache := ""
		if val, ok := extractedData["mbl_number"]; ok && val != nil {
			mblNumberForCache, _ = val.(string)
		}
		cacheDoc := &repository.MBLCacheDocument{
			FileHash:      fileHash,
			Engine:        extractionEngine,
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

	// Step 5: Lookup linked shipments
	shipmentsList, err := lookupShipmentsByMBLNumber(ctx, mblNumber, s.bookingRepo, s.shipmentRepo)
	if err != nil {
		log.Printf("Warning: shipment lookup failed for MBL %s: %v", mblNumber, err)
		shipmentsList = []mbl_schema.ShipmentListItem{}
	}

	// Step 6: Build and return response
	return &mbl_schema.ConvertMBLResponse{
		MBLNumber:       mblNumber,
		ShipmentsList: shipmentsList,
	}, nil
}

// computeFileHash returns a SHA-256 hex digest of the file bytes
func computeFileHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func normalizeExtractionEngine(model string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(model)) {
	case "grok", "groq":
		return "groq", nil
	case "gpt-oss-120b", "gpt oss 120b", "huggingface", "hf", "openai/gpt-oss-120b:novita":
		return "huggingface", nil
	case "ollama":
		return "ollama", nil
	case "":
		return "", errors.New("model is required")
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedExtractionModel, model)
	}
}
