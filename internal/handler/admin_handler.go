package handler

import (
	"net/http"
	"strconv"
	"io"
	"os"
	"path/filepath"
	"time"
	"fmt"
	"ozinse/internal/model"
	"ozinse/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "UPLOAD_ERROR", "message": "Файл не найден", "details": nil})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join("uploads", filename)

	out, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "UPLOAD_ERROR", "message": "Ошибка сохранения файла", "details": nil})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "UPLOAD_ERROR", "message": "Ошибка копирования файла", "details": nil})
		return
	}

	url := h.adminService.GetBaseURL() + "/uploads/" + filename
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// PROJECTS

func (h *AdminHandler) CreateProject(c *gin.Context) {
	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	project, err := h.adminService.CreateProject(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "CREATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusCreated, project)
}

func (h *AdminHandler) UpdateProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	if err := h.adminService.UpdateProject(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "UPDATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Проект обновлён"})
}

func (h *AdminHandler) DeleteProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.adminService.DeleteProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "DELETE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Проект удалён"})
}

// SEASONS

func (h *AdminHandler) CreateSeason(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	var req model.CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	season, err := h.adminService.CreateSeason(projectID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "CREATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusCreated, season)
}

// EPISODES

func (h *AdminHandler) CreateEpisode(c *gin.Context) {
	seasonID, _ := strconv.Atoi(c.Param("season_id"))
	var req model.CreateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	episode, err := h.adminService.CreateEpisode(seasonID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "CREATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusCreated, episode)
}

// FEATURED

func (h *AdminHandler) UpdateFeaturedOrder(c *gin.Context) {
	var req model.UpdateFeaturedOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	var items []struct {
		ProjectID int
		SortOrder int
	}
	for _, it := range req.Items {
		items = append(items, struct {
			ProjectID int
			SortOrder int
		}{it.ProjectID, it.SortOrder})
	}
	if err := h.adminService.UpdateFeaturedOrder(req.BlockType, items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "UPDATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Порядок обновлён"})
}

// CATEGORIES

func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	cat, err := h.adminService.CreateCategory(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "CREATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	if err := h.adminService.UpdateCategory(id, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "UPDATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Категория обновлена"})
}

func (h *AdminHandler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.adminService.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "DELETE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Категория удалена"})
}

// GENRES

func (h *AdminHandler) CreateGenre(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	genre, err := h.adminService.CreateGenre(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "CREATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusCreated, genre)
}

func (h *AdminHandler) UpdateGenre(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	if err := h.adminService.UpdateGenre(id, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "UPDATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Жанр обновлён"})
}

func (h *AdminHandler) DeleteGenre(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.adminService.DeleteGenre(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "DELETE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Жанр удалён"})
}

// AGE RATINGS

func (h *AdminHandler) CreateAgeRating(c *gin.Context) {
	var req struct {
		Range string `json:"range" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	ar, err := h.adminService.CreateAgeRating(req.Range)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "CREATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusCreated, ar)
}

func (h *AdminHandler) UpdateAgeRating(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Range string `json:"range" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	if err := h.adminService.UpdateAgeRating(id, req.Range); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "UPDATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Возрастной рейтинг обновлён"})
}

func (h *AdminHandler) DeleteAgeRating(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.adminService.DeleteAgeRating(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "DELETE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Возрастной рейтинг удалён"})
}

// USERS

func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	users, total, err := h.adminService.GetAllUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}
	totalPages := (total + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{"data": users, "page": page, "limit": limit, "total": total, "total_pages": totalPages})
}

func (h *AdminHandler) AssignRole(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))
	var req struct {
		RoleID int `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	if err := h.adminService.AssignRole(userID, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "ASSIGN_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Роль назначена"})
}

// ROLES

func (h *AdminHandler) GetRoles(c *gin.Context) {
	roles, err := h.adminService.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Permissions string `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	role, err := h.adminService.CreateRole(req.Name, req.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "CREATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusCreated, role)
}

func (h *AdminHandler) UpdateRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name        string `json:"name" binding:"required"`
		Permissions string `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error(), "details": nil})
		return
	}
	if err := h.adminService.UpdateRole(id, req.Name, req.Permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "UPDATE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Роль обновлена"})
}

func (h *AdminHandler) DeleteRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.adminService.DeleteRole(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "DELETE_ERROR", "message": err.Error(), "details": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Роль удалена"})
}