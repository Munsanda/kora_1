package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FormGroupRequest struct {
	GroupName string `json:"group_name"`
	GroupSpan int    `json:"group_span"`
	GroupRow  int    `json:"group_row"`
}

type FormGroupResponse struct {
	ID        uint   `json:"id"`
	GroupName string `json:"group_name"`
	GroupSpan int    `json:"group_span"`
	GroupRow  int    `json:"group_row"`
}

// CreateFormGroupHandler creates a new form group
// @Summary      Create form group
// @Description  Create a new form group
// @Tags         form-groups
// @Accept       json
// @Produce      json
// @Param        request  body      FormGroupRequest  true  "Form Group Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /form_groups [post]
func CreateFormGroupHandler(c *gin.Context) {
	var request FormGroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	fg := &models.FormGroup{
		GroupName: request.GroupName,
		GroupSpan: request.GroupSpan,
		GroupRow:  request.GroupRow,
	}

	if err := models.CreateFormGroup(database.DB, fg); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess[FormGroupResponse](FormGroupResponse{
		ID:        fg.ID,
		GroupName: fg.GroupName,
		GroupSpan: fg.GroupSpan,
		GroupRow:  fg.GroupRow,
	}, "Form group created successfully"))
}

// GetFormGroupHandler retrieves a form group by ID
// @Summary      Get form group
// @Description  Retrieve a form group by its ID
// @Tags         form-groups
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Form Group ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404  {object}  structs.ErrorResponse
// @Router       /form_groups/{id} [get]
func GetFormGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	fg, err := models.GetFormGroup(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Form group not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[FormGroupResponse](FormGroupResponse{
		ID:        fg.ID,
		GroupName: fg.GroupName,
		GroupSpan: fg.GroupSpan,
		GroupRow:  fg.GroupRow,
	}, "Form group retrieved successfully"))
}

// GetAllFormGroupsHandler retrieves all form groups
// @Summary      Get all form groups
// @Description  Retrieve all form groups
// @Tags         form-groups
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /form_groups [get]
func GetAllFormGroupsHandler(c *gin.Context) {
	fgs, err := models.GetAllFormGroups(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	var response []FormGroupResponse
	for _, fg := range fgs {
		response = append(response, FormGroupResponse{
			ID:        fg.ID,
			GroupName: fg.GroupName,
			GroupSpan: fg.GroupSpan,
			GroupRow:  fg.GroupRow,
		})
	}
	c.JSON(http.StatusOK, helpers.NewSuccess[[]FormGroupResponse](response, "Form groups retrieved successfully"))
}

// UpdateFormGroupHandler updates a form group
// @Summary      Update form group
// @Description  Update an existing form group by its ID
// @Tags         form-groups
// @Accept       json
// @Produce      json
// @Param        id       path      int               true  "Form Group ID"
// @Param        request  body      FormGroupRequest  true  "Form Group Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /form_groups/{id} [put]
func UpdateFormGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	var request FormGroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	fg, err := models.GetFormGroup(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Form group not found", http.StatusNotFound))
		return
	}

	fg.GroupName = request.GroupName
	fg.GroupSpan = request.GroupSpan
	fg.GroupRow = request.GroupRow

	if err := database.DB.Save(fg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[FormGroupResponse](FormGroupResponse{
		ID:        fg.ID,
		GroupName: fg.GroupName,
		GroupSpan: fg.GroupSpan,
		GroupRow:  fg.GroupRow,
	}, "Form group updated successfully"))
}

// DeleteFormGroupHandler deletes a form group
// @Summary      Delete form group
// @Description  Delete a form group by its ID
// @Tags         form-groups
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Form Group ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /form_groups/{id} [delete]
func DeleteFormGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := database.DB.Delete(&models.FormGroup{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "Form group deleted successfully"))
}
