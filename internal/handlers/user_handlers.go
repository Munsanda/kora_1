package handlers

import (
	"kora_1/internal/database"
	"kora_1/internal/helpers"
	"kora_1/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	Surname    string `json:"surname"`
	Dob        string `json:"dob"` // Format "YYYY-MM-DD"
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password"`
}

type UserResponse struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	Surname    string `json:"surname"`
	Dob        string `json:"dob"`
	Email      string `json:"email"`
}

// CreateUserHandler creates a new user
// @Summary      Create user
// @Description  Create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      UserRequest  true  "User Request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /users [post]
func CreateUserHandler(c *gin.Context) {
	var request UserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	var dob time.Time
	var err error
	if request.Dob != "" {
		dob, err = time.Parse("2006-01-02", request.Dob)
		if err != nil {
			c.JSON(http.StatusBadRequest, helpers.NewError("Invalid Dob format. Use YYYY-MM-DD", http.StatusBadRequest))
			return
		}
	}

	user := &models.User{
		FirstName:  request.FirstName,
		MiddleName: request.MiddleName,
		Surname:    request.Surname,
		Dob:        dob,
		Email:      request.Email,
		Password:   request.Password,
	}

	if err := models.CreateUser(database.DB, user); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewSuccess[UserResponse](userToResponse(user), "User created successfully"))
}

// GetUserHandler retrieves a user by ID
// @Summary      Get user by ID
// @Description  Retrieve a user by its ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,404  {object}  structs.ErrorResponse
// @Router       /users/{id} [get]
func GetUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	user, err := models.GetUser(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("User not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[UserResponse](userToResponse(user), "User retrieved successfully"))
}

// UpdateUserHandler updates a user
// @Summary      Update user
// @Description  Update an existing user by its ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int          true  "User ID"
// @Param        request  body      UserRequest  true  "User Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400,404,500  {object}  structs.ErrorResponse
// @Router       /users/{id} [put]
func UpdateUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	// For updates, we might want partial updates. For now assuming full update or using logic to check fields
	var request UserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	user, err := models.GetUser(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewError("User not found", http.StatusNotFound))
		return
	}

	if request.Dob != "" {
		dob, err := time.Parse("2006-01-02", request.Dob)
		if err == nil {
			user.Dob = dob
		}
	}

	user.FirstName = request.FirstName
	user.MiddleName = request.MiddleName
	user.Surname = request.Surname
	user.Email = request.Email
	if request.Password != "" {
		user.Password = request.Password
	}

	if err := models.UpdateUser(database.DB, user); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[UserResponse](userToResponse(user), "User updated successfully"))
}

// DeleteUserHandler deletes a user
// @Summary      Delete user
// @Description  Delete a user by its ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400,500  {object}  structs.ErrorResponse
// @Router       /users/{id} [delete]
func DeleteUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewError("Invalid ID", http.StatusBadRequest))
		return
	}

	if err := models.DeleteUser(database.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, helpers.NewSuccess[any](nil, "User deleted successfully"))
}

func userToResponse(user *models.User) UserResponse {
	dobStr := ""
	if !user.Dob.IsZero() {
		dobStr = user.Dob.Format("2006-01-02")
	}
	return UserResponse{
		ID:         user.ID,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		Surname:    user.Surname,
		Dob:        dobStr,
		Email:      user.Email,
	}
}
