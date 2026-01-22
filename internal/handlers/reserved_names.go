package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
}
