package town_domain

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, town *Model) error
	GetByID(ctx context.Context, id uuid.UUID) (*Model, error)
	GetByName(ctx context.Context, name string) (*Model, error)
	GetByOwner(ctx context.Context, nickname string) (*Model, error)
	Update(ctx context.Context, town *Model) error
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

func (r *repo) Create(ctx context.Context, town *Model) error {
	query, args, err := r.sb.Insert("towns").
		Columns("name", "balance", "owner_nickname", "x_coord_overworld", "y_coord_overworld", "z_coord_overworld",
			"x_coord_nether", "y_coord_nether", "z_coord_nether", "created_at").
		Values(town.Name, town.Balance, town.OwnerNickname, town.XCoordOverworld, town.YCoordOverworld, town.ZCoordOverworld,
			town.XCoordNether, town.YCoordNether, town.ZCoordNether, time.Now()).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to create town: %w", err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&town.ID)
	return err
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*Model, error) {
	query, args, err := r.sb.Select("*").
		From("towns").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to find town: %w", err)
	}

	var town Model
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&town.ID,
		&town.Name,
		&town.Balance,
		&town.OwnerNickname,
		&town.XCoordOverworld,
		&town.YCoordOverworld,
		&town.ZCoordOverworld,
		&town.XCoordNether,
		&town.YCoordNether,
		&town.ZCoordNether,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &town, nil
}

func (r *repo) GetByName(ctx context.Context, name string) (*Model, error) {
	query, args, err := r.sb.Select("*").
		From("towns").
		Where(sq.Eq{"name": name}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to find town: %w", err)
	}

	var town Model
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&town.ID,
		&town.Name,
		&town.Balance,
		&town.OwnerNickname,
		&town.XCoordOverworld,
		&town.YCoordOverworld,
		&town.ZCoordOverworld,
		&town.XCoordNether,
		&town.YCoordNether,
		&town.ZCoordNether,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &town, nil
}

func (r *repo) GetByOwner(ctx context.Context, nickname string) (*Model, error) {
	query, args, err := r.sb.Select("*").
		From("towns").
		Where(sq.Eq{"owner_nickname": nickname}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to find town: %w", err)
	}

	var town Model
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&town.ID,
		&town.Name,
		&town.Balance,
		&town.OwnerNickname,
		&town.XCoordOverworld,
		&town.YCoordOverworld,
		&town.ZCoordOverworld,
		&town.XCoordNether,
		&town.YCoordNether,
		&town.ZCoordNether,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &town, nil
}

func (r *repo) Update(ctx context.Context, town *Model) error {
	query, args, err := r.sb.Update("towns").
		Set("name", town.Name).
		Set("balance", town.Balance).
		Set("owner_nickname", town.OwnerNickname).
		Set("x_coord_overworld", town.XCoordOverworld).
		Set("y_coord_overworld", town.YCoordOverworld).
		Set("z_coord_overworld", town.ZCoordOverworld).
		Set("x_coord_nether", town.XCoordNether).
		Set("y_coord_nether", town.YCoordNether).
		Set("z_coord_nether", town.ZCoordNether).
		Where(sq.Eq{"id": town.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no town found with id %d", town.ID)
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.sb.Delete("towns").
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
		return fmt.Errorf("no town found with id %d", id)
	}

	return nil
}
