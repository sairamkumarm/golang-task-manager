package repository

import (
	"context"
	"golang-task-manager/internal/dtos"
	"golang-task-manager/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type TaskRepository interface {
    Create(ctx context.Context, t *models.Task) error
    GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Task, error)
    List(ctx context.Context, userID uuid.UUID, filter dtos.TaskFilter) ([]models.Task, int, error)
    Update(ctx context.Context, t *models.Task) error
    Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
    GetPublic(ctx context.Context, slug string) (*models.Task, error)
}
