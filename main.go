package main

import (
	"fs-backend/config"
	"fs-backend/connections"
	"fs-backend/repository"
	"fs-backend/routes"
	"fs-backend/services"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize Configuration
	config.Init()
	port := config.GetString("server.port")
	pdfBaseURL := config.GetString("pdf_service.base_url")
	mongoURI := config.GetString("mongo.uri")
	mongoDBName := config.GetString("mongo.database")
	extractionBaseURL := config.GetString("extraction_service.base_url")

	// 2. Initialize MongoDB
	db := connections.ConnectMongo(mongoURI, mongoDBName)

	// 3. Initialize Repositories
	mblRepo := repository.NewMBLRepository(db)
	mblCacheRepo := repository.NewMBLCacheRepository(db)
	hblRepo := repository.NewHBLRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	shipmentRepo := repository.NewShipmentRepository(db)
	shipperRepo := repository.NewShipperRepository(db)

	// 4. Initialize Services (Manual DI)
	pdfService := services.NewPdfGeneratorService(pdfBaseURL)
	docConvertService := services.NewDocumentConvertService(
		extractionBaseURL, mblRepo, mblCacheRepo, bookingRepo, shipmentRepo, shipperRepo,
	)
	docPreviewService := services.NewDocumentPreviewService(
		mblRepo, hblRepo, shipmentRepo, shipperRepo,
	)

	// 5. Initialize Router
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:4200"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	r.Use(cors.New(corsConfig))

	// 6. Register Routes
	routes.RegisterRoutes(r, pdfService, docConvertService, docPreviewService)

	// 7. Start Server
	log.Println("Server starting on " + port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
