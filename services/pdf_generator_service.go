package services

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "net/url"
    "time"

    "fs-backend/models"
)

type PdfGeneratorService interface {
    Generate(ctx context.Context, documentTo string, document models.PdfGeneratorDocument) (*models.PdfGeneratorResult, error)
}

type pdfGeneratorService models.PdfGeneratorService

func NewPdfGeneratorService(baseURL string) PdfGeneratorService {
    return &pdfGeneratorService{
        BaseURL: baseURL,
        Client:  &http.Client{Timeout: 30 * time.Second},
    }
}

func (s *pdfGeneratorService) Generate(ctx context.Context, documentTo string, document models.PdfGeneratorDocument) (*models.PdfGeneratorResult, error) {
    if s.BaseURL == "" {
        return nil, errors.New("pdf service base URL is not configured")
    }
    if documentTo == "" {
        return nil, errors.New("documentTo is required")
    }

    payload := models.PdfGeneratorRequest{
        DocumentTo: documentTo,
        Document:   document,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    endpoint, err := url.JoinPath(s.BaseURL, "generate", documentTo)
    if err != nil {
        return nil, err
    }

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

    return &models.PdfGeneratorResult{
        StatusCode:  resp.StatusCode,
        ContentType: resp.Header.Get("Content-Type"),
        Body:        respBody,
    }, nil
}
