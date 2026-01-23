package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceResponse struct {
	ID          uint   `json:"id"`
	ServiceName string `json:"service_name"`
}

type ServiceRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
}

// GetServiceHandler retrieves a service by ID
// @Summary      Get service by ID
// @Description  Retrieve a service by its ID
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Service ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404  {object}  structs.ErrorResponse
// @Router       /services/{id} [get]
func GetServiceHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	service, err := models.GetServiceByID(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Service not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[ServiceResponse](ServiceResponse{
		ID:          service.ID,
		ServiceName: service.ServiceName,
	}, "Service retrieved successfully"))
}

// ListServicesHandler retrieves all services
// @Summary      Get all services
// @Description  Retrieve all services
// @Tags         services
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /services [get]
func ListServicesHandler(c *gin.Context) {
	services, err := models.ListAllServices(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError("Failed to retrieve services", http.StatusInternalServerError))
		return
	}

	var response []ServiceResponse
	for _, s := range services {
		response = append(response, ServiceResponse{
			ID:          s.ID,
			ServiceName: s.ServiceName,
		})
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[[]ServiceResponse](response, "Services retrieved successfully"))
}

// AddServiceHandler creates a new service
// @Summary      Create a new service
// @Description  Create a new service with the provided name
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        request  body      ServiceRequest  true  "Service Request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /services [post]
func AddServiceHandler(c *gin.Context) {
	var request ServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	service, err := models.CreateService(database.DB, request.ServiceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess[ServiceResponse](ServiceResponse{
		ID:          service.ID,
		ServiceName: service.ServiceName,
	}, "Service created successfully"))
}

// UpdateServiceHandler updates a service
// @Summary      Update a service
// @Description  Update an existing service by its ID
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        id       path      int             true  "Service ID"
// @Param        request  body      ServiceRequest  true  "Service Request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /services/{id} [put]
func UpdateServiceHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	var request ServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	service, err := models.GetServiceByID(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Service not found", http.StatusNotFound))
		return
	}

	service.ServiceName = request.ServiceName
	if err := models.UpdateService(database.DB, service); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError("Failed to update service", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[ServiceResponse](ServiceResponse{
		ID:          service.ID,
		ServiceName: service.ServiceName,
	}, "Service updated successfully"))
}

// DeleteServiceHandler deletes a service
// @Summary      Delete a service
// @Description  Delete a service by its ID
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Service ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /services/{id} [delete]
func DeleteServiceHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := models.DeleteService(database.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError("Failed to delete service", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "Service deleted successfully"))
}
