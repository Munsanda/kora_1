package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SubmissionResponse is a Swagger-friendly representation of Submission
type SubmissionResponse struct {
	ID              uint                 `json:"id" example:"1"`
	CreatedAt       string               `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       string               `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt       *string              `json:"deleted_at,omitempty"`
	FormID          uint                 `json:"form_id" example:"1"`
	FormVersion     int                  `json:"form_version" example:"1"`
	CreatedByUserId uint                 `json:"created_by_user_id" example:"1"`
	Status          string               `json:"status" example:"submitted"`
	Answers         []FormAnswerResponse `json:"answers,omitempty"`
}

// FormAnswerResponse is a Swagger-friendly representation of FormAnswer
type FormAnswerResponse struct {
	ID           uint    `json:"id" example:"1"`
	CreatedAt    string  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    string  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt    *string `json:"deleted_at,omitempty"`
	SubmissionID uint    `json:"submission_id" example:"1"`
	QuestionID   uint    `json:"question_id" example:"1"`
	Answer       string  `json:"answer" example:"John Doe"`
	AnswerJSON   string  `json:"answer_json,omitempty"`
}

// SubmissionCreateSuccessResponse is a success response containing a submission
type SubmissionCreateSuccessResponse struct {
	Status  bool               `json:"status"`
	Message string             `json:"message,omitempty"`
	Data    SubmissionResponse `json:"data,omitempty"`
}

// SubmissionGetSuccessResponse is a success response containing a submission
type SubmissionGetSuccessResponse struct {
	Status  bool               `json:"status"`
	Message string             `json:"message,omitempty"`
	Data    SubmissionResponse `json:"data,omitempty"`
}

// SubmissionListSuccessResponse is a success response containing a list of submissions
type SubmissionListSuccessResponse struct {
	Status  bool                 `json:"status"`
	Message string               `json:"message,omitempty"`
	Data    []SubmissionResponse `json:"data,omitempty"`
}

type SubmitFormRequest struct {
	FormID  uint                `json:"form_id" binding:"required" example:"1"`
	Answers []models.FormAnswer `json:"answers" binding:"required"`
}

// SubmitFormHandler creates a new form submission
// @Summary      Submit a form
// @Description  Create a new form submission with form ID and answers
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        request  body      SubmitFormRequest  true  "Submission Request"
// @Success      201      {object}  SubmissionCreateSuccessResponse
// @Failure      400      {object}  structs.ErrorResponse
// @Failure      500      {object}  structs.ErrorResponse
// @Router       /submission [post]
func SubmitFormHandler(c *gin.Context) {
	var request SubmitFormRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	submission := &models.Submission{
		FormID:  request.FormID,
		Answers: request.Answers,
	}

	err := models.CreateSubmission(database.DB, submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess(submission, "Form submitted successfully"))
}

// GetSubmissionHandler retrieves a submission by ID
// @Summary      Get submission by ID
// @Description  Retrieve a submission by its ID
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Submission ID"
// @Success      200  {object}  SubmissionGetSuccessResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      404  {object}  structs.ErrorResponse
// @Router       /submission/{id} [get]
func GetSubmissionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("ID parameter is required", http.StatusBadRequest))
		return
	}

	parsedID, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID parameter", http.StatusBadRequest))
		return
	}

	submission, err := models.GetSubmission(database.DB, uint(parsedID))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError(err.Error(), http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(submission, "Submission retrieved successfully"))
}

// GetSubmissionsByFormIDHandler retrieves all submissions for a specific form
// @Summary      Get submissions by form ID
// @Description  Retrieve all submissions for a specific form by its form ID
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        form_id   path      int  true  "Form ID"
// @Success      200       {object}  SubmissionListSuccessResponse
// @Failure      400       {object}  structs.ErrorResponse
// @Failure      500       {object}  structs.ErrorResponse
// @Router       /submission/form/{form_id} [get]
func GetSubmissionsByFormIDHandler(c *gin.Context) {
	formIDStr := c.Param("form_id")
	if formIDStr == "" {
		c.JSON(http.StatusBadRequest, helpers.NewError("Form ID parameter is required", http.StatusBadRequest))
		return
	}

	parsedFormID, err := strconv.ParseUint(formIDStr, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid Form ID parameter", http.StatusBadRequest))
		return
	}

	submissions, err := models.GetSubmissionsByFormID(database.DB, uint(parsedFormID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(submissions, "Submissions retrieved successfully"))
}
