package handler

import (
	"gateway/internal/middleware"
	"net/http"
	"pkg/api/billing"

	"github.com/gin-gonic/gin"
)

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type BillingHandler struct {
	client billing.BillingServiceClient
}

func NewBillingHandler(client billing.BillingServiceClient) *BillingHandler {
	return &BillingHandler{client: client}
}

func (h *BillingHandler) Deposit(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int32)

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.Deposit(c.Request.Context(), &billing.DepositRequest{
		UserId: userID,
		Amount: req.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}
