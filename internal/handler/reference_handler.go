package handler

import (
	"net/http"

	"ozinse/internal/service"

	"github.com/gin-gonic/gin"
)

type ReferenceHandler struct {
	refService *service.ReferenceService
}

func NewReferenceHandler(refService *service.ReferenceService) *ReferenceHandler {
	return &ReferenceHandler{refService: refService}
}

func (h *ReferenceHandler) GetCategories(c *gin.Context) {
	list, err := h.refService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ReferenceHandler) GetGenres(c *gin.Context) {
	list, err := h.refService.GetGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ReferenceHandler) GetAgeRatings(c *gin.Context) {
	list, err := h.refService.GetAgeRatings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, list)
}