package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang-task-manager/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
    ErrUserNotFound = errors.New("user not found")
)

type PostgresUserRepository struct {
	DB *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
    return &PostgresUserRepository{DB: db}
}

func (p *PostgresUserRepository) Create(ctx context.Context, u *models.User) error{
	statement := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	return p.DB.QueryRowxContext(ctx, statement, u.Name, u.Email, u.Password).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (p *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error){
	u := new(models.User)
	statement := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	err := p.DB.GetContext(ctx, u, statement, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (p *PostgresUserRepository) GetById(ctx context.Context, id uuid.UUID) (*models.User, error) {
    u := new(models.User)
    query := `
        SELECT id, name, email, password, created_at, updated_at
        FROM users
        WHERE id = $1
        LIMIT 1
    `
    err := p.DB.GetContext(ctx, u, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return u, nil
}