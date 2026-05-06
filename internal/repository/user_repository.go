package repository

import (
	"calcula_pagamento/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(entry *model.User) error
	FindByCode(code string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(entry *model.User) error {
	return r.db.Create(entry).Error
}

func (r *userRepo) FindByCode(code string) (*model.User, error) {
	var user model.User
	result := r.db.Preload("PreferredCompany").Where("code = ?", code).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	result := r.db.Preload("PreferredCompany").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
