package image

import (
	"context"
	"errors"
	"fmt"
	"github.com/d1rtyloudx/spotiby/user-service/internal/config"
	"github.com/d1rtyloudx/spotiby/user-service/internal/domain/model"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type imageUploader interface {
	UploadProfile(ctx context.Context, id string, image model.Image) (string, error)
	UploadPlaylist(ctx context.Context, id string, image model.Image) (string, error)
	UploadTrack(ctx context.Context, id string, image model.Image) (string, error)
}

type imageProvider interface {
	GetProfile(ctx context.Context, buketName, fileName string) (*minio.Object, error)
	GetPlaylist(ctx context.Context, buketName, fileName string) (*minio.Object, error)
	GetTrack(ctx context.Context, buketName, fileName string) (*minio.Object, error)
}

type Handlers struct {
	uploader imageUploader
	provider imageProvider
	log      *zap.Logger
	cfg      *config.MinioBucketsConfig
}

func New(uploader imageUploader, provider imageProvider, cfg *config.MinioBucketsConfig, log *zap.Logger) *Handlers {
	return &Handlers{
		uploader: uploader,
		provider: provider,
		cfg:      cfg,
		log:      log,
	}
}

func (h *Handlers) getImage(bucketName string, get func(ctx context.Context, bucketName, fileName string) (*minio.Object, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fileName := ctx.Param("fileName")
		if fileName == "" {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "filename cannot be empty"})
		}

		obj, err := get(ctx.Request().Context(), bucketName, fileName)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get image"})
		}
		defer obj.Close()

		info, err := obj.Stat()
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get object info"})
		}

		ctx.Response().Header().Set("Content-Type", info.ContentType)
		ctx.Response().Header().Set("Content-Length", fmt.Sprintf("%d", info.Size))
		ctx.Response().WriteHeader(http.StatusOK)

		_, err = io.Copy(ctx.Response().Writer, obj)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to stream image"})
		}

		return nil
	}
}

func (h *Handlers) uploadImage(bucketName string, upload func(ctx context.Context, id string, image model.Image) (string, error)) echo.HandlerFunc {
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

		fileName, err := upload(ctx.Request().Context(), id, image)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to upload image",
			})
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"image": fileName,
		})
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

func (h *Handlers) GetProfile() echo.HandlerFunc {
	return h.getImage(h.cfg.ProfileBucket, h.provider.GetProfile)
}

func (h *Handlers) GetTrack() echo.HandlerFunc {
	return h.getImage(h.cfg.TrackBucket, h.provider.GetTrack)
}

func (h *Handlers) GetPlaylist() echo.HandlerFunc {
	return h.getImage(h.cfg.PlaylistBucket, h.provider.GetPlaylist)
}
