package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReservedNameResponse is a Swagger-friendly representation of ReservedName
type ReservedNameResponse struct {
	ID        uint   `json:"id" example:"1"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt *string `json:"deleted_at,omitempty"`
	Name      string `json:"name" example:"admin"`
}

// ReservedNameListSuccessResponse is a success response containing a list of reserved names
type ReservedNameListSuccessResponse struct {
	Status  bool                  `json:"status"`
	Message string                `json:"message,omitempty"`
	Data    []ReservedNameResponse `json:"data,omitempty"`
}

// ReservedNameCreateSuccessResponse is a success response containing a reserved name
type ReservedNameCreateSuccessResponse struct {
	Status  bool                `json:"status"`
	Message string              `json:"message,omitempty"`
	Data    ReservedNameResponse `json:"data,omitempty"`
}

// ReservedNameRequest is the request body for creating a reserved name
type ReservedNameRequest struct {
	Name string `json:"name" binding:"required" example:"admin"`
}

// GetReservedNameHandler retrieves reserved names similar to the provided name
// @Summary      Get reserved names by name
// @Description  Retrieve reserved names that match or are similar to the provided name parameter
// @Tags         reserved-name
// @Accept       json
// @Produce      json
// @Param        name  path      string  true  "Name to search for"
// @Success      200   {object}  ReservedNameListSuccessResponse
// @Failure      400   {object}  structs.ErrorResponse
// @Failure      500   {object}  structs.ErrorResponse
// @Router       /reserved-name/{name} [get]
func GetReservedNameHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("Name parameter is required", http.StatusBadRequest))
		return
	}

	reservedNames, err := models.GetSimilarNames(database.DB, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(reservedNames, "Reserved names retrieved successfully"))
}

// CreateReservedNameHandler creates a new reserved name
// @Summary      Create a new reserved name
// @Description  Create a new reserved name with the provided name
// @Tags         reserved-name
// @Accept       json
// @Produce      json
// @Param        request  body      ReservedNameRequest  true  "Reserved Name Request"
// @Success      201      {object}  ReservedNameCreateSuccessResponse
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /reserved-name [post]
func CreateReservedNameHandler(c *gin.Context) {
	var request models.ReservedName
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	err := models.CreateReservedName(database.DB, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess(request, "Reserved name created successfully"))
}
