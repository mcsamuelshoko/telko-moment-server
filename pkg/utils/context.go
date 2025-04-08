package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strings"
)

// GetHeaderMap extracts all headers from the Fiber context into a map.
func GetHeaderMap(c *fiber.Ctx) map[string]string {
	headerMap := make(map[string]string)
	c.Context().Request.Header.VisitAll(func(key, value []byte) {
		headerMap[string(key)] = string(value)
	})
	return headerMap
}

// GetIP gets the client IP address
func GetIP(c *fiber.Ctx) string {
	return c.IP()
}

// GetProtocol gets the request protocol (http or https)
func GetProtocol(c *fiber.Ctx) string {
	if c.Secure() {
		return "https"
	}
	return "http"
}

// GetHost gets the Host header
func GetHost(c *fiber.Ctx) string {
	return c.Hostname()
}

// GetClientIP extracts the client's IP address from the request
func GetClientIP(c *fiber.Ctx) string {
	ip := c.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.IP()
	}
	return strings.Split(ip, ", ")[0]
}

// GetUserAgent extracts the user agent
func GetUserAgent(c *fiber.Ctx) string {
	return c.Get("User-Agent")
}

// GetRequestID gets or generates a request ID
func GetRequestID(c *fiber.Ctx) string {
	requestID := c.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
		c.Set("X-Request-ID", requestID)
	}
	return requestID
}
