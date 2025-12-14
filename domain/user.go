package domain

// User represents a user in the system
type User struct {
	Username string
	DN       string // Distinguished Name
}

// Credentials represents user authentication credentials
type Credentials struct {
	Username string
	Password string
}

// PasswordChange represents a password change request
type PasswordChange struct {
	Username        string
	UserDN          string
	CurrentPassword string
	NewPassword     string
}
