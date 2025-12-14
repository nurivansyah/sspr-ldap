package config

import "os"

type Config struct {
	Port       string
	SessionKey string
	LDAP       LDAPConfig
}

type LDAPConfig struct {
	Server        string
	Port          string
	BaseDN        string
	BindDN        string
	BindPassword  string
	UserFilter    string
	UseTLS        string // Enable TLS ("true"/"false")
	TLSSkipVerify string // Skip certificate verification ("true"/"false")
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		SessionKey: getEnv("SESSION_KEY", "default-secret-key-change-this"),
		LDAP: LDAPConfig{
			Server:       getEnv("LDAP_SERVER", "localhost"),
			Port:         getEnv("LDAP_PORT", "389"),
			BaseDN:       getEnv("LDAP_BASE_DN", "dc=example,dc=com"),
			BindDN:       getEnv("LDAP_BIND_DN", ""),
			BindPassword: getEnv("LDAP_BIND_PASSWORD", ""),
			UserFilter:   getEnv("LDAP_USER_FILTER", "(userPrincipalName=%s)"),
			UseTLS:       getEnv("LDAP_USE_TLS", "false"),
			// Default to not skipping certificate verification in production
			TLSSkipVerify: getEnv("LDAP_TLS_SKIP_VERIFY", "false"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
