package handlers

import (
	"billing/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DepositRequest struct {
	UserID int32   `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type WalletHandler struct {
	service *service.BillingService
}

func NewWalletHandler(service *service.BillingService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Deposit(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deposit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deposited successfully"})
}
