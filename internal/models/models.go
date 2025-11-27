package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    Password  string    `db:"password"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type Task struct {
    ID         uuid.UUID  `db:"id"`
    UserID     uuid.UUID  `db:"user_id"`
    Title      string     `db:"title"`
    Description string    `db:"description"`
    Status     string     `db:"status"`
    DueDate    *time.Time `db:"due_date"`
    IsPublic   bool       `db:"is_public"`
    PublicSlug *string    `db:"public_slug"`
    CreatedAt  time.Time  `db:"created_at"`
    UpdatedAt  time.Time  `db:"updated_at"`
}
