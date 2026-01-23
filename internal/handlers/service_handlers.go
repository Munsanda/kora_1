package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServiceResponse is a Swagger-friendly representation of Service
type ServiceResponse struct {
	ID        uint    `json:"id" example:"1"`
	CreatedAt string  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt *string `json:"deleted_at,omitempty"`
	Name      string  `json:"name" example:"User Service"`
}

// ServiceGetSuccessResponse is a success response containing a service
type ServiceGetSuccessResponse struct {
	Status  bool            `json:"status"`
	Message string          `json:"message,omitempty"`
	Data    ServiceResponse `json:"data,omitempty"`
}

// ServiceListSuccessResponse is a success response containing a list of services
type ServiceListSuccessResponse struct {
	Status  bool              `json:"status"`
	Message string            `json:"message,omitempty"`
	Data    []ServiceResponse `json:"data,omitempty"`
}

// ServiceCreateSuccessResponse is a success response containing a service
type ServiceCreateSuccessResponse struct {
	Status  bool            `json:"status"`
	Message string          `json:"message,omitempty"`
	Data    ServiceResponse `json:"data,omitempty"`
}

type ServiceRequest struct {
	Name string `json:"name" binding:"required" example:"User Service"`
}

// GetServiceHandler retrieves a service by ID
// @Summary      Get service by ID
// @Description  Retrieve a service by its ID
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Service ID"
// @Success      200  {object}  ServiceGetSuccessResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Router       /services/{id} [get]
func GetServiceHandler(c *gin.Context) {
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

// ListServicesHandler retrieves all services
// @Summary      Get all services
// @Description  Retrieve all services
// @Tags         services
// @Accept       json
// @Produce      json
// @Success      200  {object}  ServiceListSuccessResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /services [get]
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

// AddServiceHandler creates a new service
// @Summary      Create a new service
// @Description  Create a new service with the provided name
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        request  body      ServiceRequest  true  "Service Request"
// @Success      201      {object}  ServiceCreateSuccessResponse
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /services [post]
func AddServiceHandler(c *gin.Context) {
	var request ServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	service, err := models.CreateService(database.DB, request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess(service, "Service created successfully"))
}
