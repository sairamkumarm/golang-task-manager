package service

import (
    "context"
    "errors"
    "golang-task-manager/internal/dtos"
    "golang-task-manager/internal/models"
    "golang-task-manager/internal/repository"
    "golang-task-manager/internal/utils"
)

type UserService struct {
    repo repository.UserRepository
    jwtHours int
}

func NewUserService(repo repository.UserRepository, jwtHours int) *UserService {
    return &UserService{repo: repo, jwtHours: jwtHours}
}

func (s *UserService) Register(ctx context.Context, req dtos.RegisterRequest) (*dtos.RegisterResponse, error) {
    if err := utils.ValidateStruct(req); err != nil {
        return nil, err
    }

    hashed, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    u := &models.User{
        Name: req.Name,
        Email: req.Email,
        Password: hashed,
    }

    if err := s.repo.Create(ctx, u); err != nil {
        return nil, err
    }

    return &dtos.RegisterResponse{
        ID: u.ID,
        Name: u.Name,
        Email: u.Email,
        CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
    }, nil
}

func (s *UserService) Login(ctx context.Context, req dtos.LoginRequest) (*dtos.LoginResponse, error) {
    if err := utils.ValidateStruct(req); err != nil {
        return nil, err
    }

    u, err := s.repo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }

    if !utils.CheckPassword(u.Password, req.Password) {
        return nil, errors.New("invalid credentials")
    }

    token, err := utils.GenerateJwt(u.ID.String(), s.jwtHours)
    if err != nil {
        return nil, err
    }

    return &dtos.LoginResponse{Token: token}, nil
}
