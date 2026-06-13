package service

import (
	"errors"
	"time"

	"ozinse/internal/model"
	"ozinse/internal/repository"
	"ozinse/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repository.UserRepo
	tokenRepo  *repository.TokenRepo
	jwtService *jwt.Service
	baseURL    string
}

func NewAuthService(userRepo *repository.UserRepo, tokenRepo *repository.TokenRepo, jwtService *jwt.Service, baseURL string) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
		baseURL:    baseURL,
	}
}

func (s *AuthService) Register(req model.RegisterRequest) (*model.TokenPair, error) {
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// получаем роль user
	roleID, err := s.userRepo.GetRoleByName("user")
	if err != nil {
		return nil, errors.New("роль 'user' не найдена в БД")
	}

	// дефолтный аватар
	avatarURL := s.baseURL + "/static/avatars/Avatar.svg"

	user, err := s.userRepo.CreateWithRole(req.Email, string(hash), req.FullName, avatarURL, roleID)
	if err != nil {
		return nil, err
	}

	return s.generateAndSaveTokens(user)
}

func (s *AuthService) Login(req model.LoginRequest) (*model.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	return s.generateAndSaveTokens(user)
}

func (s *AuthService) Refresh(refreshTokenStr string) (*model.TokenPair, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, errors.New("недействительный refresh токен")
	}

	stored, err := s.tokenRepo.FindByToken(refreshTokenStr)
	if err != nil || stored == nil {
		return nil, errors.New("токен не найден или уже использован")
	}

	// удаляем старый токен
	s.tokenRepo.DeleteByToken(refreshTokenStr)

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("пользователь не найден")
	}

	return s.generateAndSaveTokens(user)
}

func (s *AuthService) generateAndSaveTokens(user *model.User) (*model.TokenPair, error) {
	roleID := 0
	if user.RoleID != nil {
		roleID = *user.RoleID
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, roleID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	err = s.tokenRepo.Save(user.ID, refreshToken, expiresAt)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
