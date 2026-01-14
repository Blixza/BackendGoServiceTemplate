package user_domain

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *Model) error
	GetByNickname(ctx context.Context, nickname string) (*Model, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Model, error)
	Update(ctx context.Context, user *Model) error
	Delete(ctx context.Context, id uuid.UUID) error
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
		Columns("nickname", "password", "discord", "email", "balance", "towns", "created_at").
		Values(user.Nickname, user.Password, user.Discord, user.Email, user.Balance, user.Towns, time.Now()).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID)
	return err
}

func (r *repo) GetByNickname(ctx context.Context, nickname string) (*Model, error) {
	query, args, err := r.sb.Select("*").
		From("users").
		Where(sq.Eq{"nickname": nickname}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	var user Model
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Nickname,
		&user.Discord,
		&user.Email,
		&user.Balance,
		&user.Towns,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &user, nil
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*Model, error) {
	query, args, err := r.sb.
		Select(
			"id",
			"nickname",
			"discord",
			"email",
			"balance",
			"towns",
			"created_at",
			"updated_at",
		).
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	var user Model
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Nickname,
		&user.Discord,
		&user.Email,
		&user.Balance,
		&user.Towns,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &user, nil
}

func (r *repo) Update(ctx context.Context, user *Model) error {
	query, args, err := r.sb.Update("users").
		Set("nickname", user.Nickname).
		Set("password", user.Password).
		Set("discord", user.Discord).
		Set("email", user.Email).
		Set("balance", user.Balance).
		Set("towns", user.Towns).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %d", user.ID)
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.sb.Delete("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found with id %d", id)
	}

	return nil
}
