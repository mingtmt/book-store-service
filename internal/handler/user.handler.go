package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/dto"
	"github.com/mingtmt/book-store/internal/service"
	"github.com/mingtmt/book-store/internal/utils"
	"github.com/mingtmt/book-store/internal/utils/validation"
	"github.com/mingtmt/book-store/pkg/response"
)

type UserHandler struct {
	service service.UserService
}

type GetUserByUUIDParam struct {
	UUID string `uri:"uuid" binding:"uuid"`
}

type GetUsersParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50,search"`
	Page   int    `form:"page" binding:"omitempty,gte=1,lt=100"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	var params GetUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	users, err := uh.service.GetAllUsers(params.Search, params.Page, params.Limit)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	usersDTO := dto.MapUsersToDTO(users)

	response.Success(ctx, usersDTO)
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapCreateInputToModel()

	createdUser, err := uh.service.CreateUser(user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(createdUser)

	response.Success(ctx, userDTO)
}

func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var params GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user, err := uh.service.GetUserByUUID(params.UUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(user)

	utils.ResponseSuccess(ctx, http.StatusOK, &userDTO)
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var params GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var input dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapUpdateInputToModel()

	updatedUser, err := uh.service.UpdateUser(params.UUID, user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(updatedUser)

	response.Success(ctx, userDTO)
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var params GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := uh.service.DeleteUser(params.UUID); err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}
