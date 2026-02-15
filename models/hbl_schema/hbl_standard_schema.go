package hbl_schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HBLDocument is the top-level struct stored in MongoDB "HBL" collection
type HBLDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ShipmentID string             `bson:"shipment_id" json:"shipment_id"`
	HBLNumber  string             `bson:"hbl_number" json:"hbl_number"`
	HBL        HBLData            `bson:"hbl" json:"hbl"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

// HBLData contains all the fields of a House Bill of Lading
type HBLData struct {
	BillType           string             `bson:"bill_type" json:"bill_type"`
	SeaWaybillNo       string             `bson:"sea_waybill_no" json:"sea_waybill_no"`
	CarrierReference   string             `bson:"carrier_reference" json:"carrier_reference"`
	ExportReference    string             `bson:"export_reference" json:"export_reference"`
	ConsigneeReference string             `bson:"consignee_reference" json:"consignee_reference"`
	Carrier            HBLCarrier         `bson:"carrier" json:"carrier"`
	Shipper            HBLParty           `bson:"shipper" json:"shipper"`
	Consignee          HBLParty           `bson:"consignee" json:"consignee"`
	NotifyParty        HBLParty           `bson:"notify_party" json:"notify_party"`
	ForwardingAgent    HBLParty           `bson:"forwarding_agent" json:"forwarding_agent"`
	MovementType       string             `bson:"movement_type" json:"movement_type"`
	Routing            HBLRouting         `bson:"routing" json:"routing"`
	VesselDetails      []HBLVessel        `bson:"vessel_details" json:"vessel_details"`
	ShipmentDates      HBLShipmentDates   `bson:"shipment_dates" json:"shipment_dates"`
	ContainerDetails   []HBLContainer     `bson:"container_details" json:"container_details"`
	ReeferDetails      HBLReeferDetails   `bson:"reefer_details" json:"reefer_details"`
	ShipmentSummary    HBLShipmentSummary `bson:"shipment_summary" json:"shipment_summary"`
	FreightDetails     HBLFreightDetails  `bson:"freight_details" json:"freight_details"`
}

// HBLCarrier holds carrier name
type HBLCarrier struct {
	Name string `bson:"name" json:"name"`
}

// HBLParty represents a party with name and address
type HBLParty struct {
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
}

// HBLRouting holds the shipment routing information
type HBLRouting struct {
	PlaceOfReceipt  string `bson:"place_of_receipt" json:"place_of_receipt"`
	PortOfLoading   string `bson:"port_of_loading" json:"port_of_loading"`
	PortOfDischarge string `bson:"port_of_discharge" json:"port_of_discharge"`
	PlaceOfDelivery string `bson:"place_of_delivery" json:"place_of_delivery"`
}

// HBLVessel holds vessel and voyage info (array because HBL can have multiple legs)
type HBLVessel struct {
	VesselName string `bson:"vessel_name" json:"vessel_name"`
	VoyageNo   string `bson:"voyage_no" json:"voyage_no"`
}

// HBLShipmentDates holds date/place of issue and freight payable info
type HBLShipmentDates struct {
	PlaceAndDateOfIssue string `bson:"place_and_date_of_issue" json:"place_and_date_of_issue"`
	FreightPayableAt    string `bson:"freight_payable_at" json:"freight_payable_at"`
}

// HBLContainer holds per-container cargo details
type HBLContainer struct {
	ContainerNo        string               `bson:"container_no" json:"container_no"`
	ContainerSize      string               `bson:"container_size" json:"container_size"`
	SealNo             string               `bson:"seal_no" json:"seal_no"`
	PackageCount       int                  `bson:"package_count" json:"package_count"`
	PackageType        string               `bson:"package_type" json:"package_type"`
	MarksAndNumbers    string               `bson:"marks_and_numbers" json:"marks_and_numbers"`
	DescriptionOfGoods string               `bson:"description_of_goods" json:"description_of_goods"`
	HSCode             string               `bson:"hs_code" json:"hs_code"`
	GrossWeight        HBLWeightMeasurement `bson:"gross_weight" json:"gross_weight"`
	NetWeight          HBLWeightMeasurement `bson:"net_weight" json:"net_weight"`
	Measurement        HBLWeightMeasurement `bson:"measurement" json:"measurement"`
}

// HBLWeightMeasurement holds a value and its unit
type HBLWeightMeasurement struct {
	Value float64 `bson:"value" json:"value"`
	Unit  string  `bson:"unit" json:"unit"`
}

// HBLReeferDetails holds reefer container settings
type HBLReeferDetails struct {
	Temperature string `bson:"temperature" json:"temperature"`
	Humidity    string `bson:"humidity" json:"humidity"`
	Ventilation string `bson:"ventilation" json:"ventilation"`
}

// HBLShipmentSummary holds summary totals
type HBLShipmentSummary struct {
	TotalContainersReceived int `bson:"total_containers_received" json:"total_containers_received"`
	PackagesReceived        int `bson:"packages_received" json:"packages_received"`
}

// HBLFreightDetails holds freight status and port charge info
type HBLFreightDetails struct {
	FreightStatus         string         `bson:"freight_status" json:"freight_status"`
	FreeTimeAtDestination string         `bson:"free_time_at_destination" json:"free_time_at_destination"`
	PortCharges           HBLPortCharges `bson:"port_charges" json:"port_charges"`
}

// HBLPortCharges holds origin/destination charge types
type HBLPortCharges struct {
	Origin      string `bson:"origin" json:"origin"`
	Destination string `bson:"destination" json:"destination"`
}
