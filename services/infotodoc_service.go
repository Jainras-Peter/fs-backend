package services

import (
	"context"
	"encoding/json"
	"errors"
	"fs-backend/models"

	"fs-backend/repository"
	"time"
)

type InfoToDocService interface {
	ProcessTemplate(ctx context.Context, req models.InfoToDocCreateRequest) (string, error)
}

type infoToDocService struct {
	infoRepo   repository.InfoToDocRepository
	hblDocRepo repository.HBLDocRepository
	pdfService PdfGeneratorService // Note: We need a generic way to call generator, or use existing one
}

func NewInfoToDocService(infoRepo repository.InfoToDocRepository, hblDocRepo repository.HBLDocRepository, pdfService PdfGeneratorService) InfoToDocService {
	return &infoToDocService{
		infoRepo:   infoRepo,
		hblDocRepo: hblDocRepo,
		pdfService: pdfService,
	}
}

func (s *infoToDocService) ProcessTemplate(ctx context.Context, req models.InfoToDocCreateRequest) (string, error) {
	if req.Template != "BillOfLading" && req.Template != "CommercialInvoice" {
		return "", errors.New("invalid template type")
	}

	// Because existing generate api uses PdfGenerationRequest which expects MBL and HBLList.
	// We will repurpose it. We can map our bill_of_lading generic json into a simulated hbl_list payload 
	// Or even better adjust the pdf generation service to allow custom generator endpoint just for us
	// For now let's construct a payload that passes validations.
	
	dataBytes, _ := json.Marshal(req.Data)
	var rawData map[string]interface{}
	json.Unmarshal(dataBytes, &rawData)

	payload := models.TemplateGenerationRequest{
		TemplateType: req.Template,
		Filename:     req.Filename,
		Data:         rawData,
	}

	result, err := s.pdfService.GenerateTemplate(ctx, payload)
	if err != nil {
		return "", err
	}

    // Now we must save it to Cloudinary through existing logic if it was handled by fs-doc-generator. 
	// Wait, fs-doc-generator actually returns { uploadedFiles: [ { filename, type, url } ] }?
	// Let's parse result.Body
	var genResp struct {
		Success       bool   `json:"success"`
		UploadedFiles []struct {
			Filename string `json:"filename"`
			Type     string `json:"type"`
			Url      string `json:"url"`
		} `json:"uploadedFiles"`
	}

	if err := json.Unmarshal(result.Body, &genResp); err != nil {
		return "", err
	}

	if !genResp.Success || len(genResp.UploadedFiles) == 0 {
		return "", errors.New("pdf generation failed")
	}

	docURL := genResp.UploadedFiles[0].Url

	// 1. Save to info-to-doc
	docType := "Bill of Lading"
	if req.Template == "CommercialInvoice" {
		docType = "Commercial Invoice"
	}

	infoDoc := &models.InfoToDoc{
		Filename:  req.Filename,
		Type:      docType,
		Data:      req.Data,
		URL:       docURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.infoRepo.Create(ctx, infoDoc); err != nil {
		return "", err
	}

	// 2. Save to hbl_doc
	hblDoc := models.HBLDoc{
		Filename: req.Filename, // using reference name for now as requested
		URL:      docURL,
		Type:     docType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.hblDocRepo.InsertMany(ctx, []models.HBLDoc{hblDoc}); err != nil {
		return "", err
	}

	return docURL, nil
}
