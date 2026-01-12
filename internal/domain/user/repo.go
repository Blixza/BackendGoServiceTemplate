package user_domain

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *Model) error
	GetByID(ctx context.Context, id int) (*Model, error)
}

type repo struct {
	db *pgxpool.Pool
	sb sq.StatementBuilderType
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *repo) Create(ctx context.Context, user *Model) error {
	query, args, err := r.sb.Insert("users").
		Columns("name", "email", "created_at").
		Values(user.Name, user.Email, user.CreatedAt).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID)
	return err
}

func (r *repo) GetByID(ctx context.Context, id int) (*Model, error) {
	query, args, err := r.sb.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var user Model
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
