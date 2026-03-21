package controllers

import (
	"context"

	"fs-backend/models"
	"fs-backend/services"
)

type PdfSaveController struct {
	service services.PdfSaveService
}

func NewPdfSaveController(service services.PdfSaveService) *PdfSaveController {
	return &PdfSaveController{service: service}
}

func (c *PdfSaveController) Save(ctx context.Context, req models.PdfSaveRequest) (*models.PdfSaveResponse, error) {
	return c.service.Save(ctx, req)
}
