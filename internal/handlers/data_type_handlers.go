package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DataTypeRequest struct {
	DataType string `json:"data_type" binding:"required"`
}

type DataTypeResponse struct {
	ID       uint   `json:"id"`
	DataType string `json:"data_type"`
}

// CreateDataTypeHandler creates a new data type
// @Summary      Create data type
// @Description  Create a new data type
// @Tags         data-types
// @Accept       json
// @Produce      json
// @Param        request  body      DataTypeRequest  true  "Data Type Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /data_types [post]
func CreateDataTypeHandler(c *gin.Context) {
	var request DataTypeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	dt := &models.DataType{DataType: request.DataType}
	if err := models.CreateDataType(database.DB, dt); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess[DataTypeResponse](DataTypeResponse{
		ID:       dt.ID,
		DataType: dt.DataType,
	}, "Data type created successfully"))
}

// GetDataTypeHandler retrieves a data type by ID
// @Summary      Get data type
// @Description  Retrieve a data type by its ID
// @Tags         data-types
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Data Type ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404  {object}  structs.ErrorResponse
// @Router       /data_types/{id} [get]
func GetDataTypeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	dt, err := models.GetDataType(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Data type not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[DataTypeResponse](DataTypeResponse{
		ID:       dt.ID,
		DataType: dt.DataType,
	}, "Data type retrieved successfully"))
}

// GetAllDataTypesHandler retrieves all data types
// @Summary      Get all data types
// @Description  Retrieve all data types
// @Tags         data-types
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /data_types [get]
func GetAllDataTypesHandler(c *gin.Context) {
	dts, err := models.GetAllDataTypes(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	var response []DataTypeResponse
	for _, dt := range dts {
		response = append(response, DataTypeResponse{
			ID:       dt.ID,
			DataType: dt.DataType,
		})
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[[]DataTypeResponse](response, "Data types retrieved successfully"))
}

// UpdateDataTypeHandler updates a data type
// @Summary      Update data type
// @Description  Update an existing data type by its ID
// @Tags         data-types
// @Accept       json
// @Produce      json
// @Param        id       path      int              true  "Data Type ID"
// @Param        request  body      DataTypeRequest  true  "Data Type Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /data_types/{id} [put]
func UpdateDataTypeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	var request DataTypeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	dt, err := models.GetDataType(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Data type not found", http.StatusNotFound))
		return
	}

	dt.DataType = request.DataType
	if err := database.DB.Save(dt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[DataTypeResponse](DataTypeResponse{
		ID:       dt.ID,
		DataType: dt.DataType,
	}, "Data type updated successfully"))
}

// DeleteDataTypeHandler deletes a data type
// @Summary      Delete data type
// @Description  Delete a data type by its ID
// @Tags         data-types
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Data Type ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /data_types/{id} [delete]
func DeleteDataTypeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := database.DB.Delete(&models.DataType{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "Data type deleted successfully"))
}
