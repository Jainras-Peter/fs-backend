package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"
	"time"

	"fs-backend/models/mbl_schema"
	"fs-backend/repository"
)

// extractMBLFromDocument sends the file and MBL extraction schema to the
// Document Extraction server and returns the filled data map.
func extractMBLFromDocument(extractionBaseURL string, fileBytes []byte, filename string) (map[string]interface{}, error) {
	// Build the extraction schema
	schema := mbl_schema.GetMBLExtractionSchema()
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal extraction schema: %w", err)
	}

	// Build multipart form body
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Detect MIME type from filename extension (e.g. "application/pdf", "image/png")
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Create form file part with the correct Content-Type header
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	partHeader.Set("Content-Type", mimeType)

	filePart, err := writer.CreatePart(partHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(filePart, bytes.NewReader(fileBytes)); err != nil {
		return nil, fmt.Errorf("failed to write file bytes: %w", err)
	}

	// Add the schema field
	if err := writer.WriteField("schema", string(schemaJSON)); err != nil {
		return nil, fmt.Errorf("failed to write schema field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Send HTTP request to extraction server
	client := &http.Client{Timeout: 120 * time.Second}
	req, err := http.NewRequest(http.MethodPost, extractionBaseURL, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Printf("Sending extraction request to: %s", extractionBaseURL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("extraction server request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read extraction response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("extraction server returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse the JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse extraction response: %w", err)
	}

	return result, nil
}

// mapExtractionToMBLDocument maps the flat extracted key-value pairs to the
// nested MBLDocument struct for storage in MongoDB.
func mapExtractionToMBLDocument(data map[string]interface{}) *mbl_schema.MBLDocument {
	return &mbl_schema.MBLDocument{
		MBL: mbl_schema.MBLData{
			BillType:            getStr(data, "bill_type", "MBL"),
			BillOfLadingNo:      getStr(data, "mbl_number", ""),
			PackingListNo:       getStr(data, "packing_list_no", ""),
			NumberOfOriginalBLs: getInt(data, "number_of_original_bls"),
			TermsOfSale:         getStr(data, "terms_of_sale", ""),
			FreightPaymentType:  getStr(data, "freight_payment_type", ""),

			Carrier: mbl_schema.Carrier{
				Name:        getStr(data, "carrier_name", ""),
				SCACCode:    getStr(data, "carrier_scac_code", ""),
				ReferenceNo: getStr(data, "carrier_reference_no", ""),
			},

			Shipper: mbl_schema.Party{
				Name:    getStr(data, "shipper_name", ""),
				Address: getStr(data, "shipper_address", ""),
				Phone:   getStr(data, "shipper_phone", ""),
				Fax:     getStr(data, "shipper_fax", ""),
			},

			Consignee: mbl_schema.ConsigneeParty{
				Name:    getStr(data, "consignee_name", ""),
				Address: getStr(data, "consignee_address", ""),
				Phone:   getStr(data, "consignee_phone", ""),
				Email:   getStr(data, "consignee_email", ""),
			},

			NotifyParty: mbl_schema.NotifyParty{
				Name:    getStr(data, "notify_party_name", ""),
				Address: getStr(data, "notify_party_address", ""),
			},

			Routing: mbl_schema.Routing{
				PlaceOfReceipt:  getStr(data, "place_of_receipt", ""),
				PortOfLoading:   getStr(data, "port_of_loading", ""),
				PortOfDischarge: getStr(data, "port_of_discharge", ""),
				PlaceOfDelivery: getStr(data, "place_of_delivery", ""),
			},

			VesselDetails: mbl_schema.VesselDetails{
				VesselName: getStr(data, "vessel_name", ""),
				VoyageNo:   getStr(data, "voyage_number", ""),
			},

			ShipmentDates: mbl_schema.ShipmentDates{
				DateOfIssue:         getStr(data, "date_of_issue", ""),
				PlaceOfIssue:        getStr(data, "place_of_issue", ""),
				ShippedOnBoardDate:  getStr(data, "shipped_on_board_date", ""),
				ShippedOnBoardPlace: getStr(data, "shipped_on_board_place", ""),
			},

			Cargo: mbl_schema.Cargo{
				ContainerNo:        getStr(data, "container_number", ""),
				ContainerType:      getStr(data, "container_type", ""),
				SealNumber:         getStr(data, "seal_number", ""),
				MarksAndNumbers:    getStr(data, "marks_and_numbers", ""),
				NumberOfPackages:   getInt(data, "number_of_packages"),
				PackageType:        getStr(data, "package_type", ""),
				DescriptionOfGoods: getStr(data, "description_of_goods", ""),
				HSCode:             getStr(data, "hs_code", ""),
				GrossWeight: mbl_schema.WeightMeasurement{
					Value: getFloat(data, "gross_weight_kgs"),
					Unit:  "KGS",
				},
				NetWeight: mbl_schema.WeightMeasurement{
					Value: getFloat(data, "net_weight_kgs"),
					Unit:  "KGS",
				},
				Measurement: mbl_schema.WeightMeasurement{
					Value: getFloat(data, "measurement_cbm"),
					Unit:  "CBM",
				},
			},

			FreightCharges: mbl_schema.FreightCharges{
				OceanFreight: mbl_schema.OceanFreight{
					PrepaidAmount: getFloat(data, "ocean_freight_prepaid"),
					CollectAmount: getFloat(data, "ocean_freight_collect"),
					Currency:      getStr(data, "freight_currency", ""),
				},
			},
		},
	}
}

// lookupShippersByMBLNumber chains Booking → Shipment → Shipper DB lookups
// to find all shipper details linked to the given MBL number.
func lookupShippersByMBLNumber(
	ctx context.Context,
	mblNumber string,
	bookingRepo repository.BookingRepository,
	shipmentRepo repository.ShipmentRepository,
	shipperRepo repository.ShipperRepository,
) ([]mbl_schema.ShipperDetail, error) {
	// Step 1: Find booking by MBL number → get shipment_ids
	booking, err := bookingRepo.FindByMBLNumber(ctx, mblNumber)
	if err != nil {
		return nil, fmt.Errorf("booking not found for MBL %s: %w", mblNumber, err)
	}

	if len(booking.ShipmentIDs) == 0 {
		return []mbl_schema.ShipperDetail{}, nil
	}

	// Step 2: Find shipments by shipment_ids → get unique shipper_ids
	shipments, err := shipmentRepo.FindByShipmentIDs(ctx, booking.ShipmentIDs)
	if err != nil {
		return nil, fmt.Errorf("shipments lookup failed: %w", err)
	}

	// Collect unique shipper IDs
	shipperIDSet := make(map[string]bool)
	var shipperIDs []string
	for _, s := range shipments {
		if s.ShipperID != "" && !shipperIDSet[s.ShipperID] {
			shipperIDSet[s.ShipperID] = true
			shipperIDs = append(shipperIDs, s.ShipperID)
		}
	}

	if len(shipperIDs) == 0 {
		return []mbl_schema.ShipperDetail{}, nil
	}

	// Step 3: Find shippers by shipper_ids → get details
	shippers, err := shipperRepo.FindByShipperIDs(ctx, shipperIDs)
	if err != nil {
		return nil, fmt.Errorf("shippers lookup failed: %w", err)
	}

	// Map to response models
	result := make([]mbl_schema.ShipperDetail, 0, len(shippers))
	for _, sh := range shippers {
		result = append(result, mbl_schema.ShipperDetail{
			ShipperID:      sh.ShipperID,
			ShipperName:    sh.ShipperName,
			ShipperAddress: sh.ShipperAddress,
			ShipperContact: sh.ShipperContact,
		})
	}

	return result, nil
}

// --- Helper functions for safe type extraction from map[string]interface{} ---

func getStr(data map[string]interface{}, key string, defaultVal string) string {
	if val, ok := data[key]; ok && val != nil {
		switch v := val.(type) {
		case string:
			return v
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return defaultVal
}

func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok && val != nil {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

func getFloat(data map[string]interface{}, key string) float64 {
	if val, ok := data[key]; ok && val != nil {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0
}
