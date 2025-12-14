package ldap

import (
	"crypto/tls"
	"fmt"
	"sspr-ldap/config"
	"sspr-ldap/domain"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
)

type Repository struct {
	config config.LDAPConfig
}

func NewRepository(cfg config.LDAPConfig) *Repository {
	return &Repository{
		config: cfg,
	}
}

func (r *Repository) Authenticate(username, password string) (*domain.User, error) {

	// Decide whether to use LDAPS or plain LDAP based on config
	useTLS := r.config.UseTLS == "true"
	tlsSkip := r.config.TLSSkipVerify == "true"
	var l *ldap.Conn
	var err error
	if useTLS {
		l, err = ldap.DialURL(fmt.Sprintf("ldaps://%s:%s", r.config.Server, r.config.Port), ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: tlsSkip}))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to LDAPS: %w", err)
		}
	} else {
		l, err = ldap.DialURL(fmt.Sprintf("ldap://%s:%s", r.config.Server, r.config.Port))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to LDAP: %w", err)
		}
	}
	defer l.Close()

	// Bind with service account if configured
	if r.config.BindDN != "" {
		err = l.Bind(r.config.BindDN, r.config.BindPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to bind with service account: %w", err)
		}
	}

	// Search for user
	searchRequest := ldap.NewSearchRequest(
		r.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(r.config.UserFilter, ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search for user: %w", err)
	}
	if len(sr.Entries) != 1 {
		return nil, fmt.Errorf("user not found or too many entries returned")
	}

	// Set User DN from query
	userDN := sr.Entries[0].DN

	// Authenticate user
	err = l.Bind(userDN, password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &domain.User{
		Username: username,
		DN:       userDN,
	}, nil
}

func (r *Repository) ChangePassword(change *domain.PasswordChange) error {

	// Decide whether to use LDAPS or plain LDAP based on config
	useTLS := r.config.UseTLS == "true"
	tlsSkip := r.config.TLSSkipVerify == "true"
	var err error

	var conn *ldap.Conn
	if useTLS {
		conn, err = ldap.DialURL(fmt.Sprintf("ldaps://%s:%s", r.config.Server, r.config.Port), ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: tlsSkip}))
		if err != nil {
			return fmt.Errorf("failed to connect to LDAPS: %w", err)
		}
	} else {
		conn, err = ldap.DialURL(fmt.Sprintf("ldap://%s:%s", r.config.Server, r.config.Port))
		if err != nil {
			return fmt.Errorf("failed to connect to LDAP: %w", err)
		}
	}
	defer conn.Close()

	// Verify current password by attempting a bind as the user
	if err := conn.Bind(change.UserDN, change.CurrentPassword); err != nil {
		return fmt.Errorf("invalid current password")
	}

	// Re-bind as service account (required to perform password modifications in many LDAP servers)
	if r.config.BindDN != "" {
		if err := conn.Bind(r.config.BindDN, r.config.BindPassword); err != nil {
			return fmt.Errorf("failed to bind with service account: %w", err)
		}
	}

	// Encode and set new password
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	newEncodedPwd, encErr := utf16.NewEncoder().String("\"" + change.NewPassword + "\"")
	if encErr != nil {
		return fmt.Errorf("failed to encode new password: %w", encErr)
	}

	modifyRequest := ldap.NewModifyRequest(change.UserDN, nil)
	modifyRequest.Replace("unicodePwd", []string{newEncodedPwd})

	if err := conn.Modify(modifyRequest); err != nil {
		return fmt.Errorf("failed to modify password: %w", err)
	}

	return nil
}
