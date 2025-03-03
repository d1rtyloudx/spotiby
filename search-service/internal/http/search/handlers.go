package search

import (
	"context"
	"github.com/d1rtyloudx/spotiby-pkg/lib"
	"github.com/d1rtyloudx/spotiby/search-service/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type searcher interface {
	Search(ctx context.Context, term string, paginationQuery lib.PaginationQuery) (model.SearchResponse, error)
}

type Handlers struct {
	searcher searcher
}

func New(searcher searcher) *Handlers {
	return &Handlers{
		searcher: searcher,
	}
}

func (h *Handlers) Search() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		term := ctx.Param("term")

		paginationQuery, err := lib.ExtractPageQueryParams(ctx)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		resp, err := h.searcher.Search(ctx.Request().Context(), term, paginationQuery)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to search"})
		}

		return ctx.JSON(http.StatusOK, echo.Map{"data": resp})
	}
}
