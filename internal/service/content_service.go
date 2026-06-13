package service

import (
	"ozinse/internal/model"
	"ozinse/internal/repository"
)

type ProjectService struct {
	projectRepo *repository.ProjectRepo
}

func NewProjectService(projectRepo *repository.ProjectRepo) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

func (s *ProjectService) GetProjects(search *string, categoryID, genreID, ageRatingID *int, projectType *string, page, limit int) ([]model.Project, int, error) {
	return s.projectRepo.GetAll(search, categoryID, genreID, ageRatingID, projectType, page, limit)
}

func (s *ProjectService) GetProjectByID(id int) (*model.Project, error) {
	return s.projectRepo.GetByID(id)
}

func (s *ProjectService) GetSeasons(projectID int) ([]model.Season, error) {
	return s.projectRepo.GetSeasonsWithEpisodes(projectID)
}

func (s *ProjectService) GetFeatured() (map[string][]model.Project, error) {
	return s.projectRepo.GetFeatured()
}