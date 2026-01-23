package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FormRequest struct {
	FormName    string               `json:"form_name" binding:"required"`
	Description string               `json:"description"`
	DataTypeID  uint                 `json:"data_type_id" binding:"required"`
	Fields      []FormFieldReference `json:"fields"`
	ServiceID   *uint                `json:"service_id"`
	Status      *bool                `json:"status"`
}

type FormFieldReference struct {
	FieldID     uint   `json:"fields_id" binding:"required"`
	Validations string `json:"validations"`
	FieldSpan   int    `json:"field_span"`
	FieldRow    int    `json:"field_row"`
}

type FormResponse struct {
	ID          uint   `json:"id"`
	FormName    string `json:"form_name"`
	Description string `json:"description"`
	DataTypeID  uint   `json:"data_type_id"`
	ServiceID   *uint  `json:"service_id"`
	Status      *bool  `json:"status"`
}

// FormHandler creates a new form
// @Summary      Create a new form
// @Description  Create a new form with fields
// @Tags         form
// @Accept       json
// @Produce      json
// @Param        request  body      FormRequest  true  "Form Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /form [post]
func FormHandler(c *gin.Context) {
	var request FormRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	newForm := &models.Form{
		FormName:    request.FormName,
		Description: request.Description,
		DataTypeID:  request.DataTypeID,
		ServiceID:   request.ServiceID,
		Status:      request.Status,
	}

	createdForm, err := models.CreateForm(database.DB, newForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	if len(request.Fields) > 0 {
		for _, field := range request.Fields {
			models.CreateFormFields(database.DB, &models.FormFields{
				FormID:     createdForm.ID,
				FieldID:    field.FieldID,
				Validation: field.Validations,
				FieldSpan:  field.FieldSpan,
				FieldRow:   field.FieldRow,
			})
		}
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[FormResponse](FormResponse{
		ID:          createdForm.ID,
		FormName:    createdForm.FormName,
		Description: createdForm.Description,
		DataTypeID:  createdForm.DataTypeID,
		ServiceID:   createdForm.ServiceID,
		Status:      createdForm.Status,
	}, "Form created successfully"))
}

// GetFormWithFieldsHandler retrieves a form by ID
// @Summary      Get form by ID
// @Description  Retrieve a form by its ID
// @Tags         form
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Form ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404  {object}  structs.ErrorResponse
// @Router       /form/{id} [get]
func GetFormWithFieldsHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	form, err := models.GetForm(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError(err.Error(), http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[FormResponse](FormResponse{
		ID:          form.ID,
		FormName:    form.FormName,
		Description: form.Description,
		DataTypeID:  form.DataTypeID,
		ServiceID:   form.ServiceID,
		Status:      form.Status,
	}, "Form retrieved successfully"))
}

// UpdateFormHandler updates a form
// @Summary      Update form
// @Description  Update an existing form by its ID
// @Tags         form
// @Accept       json
// @Produce      json
// @Param        id       path      int          true  "Form ID"
// @Param        request  body      FormRequest  true  "Form Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /form/{id} [put]
func UpdateFormHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	var request FormRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	form, err := models.GetForm(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Form not found", http.StatusNotFound))
		return
	}

	form.FormName = request.FormName
	form.Description = request.Description
	form.DataTypeID = request.DataTypeID
	form.ServiceID = request.ServiceID
	form.Status = request.Status

	if err := models.UpdateForm(database.DB, form); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[FormResponse](FormResponse{
		ID:          form.ID,
		FormName:    form.FormName,
		Description: form.Description,
		DataTypeID:  form.DataTypeID,
		ServiceID:   form.ServiceID,
		Status:      form.Status,
	}, "Form updated successfully"))
}

// DeleteFormHandler deletes a form
// @Summary      Delete form
// @Description  Delete a form by its ID
// @Tags         form
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Form ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /form/{id} [delete]
func DeleteFormHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := models.DeleteForm(database.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "Form deleted successfully"))
}

// Form Field Handlers

type FormFieldRequest struct {
	FormID      uint   `json:"form_id" binding:"required"`
	FieldID     uint   `json:"field_id" binding:"required"`
	Validation  string `json:"validation"`
	FieldSpan   int    `json:"field_span"`
	FieldRow    int    `json:"field_row"`
	FormGroupID *uint  `json:"form_group_id"`
}

type FormFieldResponse struct {
	ID          uint   `json:"id"`
	FormID      uint   `json:"form_id"`
	FieldID     uint   `json:"field_id"`
	Validation  string `json:"validation"`
	FieldSpan   int    `json:"field_span"`
	FieldRow    int    `json:"field_row"`
	FormGroupID *uint  `json:"form_group_id"`
}

// CreateFormFieldsHandler creates a form field association
// @Summary      Create form field
// @Description  Create a form field association
// @Tags         form-fields
// @Accept       json
// @Produce      json
// @Param        request  body      FormFieldRequest  true  "Form Field Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /form_fields [post]
func CreateFormFieldsHandler(c *gin.Context) {
	var request FormFieldRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	ff := &models.FormFields{
		FormID:      request.FormID,
		FieldID:     request.FieldID,
		Validation:  request.Validation,
		FieldSpan:   request.FieldSpan,
		FieldRow:    request.FieldRow,
		FormGroupID: request.FormGroupID,
	}

	if err := models.CreateFormFields(database.DB, ff); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[FormFieldResponse](FormFieldResponse{
		ID:          ff.ID,
		FormID:      ff.FormID,
		FieldID:     ff.FieldID,
		Validation:  ff.Validation,
		FieldSpan:   ff.FieldSpan,
		FieldRow:    ff.FieldRow,
		FormGroupID: ff.FormGroupID,
	}, "Form field created successfully"))
}

// CreateMultipleFormFieldsHandler creates multiple form field associations
// @Summary      Create multiple form fields
// @Description  Create multiple form field associations
// @Tags         form-fields
// @Accept       json
// @Produce      json
// @Param        request  body      []FormFieldRequest  true  "Form Field Requests"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /form_fields/multiple [post]
func CreateMultipleFormFieldsHandler(c *gin.Context) {
	var requests []FormFieldRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	var responses []FormFieldResponse
	for _, req := range requests {
		ff := &models.FormFields{
			FormID:      req.FormID,
			FieldID:     req.FieldID,
			Validation:  req.Validation,
			FieldSpan:   req.FieldSpan,
			FieldRow:    req.FieldRow,
			FormGroupID: req.FormGroupID,
		}
		if err := models.CreateFormFields(database.DB, ff); err != nil {
			c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
			return
		}
		responses = append(responses, FormFieldResponse{
			ID:          ff.ID,
			FormID:      ff.FormID,
			FieldID:     ff.FieldID,
			Validation:  ff.Validation,
			FieldSpan:   ff.FieldSpan,
			FieldRow:    ff.FieldRow,
			FormGroupID: ff.FormGroupID,
		})
	}
	c.JSON(http.StatusOK, helpers.NewSuccess[[]FormFieldResponse](responses, "Multiple form fields created"))
}

// Field Handlers

type FieldRequest struct {
	Label        string `json:"label" binding:"required"`
	DataTypeID   uint   `json:"data_type_id" binding:"required"`
	GroupID      *uint  `json:"group_id"`
	CollectionID *uint  `json:"collection_id"`
	Status       *bool  `json:"status"`
}

type FieldResponse struct {
	ID           uint   `json:"id"`
	Label        string `json:"label"`
	DataTypeID   uint   `json:"data_type_id"`
	GroupID      *uint  `json:"group_id"`
	CollectionID *uint  `json:"collection_id"`
	Status       *bool  `json:"status"`
}

// CreateFieldHandler creates a new field
// @Summary      Create field
// @Description  Create a new field
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        request  body      FieldRequest  true  "Field Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /field [post]
func CreateFieldHandler(c *gin.Context) {
	var request FieldRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	field := &models.Field{
		Label:        request.Label,
		DataTypeID:   request.DataTypeID,
		GroupID:      request.GroupID,
		CollectionID: request.CollectionID,
		Status:       request.Status,
	}

	if err := models.CreateFields(database.DB, field); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess[FieldResponse](FieldResponse{
		ID:           field.ID,
		Label:        field.Label,
		DataTypeID:   field.DataTypeID,
		GroupID:      field.GroupID,
		CollectionID: field.CollectionID,
		Status:       field.Status,
	}, "Field created successfully"))
}

// GetFieldHandler retrieves a field by ID
// @Summary      Get field
// @Description  Retrieve a field by its ID
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Field ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404  {object}  structs.ErrorResponse
// @Router       /field/{id} [get]
func GetFieldHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	field, err := models.GetFields(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Field not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[FieldResponse](FieldResponse{
		ID:           field.ID,
		Label:        field.Label,
		DataTypeID:   field.DataTypeID,
		GroupID:      field.GroupID,
		CollectionID: field.CollectionID,
		Status:       field.Status,
	}, "Field retrieved successfully"))
}

// UpdateFieldHandler updates a field
// @Summary      Update field
// @Description  Update an existing field by its ID
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id       path      int           true  "Field ID"
// @Param        request  body      FieldRequest  true  "Field Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /field/{id} [put]
func UpdateFieldHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	var request FieldRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	field, err := models.GetFields(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Field not found", http.StatusNotFound))
		return
	}

	field.Label = request.Label
	field.DataTypeID = request.DataTypeID
	field.GroupID = request.GroupID
	field.CollectionID = request.CollectionID
	field.Status = request.Status

	if err := models.UpdateFields(database.DB, field); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[FieldResponse](FieldResponse{
		ID:           field.ID,
		Label:        field.Label,
		DataTypeID:   field.DataTypeID,
		GroupID:      field.GroupID,
		CollectionID: field.CollectionID,
		Status:       field.Status,
	}, "Field updated successfully"))
}

// DeleteFieldHandler deletes a field
// @Summary      Delete field
// @Description  Delete a field by its ID
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Field ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /field/{id} [delete]
func DeleteFieldHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := models.DeleteFields(database.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "Field deleted successfully"))
}

// Group Handlers

type GroupRequest struct {
	GroupName string `json:"group_name" binding:"required"`
}

type GroupResponse struct {
	ID        uint   `json:"id"`
	GroupName string `json:"group_name"`
}

// CreateGroupHandler creates a new group
// @Summary      Create group
// @Description  Create a new group
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        request  body      GroupRequest  true  "Group Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /groups [post]
func CreateGroupHandler(c *gin.Context) {
	var request GroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	group := &models.Group{GroupName: request.GroupName}
	if err := models.CreateGroup(database.DB, group); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess[GroupResponse](GroupResponse{
		ID:        group.ID,
		GroupName: group.GroupName,
	}, "Group created successfully"))
}

// GetGroupByIDHandler retrieves a group by ID
// @Summary      Get group
// @Description  Retrieve a group by its ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404  {object}  structs.ErrorResponse
// @Router       /groups/{id} [get]
func GetGroupByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	group, err := models.GetGroupByID(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Group not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[GroupResponse](GroupResponse{
		ID:        group.ID,
		GroupName: group.GroupName,
	}, "Group retrieved successfully"))
}

// GetAllGroupsHandler retrieves all groups
// @Summary      Get all groups
// @Description  Retrieve all groups
// @Tags         groups
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /groups [get]
func GetAllGroupsHandler(c *gin.Context) {
	groups, err := models.GetAllGroups(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	var response []GroupResponse
	for _, g := range groups {
		response = append(response, GroupResponse{ID: g.ID, GroupName: g.GroupName})
	}
	c.JSON(http.StatusOK, helpers.NewSuccess[[]GroupResponse](response, "Groups retrieved successfully"))
}

// UpdateGroupHandler updates a group
// @Summary      Update group
// @Description  Update an existing group by its ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        id       path      int           true  "Group ID"
// @Param        request  body      GroupRequest  true  "Group Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /groups/{id} [put]
func UpdateGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	var request GroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	group, err := models.GetGroupByID(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Group not found", http.StatusNotFound))
		return
	}

	group.GroupName = request.GroupName
	if err := models.UpdateGroup(database.DB, group); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[GroupResponse](GroupResponse{
		ID:        group.ID,
		GroupName: group.GroupName,
	}, "Group updated successfully"))
}

// DeleteGroupHandler deletes a group
// @Summary      Delete group
// @Description  Delete a group by its ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /groups/{id} [delete]
func DeleteGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := models.DeleteGroup(database.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "Group deleted successfully"))
}

// Validation helper for Unmarshal
func parseID(id string, ptr *uint) (bool, error) {
	val, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return false, err
	}
	*ptr = uint(val)
	return true, nil
}
