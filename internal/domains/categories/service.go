package categories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rasulov-emirlan/accounter-backend/internal/entities"
	"github.com/rasulov-emirlan/accounter-backend/pkg/logging"
	"github.com/rasulov-emirlan/accounter-backend/pkg/telemetry"
)

type (
	CategoriesRepository interface {
		Create(ctx context.Context, input entities.Category) (entities.Category, error)
		ReadBy(ctx context.Context, filters ReadByInput) ([]entities.Category, error)
		Update(ctx context.Context, changeset UpdateInput) (entities.Category, error)
		Delete(ctx context.Context, id string) error
	}

	Service interface {
		Create(ctx context.Context, input CreateInput) (entities.Category, error)
		ReadBy(ctx context.Context, filters ReadByInput) ([]entities.Category, error)
		Update(ctx context.Context, changeset UpdateInput) (entities.Category, error)
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo CategoriesRepository
		log  *logging.Logger
	}
)

var _ Service = (*service)(nil)

func NewService(repo CategoriesRepository, log *logging.Logger) service {
	return service{repo: repo, log: log}
}

func (s service) Create(ctx context.Context, input CreateInput) (entities.Category, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Create")).End()
	defer s.log.Sync()
	category := entities.Category{
		Name:    input.Name,
		Article: input.Article,
		Store:   &entities.Store{ID: uuid.MustParse(input.StoreID)},
	}
	if input.ParentCategoryID != nil {
		category.ParentCategory = &entities.Category{ID: uuid.MustParse(*input.ParentCategoryID)}
	}
	category, err := s.repo.Create(ctx, category)
	if err != nil {
		s.log.Debug("categories:Create - failed to create category", logging.String("stage", "repository"), logging.Error("err", err))
		return entities.Category{}, ErrDefault
	}

	s.log.Info("categories:Create - category created", logging.String("stage", "repository"), logging.String("categoryID", category.ID.String()))
	return category, nil
}

func (s service) ReadBy(ctx context.Context, filters ReadByInput) ([]entities.Category, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.ReadBy")).End()
	defer s.log.Sync()

	// validate filters
	pageNumber, ok := filters.PageNumber.Get()
	if !ok {
		filters.PageNumber.Set(1)
	} else if pageNumber < 1 {
		s.log.Debug("categories:ReadBy - pageNumber must be greater than 0", logging.String("stage", "validation"))
		return nil, errors.New("номер страницы должен быть меньше 1")
	}

	pageSize, ok := filters.PageSize.Get()
	if !ok {
		filters.PageSize.Set(10)
	} else if pageSize < 1 || pageSize > 100 {
		s.log.Debug("categories:ReadBy - pageSize must be between 1 and 100", logging.String("stage", "validation"))
		return nil, errors.New("размер страницы должен быть между 1 и 100")
	}

	text, _ := filters.Text.Get()
	if len(text) > 255 {
		s.log.Debug("categories:ReadBy - text must be less than 255 characters", logging.String("stage", "validation"))
		return nil, errors.New("текст должен быть меньше 255 символов")
	}

	sortBy, ok := filters.SortBy.Get()
	if ok {
		switch sortBy {
		case SortByName, SortByArticle, SortByCreatedAt:
		default:
			s.log.Debug("categories:ReadBy - sortBy must be one of name, article, createdAt", logging.String("stage", "validation"))
			return nil, errors.New("сортировка должна быть одной из name, article, createdAt")
		}
	} else {
		filters.SortBy.Set(SortByCreatedAt)
	}

	sortOrder, ok := filters.SortOrder.Get()
	if ok {
		switch sortOrder {
		case SortOrderAsc, SortOrderDesc:
		default:
			s.log.Debug("categories:ReadBy - sortOrder must be one of asc, desc", logging.String("stage", "validation"))
			return nil, errors.New("сортировка должна быть одной из asc, desc")
		}
	} else {
		filters.SortOrder.Set(SortOrderDesc)
	}

	categories, err := s.repo.ReadBy(ctx, filters)
	if err != nil {
		s.log.Error("categories:ReadBy - failed to read categories", logging.String("stage", "repository"), logging.Error("err", err))
		return nil, ErrDefault
	}

	s.log.Info("categories:ReadBy - categories read", logging.String("stage", "repository"), logging.Int("count", len(categories)))
	return categories, nil
}

func (s service) Update(ctx context.Context, changeset UpdateInput) (entities.Category, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Update")).End()
	defer s.log.Sync()

	// validate changeset
	name, ok := changeset.Name.Get()
	if ok && len(name) > 255 {
		s.log.Debug("categories:Update - name must be less than 255 characters", logging.String("stage", "validation"))
		return entities.Category{}, errors.New("имя должно быть меньше 255 символов")
	}

	article, ok := changeset.Article.Get()
	if ok && article != nil && len(*article) > 100 {
		s.log.Debug("categories:Update - article must be less than 100 characters", logging.String("stage", "validation"))
		return entities.Category{}, errors.New("артикул должен быть меньше 100 символов")
	}

	c, err := s.repo.Update(ctx, changeset)
	if err != nil {
		s.log.Error("categories:Update - failed to update category", logging.String("stage", "repository"), logging.Error("err", err))
		return entities.Category{}, ErrDefault
	}

	s.log.Info("categories:Update - category updated", logging.String("stage", "repository"), logging.String("categoryID", c.ID.String()))
	return c, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Delete")).End()
	defer s.log.Sync()

	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("categories:Delete - failed to delete category", logging.String("stage", "repository"), logging.Error("err", err))
		return ErrDefault
	}

	s.log.Info("categories:Delete - category deleted", logging.String("stage", "repository"), logging.String("categoryID", id))
	return nil
}
