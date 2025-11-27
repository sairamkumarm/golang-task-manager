package dtos

import "time"

// CreateTaskRequest is the incoming payload for creating a task.
type CreateTaskRequest struct {
    Title       string     `json:"title" validate:"required,min=1,max=255"`
    Description string     `json:"description" validate:"required"`
    Status      string     `json:"status" validate:"required,oneof=pending in_progress completed"`
    DueDate     *time.Time `json:"due_date" validate:"omitempty"`
    IsPublic    bool       `json:"is_public"`
}

// UpdateTaskRequest allows partial updates.
// You must use pointers to differentiate between "zero value" and "not provided".
type UpdateTaskRequest struct {
    Title       *string     `json:"title" validate:"omitempty,min=1,max=255"`
    Description *string     `json:"description"`
    Status      *string     `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
    DueDate     **time.Time `json:"due_date" validate:"omitempty"`
    IsPublic    *bool       `json:"is_public"`
}

// TaskResponse is the outward facing representation.
// Your repo populates ID, CreatedAt, UpdatedAt on insert, so they are always present here.
// PublicSlug can be nil.
type TaskResponse struct {
    ID          string     `json:"id"`
    UserID      string     `json:"user_id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      string     `json:"status"`
    DueDate     *time.Time `json:"due_date"`
    IsPublic    bool       `json:"is_public"`
    PublicSlug  *string    `json:"public_slug"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
}

// ListTasksResponse wraps paginated results.
type ListTasksResponse struct {
    Tasks []TaskResponse `json:"tasks"`
    Page  int            `json:"page"`
    Limit int            `json:"limit"`
    Total int            `json:"total"`
}


type TaskFilter struct{
	Status string
	Keyword string
	Page int
	Limit int
	DueDate time.Time
}
