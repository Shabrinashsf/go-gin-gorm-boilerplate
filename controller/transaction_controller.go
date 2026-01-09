package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/dto"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/service"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/utils/response"
	"github.com/gin-gonic/gin"
)

type (
	TransactionController interface {
		TripayWebhook(ctx *gin.Context)
	}

	transactionController struct {
		transactionService service.TransactionService
	}
)

func NewTransactionController(ts service.TransactionService) TransactionController {
	return &transactionController{
		transactionService: ts,
	}
}

func (c *transactionController) TripayWebhook(ctx *gin.Context) {
	svcCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 1. Ambil raw body sekali
	rawBody, err := ctx.GetRawData()
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// 2. Parse JSON ke struct
	var req dto.TripayWebhookRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// 3. Ambil header
	cbSignature := ctx.GetHeader("X-Callback-Signature")
	cbEvent := ctx.GetHeader("X-Callback-Event")

	// 4. Kirim semuanya ke service
	_, err = c.transactionService.TripayWebhook(svcCtx, rawBody, req, cbSignature, cbEvent)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CALLBACK_TRIPAY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_CALLBACK_TRIPAY, nil)
	ctx.JSON(http.StatusOK, res)
}
