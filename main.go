package main

import (
	"fs-backend/config"
	"fs-backend/connections"
	"fs-backend/http/controllers"
	"fs-backend/repository"
	"fs-backend/routes"
	"fs-backend/services"
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize Configuration
	config.Init()
	port := config.GetString("server.port")
	if port != "" && !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
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
	hblDocRepo := repository.NewHBLDocRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	shipmentRepo := repository.NewShipmentRepository(db)
	shipperRepo := repository.NewShipperRepository(db)
	forwarderRepo := repository.NewForwarderRepository(db)

	// 4. Initialize Services (Manual DI)
	pdfService := services.NewPdfGeneratorService(pdfBaseURL)
	pdfSaveService := services.NewPdfSaveService(hblDocRepo)
	docConvertService := services.NewDocumentConvertService(
		extractionBaseURL, mblRepo, mblCacheRepo, bookingRepo, shipmentRepo, shipperRepo,
	)
	docPreviewService := services.NewDocumentPreviewService(
		mblRepo, hblRepo, shipmentRepo, shipperRepo, mblCacheRepo,
	)
	bookingService := services.NewBookingService(shipperRepo, bookingRepo, shipmentRepo)
	shipmentService := services.NewShipmentService(shipmentRepo, bookingRepo, shipperRepo)
	dashboardService := services.NewDashboardService(hblDocRepo, hblRepo)
	forwarderService := services.NewForwarderService(forwarderRepo)

	// Initialize Controllers
	bookingController := controllers.NewBookingController(bookingService)
	shipmentController := controllers.NewShipmentController(shipmentService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	authController := controllers.NewAuthController(forwarderService)
	infoToDocRepo := repository.NewInfoToDocRepository(db)
	infoToDocService := services.NewInfoToDocService(infoToDocRepo, hblDocRepo, pdfService)
	infoToDocController := controllers.NewInfoToDocController(infoToDocService)

	// 5. Initialize Router
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:4200",
		"https://freightdocs-one.vercel.app",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	r.Use(cors.New(corsConfig))

	// 6. Register Routes
	routes.RegisterRoutes(r, pdfService, pdfSaveService, docConvertService, docPreviewService, bookingController, shipmentController, dashboardController, authController, infoToDocController)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "server is alive!")
	})

	r.HEAD("/health", func(c *gin.Context) {
	c.Status(200)
})

	// 7. Start Server
	log.Println("Server starting on " + port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
