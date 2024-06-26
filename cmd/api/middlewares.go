package main

import (
	"crypto/sha256"
	"github.com/SuhailEdu/suhail-backend/models"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
		hashString := string(hash[:])

		userToken, err := models.Tokens(qm.Where("hash = ?", hashString), qm.Load("User"), qm.Limit(1)).AllG(c.Request().Context())

		if err != nil || len(userToken) == 0 {
			return unAuthorizedError(c, "Invalid Authorization token")
		}

		if userToken[0].Expiry.Before(time.Now()) {
			return unAuthorizedError(c, "Expired Authorization token")

		}

		c.Set("user", userToken[0])

		return next(c)

	}

}
