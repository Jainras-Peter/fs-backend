package controllers

import (
	"io"
	"net/http"

	"fs-backend/services"

	"github.com/gin-gonic/gin"
)

// DocumentConvertController handles document conversion endpoints
type DocumentConvertController struct {
	service services.DocumentConvertService
}

// NewDocumentConvertController creates a new DocumentConvertController
func NewDocumentConvertController(service services.DocumentConvertService) *DocumentConvertController {
	return &DocumentConvertController{service: service}
}

// ConvertMBL handles POST /api/v1/convert/mbl
// Accepts multipart form with: file (PDF/image), from_doc, to_doc
func (ctrl *DocumentConvertController) ConvertMBL(ctx *gin.Context) {
	// 1. Validate from_doc and to_doc fields
	fromDoc := ctx.PostForm("from_doc")
	toDoc := ctx.PostForm("to_doc")

	if fromDoc != "mbl" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "from_doc must be 'mbl'"})
		return
	}
	if toDoc != "hbl" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "to_doc must be 'hbl'"})
		return
	}

	// 2. Validate mode field
	mode := ctx.PostForm("mode")
	if mode != "FCL" && mode != "LCL" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "mode must be 'FCL' or 'LCL'"})
		return
	}

	// 3. Parse uploaded file
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer file.Close()

	// 3. Read file bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file contents"})
		return
	}

	// 5. Call service
	result, err := ctrl.service.ConvertMBL(ctx.Request.Context(), fileBytes, fileHeader.Filename, mode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. Return response
	ctx.JSON(http.StatusOK, result)
}
