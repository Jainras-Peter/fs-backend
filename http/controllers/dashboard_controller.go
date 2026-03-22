package controllers

import (
	"fs-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	dashboardService services.DashboardService
}

func NewDashboardController(dashboardService services.DashboardService) *DashboardController {
	return &DashboardController{
		dashboardService: dashboardService,
	}
}

func (c *DashboardController) GetDetails(ctx *gin.Context) {
	details, err := c.dashboardService.GetDashboardDetails(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard details"})
		return
	}

	ctx.JSON(http.StatusOK, details)
}

func (c *DashboardController) DeleteDocument(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	err := c.dashboardService.DeleteDocument(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
