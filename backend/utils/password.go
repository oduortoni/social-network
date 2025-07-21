package utils

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type PasswordConfig struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
	BcryptCost     int
}

type PasswordManager struct {
	Config PasswordConfig
}

/*
* NewPasswordManager
* creates a new PasswordManager with the given config.
 */
func NewPasswordManager(config PasswordConfig) *PasswordManager {
	// Set sensible defaults if not provided
	if config.MinLength == 0 {
		config.MinLength = 8
	}
	if config.BcryptCost == 0 {
		config.BcryptCost = bcrypt.DefaultCost
	}
	config.RequireLower=true;
	config.RequireSpecial=true;
	config.RequireUpper=true;
	config.RequireNumber=true;

	return &PasswordManager{Config: config}
}

/*
* HashPassword
* hashes the plain-text password after validating strength.
 */
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	if err := pm.ValidatePasswordStrength(password); err != nil {
		return "", err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), pm.Config.BcryptCost)
	return string(hashed), err
}

/*
* ComparePassword
* compares the hashed password with a plain-text input.
 */
func (pm *PasswordManager) ComparePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

/*
* ValidatePasswordStrength
* checks if the password meets configured rules.
 */
func (pm *PasswordManager) ValidatePasswordStrength(password string) error {
	if len(password) < pm.Config.MinLength {
		return errors.New("password must be at least 8 characters")
	}
	if pm.Config.RequireUpper && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain an uppercase letter")
	}
	if pm.Config.RequireLower && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain a lowercase letter")
	}
	if pm.Config.RequireNumber && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain a number")
	}
	if pm.Config.RequireSpecial && !regexp.MustCompile(`[!@#~$%^&*()+|_.,<>?/\\-]`).MatchString(password) {
		return errors.New("password must contain a special character")
	}
	return nil
}
