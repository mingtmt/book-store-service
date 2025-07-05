package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/auths/application"
	"github.com/mingtmt/book-store/pkg/logger"
)

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	user, err := h.authService.RegisterUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	logger.Info("user registered successfully", map[string]interface{}{
		"user_id": user.ID,
	})

	c.Status(http.StatusCreated)
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	token, err := h.authService.LoginUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	logger.Info("user logged in successfully", map[string]interface{}{
		"username": req.Username,
	})

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
