package profile

import (
	"context"
	"errors"
	"github.com/d1rtyloudx/spotiby-pkg/lib"
	"github.com/labstack/echo/v4"
	"net/http"
	"user-service/internal/dto"
	"user-service/internal/http/middleware"
	"user-service/internal/service/profile"
)

type profileService interface {
	Get(ctx context.Context, pageQuery lib.PaginationQuery) (dto.PagedProfileResponse, error)
	GetByID(ctx context.Context, id string) (dto.Profile, error)
	Update(ctx context.Context, id string, req dto.UpdateProfileRequest) error
	GetFollows(ctx context.Context, followerID string, pageQuery lib.PaginationQuery) (dto.PagedProfileResponse, error)
	FollowProfile(ctx context.Context, followerID string, followeeID string) error
	UnfollowProfile(ctx context.Context, followerID string, followeeID string) error
}

type Handlers struct {
	profileService profileService
}

func New(profileService profileService) *Handlers {
	return &Handlers{
		profileService: profileService,
	}
}

func (h *Handlers) Get() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pageQuery, err := lib.ExtractPageQueryParams(ctx)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}

		resp, err := h.profileService.Get(ctx.Request().Context(), pageQuery)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to get profiles",
			})
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"items": resp,
		})
	}
}

func (h *Handlers) GetMe() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		profileID := ctx.Get(middleware.UserProfileIDCtx).(string)

		resp, err := h.profileService.GetByID(ctx.Request().Context(), profileID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to get profile",
			})
		}

		return ctx.JSON(http.StatusOK, resp)
	}
}

func (h *Handlers) GetByID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		profileID := ctx.Param("id")
		if profileID == "" {
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": "the requested profile id cannot be empty",
			})
		}

		resp, err := h.profileService.GetByID(ctx.Request().Context(), profileID)
		if err != nil {
			if errors.Is(err, profile.ErrProfileNotFound) {
				return ctx.JSON(http.StatusNotFound, echo.Map{
					"error": err.Error(),
				})
			}

			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to get profile by id",
			})
		}

		return ctx.JSON(http.StatusOK, resp)
	}
}

func (h *Handlers) Update() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		profileID := ctx.Get(middleware.UserProfileIDCtx).(string)

		var req dto.UpdateProfileRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		err := h.profileService.Update(ctx.Request().Context(), profileID, req)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to update profile",
			})
		}

		return ctx.NoContent(http.StatusOK)
	}
}

func (h *Handlers) GetFollows() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		profileID := ctx.Get(middleware.UserProfileIDCtx).(string)

		pageQuery, err := lib.ExtractPageQueryParams(ctx)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}

		followers, err := h.profileService.GetFollows(ctx.Request().Context(), profileID, pageQuery)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to get followers",
			})
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"items": followers,
		})
	}
}

func (h *Handlers) FollowProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		profileID := ctx.Get(middleware.UserProfileIDCtx).(string)

		followeeID := ctx.Param("id")
		if followeeID == "" || profileID == followeeID {
			return ctx.NoContent(http.StatusBadRequest)
		}

		err := h.profileService.FollowProfile(ctx.Request().Context(), profileID, followeeID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to follow profile",
			})
		}

		return ctx.NoContent(http.StatusOK)
	}
}

func (h *Handlers) UnfollowProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		profileID := ctx.Get(middleware.UserProfileIDCtx).(string)

		followeeID := ctx.Param("id")
		if followeeID == "" || profileID == followeeID {
			return ctx.NoContent(http.StatusBadRequest)
		}

		err := h.profileService.UnfollowProfile(ctx.Request().Context(), profileID, followeeID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to unfollow profile",
			})
		}

		return ctx.NoContent(http.StatusOK)
	}
}
