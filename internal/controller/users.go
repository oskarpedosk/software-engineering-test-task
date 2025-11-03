package controller

import (
	"cruder/internal/model"
	"cruder/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (userController *UserController) handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrUserNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrEmptyField), errors.Is(err, model.ErrInvalidEmail):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (userController *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := userController.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (userController *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := userController.userService.GetUserByUsername(username)
	if err != nil {
		userController.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (userController *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := userController.userService.GetUserByID(id)
	if err != nil {
		userController.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (userController *UserController) CreateUser(ctx *gin.Context) {
	var request model.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userController.userService.CreateUser(&request)
	if err != nil {
		userController.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (userController *UserController) UpdateUser(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var request model.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userController.userService.UpdateUser(uuid, &request)
	if err != nil {
		userController.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (userController *UserController) DeleteUser(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := userController.userService.DeleteUser(uuid)
	if err != nil {
		userController.handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
