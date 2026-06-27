package user

import (
	"fmt"
	"sportsync-api/internal/auth"
	"sportsync-api/internal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *service {
	return &service{repo, jwtService}
}

func (s *service) CreateUser(req dto.CreateRequest) (*dto.Response, error) {
	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	err := user.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil
}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.Response, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	err = user.checkPassword(req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Token:     token,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil
}
