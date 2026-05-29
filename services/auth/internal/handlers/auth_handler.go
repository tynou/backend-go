package handlers

import (
	"auth/internal/service"
	"net/http"
	_ "pkg/response"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя в базе данных
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body handlers.LoginRequest true "Данные пользователя"
// @Success      201  {object}  response.SuccessResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

// Login godoc
// @Summary      Авторизация пользователя
// @Description  Возвращает JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body handlers.LoginRequest true "Учетные данные"
// @Success      200  {object}  map[string]string "Пример: {\"token\": \"eyJhb...\"}"
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
