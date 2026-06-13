package handler

import (
	"net/http"
	"strconv"

	"ozinse/internal/service"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favService *service.FavoriteService
}

func NewFavoriteHandler(favService *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favService: favService}
}

func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID := c.GetInt("user_id")
	list, err := h.favService.GetFavorites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_ID", "message": "Некорректный ID", "details": nil})
		return
	}

	if err := h.favService.Add(userID, projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Добавлено в избранное"})
}

func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_ID", "message": "Некорректный ID", "details": nil})
		return
	}

	if err := h.favService.Remove(userID, projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Удалено из избранного"})
}