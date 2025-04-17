package handler

import (
	"biblia-be/internal/model"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler manages user account operations
type UserHandler struct {
	db *gorm.DB
}

// Initialize sets up the handler with a database connection and performs migrations
func (handler *UserHandler) Initialize(db *gorm.DB) {
	handler.db = db
	db.AutoMigrate(&model.User{})
}

// UserResponse is a user response with password field removed
type UserResponse struct {
	ID             uint           `json:"id"`
	Username       string         `json:"username"`
	FavoriteGenres []string       `json:"favorite_genres"`
	Records        []model.Record `json:"records,omitempty"`
}

// toResponse converts a User model to a UserResponse model (without password)
func toResponse(user model.User) UserResponse {
	return UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		FavoriteGenres: user.FavoriteGenres,
		Records:        user.Records,
	}
}

// toResponseArray converts an array of User models to UserResponse models
func toResponseArray(users []model.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, user := range users {
		result[i] = toResponse(user)
	}
	return result
}

// GetUsers godoc
//
//	@Summary	Get all users
//	@Schemes
//	@Description	Returns all user accounts with their reading records
//	@Tags			users
//	@Accept			json
//	@Produce		json
//
//	@Success		200	{object} Response{data=[]UserResponse} "Successfully retrieved users"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/users [get]
func (handler *UserHandler) GetUsers(c *gin.Context) {
	var users []model.User

	result := handler.db.Preload("Records").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve users: " + result.Error.Error(),
		})
		return
	}

	// Return users without password fields
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    toResponseArray(users),
		Message: "Users retrieved successfully",
	})
}

// GetUser godoc
//
//	@Summary	Get user by ID
//	@Schemes
//	@Description	Returns a specific user account and their reading records
//	@Tags			users
//	@Accept			json
//	@Produce		json
//
//	@Param			id	path	int	true	"User ID"
//	@Success		200	{object} Response{data=UserResponse} "Successfully retrieved user"
//	@Failure		400	{object} Response "Invalid user ID"
//	@Failure		404	{object} Response "User not found"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/users/{id} [get]
func (handler *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")

	// Validate ID
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID format",
		})
		return
	}

	var user model.User
	result := handler.db.Preload("Records").First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error:   "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve user: " + result.Error.Error(),
		})
		return
	}

	// Return user without password field
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    toResponse(user),
		Message: "User retrieved successfully",
	})
}

// CreateUser godoc
//
//	@Summary	Create a new user
//	@Schemes
//	@Description	Creates a new user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//
//	@Param user body model.CreateUser true "User data"
//	@Success		201	{object} Response{data=UserResponse} "User created successfully"
//	@Failure		400	{object} Response "Invalid request body or validation error"
//	@Failure		409	{object} Response "Username already exists"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/users [post]
func (handler *UserHandler) CreateUser(c *gin.Context) {
	var createUser model.CreateUser

	// Parse request body
	if err := c.ShouldBindJSON(&createUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate username length
	if len(createUser.Username) < 3 || len(createUser.Username) > 50 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Username must be between 3 and 50 characters",
		})
		return
	}

	// Validate password length
	if len(createUser.Password) < 6 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Password must be at least 6 characters long",
		})
		return
	}

	// Check if username already exists
	var existingUser model.User
	if handler.db.Where("username = ?", createUser.Username).First(&existingUser).RowsAffected > 0 {
		c.JSON(http.StatusConflict, Response{
			Success: false,
			Error:   "Username already exists",
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to process password",
		})
		return
	}

	// Create new user
	user := model.User{
		Username:       createUser.Username,
		Password:       string(hashedPassword),
		FavoriteGenres: createUser.FavoriteGenres,
	}

	// Save user to database
	if err := handler.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to create user: " + err.Error(),
		})
		return
	}

	// Return created user without password field
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    toResponse(user),
		Message: "User created successfully",
	})
}

// UpdateUser godoc
//
//	@Summary	Update an existing user
//	@Schemes
//	@Description	Updates a user's username, password, and/or favorite genres
//	@Tags			users
//	@Accept			json
//	@Produce		json
//
//	@Param			id	path	int	true	"User ID"
//	@Param user body model.UpdateUser true "Updated user data"
//	@Success		200	{object} Response{data=UserResponse} "User updated successfully"
//	@Failure		400	{object} Response "Invalid request body or validation error"
//	@Failure		404	{object} Response "User not found"
//	@Failure		409	{object} Response "Username already exists"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/users/{id} [put]
func (handler *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")

	// Validate ID
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID format",
		})
		return
	}

	var updateUser model.UpdateUser
	// Parse request body
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate username length
	if len(updateUser.Username) < 3 || len(updateUser.Username) > 50 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Username must be between 3 and 50 characters",
		})
		return
	}

	// Validate password length
	if len(updateUser.Password) < 6 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Password must be at least 6 characters long",
		})
		return
	}

	// Find user to update
	var user model.User
	result := handler.db.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error:   "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve user for update: " + result.Error.Error(),
		})
		return
	}

	// Check if the new username already exists for another user
	if updateUser.Username != user.Username {
		var existingUser model.User
		if handler.db.Where("username = ? AND id != ?", updateUser.Username, id).First(&existingUser).RowsAffected > 0 {
			c.JSON(http.StatusConflict, Response{
				Success: false,
				Error:   "Username already exists",
			})
			return
		}
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to process password",
		})
		return
	}

	// Update user fields
	user.Username = updateUser.Username
	user.Password = string(hashedPassword)
	user.FavoriteGenres = updateUser.FavoriteGenres

	// Save updated user
	if err := handler.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to update user: " + err.Error(),
		})
		return
	}

	// Return updated user without password field
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    toResponse(user),
		Message: "User updated successfully",
	})
}

// DeleteUser godoc
//
//	@Summary	Delete a user
//	@Schemes
//	@Description	Permanently removes a user account and their reading records
//	@Tags			users
//	@Accept			json
//	@Produce		json
//
//	@Param			id	path	int	true	"User ID"
//	@Success		200	{object} Response "User deleted successfully"
//	@Failure		400	{object} Response "Invalid user ID"
//	@Failure		404	{object} Response "User not found"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/users/{id} [delete]
func (handler *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")

	// Validate ID
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID format",
		})
		return
	}

	// Find user to delete
	var user model.User
	result := handler.db.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error:   "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve user for deletion: " + result.Error.Error(),
		})
		return
	}

	// Delete user
	if err := handler.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to delete user: " + err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// AuthenticateUser godoc
//
//	@Summary	Authenticate a user
//	@Schemes
//	@Description	Validates user credentials and returns user info
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//
//	@Param			credentials	body	object	true	"Login credentials"	Schema(type=object,required=["username","password"],properties={username={type=string},password={type=string}})
//	@Success		200	{object} Response{data=UserResponse} "Authentication successful"
//	@Failure		400	{object} Response "Invalid request body"
//	@Failure		401	{object} Response "Invalid credentials"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/auth/login [post]
func (handler *UserHandler) AuthenticateUser(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Parse request body
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Find user by username
	var user model.User
	result := handler.db.Where("username = ?", credentials.Username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, Response{
				Success: false,
				Error:   "Invalid username or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Authentication error: " + result.Error.Error(),
		})
		return
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid username or password",
		})
		return
	}

	// Get user records
	handler.db.Model(&user).Association("Records").Find(&user.Records)

	// Return authenticated user without password
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    toResponse(user),
		Message: "Authentication successful",
	})
}
