package services

import (
	"fmt"
	"regexp"
	"strings"

	"fs-backend/models/hbl_schema"
	"fs-backend/models/mbl_schema"
	"fs-backend/repository"
)

// generateHBLNumber creates a unique HBL number based on the MBL number and index.
// Format: HBL-{MBL_NUMBER}-001, HBL-{MBL_NUMBER}-002, etc.
func generateHBLNumber(mblNumber string, index int) string {
	return fmt.Sprintf("HBL-%s-%03d", mblNumber, index)
}

// mapMBLToHBL maps MBL data + shipment cargo data + shipper details into an HBLData struct.
//
// Mapping rules:
//   - HBL shipper        ← DB shippers collection (actual shipper for this HBL)
//   - HBL forwarding_agent ← MBL shipper (MBL shipper becomes the freight forwarder)
//   - HBL carrier         ← MBL carrier
//   - HBL consignee       ← MBL consignee
//   - HBL notify_party    ← MBL notify_party
//   - HBL routing         ← MBL routing
//   - HBL vessel_details  ← MBL vessel_details
//   - HBL shipment_dates  ← MBL shipment_dates
//   - HBL cargo details   ← DB shipments collection (per-shipper cargo)
//   - HBL container info  ← MBL cargo (container_no, seal_number)
//   - HBL freight_status  ← MBL freight_payment_type
func mapMBLToHBL(
	mbl mbl_schema.MBLData,
	shipment repository.ShipmentDocument,
	shipper repository.ShipperDocument,
	hblNumber string,
	mode string,
	validationScore float64,
	accuracyScore float64,
) hbl_schema.HBLData {
	hbl := hbl_schema.HBLData{
		BillType:         "HBL",
		SeaWaybillNo:     hblNumber,
		CarrierReference: mbl.BillOfLadingNo,
		ExportReference:  mbl.Carrier.ReferenceNo,
		MovementType:     mode,

		// Carrier from MBL
		Carrier: hbl_schema.HBLCarrier{
			Name: mbl.Carrier.Name,
		},

		// Shipper from DB (actual shipper for this HBL)
		Shipper: hbl_schema.HBLParty{
			Name:    shipper.ShipperName,
			Address: shipper.ShipperAddress,
		},

		// Consignee from MBL
		Consignee: hbl_schema.HBLParty{
			Name:    mbl.Consignee.Name,
			Address: mbl.Consignee.Address,
		},

		// Notify party from MBL
		NotifyParty: hbl_schema.HBLParty{
			Name:    mbl.NotifyParty.Name,
			Address: mbl.NotifyParty.Address,
		},

		// Forwarding agent = MBL shipper
		ForwardingAgent: hbl_schema.HBLParty{
			Name:    mbl.Shipper.Name,
			Address: mbl.Shipper.Address,
		},

		// Routing from MBL
		Routing: hbl_schema.HBLRouting{
			PlaceOfReceipt:  mbl.Routing.PlaceOfReceipt,
			PortOfLoading:   mbl.Routing.PortOfLoading,
			PortOfDischarge: mbl.Routing.PortOfDischarge,
			PlaceOfDelivery: mbl.Routing.PlaceOfDelivery,
		},

		// Vessel details from MBL
		VesselDetails: []hbl_schema.HBLVessel{
			{
				VesselName: mbl.VesselDetails.VesselName,
				VoyageNo:   mbl.VesselDetails.VoyageNo,
			},
		},

		// Shipment dates from MBL
		ShipmentDates: hbl_schema.HBLShipmentDates{
			PlaceAndDateOfIssue: fmt.Sprintf("%s, %s", mbl.ShipmentDates.PlaceOfIssue, mbl.ShipmentDates.DateOfIssue),
			FreightPayableAt:    mbl.Routing.PortOfDischarge,
		},

		// Container details: container info from MBL, cargo details from DB shipment
		ContainerDetails: []hbl_schema.HBLContainer{
			{
				ContainerNo:        mbl.Cargo.ContainerNo,
				ContainerSize:      mbl.Cargo.ContainerType,
				SealNo:             mbl.Cargo.SealNumber,
				PackageCount:       shipment.PackagesCount,
				MarksAndNumbers:    shipment.MarksAndNumbers,
				DescriptionOfGoods: shipment.GoodsDescription,
				GrossWeight: hbl_schema.HBLWeightMeasurement{
					Value: shipment.GrossWeight,
					Unit:  "KGS",
				},
				NetWeight: hbl_schema.HBLWeightMeasurement{
					Value: shipment.NetWeight,
					Unit:  "KGS",
				},
				Measurement: hbl_schema.HBLWeightMeasurement{
					Value: shipment.Volume,
					Unit:  "CBM",
				},
			},
		},

		// Shipment summary from DB shipment
		ShipmentSummary: hbl_schema.HBLShipmentSummary{
			TotalContainersReceived: 1,
			PackagesReceived:        shipment.PackagesCount,
		},

		FreightDetails: hbl_schema.HBLFreightDetails{
			FreightStatus: mbl.FreightPaymentType,
		},
	}

	hbl.ValidationScore = validationScore
	hbl.AccuracyScore = accuracyScore
	return hbl
}

// CalculateScores calculates the validation and accuracy scores based on raw extracted MBL flat data
func CalculateScores(extracted map[string]interface{}) (float64, float64) {
	totalExpected := 20.0
	validated := 0.0
	accurate := 0.0

	// Safe string extractor
	getStr := func(key string) string {
		if val, ok := extracted[key]; ok && val != nil {
			return fmt.Sprintf("%v", val)
		}
		return ""
	}

	checkValid := func(val string) bool {
		if strings.TrimSpace(val) != "" && strings.TrimSpace(val) != "null" {
			validated++
			return true
		}
		return false
	}

	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9\s\-\/\.]+$`).MatchString
	isContainerNo := regexp.MustCompile(`^[A-Z]{4}[0-9]{7}$`).MatchString
	isValidPayment := func(val string) bool {
		v := strings.ToLower(val)
		return v == "prepaid" || v == "collect"
	}

	// 1. MBL Number
	if v := getStr("mbl_number"); checkValid(v) {
		if isAlphanumeric(v) { accurate++ }
	}
	// 2. Carrier Name
	if v := getStr("carrier_name"); checkValid(v) { accurate++ }
	// 3. Shipper Name
	if v := getStr("shipper_name"); checkValid(v) { accurate++ }
	// 4. Shipper Address
	if v := getStr("shipper_address"); checkValid(v) { accurate++ }
	// 5. Consignee Name
	if v := getStr("consignee_name"); checkValid(v) { accurate++ }
	// 6. Consignee Address
	if v := getStr("consignee_address"); checkValid(v) { accurate++ }
	// 7. Notify Party Name
	if v := getStr("notify_party_name"); checkValid(v) { accurate++ }
	// 8. Place of Receipt
	if v := getStr("place_of_receipt"); checkValid(v) {
		if !strings.Contains(v, "|") && !strings.Contains(v, "=") { accurate++ }
	}
	// 9. Port of Loading
	if v := getStr("port_of_loading"); checkValid(v) {
		if !strings.Contains(v, "|") && !strings.Contains(v, "=") { accurate++ }
	}
	// 10. Port of Discharge
	if v := getStr("port_of_discharge"); checkValid(v) {
		if !strings.Contains(v, "|") && !strings.Contains(v, "=") { accurate++ }
	}
	// 11. Place of Delivery
	if v := getStr("place_of_delivery"); checkValid(v) {
		if !strings.Contains(v, "|") && !strings.Contains(v, "=") { accurate++ }
	}
	// 12. Vessel Name
	if v := getStr("vessel_name"); checkValid(v) { accurate++ }
	// 13. Voyage No
	if v := getStr("voyage_number"); checkValid(v) { accurate++ }
	// 14. Date of Issue
	if v := getStr("date_of_issue"); checkValid(v) { accurate++ }
	// 15. Place of Issue
	if v := getStr("place_of_issue"); checkValid(v) { accurate++ }
	// 16. Container No
	if v := getStr("container_number"); checkValid(v) {
		if isContainerNo(strings.ReplaceAll(v, " ", "")) { accurate++ }
	}
	// 17. Seal Number
	if v := getStr("seal_number"); checkValid(v) { accurate++ }
	// 18. Marks and Numbers
	if v := getStr("marks_and_numbers"); checkValid(v) { accurate++ }
	// 19. Freight Payment Type
	if v := getStr("freight_payment_type"); checkValid(v) {
		if isValidPayment(strings.TrimSpace(v)) { accurate++ }
	}
	// 20. Terms of Sale
	if v := getStr("terms_of_sale"); checkValid(v) { accurate++ }

	if validated > totalExpected {
		validated = totalExpected
	}
	if accurate > totalExpected {
		accurate = totalExpected
	}

	validationScore := (validated / totalExpected) * 100
	accuracyScore := (accurate / totalExpected) * 100

	return validationScore, accuracyScore
}
