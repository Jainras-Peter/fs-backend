package mbl_schema

// ConvertMBLResponse is the API response for POST /api/v1/convert/mbl
type ConvertMBLResponse struct {
	MBLNumber   string          `json:"mbl_number"`
	ShipperList []ShipperDetail `json:"shipper_list"`
}

// ShipperDetail holds individual shipper information returned in the response
type ShipperDetail struct {
	ShipperID      string `json:"shipper_id"`
	ShipperName    string `json:"shipper_name"`
	ShipperAddress string `json:"shipper_address"`
	ShipperContact string `json:"shipper_contact"`
}
