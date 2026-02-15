package mbl_schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MBLDocument is the top-level struct stored in MongoDB "MBL" collection
type MBLDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Mode      string             `bson:"mode" json:"mode"` // "FCL" or "LCL"
	MBL       MBLData            `bson:"mbl" json:"mbl"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// MBLData contains all the fields of a Master Bill of Lading
type MBLData struct {
	BillType            string         `bson:"bill_type" json:"bill_type"`
	BillOfLadingNo      string         `bson:"bill_of_lading_no" json:"bill_of_lading_no"`
	PackingListNo       string         `bson:"packing_list_no" json:"packing_list_no"`
	NumberOfOriginalBLs int            `bson:"number_of_original_bls" json:"number_of_original_bls"`
	TermsOfSale         string         `bson:"terms_of_sale" json:"terms_of_sale"`
	FreightPaymentType  string         `bson:"freight_payment_type" json:"freight_payment_type"`
	Carrier             Carrier        `bson:"carrier" json:"carrier"`
	Shipper             Party          `bson:"shipper" json:"shipper"`
	Consignee           ConsigneeParty `bson:"consignee" json:"consignee"`
	NotifyParty         NotifyParty    `bson:"notify_party" json:"notify_party"`
	Routing             Routing        `bson:"routing" json:"routing"`
	VesselDetails       VesselDetails  `bson:"vessel_details" json:"vessel_details"`
	ShipmentDates       ShipmentDates  `bson:"shipment_dates" json:"shipment_dates"`
	Cargo               Cargo          `bson:"cargo" json:"cargo"`
	FreightCharges      FreightCharges `bson:"freight_charges" json:"freight_charges"`
}

// Carrier holds carrier/shipping line information
type Carrier struct {
	Name        string `bson:"name" json:"name"`
	SCACCode    string `bson:"scac_code" json:"scac_code"`
	ReferenceNo string `bson:"reference_no" json:"reference_no"`
}

// Party represents shipper with phone and fax
type Party struct {
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
	Phone   string `bson:"phone" json:"phone"`
	Fax     string `bson:"fax" json:"fax"`
}

// ConsigneeParty represents consignee with phone and email
type ConsigneeParty struct {
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
	Phone   string `bson:"phone" json:"phone"`
	Email   string `bson:"email" json:"email"`
}

// NotifyParty holds notify party details
type NotifyParty struct {
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
}

// Routing holds the shipment routing information
type Routing struct {
	PlaceOfReceipt  string `bson:"place_of_receipt" json:"place_of_receipt"`
	PortOfLoading   string `bson:"port_of_loading" json:"port_of_loading"`
	PortOfDischarge string `bson:"port_of_discharge" json:"port_of_discharge"`
	PlaceOfDelivery string `bson:"place_of_delivery" json:"place_of_delivery"`
}

// VesselDetails holds vessel and voyage info
type VesselDetails struct {
	VesselName string `bson:"vessel_name" json:"vessel_name"`
	VoyageNo   string `bson:"voyage_no" json:"voyage_no"`
}

// ShipmentDates holds date and place of issue and on-board info
type ShipmentDates struct {
	DateOfIssue         string `bson:"date_of_issue" json:"date_of_issue"`
	PlaceOfIssue        string `bson:"place_of_issue" json:"place_of_issue"`
	ShippedOnBoardDate  string `bson:"shipped_on_board_date" json:"shipped_on_board_date"`
	ShippedOnBoardPlace string `bson:"shipped_on_board_place" json:"shipped_on_board_place"`
}

// Cargo holds container and goods information
type Cargo struct {
	ContainerNo        string            `bson:"container_no" json:"container_no"`
	ContainerType      string            `bson:"container_type" json:"container_type"`
	SealNumber         string            `bson:"seal_number" json:"seal_number"`
	MarksAndNumbers    string            `bson:"marks_and_numbers" json:"marks_and_numbers"`
	NumberOfPackages   int               `bson:"number_of_packages" json:"number_of_packages"`
	PackageType        string            `bson:"package_type" json:"package_type"`
	DescriptionOfGoods string            `bson:"description_of_goods" json:"description_of_goods"`
	HSCode             string            `bson:"hs_code" json:"hs_code"`
	GrossWeight        WeightMeasurement `bson:"gross_weight" json:"gross_weight"`
	NetWeight          WeightMeasurement `bson:"net_weight" json:"net_weight"`
	Measurement        WeightMeasurement `bson:"measurement" json:"measurement"`
}

// WeightMeasurement holds a value and its unit
type WeightMeasurement struct {
	Value float64 `bson:"value" json:"value"`
	Unit  string  `bson:"unit" json:"unit"`
}

// FreightCharges holds freight charge details
type FreightCharges struct {
	OceanFreight OceanFreight `bson:"ocean_freight" json:"ocean_freight"`
}

// OceanFreight holds prepaid/collect amounts and currency
type OceanFreight struct {
	PrepaidAmount float64 `bson:"prepaid_amount" json:"prepaid_amount"`
	CollectAmount float64 `bson:"collect_amount" json:"collect_amount"`
	Currency      string  `bson:"currency" json:"currency"`
}
