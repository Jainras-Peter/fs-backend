package controllers

import (
	"fs-backend/models"
	"fs-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetForwarderDetails(c *gin.Context)
	UpdateForwarderDetails(c *gin.Context)
}

type authController struct {
	forwarderService services.ForwarderService
}

func NewAuthController(s services.ForwarderService) AuthController {
	return &authController{forwarderService: s}
}

func (ac *authController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ac *authController) Signup(c *gin.Context) {
	var req models.Forwarder
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := ac.forwarderService.Signup(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Signup successful"})
}

func (ac *authController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	token, err := ac.forwarderService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"token":    token,
		"username": req.Username,
	})
}

func (ac *authController) GetForwarderDetails(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	forwarder, err := ac.forwarderService.GetForwarderDetails(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if forwarder == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Forwarder not found"})
		return
	}

	c.JSON(http.StatusOK, forwarder)
}

func (ac *authController) UpdateForwarderDetails(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	var req models.Forwarder
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := ac.forwarderService.UpdateForwarderDetails(c.Request.Context(), username, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Forwarder details updated successfully"})
}
