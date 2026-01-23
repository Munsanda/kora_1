package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubmitFormRequest struct {
	ServicesID *uint               `json:"services_id"`
	CreatedBy  *uint               `json:"created_by"`
	Answers    []models.FormAnswer `json:"answers" binding:"required"`
}

type SubmissionResponse struct {
	ID         uint                `json:"id"`
	ServicesID *uint               `json:"services_id"`
	CreatedBy  *uint               `json:"created_by"`
	CreatedOn  string              `json:"created_on"`
	Answers    []models.FormAnswer `json:"answers"`
}

// SubmitFormHandler creates a new form submission
// @Summary      Submit a form
// @Description  Create a new form submission
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        request  body      SubmitFormRequest  true  "Submission Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /submission [post]
func SubmitFormHandler(c *gin.Context) {
	var request SubmitFormRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	submission := &models.Submission{
		ServicesID: request.ServicesID,
		CreatedBy:  request.CreatedBy,
	}

	if err := models.CreateSubmission(database.DB, submission); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	// Create answers
	for _, ans := range request.Answers {
		ans.SubmissionID = &submission.ID // Link to created submission
		if err := models.CreateFormAnswer(database.DB, &ans); err != nil {
			c.JSON(http.StatusInternalServerError, helpers.NewError("Failed to save answers: "+err.Error(), http.StatusInternalServerError))
			return
		}
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess[models.Submission](*submission, "Form submitted successfully"))
}

// GetSubmissionHandler retrieves a submission by ID
// @Summary      Get submission by ID
// @Description  Retrieve a submission by its ID
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Submission ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /submission/{id} [get]
func GetSubmissionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	submission, err := models.GetSubmission(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Submission not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(submission, "Submission retrieved successfully"))
}

// GetSubmissionsByFormIDHandler retrieves all submissions by service ID (formerly by form ID)
// @Summary      Get submissions by Service ID
// @Description  Retrieve all submissions for a specific service by its Service ID
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        service_id   path      int  true  "Service ID"
// @Success      200       {object}  map[string]interface{}
// @Failure      400,500   {object}  structs.ErrorResponse
// @Router       /submission/service/{service_id} [get]
func GetSubmissionsByFormIDHandler(c *gin.Context) {
	serviceIDStr := c.Param("service_id")
	if serviceIDStr == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("Service ID parameter is required", http.StatusBadRequest))
		return
	}

	serviceID, err := strconv.ParseUint(serviceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	var submissions []models.Submission
	if err := database.DB.Preload("Answers").Where("services_id = ?", serviceID).Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[[]models.Submission](submissions, "Submissions retrieved successfully"))
}
