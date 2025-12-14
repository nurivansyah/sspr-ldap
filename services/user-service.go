package services

import (
	"errors"
	"sspr-ldap/domain"
	"sspr-ldap/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) ChangePassword(change *domain.PasswordChange) error {
	if change.NewPassword == "" {
		return errors.New("new password cannot be empty")
	}

	if len(change.NewPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return s.repo.ChangePassword(change)
}
