package service

import (
	"calcula_pagamento/internal/model"
	"calcula_pagamento/internal/repository"
)

type UserService interface {
	RegisterUser(entry *model.User) error
	Authenticate(code, password string) *model.User
	GetByID(id uint) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) RegisterUser(entry *model.User) error {
	return s.repo.Create(entry)
}

func (s *userService) Authenticate(code, password string) *model.User {
	user, err := s.repo.FindByCode(code)
	if err != nil || user.Password != password {
		return nil
	}
	return user
}

func (s *userService) GetByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}
