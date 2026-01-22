package server

import (
	"kora_1/internal/handlers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Database Routes
	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	reserved_name := r.Group("/reserved-name")
	{
		reserved_name.POST("", handlers.CreateReservedNameHandler)
		reserved_name.GET("/:name", handlers.GetReservedNameHandler)
	}

	// services := r.Group("/services")
	// {
	// 	services.GET("/:id")
	// 	services.GET("/")
	// 	services.POST("/")
	// 	services.PATCH("/:id")
	// 	services.DELETE("/:id")
	// }

	form := r.Group("/form")
	{
		form.POST("/", handlers.FormHandler)
		form.GET("/:id", handlers.GetFormWithFieldsHandler)
	}
	

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
		// post multiple fields to a form
		form_fields.POST("/multiple", handlers.CreateMultipleFormFieldsHandler)
		form_fields.PATCH("/:id")
	}

	answers := r.Group("/answers")
	{
		answers.GET("/:id")
		answers.POST("/")
	}

	services := r.Group("/services")
	{
		services.GET("/:id", handlers.GetserviceHandler)
		services.GET("/", handlers.ListServicesHandler)
	}

	r.POST("submission")

	return r
}

// HelloWorldHandler returns a hello world message
// @Summary      Hello World
// @Description  Returns a simple hello world message
// @Tags         general
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       / [get]
func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

// healthHandler returns the health status of the server
// @Summary      Health check
// @Description  Returns the health status of the server and database
// @Tags         general
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
