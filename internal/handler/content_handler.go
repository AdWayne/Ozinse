package handler

import (
	"net/http"
	"strconv"

	"ozinse/internal/service"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	projectService *service.ProjectService
}

func NewContentHandler(projectService *service.ProjectService) *ContentHandler {
	return &ContentHandler{projectService: projectService}
}

func (h *ContentHandler) GetProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	projects, total, err := h.projectService.GetProjects(
		strPtr(c.Query("search")),
		intPtr(c.Query("category_id")),
		intPtr(c.Query("genre_id")),
		intPtr(c.Query("age_rating_id")),
		strPtr(c.Query("project_type")),
		page, limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": "INTERNAL_ERROR",
			"message":    "Ошибка получения проектов",
			"details":    nil,
		})
		return
	}

	totalPages := (total + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"data":        projects,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	})
}

func (h *ContentHandler) GetProjectByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_ID", "message": "Некорректный ID", "details": nil})
		return
	}

	project, err := h.projectService.GetProjectByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error_code": "NOT_FOUND", "message": "Проект не найден", "details": nil})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ContentHandler) GetSeasons(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_ID", "message": "Некорректный ID", "details": nil})
		return
	}

	seasons, err := h.projectService.GetSeasons(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}

	c.JSON(http.StatusOK, seasons)
}

func (h *ContentHandler) GetFeatured(c *gin.Context) {
	blocks, err := h.projectService.GetFeatured()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}

	c.JSON(http.StatusOK, blocks)
}

// хелперы
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(s string) *int {
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}