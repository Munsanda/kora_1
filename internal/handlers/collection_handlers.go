package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models" // Required for Swagger
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Collection Handlers

type CollectionRequest struct {
	CollectionName string `json:"collection_name" binding:"required"`
}

type CollectionResponse struct {
	ID             uint   `json:"id"`
	CollectionName string `json:"collection_name"`
}

// CreateCollectionHandler creates a new collection
// @Summary      Create a new collection
// @Description  Create a new collection with the provided name
// @Tags         collections
// @Accept       json
// @Produce      json
// @Param        request  body      CollectionRequest  true  "Collection Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /collections [post]
func CreateCollectionHandler(c *gin.Context) {
	var request CollectionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	collection := &models.Collection{CollectionName: request.CollectionName}
	if err := models.CreateCollection(database.DB, collection); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess(CollectionResponse{
		ID:             collection.ID,
		CollectionName: collection.CollectionName,
	}, "Collection created successfully"))
}

// GetCollectionHandler retrieves a collection by ID
// @Summary      Get collection by ID
// @Description  Retrieve a collection by its ID
// @Tags         collections
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Collection ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404  {object}  structs.ErrorResponse
// @Router       /collections/{id} [get]
func GetCollectionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	collection, err := models.GetCollection(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Collection not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess(CollectionResponse{
		ID:             collection.ID,
		CollectionName: collection.CollectionName,
	}, "Collection retrieved successfully"))
}

// GetAllCollectionsHandler retrieves all collections
// @Summary      Get all collections
// @Description  Retrieve all collections
// @Tags         collections
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /collections [get]
func GetAllCollectionsHandler(c *gin.Context) {
	collections, err := models.GetAllCollections(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	var response []CollectionResponse
	for _, col := range collections {
		response = append(response, CollectionResponse{
			ID:             col.ID,
			CollectionName: col.CollectionName,
		})
	}
	c.JSON(http.StatusOK, helpers.NewSuccess[[]CollectionResponse](response, "Collections retrieved successfully"))
}

// UpdateCollectionHandler updates a collection
// @Summary      Update collection
// @Description  Update an existing collection by its ID
// @Tags         collections
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "Collection ID"
// @Param        request  body      CollectionRequest  true  "Collection Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /collections/{id} [put]
func UpdateCollectionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	var request CollectionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	collection, err := models.GetCollection(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Collection not found", http.StatusNotFound))
		return
	}

	collection.CollectionName = request.CollectionName
	if err := database.DB.Save(collection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[CollectionResponse](CollectionResponse{
		ID:             collection.ID,
		CollectionName: collection.CollectionName,
	}, "Collection updated successfully"))
}

// DeleteCollectionHandler deletes a collection
// @Summary      Delete collection
// @Description  Delete a collection by its ID
// @Tags         collections
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Collection ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /collections/{id} [delete]
func DeleteCollectionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := database.DB.Delete(&models.Collection{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "Collection deleted successfully"))
}

// Collection Item Handlers

type CollectionItemRequest struct {
	CollectionID              *uint  `json:"collection_id"`
	CollectionItem            string `json:"collection_item"`
	RelationCollectionItemsID *uint  `json:"relation_collection_items_id"`
}

type CollectionItemResponse struct {
	ID                        uint   `json:"id"`
	CollectionID              *uint  `json:"collection_id"`
	CollectionItem            string `json:"collection_item"`
	RelationCollectionItemsID *uint  `json:"relation_collection_items_id"`
}

// CreateCollectionItemHandler creates a new collection item
// @Summary      Create collection item
// @Description  Create a new collection item
// @Tags         collection-items
// @Accept       json
// @Produce      json
// @Param        request  body      CollectionItemRequest  true  "Collection Item Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /collection_items [post]
func CreateCollectionItemHandler(c *gin.Context) {
	var request CollectionItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	item := &models.CollectionItem{
		CollectionID:              request.CollectionID,
		CollectionItem:            request.CollectionItem,
		RelationCollectionItemsID: request.RelationCollectionItemsID,
	}

	if err := models.CreateCollectionItem(database.DB, item); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess[CollectionItemResponse](CollectionItemResponse{
		ID:                        item.ID,
		CollectionID:              item.CollectionID,
		CollectionItem:            item.CollectionItem,
		RelationCollectionItemsID: item.RelationCollectionItemsID,
	}, "Collection item created successfully"))
}

// GetCollectionItemHandler retrieves a collection item by ID
// @Summary      Get collection item
// @Description  Retrieve a collection item by its ID
// @Tags         collection-items
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Collection Item ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404  {object}  structs.ErrorResponse
// @Router       /collection_items/{id} [get]
func GetCollectionItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	item, err := models.GetCollectionItem(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Collection item not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[CollectionItemResponse](CollectionItemResponse{
		ID:                        item.ID,
		CollectionID:              item.CollectionID,
		CollectionItem:            item.CollectionItem,
		RelationCollectionItemsID: item.RelationCollectionItemsID,
	}, "Collection item retrieved successfully"))
}

// UpdateCollectionItemHandler updates a collection item
// @Summary      Update collection item
// @Description  Update an existing collection item by its ID
// @Tags         collection-items
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "Collection Item ID"
// @Param        request  body      CollectionItemRequest  true  "Collection Item Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /collection_items/{id} [put]
func UpdateCollectionItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	var request CollectionItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	item, err := models.GetCollectionItem(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("Collection item not found", http.StatusNotFound))
		return
	}

	item.CollectionID = request.CollectionID
	item.CollectionItem = request.CollectionItem
	item.RelationCollectionItemsID = request.RelationCollectionItemsID

	if err := database.DB.Save(item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[CollectionItemResponse](CollectionItemResponse{
		ID:                        item.ID,
		CollectionID:              item.CollectionID,
		CollectionItem:            item.CollectionItem,
		RelationCollectionItemsID: item.RelationCollectionItemsID,
	}, "Collection item updated successfully"))
}

// DeleteCollectionItemHandler deletes a collection item
// @Summary      Delete collection item
// @Description  Delete a collection item by its ID
// @Tags         collection-items
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Collection Item ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /collection_items/{id} [delete]
func DeleteCollectionItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := database.DB.Delete(&models.CollectionItem{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "Collection item deleted successfully"))
}
