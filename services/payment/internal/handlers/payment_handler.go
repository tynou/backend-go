package handlers

import (
	"net/http"
	"payment/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	UserId int32   `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) Pay(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Pay(c.Request.Context(), req.UserId, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "payment pending"})
}
