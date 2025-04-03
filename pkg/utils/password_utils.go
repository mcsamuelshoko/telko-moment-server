package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string, cost ...int) (string, error) {
	// Default cost is 12 if not specified
	// 12 is a good work factor, adjust based on security needs
	assignedCost := 12

	// Use provided cost if available
	if len(cost) > 0 && cost[0] > 0 {
		assignedCost = cost[0]
	}

	// Make sure cost is within Bcrypt's allowed range (4-31)
	if assignedCost < bcrypt.MinCost {
		assignedCost = bcrypt.MinCost
	} else if assignedCost > bcrypt.MaxCost {
		assignedCost = bcrypt.MaxCost
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), assignedCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
