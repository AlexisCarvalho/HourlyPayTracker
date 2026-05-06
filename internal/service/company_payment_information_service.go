package service

import (
	"calcula_pagamento/internal/model"
	"calcula_pagamento/internal/repository"
)

type CompanyPaymentInformationService interface {
	Create(company *model.CompanyPaymentInformation) error
	GetByID(id uint) (*model.CompanyPaymentInformation, error)
	GetAll() ([]model.CompanyPaymentInformation, error)
}

type companyPaymentInfoService struct {
	repo repository.CompanyPaymentInformationRepository
}

func NewCompanyPaymentInformationService(repo repository.CompanyPaymentInformationRepository) CompanyPaymentInformationService {
	return &companyPaymentInfoService{repo}
}

func (s *companyPaymentInfoService) Create(company *model.CompanyPaymentInformation) error {
	return s.repo.Create(company)
}

func (s *companyPaymentInfoService) GetByID(id uint) (*model.CompanyPaymentInformation, error) {
	return s.repo.FindByID(id)
}

func (s *companyPaymentInfoService) GetAll() ([]model.CompanyPaymentInformation, error) {
	return s.repo.FindAll()
}
