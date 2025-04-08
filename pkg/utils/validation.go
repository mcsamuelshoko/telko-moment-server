package utils

import "regexp"

// IsValidEmail checks if a string is a valid email
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// IsValidPhoneNumber validates phone numbers (customize for your format)
func IsValidPhoneNumber(phone string) bool {
	pattern := `^\+?[0-9]{10,15}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(phone)
}

// IsStrongPassword checks password strength
func IsStrongPassword(password string) bool {
	// At least 8 chars, 1 uppercase, 1 lowercase, 1 number, 1 special char
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}
