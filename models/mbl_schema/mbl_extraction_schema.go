package mbl_schema

// GetMBLExtractionSchema returns the flat key→null schema sent to the
// Document Extraction server. The server fills in the null values from the
// uploaded MBL document (PDF/image).
func GetMBLExtractionSchema() map[string]interface{} {
	return map[string]interface{}{
		// Bill basics
		"mbl_number":             nil,
		"bill_type":              nil,
		"packing_list_no":        nil,
		"number_of_original_bls": nil,
		"terms_of_sale":          nil,
		"freight_payment_type":   nil,

		// Carrier
		"carrier_name":         nil,
		"carrier_scac_code":    nil,
		"carrier_reference_no": nil,

		// Shipper
		"shipper_name":    nil,
		"shipper_address": nil,
		"shipper_phone":   nil,
		"shipper_fax":     nil,

		// Consignee
		"consignee_name":    nil,
		"consignee_address": nil,
		"consignee_phone":   nil,
		"consignee_email":   nil,

		// Notify Party
		"notify_party_name":    nil,
		"notify_party_address": nil,

		// Routing
		"place_of_receipt":  nil,
		"port_of_loading":   nil,
		"port_of_discharge": nil,
		"place_of_delivery": nil,

		// Vessel
		"vessel_name":   nil,
		"voyage_number": nil,

		// Shipment dates
		"date_of_issue":          nil,
		"place_of_issue":         nil,
		"shipped_on_board_date":  nil,
		"shipped_on_board_place": nil,

		// Cargo
		"container_number":     nil,
		"container_type":       nil,
		"seal_number":          nil,
		"marks_and_numbers":    nil,
		"number_of_packages":   nil,
		"package_type":         nil,
		"description_of_goods": nil,
		"hs_code":              nil,
		"gross_weight_kgs":     nil,
		"net_weight_kgs":       nil,
		"measurement_cbm":      nil,

		// Freight
		"ocean_freight_prepaid": nil,
		"ocean_freight_collect": nil,
		"freight_currency":      nil,
	}
}
