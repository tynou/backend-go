package handler

import (
	"gateway/internal/api/response"
	"log/slog"
	"net/http"
	"pkg/api/auth"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	client auth.AuthServiceClient
	log    *slog.Logger
}

func NewAuthHandler(client auth.AuthServiceClient, log *slog.Logger) *AuthHandler {
	return &AuthHandler{client: client, log: log}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to decode request", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, response.Error("failed to decode request"))
		return
	}

	resp, err := h.client.Register(c.Request.Context(), &auth.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		h.log.Error("failed to register user", slog.Any("err", err))
		if st, ok := status.FromError(err); ok {
			c.JSON(http.StatusInternalServerError, response.Error(st.Message()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("internal error"))
		return
	}

	c.JSON(http.StatusCreated, response.OK(resp.Message))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to decode request", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, response.Error("failed to decode request"))
		return
	}

	resp, err := h.client.Login(c.Request.Context(), &auth.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		h.log.Error("failed to login", slog.Any("err", err))
		if st, ok := status.FromError(err); ok {
			c.JSON(http.StatusUnauthorized, response.Error(st.Message()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("internal error"))
		return
	}

	c.JSON(http.StatusOK, &LoginResponse{Token: resp.Token})
}
