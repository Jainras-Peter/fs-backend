package routes

import (
    "fs-backend/http/controllers"
    "fs-backend/services"

    "github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, pdfService services.PdfGeneratorService) {
    pdfController := controllers.NewPdfGeneratorController(pdfService)

    api := router.Group("/api/v1")
    {
        api.POST("/pdf-generator", pdfController.Generate)
    }
}
