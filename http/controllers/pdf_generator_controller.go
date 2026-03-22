package controllers

import (
	"encoding/json"
	"fmt"
	"fs-backend/models"
	"fs-backend/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PdfGeneratorController struct {
	service        services.PdfGeneratorService
	saveController *PdfSaveController
}

func NewPdfGeneratorController(service services.PdfGeneratorService, saveController *PdfSaveController) *PdfGeneratorController {
	return &PdfGeneratorController{service: service, saveController: saveController}
}

func (c *PdfGeneratorController) Generate(ctx *gin.Context) {
	documentTo := strings.TrimSpace(ctx.Query("documentTo"))
	if documentTo == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "documentTo query parameter is required"})
		return
	}

	var req models.PdfGenerationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("pdf-generator request received: mbl_number=%s total_count=%d hbl_count=%d documentTo=%s", req.MBLNumber, req.TotalCount, len(req.HBLList), documentTo)

	result, err := c.service.Generate(ctx.Request.Context(), req, documentTo)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if result.StatusCode < http.StatusOK || result.StatusCode >= http.StatusMultipleChoices {
		contentType := result.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		ctx.Data(result.StatusCode, contentType, result.Body)
		return
	}

	var uploadResponse models.PdfGeneratorUploadResponse
	if err := json.Unmarshal(result.Body, &uploadResponse); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to parse pdf-generator response: %v", err)})
		return
	}

	saveResult, err := c.saveController.Save(ctx.Request.Context(), models.PdfSaveRequest{
		UploadedFiles: uploadResponse.UploadedFiles,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success":       uploadResponse.Success,
		"message":       uploadResponse.Message,
		"uploadedFiles": uploadResponse.UploadedFiles,
		"savedCount":    saveResult.SavedCount,
	})
}
