package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sspr-ldap/config"
	"sspr-ldap/handlers"
	"sspr-ldap/infra/ldap"
	"sspr-ldap/infra/session"
	tmpl "sspr-ldap/infra/template"
	"sspr-ldap/services"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: .env file not found, using system environment variables")
	} else {
		log.Println("✅ Loaded configuration from .env file")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize infrastructure
	sessionStore := session.NewStore(cfg.SessionKey)
	templateEngine, err := tmpl.NewEngine("templates/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}
	ldapRepo := ldap.NewRepository(cfg.LDAP)

	// Initialize services
	authService := services.NewAuthService(ldapRepo)
	userService := services.NewUserService(ldapRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, sessionStore, templateEngine)
	userHandler := handlers.NewUserHandler(userService, sessionStore, templateEngine)

	// Setup routes
	http.HandleFunc("/", authHandler.Home)
	http.HandleFunc("/login", authHandler.Login)
	http.HandleFunc("/logout", authHandler.Logout)
	http.HandleFunc("/dashboard", userHandler.Dashboard)
	http.HandleFunc("/change-password", userHandler.ChangePassword)

	// Start server with timeouts and graceful shutdown
	port := cfg.Port
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Run server
	go func() {
		log.Printf("🚀 Server starting on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown on interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped gracefully")
}
