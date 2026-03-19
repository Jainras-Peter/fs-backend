package controllers

import (
	"fs-backend/repository"
	"fs-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShipmentController struct {
	shipmentService services.ShipmentService
}

func NewShipmentController(shipmentService services.ShipmentService) *ShipmentController {
	return &ShipmentController{
		shipmentService: shipmentService,
	}
}

func (c *ShipmentController) CreateShipment(ctx *gin.Context) {
	var input repository.ShipmentDocument
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.shipmentService.InsertShipment(ctx.Request.Context(), &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shipment"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":     "Shipment created successfully",
		"shipment_id": id,
	})
}

func (c *ShipmentController) GetShipmentList(ctx *gin.Context) {
	shipments, err := c.shipmentService.GetAllShipments(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shipments"})
		return
	}

	ctx.JSON(http.StatusOK, shipments)
}

func (c *ShipmentController) UpdateShipment(ctx *gin.Context) {
	id := ctx.Param("id")

	var updates repository.ShipmentDocument
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.shipmentService.UpdateShipment(ctx.Request.Context(), id, &updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shipment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Shipment updated successfully"})
}

func (c *ShipmentController) DeleteShipment(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.shipmentService.DeleteShipment(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete shipment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Shipment deleted successfully"})
}
