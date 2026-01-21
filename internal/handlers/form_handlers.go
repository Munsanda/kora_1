package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/models"

	"github.com/gin-gonic/gin"
)

type FormRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ServiceID   int    `json:"service_id" binding:"required"`
}

type FormFieldRequest struct {
	FormID  uint `json:"form_id" binding:"required"`
	GroupId int
}

func FormHandler(c *gin.Context) {
	var request FormRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	models.CreateForm(database.DB, &models.Form{
		Title:       request.Title,
		Description: request.Description,
		ServiceId:   request.ServiceID,
	})

	// TODO:: Return PACRA Response
	c.JSON(200, gin.H{
		"message": "Hello, World!",
		"data":    request,
	})
}

func CreateFormFieldHandler(c *gin.Context) {

}

func CreateFormFieldsHandler(c *gin.Context) {

}
