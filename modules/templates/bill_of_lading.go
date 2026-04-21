package modules



type BillOfLadingTemplate struct {
	Shipper                 string    `json:"shipper" bson:"shipper"`
	ShipperReference        string    `json:"shipper_reference" bson:"shipper_reference"`
	CarrierReference        string    `json:"carrier_reference" bson:"carrier_reference"`
	BillOfLadingNumber      string    `json:"bill_of_lading_number" bson:"bill_of_lading_number"`
	UniqueConsignmentRef    string    `json:"unique_consignment_reference" bson:"unique_consignment_reference"`
	Consignee               string    `json:"consignee" bson:"consignee"`
	CarrierName             string    `json:"carrier_name" bson:"carrier_name"`
	NotifyParty             string    `json:"notify_party" bson:"notify_party"`
	AdditionalNotifyParty   string    `json:"additional_notify_party" bson:"additional_notify_party"`
	PreCarriageBy           string    `json:"pre_carriage_by" bson:"pre_carriage_by"`
	PlaceOfReceipt          string    `json:"place_of_receipt" bson:"place_of_receipt"`
	VesselAircraft          string    `json:"vessel_aircraft" bson:"vessel_aircraft"`
	VoyageNo                string    `json:"voyage_no" bson:"voyage_no"`
	PortOfLoading           string    `json:"port_of_loading" bson:"port_of_loading"`
	PortOfDischarge         string    `json:"port_of_discharge" bson:"port_of_discharge"`
	PlaceOfDelivery         string    `json:"place_of_delivery" bson:"place_of_delivery"`
	FinalDestination        string    `json:"final_destination" bson:"final_destination"`
	AdditionalInformation   string    `json:"additional_information" bson:"additional_information"`
	MarksAndNumbers         string    `json:"marks_and_numbers" bson:"marks_and_numbers"`
	KindAndNoOfPackages     string    `json:"kind_and_no_of_packages" bson:"kind_and_no_of_packages"`
	DescriptionOfGoods      string    `json:"description_of_goods" bson:"description_of_goods"`
	NetWeight               string    `json:"net_weight" bson:"net_weight"`
	GrossWeight             string    `json:"gross_weight" bson:"gross_weight"`
	Measurements            string    `json:"measurements" bson:"measurements"`
	TotalNoOfContainers     string    `json:"total_no_of_containers" bson:"total_no_of_containers"`
	NoOfOriginalBills       string    `json:"no_of_original_bills" bson:"no_of_original_bills"`
	Incoterms               string    `json:"incoterms" bson:"incoterms"`
	PayableAt               string    `json:"payable_at" bson:"payable_at"`
	FreightCharges          string    `json:"freight_charges" bson:"freight_charges"`
	ShippedOnBoardDate      string    `json:"shipped_on_board_date" bson:"shipped_on_board_date"`
	PlaceAndDateOfIssue     string    `json:"place_and_date_of_issue" bson:"place_and_date_of_issue"`
	SignatoryCompany        string    `json:"signatory_company" bson:"signatory_company"`
	NameOfAuthorizedSign    string    `json:"name_of_authorized_signatory" bson:"name_of_authorized_signatory"`
}
