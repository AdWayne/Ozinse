package service

import (
	"ozinse/internal/model"
	"ozinse/internal/repository"
)

type AdminService struct {
	adminRepo *repository.AdminRepo
	baseURL   string
}


func NewAdminService(adminRepo *repository.AdminRepo, baseURL string) *AdminService {
	return &AdminService{adminRepo: adminRepo, baseURL: baseURL}
}

func (s *AdminService) CreateProject(req model.CreateProjectRequest) (*model.Project, error) {
	return s.adminRepo.CreateProject(req)
}

func (s *AdminService) UpdateProject(id int, req model.CreateProjectRequest) error {
	return s.adminRepo.UpdateProject(id, req)
}

func (s *AdminService) DeleteProject(id int) error {
	return s.adminRepo.DeleteProject(id)
}

func (s *AdminService) CreateSeason(projectID int, req model.CreateSeasonRequest) (*model.Season, error) {
	return s.adminRepo.CreateSeason(projectID, req)
}

func (s *AdminService) CreateEpisode(seasonID int, req model.CreateEpisodeRequest) (*model.Episode, error) {
	return s.adminRepo.CreateEpisode(seasonID, req)
}

func (s *AdminService) UpdateFeaturedOrder(blockType string, items []struct {
	ProjectID int
	SortOrder int
}) error {
	return s.adminRepo.UpdateFeaturedOrder(blockType, items)
}

func (s *AdminService) CreateCategory(name string) (*model.Category, error) {
	return s.adminRepo.CreateCategory(name)
}

func (s *AdminService) UpdateCategory(id int, name string) error {
	return s.adminRepo.UpdateCategory(id, name)
}

func (s *AdminService) DeleteCategory(id int) error {
	return s.adminRepo.DeleteCategory(id)
}

func (s *AdminService) CreateGenre(name string) (*model.Genre, error) {
	return s.adminRepo.CreateGenre(name)
}

func (s *AdminService) UpdateGenre(id int, name string) error {
	return s.adminRepo.UpdateGenre(id, name)
}

func (s *AdminService) DeleteGenre(id int) error {
	return s.adminRepo.DeleteGenre(id)
}

func (s *AdminService) CreateAgeRating(rng string) (*model.AgeRating, error) {
	return s.adminRepo.CreateAgeRating(rng)
}

func (s *AdminService) UpdateAgeRating(id int, rng string) error {
	return s.adminRepo.UpdateAgeRating(id, rng)
}

func (s *AdminService) DeleteAgeRating(id int) error {
	return s.adminRepo.DeleteAgeRating(id)
}

func (s *AdminService) GetAllUsers(page, limit int) ([]model.User, int, error) {
	return s.adminRepo.GetAllUsers(page, limit)
}

func (s *AdminService) AssignRole(userID, roleID int) error {
	return s.adminRepo.AssignRole(userID, roleID)
}

func (s *AdminService) GetRoles() ([]model.Role, error) {
	return s.adminRepo.GetRoles()
}

func (s *AdminService) CreateRole(name, permissions string) (*model.Role, error) {
	return s.adminRepo.CreateRole(name, permissions)
}

func (s *AdminService) UpdateRole(id int, name, permissions string) error {
	return s.adminRepo.UpdateRole(id, name, permissions)
}

func (s *AdminService) DeleteRole(id int) error {
	return s.adminRepo.DeleteRole(id)
}

func (s *AdminService) GetBaseURL() string {
	return s.baseURL
}