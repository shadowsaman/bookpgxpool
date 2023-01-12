package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"app/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"

	"github.com/google/uuid"
)

type categoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *categoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) Insert(ctx context.Context, category *models.CreateCategory) (string, error) {

	var (
		id = uuid.New().String()
	)

	query := `
		INSERT INTO category (
			id,
			name,
			updated_at
		) VALUES ($1, $2, now())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		category.Name,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *categoryRepo) GetByID(ctx context.Context, req *models.CategoryPrimeryKey) (*models.Category, error) {

	query := `
		SELECT
			c.id,
			c.name,
			c.created_at,
			c.updated_at,
			(
				SELECT
					ARRAY_AGG(books_id)
				FROM book_category AS bc 
				WHERE bc.category_id = $1
			) AS book_ids
		FROM
			category AS c
		WHERE c.id = $1
	`

	var (
		id        sql.NullString
		name      sql.NullString
		createdAt sql.NullString
		updatedAt sql.NullString
		booksIds  []string
	)

	err := r.db.QueryRow(ctx, query, req.Id).
		Scan(
			&id,
			&name,
			&createdAt,
			&updatedAt,
			(*pq.StringArray)(&booksIds),
		)

	if err != nil {
		return nil, err
	}

	category := &models.Category{
		Id:        id.String,
		Name:      name.String,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
	}

	if len(booksIds) > 0 {
		bookQuery := `
			SELECT
				id,
				name,
				price,
				description,
				created_at,
				updated_at
			FROM
				books
			WHERE id IN (`

		for _, bookId := range booksIds {
			bookQuery += fmt.Sprintf("'%s',", bookId)
		}
		bookQuery = bookQuery[:len(bookQuery)-1]
		bookQuery += ")"

		rows, err := r.db.Query(ctx, bookQuery)
		if err != nil {
			return nil, err
		}

		for rows.Next() {

			var (
				id          sql.NullString
				name        sql.NullString
				price       sql.NullFloat64
				description sql.NullString
				createdAt   sql.NullString
				updatedAt   sql.NullString
			)

			err = rows.Scan(
				&id,
				&name,
				&price,
				&description,
				&createdAt,
				&updatedAt,
			)

			category.Books = append(category.Books, models.Book1{
				Id:          id.String,
				Name:        name.String,
				Price:       price.Float64,
				Description: description.String,
				CreatedAt:   createdAt.String,
				UpdatedAt:   updatedAt.String,
			})
		}
	}

	return category, nil
}

func (r *categoryRepo) GetList(ctx context.Context, req *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error) {

	var (
		resp   models.GetListCategoryResponse
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query := `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			created_at,
			updated_at
		FROM category
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return &models.GetListCategoryResponse{}, err
	}

	var (
		id        sql.NullString
		name      sql.NullString
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	for rows.Next() {

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&createdAt,
			&updatedAt,
		)

		category := models.Category1{
			Id:        id.String,
			Name:      name.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		}

		if err != nil {
			return &models.GetListCategoryResponse{}, err
		}

		resp.Categories = append(resp.Categories, &category)
	}

	return &resp, nil
}

func (r *categoryRepo) Update(ctx context.Context, category *models.UpdateCategory) error {

	query := `
	UPDATE
		category
	SET
		name = $2,
		updated_at = now()
	WHERE id = $1
`

	_, err := r.db.Exec(ctx, query,
		category.Id,
		category.Name,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepo) Delete(ctx context.Context, req *models.CategoryPrimeryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM book_category WHERE id = $1", req.Id)

	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, "delete from category where id = $1")
	if err != nil {
		return err
	}

	return nil
}
