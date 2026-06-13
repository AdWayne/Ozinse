package service

import (
	"ozinse/internal/model"
	"ozinse/internal/repository"
)

type ReferenceService struct {
	refRepo *repository.ReferenceRepo
}

func NewReferenceService(refRepo *repository.ReferenceRepo) *ReferenceService {
	return &ReferenceService{refRepo: refRepo}
}

func (s *ReferenceService) GetCategories() ([]model.Category, error) {
	return s.refRepo.GetCategories()
}

func (s *ReferenceService) GetGenres() ([]model.Genre, error) {
	return s.refRepo.GetGenres()
}

func (s *ReferenceService) GetAgeRatings() ([]model.AgeRating, error) {
	return s.refRepo.GetAgeRatings()
}