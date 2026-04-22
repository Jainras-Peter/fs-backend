package hbl_schema

// PreviewHBLRequest is the JSON payload for POST /api/v1/preview/hbl
type PreviewHBLRequest struct {
	MBLNumber   string   `json:"mbl_number" binding:"required"`
	ShipmentList []string `json:"shipment_list" binding:"required"` // array of shipment_ids
}

// PreviewHBLResponse is the response from POST /api/v1/preview/hbl
type PreviewHBLResponse struct {
	MBLNumber  string    `json:"mbl_number"`
	TotalCount int       `json:"total_count"`
	HBLList    []HBLData `json:"hbl_list"`
}
