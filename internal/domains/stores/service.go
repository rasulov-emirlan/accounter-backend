package stores

import (
	"context"
	"errors"

	"github.com/SanaripEsep/esep-backend/internal/entities"
	"github.com/SanaripEsep/esep-backend/pkg/logging"
	"github.com/SanaripEsep/esep-backend/pkg/telemetry"
	"github.com/google/uuid"
)

type (
	// TODO: it is actually a bad practice to use types that are not buisness entities in repos, but for now fuck it
	StoresRepository interface {
		Create(ctx context.Context, store entities.Store) (entities.Store, error)
		ReadBy(ctx context.Context, filter ReadByInput) ([]entities.Store, error)
		Update(ctx context.Context, id string, changeset UpdateInput) (entities.Store, error)
		Delete(ctx context.Context, id string) error
	}

	Service interface {
		Create(ctx context.Context, input CreateInput) (entities.Store, error)
		ReadBy(ctx context.Context, filter ReadByInput) ([]entities.Store, error)
		Update(ctx context.Context, id string, input UpdateInput) (entities.Store, error)
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo StoresRepository
		log  *logging.Logger
	}
)

var _ Service = (*service)(nil)

func NewService(repo StoresRepository, log *logging.Logger) service {
	return service{repo: repo, log: log}
}

func (s service) Create(ctx context.Context, input CreateInput) (entities.Store, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Create")).End()
	defer s.log.Sync()
	ownerID, err := uuid.Parse(input.OwnerID)
	if err != nil {
		s.log.Debug("stores:Create - failed to parse owner id", logging.String("stage", "validation"), logging.Error("err", err))
		return entities.Store{}, errors.New("id владельца не валиден")
	}
	store := entities.Store{
		Name:        input.Name,
		Description: input.Description,
		Owner:       &entities.Owner{ID: ownerID},
	}

	store, err = s.repo.Create(ctx, store)
	if err != nil {
		s.log.Debug("stores:Create - failed to create store", logging.String("stage", "repository"), logging.Error("err", err))
		return entities.Store{}, ErrDefault
	}

	s.log.Info("stores:Create - store created", logging.String("stage", "repository"), logging.String("storeID", store.ID.String()), logging.String("ownerID", store.Owner.ID.String()))
	return store, nil
}

func (s service) ReadBy(ctx context.Context, filter ReadByInput) ([]entities.Store, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.ReadBy")).End()
	defer s.log.Sync()

	// validate filters
	val, ok := filter.PageNumber.Get()
	if !ok {
		filter.PageNumber.Set(1)
	} else if val < 1 {
		s.log.Debug("stores:ReadBy - invalid page number", logging.String("stage", "validation"), logging.Uint64("pageNumber", val))
		return nil, errors.New("номер страницы не может быть меньше 1")
	}

	val1, ok := filter.PageSize.Get()
	if !ok {
		filter.PageSize.Set(10)
	} else if val < 1 || val > 100 {
		s.log.Debug("stores:ReadBy - invalid page size", logging.String("stage", "validation"), logging.Uint("pageSize", val1))
		return nil, errors.New("размер страницы должен быть в диапазоне от 1 до 100")
	}

	// filter
	stores, err := s.repo.ReadBy(ctx, filter)
	if err != nil {
		s.log.Debug("stores:ReadBy - failed to read stores", logging.String("stage", "repository"), logging.Error("err", err))
		return nil, ErrDefault
	}

	s.log.Info("stores:ReadBy - stores read", logging.String("stage", "repository"), logging.Int("count", len(stores)))
	return stores, nil
}

func (s service) Update(ctx context.Context, id string, input UpdateInput) (entities.Store, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Update")).End()
	defer s.log.Sync()

	// validate
	countChanges := 0
	val, ok := input.Name.Get()
	if ok {
		countChanges++
		if len(val) < 3 {
			s.log.Debug("stores:Update - invalid name", logging.String("stage", "validation"), logging.String("name", val))
			return entities.Store{}, errors.New("название магазина должно содержать минимум 3 символа")
		}
	}

	val, ok = input.Description.Get()
	if ok {
		countChanges++
		if len(val) < 3 {
			s.log.Debug("stores:Update - invalid description", logging.String("stage", "validation"), logging.String("description", val))
			return entities.Store{}, errors.New("описание магазина должно содержать минимум 3 символа")
		}
	}

	if countChanges == 0 {
		s.log.Debug("stores:Update - no changes", logging.String("stage", "validation"))
		return entities.Store{}, errors.New("не переданы изменения")
	}

	// update
	store, err := s.repo.Update(ctx, id, input)
	if err != nil {
		s.log.Debug("stores:Update - failed to update store", logging.String("stage", "repository"), logging.Error("err", err))
		return entities.Store{}, ErrDefault
	}

	s.log.Info("stores:Update - store updated", logging.String("stage", "repository"), logging.String("storeID", store.ID.String()))
	return store, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Delete")).End()
	defer s.log.Sync()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Debug("stores:Delete - failed to delete store", logging.String("stage", "repository"), logging.Error("err", err))
		return ErrDefault
	}

	s.log.Info("stores:Delete - store deleted", logging.String("stage", "repository"), logging.String("storeID", id))
	return nil
}
