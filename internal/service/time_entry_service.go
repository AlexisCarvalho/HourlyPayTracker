package service

import (
	"calcula_pagamento/internal/model"
	"calcula_pagamento/internal/repository"
	"errors"
)

type TimeEntryService interface {
	RegisterTime(entry *model.TimeEntry) error
	Delete(ids []uint) error
	MarkAsPaid(id uint) error
	MarkMultipleAsPaid(ids []uint) error
	GetDurations(idUser uint, paid bool, desc bool, limit int, page int) ([]model.TimeEntryDuration, error)
	GetDurationsMonth(idUser uint, year int, month int) ([]model.TimeEntryDurationMonth, error)
	Update(entry *model.TimeEntry) error
	CountEntries(idUser uint, paid bool) (int, error)
}

type timeEntryService struct {
	repo repository.TimeEntryRepository
}

func NewTimeEntryService(repo repository.TimeEntryRepository) TimeEntryService {
	return &timeEntryService{repo: repo}
}

func (s *timeEntryService) RegisterTime(entry *model.TimeEntry) error {
	if entry.ClockOut.Before(entry.ClockIn) {
		return errors.New("ClockOut must be after ClockIn")
	}
	return s.repo.Create(entry)
}

func (s *timeEntryService) MarkAsPaid(id uint) error {
	return s.repo.MarkAsPaid(id)
}

func (s *timeEntryService) MarkMultipleAsPaid(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("no IDs provided")
	}
	return s.repo.MarkMultipleAsPaid(ids)
}

func (s *timeEntryService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("no IDs provided")
	}
	return s.repo.Delete(ids)
}

func (s *timeEntryService) GetDurations(idUser uint, paid bool, desc bool, limit int, page int) ([]model.TimeEntryDuration, error) {
	return s.repo.GetDurations(idUser, paid, desc, limit, page)
}

func (s *timeEntryService) GetDurationsMonth(idUser uint, year int, month int) ([]model.TimeEntryDurationMonth, error) {
	return s.repo.GetDurationsMonth(idUser, year, month)
}

func (s *timeEntryService) Update(entry *model.TimeEntry) error {
	if entry.ClockOut.Before(entry.ClockIn) {
		return errors.New("ClockOut must be after ClockIn")
	}
	return s.repo.Update(entry)
}

func (s *timeEntryService) CountEntries(idUser uint, paid bool) (int, error) {
	return s.repo.CountEntries(idUser, paid)
}
