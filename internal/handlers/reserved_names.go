package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReservedNameRequest struct {
	ReservedName string `json:"reserved_name" binding:"required"`
}

type ReservedNameResponse struct {
	ID           uint   `json:"id"`
	ReservedName string `json:"reserved_name"`
}

// GetReservedNameHandler retrieves reserved names similar to the provided name
// @Summary      Get reserved names
// @Description  Retrieve reserved names that match or are similar to the provided name parameter
// @Tags         reserved-name
// @Accept       json
// @Produce      json
// @Param        name  path      string  true  "Name to search for"
// @Success      200   {object}  map[string]interface{}
// @Failure      400,500   {object}  structs.ErrorResponse
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

	// Helper to convert to response struct if needed, or return directly
	c.JSON(http.StatusOK, helpers.NewSuccess[[]models.ReservedName](reservedNames, "Reserved names retrieved successfully"))
}

// CreateReservedNameHandler creates a new reserved name
// @Summary      Create reserved name
// @Description  Create a new reserved name
// @Tags         reserved-name
// @Accept       json
// @Produce      json
// @Param        request  body      ReservedNameRequest  true  "Reserved Name Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500      {object}  structs.ErrorResponse
// @Router       /reserved-name [post]
func CreateReservedNameHandler(c *gin.Context) {
	var request ReservedNameRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	rn := &models.ReservedName{ReservedName: request.ReservedName}
	err := models.CreateReservedName(database.DB, rn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess[ReservedNameResponse](ReservedNameResponse{
		ID:           rn.ID,
		ReservedName: rn.ReservedName,
	}, "Reserved name created successfully"))
}

// DeleteReservedNameHandler deletes a reserved name
// @Summary      Delete reserved name
// @Description  Delete a reserved name by its ID
// @Tags         reserved-name
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Reserved Name ID"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,500      {object}  structs.ErrorResponse
// @Router       /reserved-name/{id} [delete]
func DeleteReservedNameHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := database.DB.Delete(&models.ReservedName{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "Reserved name deleted"))
}
