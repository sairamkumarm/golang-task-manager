package service

import (
    "context"
    "errors"
    "golang-task-manager/internal/dtos"
    "golang-task-manager/internal/models"
    "golang-task-manager/internal/repository"
    "golang-task-manager/internal/utils"
    "strings"
    "time"

    "github.com/google/uuid"
)

type TaskService struct {
    repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

func slugify(title string) string {
    s := strings.ToLower(title)
    s = strings.ReplaceAll(s, " ", "-")
    s = strings.ReplaceAll(s, ".", "")
    s = strings.ReplaceAll(s, "/", "")
    return s + "-" + uuid.New().String()[:8]
}

func toResponse(t *models.Task) dtos.TaskResponse {
    return dtos.TaskResponse{
        ID:         t.ID.String(),
		UserID: 	t.UserID.String(),
        Title:      t.Title,
        Description: t.Description,
        Status:     t.Status,
        DueDate:    t.DueDate,
        IsPublic:   t.IsPublic,
        PublicSlug: t.PublicSlug,
        CreatedAt:  t.CreatedAt.Format(time.RFC3339),
        UpdatedAt:  t.UpdatedAt.Format(time.RFC3339),
    }
}

func (s *TaskService) Create(ctx context.Context, userID uuid.UUID, req dtos.CreateTaskRequest) (*dtos.TaskResponse, error) {
    if err := utils.ValidateStruct(req); err != nil {
        return nil, err
    }

    slug := (*string)(nil)
    if req.IsPublic {
        generated := slugify(req.Title)
        slug = &generated
    }

    t := &models.Task{
        UserID:     userID,
        Title:      req.Title,
        Description: req.Description,
        Status:     req.Status,
        DueDate:    req.DueDate,
        IsPublic:   req.IsPublic,
        PublicSlug: slug,
    }

    if err := s.repo.Create(ctx, t); err != nil {
        return nil, err
    }

    resp := toResponse(t)
    return &resp, nil
}

func (s *TaskService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dtos.TaskResponse, error) {
    t, err := s.repo.GetByID(ctx, id, userID)
    if err != nil {
        return nil, err
    }
    resp := toResponse(t)
    return &resp, nil
}

func (s *TaskService) List(ctx context.Context, userID uuid.UUID, filter dtos.TaskFilter) (*dtos.ListTasksResponse, error) {
    tasks, total, err := s.repo.List(ctx, userID, filter)
    if err != nil {
        return nil, err
    }

    out := make([]dtos.TaskResponse, len(tasks))
    for i, t := range tasks {
        out[i] = toResponse(&t)
    }

    return &dtos.ListTasksResponse{
        Tasks: out,
        Page:  filter.Page,
        Limit: filter.Limit,
        Total: total,
    }, nil
}

func (s *TaskService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dtos.UpdateTaskRequest) (*dtos.TaskResponse, error) {
    if err := utils.ValidateStruct(req); err != nil {
        return nil, err
    }

    t, err := s.repo.GetByID(ctx, id, userID)
    if err != nil {
        return nil, err
    }

    if req.Title != nil {
        t.Title = *req.Title
    }
    if req.Description != nil {
        t.Description = *req.Description
    }
    if req.Status != nil {
        t.Status = *req.Status
    }
    if req.DueDate != nil {
        t.DueDate = *req.DueDate
    }
    if req.IsPublic != nil {
        t.IsPublic = *req.IsPublic
        if t.IsPublic {
            if t.PublicSlug == nil {
                slug := slugify(t.Title)
                t.PublicSlug = &slug
            }
        } else {
            t.PublicSlug = nil
        }
    }

    if err := s.repo.Update(ctx, t); err != nil {
        return nil, err
    }

    resp := toResponse(t)
    return &resp, nil
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
    return s.repo.Delete(ctx, id, userID)
}

func (s *TaskService) GetPublic(ctx context.Context, slug string) (*dtos.TaskResponse, error) {
    t, err := s.repo.GetPublic(ctx, slug)
    if err != nil {
        return nil, err
    }

    if !t.IsPublic {
        return nil, errors.New("task is not public")
    }

    resp := toResponse(t)
    return &resp, nil
}
