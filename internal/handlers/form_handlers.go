package handlers

import (
	"fmt"
	"kora_1/internal/database"
	"kora_1/internal/models"
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

// Create form handler
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

//create an array of form fields
func CreateMultipleFormFieldsHandler(c *gin.Context) {
	var requests []FormFieldRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for _, request := range requests {
		models.CreateFormFields(database.DB, &models.FormFields{
			FormID:      request.FormID,
			FieldsID:    request.FieldsID,
			Validations: request.Validations,
		})
	}

	c.JSON(200, gin.H{
		"message": "Multiple form fields created successfully",
		"data":    requests,
	})
}

// Field Handlers

type FieldRequest struct {
	Label string         `json:"label" binding:"required"`
	Type  string         `json:"type" binding:"required"`
	Meta  datatypes.JSON `json:"meta"`
}

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
