package controllers

import (
	"fs-backend/repository"
	"fs-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingController struct {
	bookingService services.BookingService
}

func NewBookingController(bookingService services.BookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
	}
}

func (c *BookingController) AddShipper(ctx *gin.Context) {
	var input repository.ShipperDocument
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.bookingService.AddShipper(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add shipper"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Shipper added successfully",
		"id":      id.Hex(),
	})
}

func (c *BookingController) GetShipperList(ctx *gin.Context) {
	shippers, err := c.bookingService.GetShipperList(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shippers"})
		return
	}

	ctx.JSON(http.StatusOK, shippers)
}

func (c *BookingController) UpdateShipper(ctx *gin.Context) {
	idParam := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.bookingService.UpdateShipper(ctx.Request.Context(), objID, updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shipper"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Shipper updated successfully"})
}

func (c *BookingController) DeleteShipper(ctx *gin.Context) {
	idParam := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = c.bookingService.DeleteShipper(ctx.Request.Context(), objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete shipper"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Shipper deleted successfully"})
}

func (c *BookingController) SyncBooking(ctx *gin.Context) {
	var input struct {
		MBLNumber  string `json:"mbl_number" binding:"required"`
		ShipmentID string `json:"shipment_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.bookingService.SyncBooking(ctx.Request.Context(), input.MBLNumber, input.ShipmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync booking"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Booking synced successfully"})
}
