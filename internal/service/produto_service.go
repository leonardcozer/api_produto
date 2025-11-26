package service

import (
	"context"

	"api-go/internal/model"
	"api-go/internal/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Create(user model.User) error {
	_, err := s.Repo.Create(context.Background(), user)
	return err
}

func (s *UserService) GetAll() ([]model.User, error) {
	return s.Repo.FindAll(context.Background())
}
