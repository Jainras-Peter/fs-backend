package routes

import (
	"fs-backend/http/controllers"
	"fs-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, pdfService services.PdfGeneratorService, docConvertService services.DocumentConvertService, docPreviewService services.DocumentPreviewService, bookingController *controllers.BookingController, shipmentController *controllers.ShipmentController) {
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
}
