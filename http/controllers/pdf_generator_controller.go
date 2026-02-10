package controllers

import (
    "net/http"

    "fs-backend/models"
    "fs-backend/services"
    "github.com/gin-gonic/gin"
)

type PdfGeneratorController struct {
    service services.PdfGeneratorService
}

func NewPdfGeneratorController(service services.PdfGeneratorService) *PdfGeneratorController {
    return &PdfGeneratorController{service: service}
}

func (c *PdfGeneratorController) Generate(ctx *gin.Context) {
    var req models.PdfGeneratorRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := c.service.Generate(ctx.Request.Context(), req.DocumentTo, req.Document)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
        return
    }

    contentType := result.ContentType
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    ctx.Data(result.StatusCode, contentType, result.Body)
}
