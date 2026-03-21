package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"fs-backend/models"
)

type PdfGeneratorService interface {
	Generate(ctx context.Context, payload models.PdfGenerationRequest) (*models.PdfGeneratorResult, error)
}

type pdfGeneratorService models.PdfGeneratorService

func NewPdfGeneratorService(baseURL string) PdfGeneratorService {
	return &pdfGeneratorService{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *pdfGeneratorService) Generate(ctx context.Context, payload models.PdfGenerationRequest) (*models.PdfGeneratorResult, error) {
	if s.BaseURL == "" {
		return nil, errors.New("pdf service base URL is not configured")
	}
	if payload.MBLNumber == "" {
		return nil, errors.New("mbl_number is required")
	}
	if len(payload.HBLList) == 0 {
		return nil, errors.New("hbl_list is required")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.JoinPath(s.BaseURL, "api", "v1", "generate")
	if err != nil {
		return nil, err
	}

	log.Printf("forwarding pdf generation request to %s (mbl_number=%s, hbl_count=%d)", endpoint, payload.MBLNumber, len(payload.HBLList))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("pdf generation response status=%d content_type=%s", resp.StatusCode, resp.Header.Get("Content-Type"))

	return &models.PdfGeneratorResult{
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Body:        respBody,
	}, nil
}
