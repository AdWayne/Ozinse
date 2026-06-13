package handler

import (
	"net/http"

	"ozinse/internal/model"
	"ozinse/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	tokens, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "REGISTRATION_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	tokens, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_code": "AUTH_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	tokens, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_code": "REFRESH_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Инструкция отправлена на ваш email.",
	})
}
