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

	fields := r.Group("/field")
	{
		fields.GET("/:id", handlers.GetFieldHandler)
		fields.POST("/", handlers.CreateFieldHandler)
		fields.PATCH("/:id", handlers.UpdateFieldHandler)
		fields.DELETE("/:id", handlers.DeleteFieldHandler)
	}

	groups := r.Group("/groups")
	{
		groups.GET("/:id", handlers.GetGroupByIDHandler)
		groups.GET("/", handlers.GetAllGroupsHandler)
		groups.POST("/", handlers.CreateGroupHandler)
		groups.PATCH("/:id", handlers.UpdateGroupHandler)
		groups.DELETE("/:id", handlers.DeleteGroupHandler)
	}

	form_fields := r.Group("/form_fields")
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
