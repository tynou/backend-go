package handlers

import (
	"net/http"
	"payment/internal/service"
	_ "pkg/response"

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

// Pay godoc
// @Summary      Оплатить со счёта пользователя
// @Description  Регистрирует платёж и отправляет его на обработку
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        request body handlers.PaymentRequest true "Данные платежа"
// @Success      201  {object}  response.SuccessResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /pay [post]
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
