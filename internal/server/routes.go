package server

import (
	_ "kora_1/docs"
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
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	// Reserved Names
	reservedName := r.Group("/reserved-name")
	{
		reservedName.POST("", handlers.CreateReservedNameHandler)
		reservedName.GET("/:name", handlers.GetReservedNameHandler)
		reservedName.DELETE("/:id", handlers.DeleteReservedNameHandler)
	}

	// Services
	services := r.Group("/services")
	{
		services.GET("/:id", handlers.GetServiceHandler)
		services.GET("/", handlers.ListServicesHandler)
		services.POST("/", handlers.AddServiceHandler)
		services.PUT("/:id", handlers.UpdateServiceHandler)
		services.DELETE("/:id", handlers.DeleteServiceHandler)
	}

	// Forms
	form := r.Group("/form")
	{
		form.POST("/", handlers.FormHandler)
		form.GET("/:id", handlers.GetFormWithFieldsHandler)
		form.PUT("/:id", handlers.UpdateFormHandler)
		form.DELETE("/:id", handlers.DeleteFormHandler)
	}

	// Form Fields
	formFields := r.Group("/form_fields")
	{
		formFields.POST("/", handlers.CreateFormFieldsHandler)
		formFields.POST("/multiple", handlers.CreateMultipleFormFieldsHandler)
	}

	// Form Groups
	formGroups := r.Group("/form_groups")
	{
		formGroups.POST("/", handlers.CreateFormGroupHandler)
		formGroups.GET("/:id", handlers.GetFormGroupHandler)
		formGroups.GET("/", handlers.GetAllFormGroupsHandler)
		formGroups.PUT("/:id", handlers.UpdateFormGroupHandler)
		formGroups.DELETE("/:id", handlers.DeleteFormGroupHandler)
	}

	// Fields
	fields := r.Group("/field")
	{
		fields.GET("/:id", handlers.GetFieldHandler)
		fields.POST("/", handlers.CreateFieldHandler)
		fields.PUT("/:id", handlers.UpdateFieldHandler) // Changed PATCH to PUT for consistency, check Handler
		fields.DELETE("/:id", handlers.DeleteFieldHandler)
	}

	// Groups
	groups := r.Group("/groups")
	{
		groups.GET("/:id", handlers.GetGroupByIDHandler)
		groups.GET("/", handlers.GetAllGroupsHandler)
		groups.POST("/", handlers.CreateGroupHandler)
		groups.PUT("/:id", handlers.UpdateGroupHandler)
		groups.DELETE("/:id", handlers.DeleteGroupHandler)
	}

	// Collections
	collections := r.Group("/collections")
	{
		collections.GET("/:id", handlers.GetCollectionHandler)
		collections.GET("/", handlers.GetAllCollectionsHandler)
		collections.POST("/", handlers.CreateCollectionHandler)
		collections.PUT("/:id", handlers.UpdateCollectionHandler)
		// collections.DELETE("/:id", handlers.DeleteCollectionHandler)
	}

	// Collection Items
	collectionItems := r.Group("/collection_items")
	{
		collectionItems.GET("/:id", handlers.GetCollectionItemHandler)
		collectionItems.POST("/", handlers.CreateCollectionItemHandler)
		collectionItems.PUT("/:id", handlers.UpdateCollectionItemHandler)
		// collectionItems.DELETE("/:id", handlers.DeleteCollectionItemHandler)
	}

	// Data Types
	dataTypes := r.Group("/data_types")
	{
		dataTypes.GET("/:id", handlers.GetDataTypeHandler)
		dataTypes.GET("/", handlers.GetAllDataTypesHandler)
		dataTypes.POST("/", handlers.CreateDataTypeHandler)
		dataTypes.PUT("/:id", handlers.UpdateDataTypeHandler)
		dataTypes.DELETE("/:id", handlers.DeleteDataTypeHandler)
	}

	// Users
	users := r.Group("/users")
	{
		users.GET("/:id", handlers.GetUserHandler)
		users.POST("/", handlers.CreateUserHandler)
		users.PUT("/:id", handlers.UpdateUserHandler)
		users.DELETE("/:id", handlers.DeleteUserHandler)
	}

	// Submissions
	submissions := r.Group("/submission")
	{
		submissions.POST("/", handlers.SubmitFormHandler)
		submissions.GET("/:id", handlers.GetSubmissionHandler)
		// Changed path to service/:service_id as discussed in handler update logic
		submissions.GET("/service/:service_id", handlers.GetSubmissionsByFormIDHandler)
	}

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
