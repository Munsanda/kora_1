package handlers

import (
	"kora_1/internal/database"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	port int

	db database.Service
}

func (s *Handler) HelloWorldHandler(c *gin.Context) {
	
}
