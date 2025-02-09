package image

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"image-service/internal/config"
	"image-service/internal/domain/model"
	"net/http"
)

type imageUploader interface {
	UploadProfile(ctx context.Context, id string, image model.Image) error
	UploadPlaylist(ctx context.Context, id string, image model.Image) error
	UploadTrack(ctx context.Context, id string, image model.Image) error
}

type Handlers struct {
	uploader imageUploader
	log      *zap.Logger
	cfg      *config.MinioBucketsConfig
}

func New(uploader imageUploader, cfg *config.MinioBucketsConfig, log *zap.Logger) *Handlers {
	return &Handlers{
		uploader: uploader,
		cfg:      cfg,
		log:      log,
	}
}

func (h *Handlers) uploadImage(bucketName string, upload func(ctx context.Context, id string, image model.Image) error) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")
		if id == "" {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": errors.New("id cannot be empty")})
		}

		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}
		defer file.Close()

		image := model.Image{
			File:        file,
			Name:        fileHeader.Filename,
			Size:        fileHeader.Size,
			ContentType: fileHeader.Header.Get("Content-Type"),
			BucketName:  bucketName,
		}

		err = upload(ctx.Request().Context(), id, image)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to upload image",
			})
		}

		return ctx.NoContent(http.StatusOK)
	}
}

func (h *Handlers) UploadProfile() echo.HandlerFunc {
	return h.uploadImage(h.cfg.ProfileBucket, h.uploader.UploadProfile)
}

func (h *Handlers) UploadTrack() echo.HandlerFunc {
	return h.uploadImage(h.cfg.TrackBucket, h.uploader.UploadTrack)
}

func (h *Handlers) UploadPlaylist() echo.HandlerFunc {
	return h.uploadImage(h.cfg.PlaylistBucket, h.uploader.UploadPlaylist)
}
