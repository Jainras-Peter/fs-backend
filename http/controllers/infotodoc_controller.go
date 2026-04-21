package controllers

import (
	"fs-backend/models"
	"fs-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoToDocController struct {
	service services.InfoToDocService
}

func NewInfoToDocController(service services.InfoToDocService) *InfoToDocController {
	return &InfoToDocController{service: service}
}

func (c *InfoToDocController) HandleBillOfLading(ctx *gin.Context) {
	var req models.InfoToDocCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	if req.Template == "" || req.Filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Template and Filename are required"})
		return
	}

	docURL, err := c.service.ProcessBillOfLading(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process document", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Document generated successfully",
		"uploadedFiles": []map[string]interface{}{
			{
				"filename": req.Filename,
				"url":      docURL,
			},
		},
	})
}
