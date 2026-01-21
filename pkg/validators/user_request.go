package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	hasUpperRegex   = regexp.MustCompile(`[A-Z]`)
	hasLowerRegex   = regexp.MustCompile(`[a-z]`)
	hasNumberRegex  = regexp.MustCompile(`[0-9]`)
	hasSpecialRegex = regexp.MustCompile(`[@#$%!^&*()_+]`)
)

// PasswordStrength is a custom validator function that checks if a password
// meets strength requirements: at least one uppercase letter, one lowercase letter,
// one number, and one special character.
//
// Parameters:
//   - fl: The field level information provided by the validator
//
// Returns:
//   - bool: true if the password meets the strength requirements, false otherwise
func PasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasUpper := hasUpperRegex.MatchString(password)
	hasLower := hasLowerRegex.MatchString(password)
	hasNumber := hasNumberRegex.MatchString(password)
	hasSpecial := hasSpecialRegex.MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}
