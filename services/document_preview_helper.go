package services

import (
	"fmt"

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
) hbl_schema.HBLData {
	return hbl_schema.HBLData{
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

		// Freight details from MBL
		FreightDetails: hbl_schema.HBLFreightDetails{
			FreightStatus: mbl.FreightPaymentType,
		},
	}
}
