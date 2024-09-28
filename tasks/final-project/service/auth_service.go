package service

import (
	"final-project/domain"
	"os"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// GenerateJWT generates a JWT token for the admin
func GenerateJWT(admin domain.Admin) (string, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"admin_id": admin.ID,
		"email":    admin.Email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}

// HashPassword hashes a password using bcrypt/crypto
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a plain-text password with its hashed version
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Custom password validator
func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Min. 8 characters long
	if len(password) < 8 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)                         // Min. 1 lowercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)                         // Min. 1 uppercase letter
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*()_+{}":;'?/>.<,]`).MatchString(password) // Min. 1 special character

	return hasLower && hasUpper && hasSpecial
}
