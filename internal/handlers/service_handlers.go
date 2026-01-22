package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"

	"github.com/gin-gonic/gin"
)

type ServiceRequest struct {
	Name string `json:"name" binding:"required"`
}

func GetserviceHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var serviceID uint
	if _, err := parseID(id, &serviceID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	service, err := models.GetServiceByID(database.DB, serviceID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(200, gin.H{
		"data": helpers.NewSuccess(service, "Service retrieved successfully"),
	})
}

func ListServicesHandler(c *gin.Context) {
	services, err := models.ListAllServices(database.DB)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve services"})
		return
	}

	c.JSON(200, gin.H{
		"data": helpers.NewSuccess(services, "Services retrieved successfully"),
	})
}