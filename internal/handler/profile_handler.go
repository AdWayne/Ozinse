package handler

import (
	"net/http"

	"ozinse/internal/model"
	"ozinse/internal/service"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id")

	user, err := h.profileService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code": "NOT_FOUND",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	if err := h.profileService.UpdateProfile(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": "UPDATE_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Профиль обновлён"})
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	if err := h.profileService.ChangePassword(userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "PASSWORD_ERROR",
			"message":    err.Error(),
			"details":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль изменён"})
}