package handler

import (
	"gateway/internal/api/response"
	"gateway/internal/middleware"
	"log/slog"
	"net/http"
	"pkg/api/billing"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type BillingHandler struct {
	client billing.BillingServiceClient
	log    *slog.Logger
}

func NewBillingHandler(client billing.BillingServiceClient, log *slog.Logger) *BillingHandler {
	return &BillingHandler{client: client, log: log}
}

func (h *BillingHandler) Deposit(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int32)

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to decode request", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, response.Error("failed to decode request"))
		return
	}

	resp, err := h.client.Deposit(c.Request.Context(), &billing.DepositRequest{
		UserId: userID,
		Amount: req.Amount,
	})
	if err != nil {
		h.log.Error("failed to deposit", slog.Any("err", err))
		if st, ok := status.FromError(err); ok {
			c.JSON(http.StatusInternalServerError, response.Error(st.Message()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("internal error"))
		return
	}

	c.JSON(http.StatusOK, response.OK(resp.Message))
}
