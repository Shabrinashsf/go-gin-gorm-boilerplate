package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/api/service"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/dto"
	myerror "github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/pkg/error"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	UserController interface {
		RegisterUser(ctx *gin.Context)
		Login(ctx *gin.Context)
		SendVerificationEmail(ctx *gin.Context)
		VerifyEmail(ctx *gin.Context)
		ForgotPassword(ctx *gin.Context)
		ResetPassword(ctx *gin.Context)
		MeAuth(ctx *gin.Context)
		UpdateUser(ctx *gin.Context)
	}

	userController struct {
		userService service.UserService
	}
)

func NewUserController(us service.UserService) UserController {
	return &userController{
		userService: us,
	}
}

func (c *userController) RegisterUser(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 20*time.Second)
	defer cancel()

	var req dto.UserRegistrationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, myerror.New(myerror.FormatValidationError(err), http.StatusBadRequest), nil).Send(ctx)
		return
	}

	result, err := c.userService.RegisterUser(reqCtx, req)
	if err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_REGISTER_USER, myerror.FromDBError(err), nil).Send(ctx)
		return
	}

	response.NewSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, result).Send(ctx)
}

func (c *userController) Login(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 20*time.Second)
	defer cancel()

	var req dto.UserLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, myerror.New(myerror.FormatValidationError(err), http.StatusBadRequest), nil).Send(ctx)
		return
	}

	result, err := c.userService.Login(reqCtx, req)
	if err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_LOGIN_USER, myerror.FromDBError(err), nil).Send(ctx)
		return
	}

	response.NewSuccess(dto.MESSAGE_SUCCESS_LOGIN_USER, result).Send(ctx)
}

func (c *userController) SendVerificationEmail(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 20*time.Second)
	defer cancel()

	var req dto.SendVerificationEmailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, myerror.New(myerror.FormatValidationError(err), http.StatusBadRequest), nil).Send(ctx)
		return
	}

	err := c.userService.SendVerificationEmail(reqCtx, req)
	if err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, myerror.FromDBError(err), nil)
		return
	}

	response.NewSuccess(dto.MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS, nil)
}

func (c *userController) VerifyEmail(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 20*time.Second)
	defer cancel()

	token := ctx.Query("token")

	if token == "" {
		response.NewFailed(dto.MESSAGE_FAILED_TOKEN_NOT_FOUND, dto.ErrTokenNotFound, nil).Send(ctx)
		return
	}

	req := dto.VerifyEmailRequest{
		Token: token,
	}

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, myerror.New(myerror.FormatValidationError(err), http.StatusBadRequest), nil).Send(ctx)
		return
	}

	result, err := c.userService.VerifyEmail(reqCtx, req)
	if err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_VERIFY_EMAIL, myerror.FromDBError(err), nil).Send(ctx)
		return
	}

	response.NewSuccess(dto.MESSAGE_SUCCESS_VERIFY_EMAIL, result).Send(ctx)
}

func (c *userController) ForgotPassword(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 20*time.Second)
	defer cancel()

	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, myerror.New(myerror.FormatValidationError(err), http.StatusBadRequest), nil).Send(ctx)
		return
	}

	if err := c.userService.ForgotPassword(reqCtx, req); err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_FORGET_PASSWORD, myerror.FromDBError(err), nil).Send(ctx)
		return
	}

	response.NewSuccess(dto.MESSAGE_SUCCESS_FORGET_PASSWORD, nil).Send(ctx)
}

func (c *userController) ResetPassword(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 20*time.Second)
	defer cancel()

	token := ctx.Query("token")
	var req dto.ResetPasswordRequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, myerror.New(myerror.FormatValidationError(err), http.StatusBadRequest), nil).Send(ctx)
		return
	}

	if err := c.userService.ResetPassword(reqCtx, token, req.Password); err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_RESET_PASSWORD, myerror.FromDBError(err), nil)
		return
	}

	response.NewSuccess(dto.MESSAGE_SUCCESS_RESET_PASSWORD, nil).Send(ctx)
}

func (c *userController) MeAuth(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	result, err := c.userService.GetUserByID(ctx.Request.Context(), uuid.MustParse(userId))
	if err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, myerror.FromDBError(err), nil).Send(ctx)
		return
	}

	response.NewSuccess(dto.MESSAGE_SUCCESS_GET_USER, result).Send(ctx)
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 20*time.Second)
	defer cancel()

	userId := ctx.MustGet("user_id").(string)
	var req dto.UserUpdateRequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, myerror.New(myerror.FormatValidationError(err), http.StatusBadRequest), nil).Send(ctx)
		return
	}

	result, err := c.userService.UpdateUser(reqCtx, uuid.MustParse(userId), req)
	if err != nil {
		response.NewFailed(dto.MESSAGE_FAILED_UPDATE_USER, myerror.FromDBError(err), nil).Send(ctx)
		return
	}

	response.NewSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, result).Send(ctx)
}
