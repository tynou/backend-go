package handler

import (
	"gateway/internal/api/response"
	"gateway/internal/middleware"
	"log/slog"
	"net/http"
	"pkg/api/payment"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type PaymentRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type PaymentHandler struct {
	client payment.PaymentServiceClient
	log    *slog.Logger
}

func NewPaymentHandler(client payment.PaymentServiceClient, log *slog.Logger) *PaymentHandler {
	return &PaymentHandler{client: client, log: log}
}

func (h *PaymentHandler) Pay(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int32)

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to decode request", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	resp, err := h.client.Pay(c.Request.Context(), &payment.PaymentRequest{
		UserId: userID,
		Amount: req.Amount,
	})
	if err != nil {
		h.log.Error("failed to create payment", slog.Any("err", err))
		if st, ok := status.FromError(err); ok {
			c.JSON(http.StatusInternalServerError, response.Error(st.Message()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.OK(resp.Message))
}
