package handler

import (
	"gateway/internal/middleware"
	"net/http"
	"pkg/api/payment"

	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type PaymentHandler struct {
	client payment.PaymentServiceClient
}

func NewPaymentHandler(client payment.PaymentServiceClient) *PaymentHandler {
	return &PaymentHandler{client: client}
}

func (h *PaymentHandler) Pay(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int32)

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.Pay(c.Request.Context(), &payment.PaymentRequest{
		UserId: userID,
		Amount: req.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": resp.Message})
}
