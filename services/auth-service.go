package services

import (
	"errors"
	"sspr-ldap/domain"
	"sspr-ldap/ports"
)

type AuthService struct {
	repo ports.AuthRepository
}

func NewAuthService(repo ports.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Authenticate(creds *domain.Credentials) (*domain.User, error) {
	if creds.Username == "" || creds.Password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.repo.Authenticate(creds.Username, creds.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
