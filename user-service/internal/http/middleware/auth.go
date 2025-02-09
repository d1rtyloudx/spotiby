package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"user-service/internal/service/auth"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
)

type tokenParser interface {
	ParseAccessToken(tokenString string) (auth.AccessClaims, error)
}

type AuthMiddleware struct {
	parser tokenParser
	log    *zap.Logger
}

const (
	UserCredentialIDCtx = "credential_id"
	UserEmailCtx        = "user_email"
	UserProfileIDCtx    = "profile_id"
	UserRoleCtx         = "user_role"
)

func NewAuthMiddleware(parser tokenParser, log *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		parser: parser,
		log:    log,
	}
}

func (m *AuthMiddleware) ParseAndVerifyAccessToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")

			if header == "" {
				return c.JSON(http.StatusUnauthorized, errors.New("authorization header is required"))
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, errors.New("authorization header is invalid"))
			}

			if headerParts[1] == "" {
				return c.JSON(http.StatusUnauthorized, errors.New("token is empty"))
			}

			tokenString := headerParts[1]

			claims, err := m.parser.ParseAccessToken(tokenString)
			if err != nil {
				m.log.Warn(
					"failed to parse access token",
					zap.String("op", "middleware.Auth.ParseAndVerifyAccessToken"),
					zap.Error(err))
				return c.JSON(http.StatusUnauthorized, ErrUnauthenticated)
			}

			c.Set(UserCredentialIDCtx, claims.Subject)
			c.Set(UserEmailCtx, claims.Email)
			c.Set(UserProfileIDCtx, claims.ProfileID)
			c.Set(UserRoleCtx, claims.Role)

			return next(c)
		}
	}
}
