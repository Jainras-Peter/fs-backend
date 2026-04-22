package mbl_schema

// ConvertMBLResponse is the API response for POST /api/v1/convert/mbl
type ConvertMBLResponse struct {
	MBLNumber       string             `json:"mbl_number"`
	ShipmentsList []ShipmentListItem `json:"shipments_list"`
}

// ShipmentListItem holds individual shipment information returned in the response.
type ShipmentListItem struct {
	ShipmentID       string  `json:"shipment_id"`
	ShipperID        string  `json:"shipper_id"`
	GoodsDescription string  `json:"goods_description"`
	PackagesCount    int     `json:"packages_count"`
	GrossWeight      float64 `json:"gross_weight"`
	NetWeight        float64 `json:"net_weight"`
	Volume           float64 `json:"volume"`
	MarksAndNumbers  string  `json:"marks_and_numbers"`
	Measurement      string  `json:"measurement"`
}
