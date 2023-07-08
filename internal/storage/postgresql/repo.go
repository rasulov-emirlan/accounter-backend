package postgresql

import (
	"context"
	"fmt"

	"github.com/SanaripEsep/esep-backend/config"
	"github.com/SanaripEsep/esep-backend/internal/storage/postgresql/migrations"
	"github.com/SanaripEsep/esep-backend/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

const PackageName = "internal/storage/postgresql/"

type RepositoryCombiner struct {
	ownersRepo     ownersRepository
	storesRepo     storesRepository
	categoriesRepo categoriesRepository
}

func NewRepositories(ctx context.Context, cfg config.Config, log *logging.Logger) (RepositoryCombiner, error) {
	c, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return RepositoryCombiner{}, fmt.Errorf("could not parse config: %w", err)
	}

	c.ConnConfig.Tracer = log

	conn, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return RepositoryCombiner{}, fmt.Errorf("could not connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return RepositoryCombiner{}, fmt.Errorf("could not ping database: %w", err)
	}

	if cfg.Flags.WithMigrations {
		if err := migrations.Up(ctx, cfg.DatabaseURL, log.Goosed()); err != nil {
			return RepositoryCombiner{}, fmt.Errorf("could not migrate database: %w", err)
		}
	}

	return RepositoryCombiner{
		ownersRepo:     ownersRepository{conn},
		storesRepo:     storesRepository{conn},
		categoriesRepo: categoriesRepository{conn},
	}, nil
}

func (r RepositoryCombiner) Owners() ownersRepository {
	return r.ownersRepo
}

func (r RepositoryCombiner) Stores() storesRepository {
	return r.storesRepo
}

func (r RepositoryCombiner) Categories() categoriesRepository {
	return r.categoriesRepo
}

func (r RepositoryCombiner) Close(ctx context.Context) error {
	r.ownersRepo.conn.Close()
	return nil
}
