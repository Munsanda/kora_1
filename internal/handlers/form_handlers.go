package handlers

import (
	"fmt"
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	_ "kora_1/internal/structs" // Required for Swagger annotations
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FormRequest struct {
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description"`
	Fields      []models.FormFields `json:"fields"`
	ServiceID   int                 `json:"service_id" binding:"required"`
}

// FormHandler creates a new form
// @Summary      Create a new form
// @Description  Create a new form with title, description, fields, and service ID
// @Tags         forms
// @Accept       json
// @Produce      json
// @Param        request  body      FormRequest  true  "Form Request"
// @Success      200      {object}  structs.SuccessResponse
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /form [post]
func FormHandler(c *gin.Context) {
	var request FormRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	newForm, err := models.CreateForm(database.DB, &models.Form{
		Title:       request.Title,
		Description: request.Description,
		ServiceId:   request.ServiceID,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	// TODO:: Return PACRA Response
	c.JSON(200, gin.H{
		"message": "Hello, World!",
		"data":    request,
	})

	if request.Fields != nil {
		for _, field := range request.Fields {
			models.CreateFormFields(database.DB, &models.FormFields{
				FormID:      newForm.ID,
				FieldsID:    field.ID,
				Validations: field.Validations,
			})
		}
	}
}

type FormFieldRequest struct {
	FormID      uint           `json:"form_id" binding:"required"`
	FieldsID    uint           `json:"fields_id" binding:"required"`
	Validations datatypes.JSON `json:"validations" binding:"required"`
}

// CreateFormFieldsHandler creates a form field association
// @Summary      Create form field
// @Description  Create an association between a form and a field with validations
// @Tags         form-fields
// @Accept       json
// @Produce      json
// @Param        request  body      FormFieldRequest  true  "Form Field Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /form_fields [post]
func CreateFormFieldsHandler(c *gin.Context) {
	var request FormFieldRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	models.CreateFormFields(database.DB, &models.FormFields{
		FormID:      request.FormID,
		FieldsID:    request.FieldsID,
		Validations: request.Validations,
	})

	c.JSON(200, gin.H{
		"message": "Form fields created successfully",
		"data":    request,
	})
}

// CreateMultipleFormFieldsHandler creates multiple form field associations
// @Summary      Create multiple form fields
// @Description  Create multiple associations between a form and fields with validations
// @Tags         form-fields
// @Accept       json
// @Produce      json
// @Param        request  body      []FormFieldRequest  true  "Form Field Requests"
// @Success      202      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /form_fields/multiple [post]
func CreateMultipleFormFieldsHandler(c *gin.Context) {
	var requests []FormFieldRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	for _, request := range requests {
		models.CreateFormFields(database.DB, &models.FormFields{
			FormID:      request.FormID,
			FieldsID:    request.FieldsID,
			Validations: request.Validations,
		})
	}

	c.JSON(http.StatusAccepted, helpers.NewSuccess(requests, "Form fields created successfully"))
}

// Field Handlers

type FieldRequest struct {
	Label string         `json:"label" binding:"required"`
	Type  string         `json:"type" binding:"required"`
	Meta  datatypes.JSON `json:"meta"`
}

// CreateFieldHandler creates a new field
// @Summary      Create a new field
// @Description  Create a new field with label, type, and metadata
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        request  body      FieldRequest  true  "Field Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /field [post]
func CreateFieldHandler(c *gin.Context) {
	var request FieldRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := models.CreateFields(database.DB, &models.Fields{
		Label: request.Label,
		Type:  request.Type,
		Meta:  request.Meta,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Field created successfully",
		"data":    request,
	})
}

// GetFieldHandler retrieves a field by ID
// @Summary      Get field by ID
// @Description  Retrieve a field by its ID
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Field ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /field/{id} [get]
func GetFieldHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var fieldID uint
	if _, err := parseID(id, &fieldID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	field, err := models.GetFields(database.DB, fieldID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Field not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Field retrieved successfully",
		"data":    field,
	})
}

// UpdateFieldHandler updates a field by ID
// @Summary      Update field
// @Description  Update an existing field by its ID
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id       path      int           true  "Field ID"
// @Param        request  body      FieldRequest  true  "Field Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /field/{id} [patch]
func UpdateFieldHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var fieldID uint
	if _, err := parseID(id, &fieldID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	var request FieldRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := models.UpdateFields(database.DB, &models.Fields{
		Model: gorm.Model{ID: fieldID},
		Label: request.Label,
		Type:  request.Type,
		Meta:  request.Meta,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Field updated successfully",
		"data":    request,
	})
}

// DeleteFieldHandler deletes a field by ID
// @Summary      Delete field
// @Description  Delete a field by its ID
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Field ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /field/{id} [delete]
func DeleteFieldHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var fieldID uint
	if _, err := parseID(id, &fieldID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	err := models.DeleteFields(database.DB, fieldID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Field deleted successfully",
	})
}

func parseID(id string, fieldID *uint) (bool, error) {
	numID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return false, fmt.Errorf("invalid ID format: %w", err)
	}
	*fieldID = uint(numID)
	return true, nil
}

// Group Handlers

type GroupRequest struct {
	GroupName string `json:"group_name" binding:"required"`
}

// CreateGroupHandler creates a new group
// @Summary      Create a new group
// @Description  Create a new group with a group name
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        request  body      GroupRequest  true  "Group Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /groups [post]
func CreateGroupHandler(c *gin.Context) {
	var request GroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := models.CreateGroup(database.DB, &models.Group{
		GroupName: request.GroupName,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Group created successfully",
		"data":    request,
	})
}

// GetGroupByIDHandler retrieves a group by ID
// @Summary      Get group by ID
// @Description  Retrieve a group by its ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /groups/{id} [get]
func GetGroupByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var groupID uint
	if _, err := parseID(id, &groupID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	group, err := models.GetGroupByID(database.DB, groupID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Group retrieved successfully",
		"data":    group,
	})
}

// GetAllGroupsHandler retrieves all groups
// @Summary      Get all groups
// @Description  Retrieve all groups
// @Tags         groups
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /groups [get]
func GetAllGroupsHandler(c *gin.Context) {
	groups, err := models.GetAllGroups(database.DB)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Groups retrieved successfully",
		"data":    groups,
	})
}

// UpdateGroupHandler updates a group by ID
// @Summary      Update group
// @Description  Update an existing group by its ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        id       path      int           true  "Group ID"
// @Param        request  body      GroupRequest  true  "Group Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /groups/{id} [patch]
func UpdateGroupHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var groupID uint
	if _, err := parseID(id, &groupID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	var request GroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := models.UpdateGroup(database.DB, &models.Group{
		Model:     gorm.Model{ID: groupID},
		GroupName: request.GroupName,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Group updated successfully",
		"data":    request,
	})
}

// DeleteGroupHandler deletes a group by ID
// @Summary      Delete group
// @Description  Delete a group by its ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /groups/{id} [delete]
func DeleteGroupHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var groupID uint
	if _, err := parseID(id, &groupID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	err := models.DeleteGroup(database.DB, groupID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Group deleted successfully",
	})
}
