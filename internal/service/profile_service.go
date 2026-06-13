package service

import (
	"errors"

	"ozinse/internal/model"
	"ozinse/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type ProfileService struct {
	userRepo *repository.UserRepo
}

func NewProfileService(userRepo *repository.UserRepo) *ProfileService {
	return &ProfileService{userRepo: userRepo}
}

func (s *ProfileService) GetProfile(userID int) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}
	return user, nil
}

func (s *ProfileService) UpdateProfile(userID int, req model.UpdateProfileRequest) error {
	return s.userRepo.UpdateProfile(userID, req.FullName, req.Phone, req.BirthDate)
}

func (s *ProfileService) ChangePassword(userID int, req model.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		return errors.New("неверный текущий пароль")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, string(hash))
}