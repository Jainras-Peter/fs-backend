package main

import (
    "fs-backend/config"
    "fs-backend/routes"
    "fs-backend/services"
    "log"

    "github.com/gin-gonic/gin"
)

func main() {
    // 1. Initialize Configuration
    config.Init()
    port := config.GetString("server.port")
    pdfBaseURL := config.GetString("pdf_service.base_url")

    // 2. Initialize Dependency Injection (Manual)
    pdfService := services.NewPdfGeneratorService(pdfBaseURL)

    // 3. Initialize Router
    r := gin.Default()

    // 4. Register Routes
    routes.RegisterRoutes(r, pdfService)

    // 5. Start Server
    log.Println("Server starting on " + port)
    if err := r.Run(port); err != nil {
        log.Fatal(err)
    }
}
