package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rasulov-emirlan/esep-backend/config"
	"github.com/rasulov-emirlan/esep-backend/internal/storage/postgresql/migrations"
	"github.com/rasulov-emirlan/esep-backend/pkg/logging"
)

type RepositoryCombiner struct {
	ownersRepo ownersRepository
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

	ownersRepo := ownersRepository{conn}
	return RepositoryCombiner{
		ownersRepo: ownersRepo,
	}, nil
}

func (r RepositoryCombiner) Owners() ownersRepository {
	return r.ownersRepo
}

func (r RepositoryCombiner) Close() {
	r.ownersRepo.conn.Close()
}
