package main

import (
	"crypto/sha256"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

func (config *Config) checkAuthToken(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return unAuthorizedError(c, "Invalid Authorization header")
		}

		authHeaderParts := strings.SplitN(authHeader, " ", 2)

		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			return unAuthorizedError(c, "Invalid Authorization header")
		}
		hash := sha256.Sum256([]byte(authHeaderParts[1]))

		userToken, err := config.db.GetUserByToken(c.Request().Context(), hash[:])

		if err != nil {
			return unAuthorizedError(c, "Invalid Authorization token")
		}

		if userToken.Expiry.Time.Before(time.Now()) {
			return unAuthorizedError(c, "Expired Authorization token")
		}

		c.Set("user", userToken)

		return next(c)

	}

}
