package handlers

import (
	"billing/internal/service"
	"net/http"
	"pkg/middleware"

	"github.com/gin-gonic/gin"
)

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type WalletHandler struct {
	service *service.BillingService
}

func NewWalletHandler(service *service.BillingService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int32)

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Deposit(c.Request.Context(), userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deposit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deposited successfully"})
}
