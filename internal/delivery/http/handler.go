package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dchlong/billing-be/internal/presenter"
	"github.com/dchlong/billing-be/internal/services"
	"github.com/dchlong/billing-be/pkg/logger"
)

type Handler interface {
	Register(router gin.IRouter)
}

type handler struct {
	billService services.BillService
	ilogger     logger.ILogger
}

func ProvideHandler(billService services.BillService, ilogger logger.ILogger) Handler {
	return &handler{
		billService: billService,
		ilogger:     ilogger,
	}
}

const maxUserNameLength = 32

func (h *handler) Register(router gin.IRouter) {
	router.PUT("/mobile/:user_name/call", h.createCallHistory)
	router.GET("/mobile/:user_name/billing", h.getBill)
}

func (h *handler) createCallHistory(context *gin.Context) {
	ctx := context.Request.Context()
	log := h.ilogger.GetLogger(ctx)
	userName := context.Param("user_name")
	if len(userName) >= maxUserNameLength || userName == "" {
		context.JSON(http.StatusBadRequest, presenter.NewInvalidUserNameError(maxUserNameLength))
		return
	}

	input := &presenter.CreateCallHistoryInput{}
	if err := context.ShouldBindJSON(input); err != nil {
		context.JSON(http.StatusBadRequest, presenter.NewInvalidInputError(input, err))
		return
	}

	createCallHistoryResp, err := h.billService.CreateCallHistory(ctx, userName, input.CallDuration)
	if err != nil {
		log.Errorf("could not create call history, username: %s, error: %+v", userName, err)
		context.JSON(http.StatusInternalServerError, presenter.NewInternalServerErrorError(err))
		return
	}

	context.JSON(http.StatusOK, createCallHistoryResp)
}

func (h *handler) getBill(context *gin.Context) {
	ctx := context.Request.Context()
	log := h.ilogger.GetLogger(ctx)
	userName := context.Param("user_name")
	if len(userName) >= maxUserNameLength || userName == "" {
		context.JSON(http.StatusBadRequest, presenter.NewInvalidUserNameError(maxUserNameLength))
		return
	}

	getBillResp, err := h.billService.GetBill(ctx, userName)
	if err != nil {
		log.Errorf("could not get bill, username: %s, error: %+v", userName, err)
		context.JSON(http.StatusInternalServerError, presenter.NewInternalServerErrorError(err))
		return
	}

	context.JSON(http.StatusOK, getBillResp)
}
