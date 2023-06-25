package postgresql

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rasulov-emirlan/esep-backend/internal/domains/stores"
	"github.com/rasulov-emirlan/esep-backend/internal/entities"
)

type storesRepository struct {
	conn *pgxpool.Pool
}

func (r storesRepository) Create(ctx context.Context, store entities.Store) (entities.Store, error) {
	if store.Owner == nil {
		return entities.Store{}, errors.New("owner is required")
	}
	store.CreatedAt = time.Now()
	sql, args, err := sq.Insert("stores").
		Columns("owner_id", "name", "description", "created_at", "tsv").
		Values(
			store.Owner.ID, store.Name, store.Description, store.CreatedAt,
			sq.Expr(
				`setweight(to_tsvector(?), 'A') || setweight(to_tsvector(?), 'B')`,
				store.Name, store.Description,
			)).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).ToSql()
	// dont forget to add 'tsv' column to the list of columns

	if err != nil {
		return entities.Store{}, err
	}

	row := r.conn.QueryRow(ctx, sql, args...)
	if err := row.Scan(&store.ID); err != nil {
		return entities.Store{}, err
	}

	return store, nil
}

var storeSortingFields = map[string]string{
	stores.SortByCreatedAt: "stores.created_at",
	stores.SortByName:      "stores.name",
}

func (r storesRepository) ReadBy(ctx context.Context, filter stores.ReadByInput) ([]entities.Store, error) {
	query := sq.Select("stores.id", "owner_id", "owners.full_name", "owners.username", "owners.created_at", "name", "description", "stores.created_at").
		LeftJoin("owners ON owners.id = stores.owner_id").
		From("stores").
		PlaceholderFormat(sq.Dollar)

	val, ok := filter.ID.Get()
	if !ok {
		val, ok := filter.OwnerID.Get()
		if ok {
			query = query.Where(sq.Eq{"owner_id": val})
		}

		val, ok = filter.Text.Get()
		if ok {
			// full text search on 'tsv' column
			query = query.Where(sq.Expr("tsv @@ plainto_tsquery(?)", val))
		}
	} else {
		query = query.Where(sq.Eq{"stores.id": val})
	}

	sortBy, ok := filter.SortBy.Get()
	if ok {
		sortBy, ok := storeSortingFields[sortBy]
		if !ok {
			// TODO: return error or something
			sortBy = "stores.created_at"
		}
		sortOrder, ok := filter.SortOrder.Get()
		if !ok {
			sortOrder = "asc"
		}
		query = query.OrderBy(sortBy + " " + sortOrder)
	} else {
		query = query.OrderBy("stores.created_at desc")
	}

	pageSize, ok := filter.PageSize.Get()
	if !ok {
		pageSize = 10
	}

	page, ok := filter.PageNumber.Get()
	if !ok {
		page = 1
	}

	query = query.Limit(uint64(pageSize)).Offset(uint64((page - 1) * uint64(pageSize)))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := make([]entities.Store, 0)
	for rows.Next() {
		var store entities.Store
		var owner entities.Owner
		if err := rows.Scan(&store.ID, &owner.ID, &owner.FullName, &owner.Username, &owner.CreatedAt, &store.Name, &store.Description, &store.CreatedAt); err != nil {
			return nil, err
		}
		store.Owner = &owner
		stores = append(stores, store)
	}

	return stores, nil
}

func (r storesRepository) Update(ctx context.Context, id string, changeset stores.UpdateInput) (entities.Store, error) {
	store := entities.Store{}
	query := sq.Update("stores").
		Where(sq.Eq{"id": id})

	val, ok := changeset.Name.Get()
	if ok {
		query = query.Set("name", val)
		store.Name = val
	}

	val, ok = changeset.Description.Get()
	if ok {
		query = query.Set("description", val)
		store.Description = val
	}

	if ok { // if at least one of the fields is changed
		query = query.Set("tsv", sq.Expr(
			`setweight(to_tsvector(name), 'A') || setweight(to_tsvector(description), 'B')`,
		))
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return entities.Store{}, err
	}

	_, err = r.conn.Exec(ctx, sql, args...)
	if err != nil {
		return entities.Store{}, err
	}
	store.ID = uuid.MustParse(id)

	return store, nil
}

func (r storesRepository) Delete(ctx context.Context, id string) error {
	sql, args, err := sq.Delete("stores").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(ctx, sql, args...)
	return err
}
