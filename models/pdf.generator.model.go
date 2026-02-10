package models

import "net/http"

type PdfGeneratorRequest struct {
    DocumentTo string               `json:"documentTo" binding:"required"`
    Document   PdfGeneratorDocument `json:"document" binding:"required"`
}

type PdfGeneratorDocument struct {
    Shipper               string                           `json:"shipper,omitempty"`
    Consignee             string                           `json:"consignee,omitempty"`
    NotifyParty           string                           `json:"notifyParty,omitempty"`
    VesselName            string                           `json:"vesselName,omitempty"`
    VoyageNumber          string                           `json:"voyageNumber,omitempty"`
    PortOfLoading         string                           `json:"portOfLoading,omitempty"`
    PortOfDischarge       string                           `json:"portOfDischarge,omitempty"`
    PlaceOfReceipt        string                           `json:"placeOfReceipt,omitempty"`
    PlaceOfDelivery       string                           `json:"placeOfDelivery,omitempty"`
    MasterBLNumber        string                           `json:"masterBLNumber,omitempty"`
    BookingNumber         string                           `json:"bookingNumber,omitempty"`
    ContainerNumber       string                           `json:"containerNumber,omitempty"`
    SealNumber            string                           `json:"sealNumber,omitempty"`
    ContainerSizeAndType  string                           `json:"containerSizeAndType,omitempty"`
    NumberOfPackages      string                           `json:"numberOfPackages,omitempty"`
    DescriptionOfGoods    string                           `json:"descriptionOfGoods,omitempty"`
    GrossWeight           string                           `json:"grossWeight,omitempty"`
    NetWeight             string                           `json:"netWeight,omitempty"`
    Measurement           string                           `json:"measurement,omitempty"`
    MarksAndNumbers       string                           `json:"marksAndNumbers,omitempty"`
    HSCode                string                           `json:"hsCode,omitempty"`
    FreightTerms          string                           `json:"freightTerms,omitempty"`
    NumberOfOriginalBL    string                           `json:"numberOfOriginalBL,omitempty"`
    PlaceOfIssue          string                           `json:"placeOfIssue,omitempty"`
    DateOfIssue           string                           `json:"dateOfIssue,omitempty"`
    ShippedOnBoardDate    string                           `json:"shippedOnBoardDate,omitempty"`
    MovementType          string                           `json:"movementType,omitempty"`
    ShippersDeclaredValue string                           `json:"shippersDeclaredValue,omitempty"`
    FreightCharges        *PdfGeneratorFreightCharges      `json:"freightCharges,omitempty"`
    TermsAndConditions    string                           `json:"termsAndConditions,omitempty"`
    SpecialInstructions   *PdfGeneratorSpecialInstructions `json:"specialInstructions,omitempty"`
    HouseBLNumber         string                           `json:"houseBLNumber,omitempty"`
    TemplateNumber        string                           `json:"templateNumber,omitempty"`
}

type PdfGeneratorFreightCharges struct {
    OceanFreight            string `json:"oceanFreight,omitempty"`
    BunkerAdjustmentFactor  string `json:"bunkerAdjustmentFactor,omitempty"`
    CurrencyAdjustmentFactor string `json:"currencyAdjustmentFactor,omitempty"`
    Total                   string `json:"total,omitempty"`
}

type PdfGeneratorSpecialInstructions struct {
    Reefer    *PdfGeneratorReeferInstructions    `json:"reefer,omitempty"`
    Hazardous *PdfGeneratorHazardousInstructions `json:"hazardous,omitempty"`
    Slac      string                             `json:"slac,omitempty"`
}

type PdfGeneratorReeferInstructions struct {
    Temperature string `json:"temperature,omitempty"`
    Humidity    string `json:"humidity,omitempty"`
    Ventilation string `json:"ventilation,omitempty"`
}

type PdfGeneratorHazardousInstructions struct {
    UNNumber   string `json:"unNumber,omitempty"`
    IMOClass   string `json:"imoClass,omitempty"`
    FlashPoint string `json:"flashPoint,omitempty"`
}

type PdfGeneratorResult struct {
    StatusCode  int
    ContentType string
    Body        []byte
}

type PdfGeneratorService struct {
    BaseURL string
    Client  *http.Client
}
