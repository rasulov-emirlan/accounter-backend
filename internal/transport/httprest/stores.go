package httprest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/esep-backend/internal/domains/stores"
)

type (
	StoresReadRequest struct {
		Text    string `query:"text"`
		OwnerID string `query:"ownerID"`

		// Pagination
		PageNumber uint64 `query:"pageNumber"`
		PageSize   int    `query:"pageSize"`

		// Sorting
		SortBy    string `query:"sortBy"`    // name, createdAt
		SortOrder string `query:"sortOrder"` // asc, desc
	}
)

type StoresHandler struct {
	storesService stores.Service
}

func (h StoresHandler) Create(ctx echo.Context) error {
	req := new(stores.CreateInput)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	if err := ctx.Validate(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	store, err := h.storesService.Create(ctx.Request().Context(), *req)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, store)
}

func (h StoresHandler) Read(ctx echo.Context) error {
	req := new(StoresReadRequest)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	in := stores.ReadByInput{}

	if req.Text != "" {
		in.Text.Set(req.Text)
	}
	if req.OwnerID != "" {
		in.OwnerID.Set(req.OwnerID)
	}
	if req.PageNumber != 0 {
		in.PageNumber.Set(req.PageNumber)
	}
	if req.PageSize != 0 {
		in.PageSize.Set(req.PageSize)
	}
	if req.SortBy != "" {
		in.SortBy.Set(req.SortBy)
	}
	if req.SortOrder != "" {
		in.SortOrder.Set(req.SortOrder)
	}

	store, err := h.storesService.ReadBy(ctx.Request().Context(), in)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, store)
}

func (h StoresHandler) Update(ctx echo.Context) error {
	req := make(map[string]any)
	if err := ctx.Bind(&req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	in := h.mapToUpdateInput(req)
	id := ctx.Param("id")

	s, err := h.storesService.Update(ctx.Request().Context(), id, in)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, s)
}

func (h StoresHandler) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := h.storesService.Delete(ctx.Request().Context(), id); err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h StoresHandler) mapToUpdateInput(in map[string]any) stores.UpdateInput {
	out := stores.UpdateInput{}

	if v, ok := in["name"]; ok {
		out.Name.Set(v.(string))
	}
	if v, ok := in["description"]; ok {
		out.Description.Set(v.(string))
	}

	return out
}
