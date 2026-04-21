package routes

import (
	"fs-backend/http/controllers"
	"fs-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, pdfService services.PdfGeneratorService, pdfSaveService services.PdfSaveService, docConvertService services.DocumentConvertService, docPreviewService services.DocumentPreviewService, bookingController *controllers.BookingController, shipmentController *controllers.ShipmentController, dashboardController *controllers.DashboardController, authController controllers.AuthController, infoToDocController *controllers.InfoToDocController) {
	pdfSaveController := controllers.NewPdfSaveController(pdfSaveService)
	pdfController := controllers.NewPdfGeneratorController(pdfService, pdfSaveController)
	docConvertController := controllers.NewDocumentConvertController(docConvertService)
	docPreviewController := controllers.NewDocumentPreviewController(docPreviewService)

	api := router.Group("/api/v1")
	{
		api.POST("/pdf-generator", pdfController.Generate)
		api.POST("/convert/mbl", docConvertController.ConvertMBL)
		api.POST("/preview/hbl", docPreviewController.PreviewHBL)
		api.PUT("/hbl/:hbl_number", docPreviewController.UpdateHBL)
		api.POST("/hbl-docs/download-archive", controllers.DownloadHBLDocsArchive)
	}

	usersAPI := router.Group("/api/users")
	{
		usersAPI.POST("/signup", authController.Signup)
		usersAPI.POST("/login", authController.Login)
		usersAPI.POST("/logout", authController.Logout)
		usersAPI.GET("/getforwarderdetails", authController.GetForwarderDetails)
		usersAPI.PUT("/updateforwarderdetails/:username", authController.UpdateForwarderDetails)
	}

	bookingApi := router.Group("/api/booking")
	{
		//Shippers
		bookingApi.POST("/addshipper", bookingController.AddShipper)
		bookingApi.GET("/shipperlist", bookingController.GetShipperList)
		bookingApi.PUT("/updateshipper/:id", bookingController.UpdateShipper)
		bookingApi.DELETE("/deleteshipper/:id", bookingController.DeleteShipper)

		//Status
		bookingApi.GET("/statusdetails", bookingController.GetStatusDetails)
		bookingApi.PUT("/updatestatus/:id", bookingController.UpdateStatus)

		//Shipments
		bookingApi.GET("/shipments", shipmentController.GetShipmentList)
		bookingApi.POST("/shipments", shipmentController.CreateShipment)
		bookingApi.PUT("/shipments/:id", shipmentController.UpdateShipment)
		bookingApi.DELETE("/shipments/:id", shipmentController.DeleteShipment)

		//Sync MBL Number
		bookingApi.POST("/syncBooking", bookingController.SyncBooking)
	}

	dashboardApi := router.Group("/api/dashboard")
	{
		dashboardApi.GET("/details", dashboardController.GetDetails)
		dashboardApi.DELETE("/delete/:id", dashboardController.DeleteDocument)
	}

	infotodocApi := router.Group("/api/infotodoc")
	{
		infotodocApi.POST("/billoflading", infoToDocController.HandleBillOfLading)
	}
}
