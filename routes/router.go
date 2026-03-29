package routes

import (
	"fs-backend/http/controllers"
	"fs-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, pdfService services.PdfGeneratorService, docConvertService services.DocumentConvertService, docPreviewService services.DocumentPreviewService, forwarderService services.ForwarderService) {
	pdfController := controllers.NewPdfGeneratorController(pdfService)
	docConvertController := controllers.NewDocumentConvertController(docConvertService)
	docPreviewController := controllers.NewDocumentPreviewController(docPreviewService)
	authController := controllers.NewAuthController(forwarderService)

	api := router.Group("/api/v1")
	{
		api.POST("/pdf-generator", pdfController.Generate)
		api.POST("/convert/mbl", docConvertController.ConvertMBL)
		api.POST("/preview/hbl", docPreviewController.PreviewHBL)
		api.PUT("/hbl/:hbl_number", docPreviewController.UpdateHBL)
	}

	usersAPI := router.Group("/api/users")
	{
		usersAPI.POST("/signup", authController.Signup)
		usersAPI.POST("/login", authController.Login)
		usersAPI.POST("/logout", authController.Logout)
	}
}
