package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/auth/application"
	"github.com/mingtmt/book-store/pkg/logger"
)

// Request payloads
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type ReAuthRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Response payloads
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

	user, accessToken, refreshToken, err := h.authService.RegisterUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	logger.Info("user registered successfully", map[string]interface{}{
		"user_id": user.ID,
	})

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
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

	accessToken, refreshToken, err := h.authService.LoginUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	logger.Info("user logged in successfully", map[string]interface{}{
		"username": req.Username,
	})

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body controller.ReAuthRequest true "Refresh token"
// @Success 200 {object} controller.AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req ReAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	access, refresh, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// LogoutUser godoc
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param logout body controller.ReAuthRequest true "Logout request"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	var req ReAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	err := h.authService.LogoutUser(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
