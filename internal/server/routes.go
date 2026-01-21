package server

import (
	"kora_1/internal/handlers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	//Database Routes
	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	fields := r.Group("/question")
	{
		fields.GET("/:id")
		fields.POST("/")
		fields.PATCH("/:id")
	}

	forms := r.Group("/forms")
	{
		forms.GET("/:id")
		forms.POST("/", handlers.FormHandler)
		forms.PATCH("/:id")
	}

	form_fields := r.Group("/form_questions")
	{
		form_fields.GET("/:id")
		form_fields.POST("/")
		form_fields.PATCH("/:id")
	}

	answers := r.Group("/answers")
	{
		answers.GET("/:id")
		answers.POST("/")
	}
	r.POST("submission")

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
