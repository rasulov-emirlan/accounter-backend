package httprest

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/accounter-backend/internal/domains/categories"
)

type (
	CategoriesReadByRequest struct {
		PageNumber       uint64 `query:"pageNumber"`
		PageSize         uint   `query:"pageSize"`
		Text             string `query:"text"`
		StoreID          string `query:"storeID"`
		ParentCategoryID string `query:"parentCategoryID"`
		SortBy           string `query:"sortBy"`
		SortOrder        string `query:"sortOrder"`
	}
)

type CategoriesHandler struct {
	categoriesService categories.Service
}

func (h CategoriesHandler) Create(ctx echo.Context) error {
	req := new(categories.CreateInput)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}
	if err := ctx.Validate(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	category, err := h.categoriesService.Create(ctx.Request().Context(), *req)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, category)
}

func (h CategoriesHandler) ReadBy(ctx echo.Context) error {
	req := new(CategoriesReadByRequest)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	in := categories.ReadByInput{}
	if req.PageNumber != 0 {
		in.PageNumber.Set(req.PageNumber)
	}
	if req.PageSize != 0 {
		in.PageSize.Set(req.PageSize)
	}
	if req.Text != "" {
		in.Text.Set(req.Text)
	}
	if req.StoreID != "" {
		in.StoreID.Set(req.StoreID)
	}
	if req.ParentCategoryID != "" {
		in.ParentCategoryID.Set(req.ParentCategoryID)
	}
	if req.SortBy != "" {
		in.SortBy.Set(req.SortBy)
	}

	res, err := h.categoriesService.ReadBy(ctx.Request().Context(), in)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h CategoriesHandler) Read(ctx echo.Context) error {
	id := ctx.Param("id")

	in := categories.ReadByInput{}
	in.ID.Set(id)

	res, err := h.categoriesService.ReadBy(ctx.Request().Context(), in)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h CategoriesHandler) Update(ctx echo.Context) error {
	req := make(echo.Map)
	if err := ctx.Bind(&req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	id := ctx.Param("id")

	in := categories.UpdateInput{ID: id}

	if v, ok := req["name"]; ok {
		tmp, ok := v.(string)
		if !ok {
			return respondErr(ctx, http.StatusBadRequest, errors.New("поле name должно быть строкой"))
		}
		in.Name.Set(tmp)
	}
	if v, ok := req["article"]; ok {
		tmp, ok := v.(string)
		if !ok {
			return respondErr(ctx, http.StatusBadRequest, errors.New("поле article должно быть строкой"))
		}
		in.Article.Set(&tmp)
	}
	if v, ok := req["parentCategoryID"]; ok {
		tmp, ok := v.(string)
		if !ok {
			return respondErr(ctx, http.StatusBadRequest, errors.New("поле parentCategoryID должно быть строкой"))
		}
		in.ParentCategoryID.Set(&tmp)
	}

	category, err := h.categoriesService.Update(ctx.Request().Context(), in)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, category)
}

func (h CategoriesHandler) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := h.categoriesService.Delete(ctx.Request().Context(), id); err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.NoContent(http.StatusOK)
}
