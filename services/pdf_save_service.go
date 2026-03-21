package services

import (
	"context"
	"fmt"

	"fs-backend/models"
	"fs-backend/repository"
)

type PdfSaveService interface {
	Save(ctx context.Context, req models.PdfSaveRequest) (*models.PdfSaveResponse, error)
}

type pdfSaveService struct {
	repo repository.HBLDocRepository
}

func NewPdfSaveService(repo repository.HBLDocRepository) PdfSaveService {
	return &pdfSaveService{repo: repo}
}

func (s *pdfSaveService) Save(ctx context.Context, req models.PdfSaveRequest) (*models.PdfSaveResponse, error) {
	docs := make([]models.HBLDoc, 0, len(req.UploadedFiles))
	for _, file := range req.UploadedFiles {
		if file.Filename == "" || file.URL == "" {
			continue
		}

		docs = append(docs, models.HBLDoc{
			Filename: file.Filename,
			URL:      file.URL,
		})
	}

	if err := s.repo.InsertMany(ctx, docs); err != nil {
		return nil, fmt.Errorf("failed to store generated HBL documents: %w", err)
	}

	return &models.PdfSaveResponse{
		Success:    true,
		Message:    "PDFs saved successfully",
		SavedCount: len(docs),
	}, nil
}
