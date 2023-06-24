package postgresql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rasulov-emirlan/esep-backend/internal/entities"
)

type ownersRepository struct {
	conn *pgxpool.Pool
}

func (r ownersRepository) Create(ctx context.Context, owner entities.Owner) (entities.Owner, error) {
	sql, args, err := sq.Insert("owners").
		Columns("full_name", "username", "password_hash", "phone_number", "created_at").
		Values(owner.FullName, owner.Username, owner.Password, owner.PhoneNumber, owner.CreatedAt).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return entities.Owner{}, fmt.Errorf("could not construct sql: %w", err)
	}

	row := r.conn.QueryRow(ctx, sql, args...)
	if err := row.Scan(&owner.ID); err != nil {
		return entities.Owner{}, fmt.Errorf("could not scan row: %w", err)
	}

	return owner, nil
}

func (r ownersRepository) Read(ctx context.Context, id string) (entities.Owner, error) {
	var owner entities.Owner

	sql, args, err := sq.Select("id", "full_name", "username", "password_hash", "phone_number", "created_at").
		From("owners").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return entities.Owner{}, fmt.Errorf("could not construct sql: %w", err)
	}

	row := r.conn.QueryRow(ctx, sql, args...)
	if err := row.Scan(&owner.ID, &owner.FullName, &owner.Username, &owner.Password, &owner.PhoneNumber, &owner.CreatedAt); err != nil {
		return entities.Owner{}, fmt.Errorf("could not scan row: %w", err)
	}

	return owner, nil
}

func (r ownersRepository) ReadByUsername(ctx context.Context, username string) (entities.Owner, error) {
	var owner entities.Owner

	sql, args, err := sq.Select("id", "full_name", "username", "password_hash", "phone_number", "created_at").
		From("owners").
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return entities.Owner{}, fmt.Errorf("could not construct sql: %w", err)
	}

	row := r.conn.QueryRow(ctx, sql, args...)
	if err := row.Scan(&owner.ID, &owner.FullName, &owner.Username, &owner.Password, &owner.PhoneNumber, &owner.CreatedAt); err != nil {
		return entities.Owner{}, fmt.Errorf("could not scan row: %w", err)
	}

	return owner, nil
}

func (r ownersRepository) ReadAll(ctx context.Context) ([]entities.Owner, error) {
	var owners []entities.Owner

	sql, args, err := sq.Select("id", "full_name", "username", "password_hash", "phone_number", "created_at").
		From("owners").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("could not construct sql: %w", err)
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("could not scan row: %w", err)
	}

	for rows.Next() {
		var owner entities.Owner
		if err := rows.Scan(&owner.ID, &owner.FullName, &owner.Username, &owner.Password, &owner.PhoneNumber, &owner.CreatedAt); err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		owners = append(owners, owner)
	}

	return owners, nil
}

func (r ownersRepository) Update(ctx context.Context, owner entities.Owner) (entities.Owner, error) {
	sql, args, err := sq.Update("owners").
		Set("full_name", owner.FullName).
		Set("username", owner.Username).
		Set("password_hash", owner.Password).
		Set("phone_number", owner.PhoneNumber).
		Where(sq.Eq{"id": owner.ID}).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return entities.Owner{}, fmt.Errorf("could not construct sql: %w", err)
	}

	row := r.conn.QueryRow(ctx, sql, args...)
	if err := row.Scan(&owner.ID); err != nil {
		return entities.Owner{}, fmt.Errorf("could not scan row: %w", err)
	}

	return owner, nil
}

func (r ownersRepository) Delete(ctx context.Context, id string) error {
	sql, args, err := sq.Delete("owners").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("could not construct sql: %w", err)
	}

	if _, err := r.conn.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("could not scan row: %w", err)
	}

	return nil
}
