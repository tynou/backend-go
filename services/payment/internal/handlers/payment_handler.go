package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) Pay(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "payment pending"})
}
