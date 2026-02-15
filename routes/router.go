package routes

import (
	"fs-backend/http/controllers"
	"fs-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, pdfService services.PdfGeneratorService, docConvertService services.DocumentConvertService, docPreviewService services.DocumentPreviewService) {
	pdfController := controllers.NewPdfGeneratorController(pdfService)
	docConvertController := controllers.NewDocumentConvertController(docConvertService)
	docPreviewController := controllers.NewDocumentPreviewController(docPreviewService)

	api := router.Group("/api/v1")
	{
		api.POST("/pdf-generator", pdfController.Generate)
		api.POST("/convert/mbl", docConvertController.ConvertMBL)
		api.POST("/preview/hbl", docPreviewController.PreviewHBL)
		api.PUT("/hbl/:hbl_number", docPreviewController.UpdateHBL)
	}
}
