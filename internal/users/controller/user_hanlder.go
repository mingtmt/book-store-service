package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/users/application"
	"github.com/mingtmt/book-store/pkg/logger"
)

type UserHandler struct {
	authService *application.AuthService
}

func NewUserHandler(authService *application.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
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

func (h *UserHandler) LoginUser(c *gin.Context) {
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
