package service

import (
	"ozinse/internal/model"
	"ozinse/internal/repository"
)

type FavoriteService struct {
	favRepo *repository.FavoriteRepo
}

func NewFavoriteService(favRepo *repository.FavoriteRepo) *FavoriteService {
	return &FavoriteService{favRepo: favRepo}
}

func (s *FavoriteService) GetFavorites(userID int) ([]model.Favorite, error) {
	return s.favRepo.GetByUserID(userID)
}

func (s *FavoriteService) Add(userID, projectID int) error {
	return s.favRepo.Add(userID, projectID)
}

func (s *FavoriteService) Remove(userID, projectID int) error {
	return s.favRepo.Remove(userID, projectID)
}