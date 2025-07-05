package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/auth/application"
	"github.com/mingtmt/book-store/pkg/logger"
)

// AuthRequest represents the login payload
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the login response payload
type AuthResponse struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body controller.AuthRequest true "User credentials"
// @Success 201 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req AuthRequest
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

// LoginUser godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body controller.AuthRequest true "Login credentials"
// @Success 200 {object} controller.AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req AuthRequest
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
