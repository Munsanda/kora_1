package handlers

import (
	"encoding/json"
	"fmt"
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FormFieldReference represents a field reference in a form creation request
type FormFieldReference struct {
	FieldsID    uint           `json:"fields_id" binding:"required"`
	Validations datatypes.JSON `json:"validations" swaggertype:"object"`
}

type FormRequest struct {
	Title       string               `json:"title" binding:"required"`
	Description string               `json:"description"`
	Fields      []FormFieldReference `json:"fields"`
	ServiceID   int                  `json:"service_id" binding:"required"`
}

type FormCreateResponse struct {
	ID          uint   `json:"id" example:"1"`
	Title       string `json:"title" example:"Contact Form"`
	Description string `json:"description" example:"A form to collect contact information"`
	ServiceID   int    `json:"service_id" example:"1"`
}

// CreateReservedNameHandler creates a new reserved name
// @Summary      Create a new form name
// @Description  Create a new form name with the provided fields
// @Tags         form
// @Accept       json
// @Produce      json
// @Param        request  body      FormRequest  true  "Form Request"
// @Success      201      {object}  FormCreateResponse
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /form [post]
func FormHandler(c *gin.Context) {
	var request FormRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	newForm, err := models.CreateForm(database.DB, &models.Form{
		Title:       request.Title,
		Description: request.Description,
		ServiceId:   request.ServiceID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	if request.Fields != nil {
		for _, field := range request.Fields {
			models.CreateFormFields(database.DB, &models.FormFields{
				FormID:      newForm.ID,
				FieldsID:    field.FieldsID,
				Validations: field.Validations,
			})
		}
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(request, "Form created successfully"))
}

// patch form handler here...
// UpdateFormStatusHandler updates a form's status by ID
// @Summary      Update form status
// @Description  Update the status of an existing form by its ID
// @Tags         forms
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "Form ID"
// @Param        request  body      UpdateFormStatusRequest  true  "Form Status Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /form/{id}/status [patch]
func UpdateFormStatusHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var formID uint
	if _, err := parseID(id, &formID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	var request UpdateFormStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	statusBool, err := strconv.ParseBool(request.Status)
	err = models.UpdateFormStatus(database.DB, formID, statusBool)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Form status updated successfully",
		"data":    request,
	})
}

func GetFormWithFieldsHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID parameter is required"})
		return
	}

	var formID uint
	if _, err := parseID(id, &formID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	form, err := models.GetForm(database.DB, formID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Form not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Form retrieved successfully",
		"data":    form,
	})
}

type UpdateFormStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type FormFieldRequest struct {
	FormID      uint           `json:"form_id" binding:"required"`
	FieldsID    uint           `json:"fields_id" binding:"required"`
	Validations datatypes.JSON `json:"validations" binding:"required" swaggertype:"object"`
}

// FormFieldResponse is a Swagger-friendly representation of FormFields (without gorm.Model)
type FormFieldResponse struct {
	ID          uint                   `json:"id" example:"1"`
	CreatedAt   string                 `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   string                 `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *string                `json:"deleted_at,omitempty"`
	FormID      uint                   `json:"form_id" example:"1"`
	FieldsID    uint                   `json:"fields_id" example:"1"`
	Validations map[string]interface{} `json:"validations" swaggertype:"object"`
}

// GroupFieldResponse represents a field in a group with complete field information
type GroupFieldResponse struct {
	FormFieldID uint                   `json:"form_field_id" example:"1"`
	FormID      uint                   `json:"form_id" example:"1"`
	FieldID     uint                   `json:"field_id" example:"1"`
	GroupID     *uint                  `json:"group_id,omitempty" example:"2"`
	Label       string                 `json:"label" example:"First Name"`
	Type        string                 `json:"type" example:"text"`
	Meta        map[string]interface{} `json:"meta" swaggertype:"object"`
	IsRequired  bool                   `json:"is_required" example:"true"`
	Validations map[string]interface{} `json:"validations" swaggertype:"object"`
	CreatedAt   string                 `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   string                 `json:"updated_at" example:"2024-01-01T00:00:05Z"`
}

// CreateFormFieldsHandler creates a form field association
// @Summary      Create form field
// @Description  Create an association between a form and a field with validations
// @Tags         form-fields
// @Accept       json
// @Produce      json
// @Param        request  body      FormFieldRequest  true  "Form Field Request"
// @Success      200      {object}  object
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /form_fields [post]
func CreateFormFieldsHandler(c *gin.Context) {
	var request FormFieldRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	formField := &models.FormFields{
		FormID:      request.FormID,
		FieldsID:    request.FieldsID,
		Validations: request.Validations,
	}
	err := models.CreateFormFields(database.DB, formField)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	// Convert validations from datatypes.JSON to map[string]interface{}
	var validationsMap map[string]interface{}
	if formField.Validations != nil {
		if err := json.Unmarshal(formField.Validations, &validationsMap); err != nil {
			validationsMap = make(map[string]interface{})
		}
	}

	response := FormFieldResponse{
		ID:          formField.ID,
		CreatedAt:   formField.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   formField.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		DeletedAt:   nil,
		FormID:      formField.FormID,
		FieldsID:    formField.FieldsID,
		Validations: validationsMap,
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(response, "Form fields created successfully"))
}

// CreateMultipleFormFieldsHandler creates multiple form field associations
// @Summary      Create multiple form fields
// @Description  Create multiple associations between a form and fields with validations
// @Tags         form-fields
// @Accept       json
// @Produce      json
// @Param        request  body      []FormFieldRequest  true  "Form Field Requests"
// @Success      202      {object}  object
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /form_fields/multiple [post]
func CreateMultipleFormFieldsHandler(c *gin.Context) {
	var requests []FormFieldRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	var responses []FormFieldResponse
	for _, request := range requests {
		formField := &models.FormFields{
			FormID:      request.FormID,
			FieldsID:    request.FieldsID,
			Validations: request.Validations,
		}
		if err := models.CreateFormFields(database.DB, formField); err != nil {
			c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
			return
		}

		// Convert validations from datatypes.JSON to map[string]interface{}
		var validationsMap map[string]interface{}
		if formField.Validations != nil {
			if err := json.Unmarshal(formField.Validations, &validationsMap); err != nil {
				validationsMap = make(map[string]interface{})
			}
		}

		responses = append(responses, FormFieldResponse{
			ID:          formField.ID,
			CreatedAt:   formField.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   formField.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			DeletedAt:   nil,
			FormID:      formField.FormID,
			FieldsID:    formField.FieldsID,
			Validations: validationsMap,
		})
	}

	c.JSON(http.StatusAccepted, helpers.NewSuccess(responses, "Form fields created successfully"))
}

// Field Handlers

type FieldRequest struct {
	Label      string         `json:"label" binding:"required"`
	Type       string         `json:"type" binding:"required"`
	Meta       datatypes.JSON `json:"meta"`
	IsRequired bool           `json:"is_required"`
	// Label string         `json:"label" binding:"required" example:"First Name"`
	// Type  string         `json:"type" binding:"required" example:"text"`
	// Meta  datatypes.JSON `json:"meta" swaggertype:"object"`
}

// FieldResponse is a Swagger-friendly representation of Fields (without gorm.Model)
type FieldResponse struct {
	ID        uint                   `json:"id" example:"1"`
	CreatedAt string                 `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string                 `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt *string                `json:"deleted_at,omitempty"`
	Label     string                 `json:"label" example:"First Name"`
	Type      string                 `json:"type" example:"text"`
	Meta      map[string]interface{} `json:"meta" swaggertype:"object"`
}

// FormResponse is a Swagger-friendly representation of Form (without gorm.Model)
type FormResponse struct {
	ID          uint    `json:"id" example:"1"`
	CreatedAt   string  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   string  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
	Title       string  `json:"title" example:"Contact Form"`
	Description string  `json:"description" example:"A form to collect contact information"`
	ServiceId   int     `json:"service_id" example:"1"`
	Status      int     `json:"status" example:"1"`
	Version     int     `json:"version" example:"1"`
}

// GroupResponse is a Swagger-friendly representation of Group (without gorm.Model)
type GroupResponse struct {
	ID        uint    `json:"id" example:"1"`
	CreatedAt string  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt *string `json:"deleted_at,omitempty"`
	GroupName string  `json:"group_name" example:"Personal Information"`
}

// CreateFieldHandler creates a new field
// @Summary      Create a new field
// @Description  Create a new field with label, type, and metadata
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        request  body      FieldRequest  true  "Field Request"
// @Success      201      {object}  object
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /field [post]
func CreateFieldHandler(c *gin.Context) {
	var request FieldRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	err := models.CreateFields(database.DB, &models.Field{
		Label: request.Label,
		Type:  request.Type,
		Meta:  request.Meta,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess(request, "Field created successfully"))
}

// GetFieldHandler retrieves a field by ID
// @Summary      Get field by ID
// @Description  Retrieve a field by its ID
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Field ID"
// @Success      200  {object}  object
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Router       /field/{id} [get]
func GetFieldHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("ID parameter is required", http.StatusBadRequest))
		return
	}

	var fieldID uint
	if _, err := parseID(id, &fieldID); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID format", http.StatusBadRequest))
		return
	}

	field, err := models.GetFields(database.DB, fieldID)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Field not found", http.StatusNotFound))
		return
	}

	// Convert model to response type
	metaMap := make(map[string]interface{})
	if len(field.Meta) > 0 {
		if err := json.Unmarshal(field.Meta, &metaMap); err != nil {
			metaMap = make(map[string]interface{})
		}
	}

	fieldResponse := FieldResponse{
		ID:        field.ID,
		CreatedAt: field.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: field.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Label:     field.Label,
		Type:      field.Type,
		Meta:      metaMap,
	}
	if field.DeletedAt.Valid {
		deletedAt := field.DeletedAt.Time.Format("2006-01-02T15:04:05Z")
		fieldResponse.DeletedAt = &deletedAt
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(fieldResponse, "Field retrieved successfully"))
}

// UpdateFieldHandler updates a field by ID
// @Summary      Update field
// @Description  Update an existing field by its ID
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id       path      int           true  "Field ID"
// @Param        request  body      FieldRequest  true  "Field Request"
// @Success      200      {object}  object
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /field/{id} [patch]
func UpdateFieldHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("ID parameter is required", http.StatusBadRequest))
		return
	}

	var fieldID uint
	if _, err := parseID(id, &fieldID); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID format", http.StatusBadRequest))
		return
	}

	var request FieldRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	err := models.UpdateFields(database.DB, &models.Field{
		Model: gorm.Model{ID: fieldID},
		Label: request.Label,
		Type:  request.Type,
		Meta:  request.Meta,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(request, "Field updated successfully"))
}

// DeleteFieldHandler deletes a field by ID
// @Summary      Delete field
// @Description  Delete a field by its ID
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Field ID"
// @Success      200  {object}  object
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /field/{id} [delete]
func DeleteFieldHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("ID parameter is required", http.StatusBadRequest))
		return
	}

	var fieldID uint
	if _, err := parseID(id, &fieldID); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID format", http.StatusBadRequest))
		return
	}

	err := models.DeleteFields(database.DB, fieldID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[interface{}](nil, "Field deleted successfully"))
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

type AddFieldsToGroupRequest struct {
	FormID   uint   `json:"form_id" binding:"required"`
	GroupID  uint   `json:"group_id" binding:"required"`
	FieldIDs []uint `json:"field_ids" binding:"required"`
}

// CreateGroupHandler creates a new group
// @Summary      Create a new group
// @Description  Create a new group with a group name
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        request  body      GroupRequest  true  "Group Request"
// @Success      201      {object}  object
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /groups [post]
func CreateGroupHandler(c *gin.Context) {
	var request GroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	err := models.CreateGroup(database.DB, &models.Group{
		GroupName: request.GroupName,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess(request, "Group created successfully"))
}

// GetGroupByIDHandler retrieves a group by ID
// @Summary      Get group by ID
// @Description  Retrieve a group by its ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {object}  object
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Router       /groups/{id} [get]
func GetGroupByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("ID parameter is required", http.StatusBadRequest))
		return
	}

	var groupID uint
	if _, err := parseID(id, &groupID); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID format", http.StatusBadRequest))
		return
	}

	group, err := models.GetGroupByID(database.DB, groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Group not found", http.StatusNotFound))
		return
	}

	// Convert model to response type
	groupResponse := GroupResponse{
		ID:        group.ID,
		CreatedAt: group.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: group.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		GroupName: group.GroupName,
	}

	// write a helpfull comment here explaining the purpose of this code block
	// This code block checks if the DeletedAt field of the group is valid (i.e., if the group has been deleted).
	if group.DeletedAt.Valid {
		deletedAt := group.DeletedAt.Time.Format("2006-01-02T15:04:05Z")
		groupResponse.DeletedAt = &deletedAt
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(groupResponse, "Group retrieved successfully"))
}

// GetAllGroupsHandler retrieves all groups
// @Summary      Get all groups
// @Description  Retrieve all groups
// @Tags         groups
// @Accept       json
// @Produce      json
// @Success      200  {object}  object
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /groups [get]
func GetAllGroupsHandler(c *gin.Context) {
	groups, err := models.GetAllGroups(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	// Convert models to response types
	groupsResponse := make([]GroupResponse, len(groups))
	for i, group := range groups {
		groupResp := GroupResponse{
			ID:        group.ID,
			CreatedAt: group.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: group.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			GroupName: group.GroupName,
		}
		if group.DeletedAt.Valid {
			deletedAt := group.DeletedAt.Time.Format("2006-01-02T15:04:05Z")
			groupResp.DeletedAt = &deletedAt
		}
		groupsResponse[i] = groupResp
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(groupsResponse, "Groups retrieved successfully"))
}

// UpdateGroupHandler updates a group by ID
// @Summary      Update group
// @Description  Update an existing group by its ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        id       path      int           true  "Group ID"
// @Param        request  body      GroupRequest  true  "Group Request"
// @Success      200      {object}  object
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /groups/{id} [patch]
func UpdateGroupHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("ID parameter is required", http.StatusBadRequest))
		return
	}

	var groupID uint
	if _, err := parseID(id, &groupID); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID format", http.StatusBadRequest))
		return
	}

	var request GroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	err := models.UpdateGroup(database.DB, &models.Group{
		Model:     gorm.Model{ID: groupID},
		GroupName: request.GroupName,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(request, "Group updated successfully"))
}

// DeleteGroupHandler deletes a group by ID
// @Summary      Delete group
// @Description  Delete a group by its ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {object}  object
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /groups/{id} [delete]
func DeleteGroupHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("ID parameter is required", http.StatusBadRequest))
		return
	}

	var groupID uint
	if _, err := parseID(id, &groupID); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID format", http.StatusBadRequest))
		return
	}

	err := models.DeleteGroup(database.DB, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[interface{}](nil, "Group deleted successfully"))
}

// CreateGroupFieldsHandler adds fields to a group in a specific form
// @Summary      Add fields to group
// @Description  Add fields to a group within a specific form
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        request  body      AddFieldsToGroupRequest  true  "Add Fields to Group Request"
// @Success      200      {object}  object
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /groups/add-fields [post]
func CreateGroupFieldsHandler(c *gin.Context) {
	var request AddFieldsToGroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	err := models.AddFieldsToGroup(database.DB, request.FormID, request.GroupID, request.FieldIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(request, "Fields added to group successfully"))
}

// GetGroupFieldsHandler retrieves all fields for a group in a specific form
// @Summary      Get group fields
// @Description  Retrieve all fields for a group within a specific form
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        form_id   query     int  true  "Form ID"
// @Param        group_id  query     int  true  "Group ID"
// @Success      200       {object}  object
// @Failure      400       {object}  structs.ErrorResponse
// @Failure      404       {object}  structs.ErrorResponse
// @Failure      500       {object}  structs.ErrorResponse
// @Router       /groups/get-fields [get]
func GetGroupFieldsHandler(c *gin.Context) {
	formIDStr := c.Query("form_id")
	groupIDStr := c.Query("group_id")

	if formIDStr == "" || groupIDStr == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("form_id and group_id query parameters are required", http.StatusBadRequest))
		return
	}

	var formID uint
	if _, err := parseID(formIDStr, &formID); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid form_id format", http.StatusBadRequest))
		return
	}

	var groupID uint
	if _, err := parseID(groupIDStr, &groupID); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid group_id format", http.StatusBadRequest))
		return
	}

	formFields, err := models.GetAllGroupFields(database.DB, formID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	// Transform FormFields to response format with complete Field details
	var responses []GroupFieldResponse
	for _, formField := range formFields {
		// Get the associated Field
		field, err := models.GetFields(database.DB, formField.FieldsID)
		if err != nil {
			// Skip fields that can't be found
			continue
		}

		// Convert validations from datatypes.JSON to map[string]interface{}
		var validationsMap map[string]interface{}
		if formField.Validations != nil {
			if err := json.Unmarshal(formField.Validations, &validationsMap); err != nil {
				validationsMap = make(map[string]interface{})
			}
		}

		// Convert meta from datatypes.JSON to map[string]interface{}
		var metaMap map[string]interface{}
		if len(field.Meta) > 0 {
			if err := json.Unmarshal(field.Meta, &metaMap); err != nil {
				metaMap = make(map[string]interface{})
			}
		}

		response := GroupFieldResponse{
			FormFieldID: formField.ID,
			FormID:      formField.FormID,
			FieldID:     formField.FieldsID,
			GroupID:     formField.GroupID,
			Label:       field.Label,
			Type:        field.Type,
			Meta:        metaMap,
			IsRequired:  field.IsRequired,
			Validations: validationsMap,
			CreatedAt:   formField.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   formField.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(responses, "Group fields retrieved successfully"))
}
