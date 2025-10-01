package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("email inválido")
	ErrInvalidUsername = errors.New("username debe tener entre 3 y 30 caracteres alfanuméricos")
	ErrPasswordTooWeak = errors.New("la contraseña debe tener al menos 8 caracteres")
)

type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	AvatarURL    string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Bio          string    `json:"bio,omitempty" db:"bio"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Estructura para registro de usuario
type RegisterUser struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Estructura para login de usuario
type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HashPassword genera un hash bcrypt de la contraseña
func (u *User) HashPassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooWeak
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifica si la contraseña es correcta
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// ValidateEmail verifica si el email es válido
func (u *User) ValidateEmail() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateUsername verifica si el username es válido
func (u *User) ValidateUsername() error {
	username := strings.TrimSpace(u.Username)
	if len(username) < 3 || len(username) > 30 {
		return ErrInvalidUsername
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}

	return nil
}

// Validate ejecuta todas las validaciones
func (u *User) Validate() error {
	if err := u.ValidateEmail(); err != nil {
		return err
	}
	if err := u.ValidateUsername(); err != nil {
		return err
	}
	return nil
}