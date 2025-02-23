package handler

import (
	"biblia-be/internal/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func (handler *UserHandler) Initialize(db *gorm.DB) {
	handler.db = db
	db.AutoMigrate(&model.User{})
}

// GetAllUsers godoc
//
//	@Summary	Get all users
//	@Schemes
//	@Description	return all user objects in the database
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	model.User
//	@Router			/user [get]
func (handler *UserHandler) GetUsers(c *gin.Context) {
	users := []model.User{}

	handler.db.Preload("Records").Find(&users)

	c.JSON(http.StatusOK, users)
}

// GetUserById godoc
//
//	@Summary	Get user with its id
//	@Schemes
//	@Description	return one user with coresponding id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	model.User
//	@Router			/user/{id} [get]
func (handler *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user := model.User{}

	if err := handler.db.Preload("Records").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
//
//	@Summary	Create a new user
//	@Schemes
//	@Description	Create a new user account
//	@Tags			user
//	@Accept			json
//	@Produce		json
//
// @Param user body model.CreateUser true "User data"
//
//	@Success		201	{object}	model.User
//	@Router			/user [post]
func (handler *UserHandler) CreateUser(c *gin.Context) {
	createUser := model.CreateUser{}

	if err := c.ShouldBindJSON(&createUser); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user := model.User{Username: createUser.Username, Password: createUser.Password, Email: createUser.Email}
	log.Print(user)
	if err := handler.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUserById godoc
//
//	@Summary	Update user with the id
//	@Schemes
//	@Description	update one user with coresponding id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//
// @Param user body model.UpdateUser true "Updated User data"
//
//	@Success		200	{object}	model.User
//	@Router			/user/{id} [put]
func (handler *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	user := model.User{}
	updateUser := model.UpdateUser{}

	if err := handler.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user.Username = updateUser.Username
	user.Password = updateUser.Password
	if err := handler.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUserById godoc
//
//	@Summary	Delete user with its id
//	@Schemes
//	@Description	Delete one user with coresponding id
//	@Tags			user
//	@Param			id	path	int	true	"User ID"
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Router			/user/{id} [delete]
func (handler *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	log.Printf("deleting user with id %s", id)
	user := model.User{}

	if err := handler.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	if err := handler.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
