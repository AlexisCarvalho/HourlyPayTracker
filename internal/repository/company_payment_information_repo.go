package repository

import (
	"calcula_pagamento/internal/model"

	"gorm.io/gorm"
)

type CompanyPaymentInformationRepository interface {
	Create(company *model.CompanyPaymentInformation) error
	FindByID(id uint) (*model.CompanyPaymentInformation, error)
	FindAll() ([]model.CompanyPaymentInformation, error)
}

type companyPaymentInfoRepo struct {
	db *gorm.DB
}

func NewCompanyPaymentInformationRepository(db *gorm.DB) CompanyPaymentInformationRepository {
	return &companyPaymentInfoRepo{db}
}

func (r *companyPaymentInfoRepo) Create(company *model.CompanyPaymentInformation) error {
	return r.db.Create(company).Error
}

func (r *companyPaymentInfoRepo) FindByID(id uint) (*model.CompanyPaymentInformation, error) {
	var company model.CompanyPaymentInformation
	result := r.db.First(&company, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &company, nil
}

func (r *companyPaymentInfoRepo) FindAll() ([]model.CompanyPaymentInformation, error) {
	var companies []model.CompanyPaymentInformation
	result := r.db.Find(&companies)
	if result.Error != nil {
		return nil, result.Error
	}
	return companies, nil
}
