package controllers

import (
	"net/http"

	"fs-backend/models/hbl_schema"
	"fs-backend/services"

	"github.com/gin-gonic/gin"
)

// DocumentPreviewController handles document preview endpoints
type DocumentPreviewController struct {
	service services.DocumentPreviewService
}

// NewDocumentPreviewController creates a new DocumentPreviewController
func NewDocumentPreviewController(service services.DocumentPreviewService) *DocumentPreviewController {
	return &DocumentPreviewController{service: service}
}

// PreviewHBL handles POST /api/v1/preview/hbl
// Accepts JSON body with mbl_number and shipper_list (array of shipper_ids)
func (ctrl *DocumentPreviewController) PreviewHBL(ctx *gin.Context) {
	var req hbl_schema.PreviewHBLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.MBLNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "mbl_number is required"})
		return
	}
	if len(req.ShipperList) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "shipper_list must contain at least one shipper_id"})
		return
	}

	result, err := ctrl.service.PreviewHBL(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// UpdateHBL handles PUT /api/v1/hbl/:hbl_number
func (ctrl *DocumentPreviewController) UpdateHBL(ctx *gin.Context) {
	hblNumber := ctx.Param("hbl_number")
	if hblNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "hbl_number is required"})
		return
	}

	var data hbl_schema.HBLData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.UpdateHBL(ctx.Request.Context(), hblNumber, data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "HBL updated successfully"})
}
