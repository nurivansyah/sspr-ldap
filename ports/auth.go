package ports

import "sspr-ldap/domain"

// AuthRepository defines the interface for authentication operations
type AuthRepository interface {
	Authenticate(username, password string) (*domain.User, error)
}

// UserRepository defines the interface for user operations
type UserRepository interface {
	ChangePassword(change *domain.PasswordChange) error
}
