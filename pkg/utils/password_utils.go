package utils

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
	// 12 is a good work factor, adjust based on your security needs
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
