package utils

import "github.com/gofiber/fiber/v2"

// SuccessResponse creates a standardized success response
func SuccessResponse(data interface{}, message string) fiber.Map {
	return fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	}
}

// ErrorResponse creates a standardized error response
func ErrorResponse(message string, details ...interface{}) fiber.Map {
	resp := fiber.Map{
		"success": false,
		"message": message,
	}

	if len(details) > 0 {
		resp["details"] = details[0]
	}

	return resp
}

// PaginatedResponse creates a response with pagination info
func PaginatedResponse(data interface{}, page, pageSize, total int) fiber.Map {
	return fiber.Map{
		"success": true,
		"data":    data,
		"pagination": fiber.Map{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + pageSize - 1) / pageSize,
		},
	}
}
