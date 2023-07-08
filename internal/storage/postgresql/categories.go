package postgresql

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/SanaripEsep/esep-backend/internal/domains/categories"
	"github.com/SanaripEsep/esep-backend/internal/entities"
	"github.com/SanaripEsep/esep-backend/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoriesRepository struct {
	conn *pgxpool.Pool
}

func (c categoriesRepository) Create(ctx context.Context, input entities.Category) (entities.Category, error) {
	defer telemetry.NewSpan(ctx, PackageName+"categoriesRepository.Create").End()

	var (
		storeID          *string
		parentCategoryID *string
	)
	if input.Store != nil {
		tmp := input.Store.ID.String()
		storeID = &tmp
	}
	if input.ParentCategory != nil {
		tmp := input.ParentCategory.ID.String()
		parentCategoryID = &tmp
	}
	const sql = "INSERT INTO categories (store_id, parent_category_id, name, article, icon_url)" +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id"

	res := c.conn.QueryRow(ctx, sql,
		storeID, parentCategoryID, input.Name, input.Article, input.IconURL,
	)
	return input, res.Scan(&input.ID)
}

func (c categoriesRepository) ReadBy(ctx context.Context, filters categories.ReadByInput) ([]entities.Category, error) {
	defer telemetry.NewSpan(ctx, PackageName+"categoriesRepository.ReadBy").End()

	query := sq.Select("id", "store_id", "parent_category_id", "name", "article", "icon_url", "created_at").
		From("categories").
		PlaceholderFormat(sq.Dollar)

	id, ok := filters.ID.Get()
	if ok {
		query = query.Where(sq.Eq{"id": id})
	} else {
		text, ok := filters.Text.Get()
		if ok {
			query = query.Where(sq.Like{"name": text})
		}
		storeID, ok := filters.StoreID.Get()
		if ok {
			query = query.Where(sq.Eq{"store_id": storeID})
		}
		parentCategoryID, ok := filters.ParentCategoryID.Get()
		if ok {
			query = query.Where(sq.Eq{"parent_category_id": parentCategoryID})
		}

		pageNumber, ok := filters.PageNumber.Get()
		if ok {
			query = query.Offset(pageNumber)
		}
		pageSize, ok := filters.PageSize.Get()
		if ok {
			query = query.Limit(uint64(pageSize))
		}

		sortBy, ok := filters.SortBy.Get()
		if ok {
			sortOrder, ok := filters.SortOrder.Get()
			if ok {
				query = query.OrderBy(sortBy + " " + sortOrder)
			}
		}
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.Category
	for rows.Next() {
		var (
			category         entities.Category
			storeID          *string
			parentCategoryID *string
		)
		err := rows.Scan(
			&category.ID,
			&storeID,
			&parentCategoryID,
			&category.Name,
			&category.Article,
			&category.IconURL,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// TODO: read more than just ids
		if storeID != nil {
			category.Store = &entities.Store{ID: uuid.MustParse(*storeID)}
		}
		if parentCategoryID != nil {
			category.ParentCategory = &entities.Category{ID: uuid.MustParse(*parentCategoryID)}
		}

		result = append(result, category)
	}
	return result, nil
}

func (c categoriesRepository) Update(ctx context.Context, changeset categories.UpdateInput) (entities.Category, error) {
	defer telemetry.NewSpan(ctx, PackageName+"categoriesRepository.Update").End()

	query := sq.Update("categories").
		Where(sq.Eq{"id": changeset.ID}).
		PlaceholderFormat(sq.Dollar)
	var category entities.Category

	name, ok := changeset.Name.Get()
	if ok {
		query = query.Set("name", name)
		category.Name = name
	}
	article, ok := changeset.Article.Get()
	if ok {
		query = query.Set("article", article)
		category.Article = article
	}
	parent, ok := changeset.ParentCategoryID.Get()
	if ok {
		query = query.Set("parent_category_id", parent)
		category.ParentCategory = &entities.Category{ID: uuid.MustParse(*parent)}
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return entities.Category{}, err
	}

	res := c.conn.QueryRow(ctx, sql, args...)
	return category, res.Scan(&category.ID)
}

func (c categoriesRepository) Delete(ctx context.Context, id string) error {
	defer telemetry.NewSpan(ctx, PackageName+"categoriesRepository.Delete").End()

	const sql = "DELETE FROM categories WHERE id = $1"
	_, err := c.conn.Exec(ctx, sql, id)
	return err
}
