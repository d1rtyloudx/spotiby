package auth

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"user-service/internal/dto"
	"user-service/internal/http/middleware"
	"user-service/internal/service/auth"
)

type authService interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, token string) (dto.TokenPair, error)
	UpdatePassword(ctx context.Context, id string, password string) error
	UpdateUsername(ctx context.Context, id string, username string) error
}

type Handlers struct {
	authService authService
}

func New(authService authService) *Handlers {
	return &Handlers{
		authService: authService,
	}
}

func (h *Handlers) Login() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var req dto.LoginRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		resp, err := h.authService.Login(ctx.Request().Context(), req)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				return ctx.JSON(http.StatusBadRequest, echo.Map{
					"error": err.Error(),
				})
			}

			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to login",
			})
		}

		return ctx.JSON(http.StatusOK, resp)
	}
}

func (h *Handlers) Register() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var req dto.RegisterRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		resp, err := h.authService.Register(ctx.Request().Context(), req)
		if err != nil {
			if errors.Is(err, auth.ErrUserAlreadyRegistered) {
				return ctx.JSON(http.StatusBadRequest, echo.Map{
					"error": err.Error(),
				})
			}

			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to register",
			})
		}

		return ctx.JSON(http.StatusOK, resp)
	}
}

func (h *Handlers) Logout() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var req dto.LogoutRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		err := h.authService.Logout(ctx.Request().Context(), req.RefreshToken)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, echo.Map{
				"error": "failed to logout",
			})
		}

		return ctx.NoContent(http.StatusOK)
	}
}

func (h *Handlers) UpdatePassword() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		credentialID := ctx.Get(middleware.UserCredentialIDCtx).(string)

		var req dto.UpdatePasswordRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		err := h.authService.UpdatePassword(ctx.Request().Context(), credentialID, req.Password)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to update password",
			})
		}

		return ctx.NoContent(http.StatusOK)
	}
}

func (h *Handlers) UpdateUsername() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		credentialID := ctx.Get(middleware.UserCredentialIDCtx).(string)

		var req dto.UpdateUsernameRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		err := h.authService.UpdateUsername(ctx.Request().Context(), credentialID, req.Username)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to update username",
			})
		}

		return ctx.NoContent(http.StatusOK)
	}
}

func (h *Handlers) IntrospectToken() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	}
}

func (h *Handlers) RefreshToken() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var req dto.RefreshTokenRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		tokens, err := h.authService.RefreshToken(ctx.Request().Context(), req.RefreshToken)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, echo.Map{
				"error": "failed to refresh token",
			})
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"tokens": tokens,
		})
	}
}
