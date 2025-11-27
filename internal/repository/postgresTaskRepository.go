package repository

import (
	"context"
	"golang-task-manager/internal/dtos"
	"golang-task-manager/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type taskRepository struct {
	db *sqlx.DB
}

func NewPostgresTaskRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{db: db}
}

// Create inserts a new task and populates ID, CreatedAt, UpdatedAt
func (r *taskRepository) Create(ctx context.Context, t *models.Task) error {
	query := `
		INSERT INTO tasks (user_id, title, description, status, due_date, is_public, public_slug)
		VALUES (:user_id, :title, :description, :status, :due_date, :is_public, :public_slug)
		RETURNING id, created_at, updated_at
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, t, t)
}

// GetByID returns a task for a given user
func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Task, error) {
	query := `
		SELECT * FROM tasks
		WHERE id = $1 AND user_id = $2
	`
	t:= new(models.Task)
	err := r.db.GetContext(ctx, t, query, id, userID)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// GetPublic fetches a public task by slug
func (r *taskRepository) GetPublic(ctx context.Context, slug string) (*models.Task, error) {
	query := `
		SELECT * FROM tasks
		WHERE public_slug = $1 AND is_public = true
	`
	t:= new(models.Task)
	err := r.db.GetContext(ctx, t, query, slug)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// List returns tasks with optional filters and pagination
func (r *taskRepository) List(ctx context.Context, userID uuid.UUID, filter dtos.TaskFilter) ([]models.Task, int, error) {
	baseQuery := `SELECT * FROM tasks WHERE user_id = :user_id`
	countQuery := `SELECT COUNT(*) FROM tasks WHERE user_id = :user_id`

	params := map[string]interface{}{
		"user_id": userID,
	}

	// Filters
	if filter.Status != "" {
		baseQuery += " AND status = :status"
		countQuery += " AND status = :status"
		params["status"] = filter.Status
	}
	if filter.Keyword != "" {
		baseQuery += " AND (title ILIKE :keyword OR description ILIKE :keyword)"
		countQuery += " AND (title ILIKE :keyword OR description ILIKE :keyword)"
		params["keyword"] = "%" + filter.Keyword + "%"
	}
	if !filter.DueDate.IsZero() {
		baseQuery += " AND DATE(due_date) = CAST(:due_date AS DATE)"
		countQuery += " AND DATE(due_date) = CAST(:due_date AS DATE)"
		params["due_date"] = filter.DueDate.Format(time.DateOnly)
	}

	// Pagination
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Page == 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit
	baseQuery += " ORDER BY created_at DESC LIMIT :limit OFFSET :offset"
	params["limit"] = filter.Limit
	params["offset"] = offset

	// Total count
	var total int
	nstmtCount, err := r.db.PrepareNamedContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	err = nstmtCount.GetContext(ctx, &total, params)
	if err != nil {
		return nil, 0, err
	}

	// Fetch tasks
	var tasks []models.Task
	nstmtSelect, err := r.db.PrepareNamedContext(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	err = nstmtSelect.SelectContext(ctx, &tasks, params)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}


// Update modifies a task
func (r *taskRepository) Update(ctx context.Context, t *models.Task) error {
	query := `
		UPDATE tasks
		SET title=:title, description=:description, status=:status, due_date=:due_date, is_public=:is_public, public_slug=:public_slug, updated_at=NOW()
		WHERE id=:id AND user_id=:user_id
		RETURNING updated_at
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, t, t)
}

// Delete removes a task
func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, id, userID)
	return err
}
