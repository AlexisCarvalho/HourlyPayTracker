package repository

import (
	"calcula_pagamento/internal/model"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TimeEntryRepository interface {
	Create(entry *model.TimeEntry) error
	Delete(ids []uint) error
	MarkAsPaid(id uint) error
	MarkMultipleAsPaid(ids []uint) error
	GetByID(id uint) (*model.TimeEntry, error)
	GetDurations(idUser uint, paid bool, desc bool, limit int, page int) ([]model.TimeEntryDuration, error)
	GetDurationsMonth(idUser uint, year int, month int) ([]model.TimeEntryDurationMonth, error)
	Update(entry *model.TimeEntry) error
	CountEntries(idUser uint, paid bool) (int, error)
}

type timeEntryRepo struct {
	db *gorm.DB
}

func NewTimeEntryRepository(db *gorm.DB) TimeEntryRepository {
	return &timeEntryRepo{db: db}
}

func (r *timeEntryRepo) Create(entry *model.TimeEntry) error {
	return r.db.Create(entry).Error
}

func (r *timeEntryRepo) MarkAsPaid(id uint) error {
	return r.db.Model(&model.TimeEntry{}).Where("id = ?", id).Update("paid", true).Error
}

func (r *timeEntryRepo) MarkMultipleAsPaid(ids []uint) error {
	return r.db.Model(&model.TimeEntry{}).
		Where("id IN ?", ids).
		Update("paid", true).Error
}

func (r *timeEntryRepo) Delete(ids []uint) error {
	return r.db.Model(&model.TimeEntry{}).
		Delete("id IN ?", ids).Error
}

func (r *timeEntryRepo) GetByID(id uint) (*model.TimeEntry, error) {
	var entry model.TimeEntry
	err := r.db.
		Preload("CompanyPaymentInfo").
		First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// Get the durations, if limit = -1 offset will be as well
func (r *timeEntryRepo) GetDurations(idUser uint, paid bool, desc bool, limit int, page int) ([]model.TimeEntryDuration, error) {
	var offset int

	if limit < 0 {
		offset = -1
	} else {
		offset = (page - 1) * limit
	}

	var entries []model.TimeEntry
	err := r.db.
		Preload("CompanyPaymentInfo").
		Where("id_user = ? AND paid = ?", idUser, paid).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: "clock_in"},
			Desc:   desc,
		}).
		Limit(limit).
		Offset(offset).
		Find(&entries).Error
	if err != nil {
		return nil, err
	}

	results := make([]model.TimeEntryDuration, 0, len(entries))

	// Compute durations
	for _, entry := range entries {
		totalMinutes := uint32(entry.ClockOut.Sub(entry.ClockIn).Minutes())
		hours := totalMinutes / 60
		minutes := totalMinutes % 60

		hourlyRate := 0.0
		if entry.CompanyPaymentInfo != nil {
			hourlyRate = entry.CompanyPaymentInfo.HourlyRate
		}

		results = append(results, model.TimeEntryDuration{
			ID:                   entry.ID,
			ClockIn:              entry.ClockIn,
			ClockOut:             entry.ClockOut,
			TotalDurationMinutes: totalMinutes,
			Hours:                hours,
			Minutes:              minutes,
			HourlyRate:           hourlyRate,
		})
	}

	return results, nil
}

func (r *timeEntryRepo) GetDurationsMonth(idUser uint, year int, month int) ([]model.TimeEntryDurationMonth, error) {
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("mês inválido: %d", month)
	}

	loc := time.UTC // ou injete/configure isso

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0) // início do próximo mês

	var entries []model.TimeEntry
	if err := r.db.
		Preload("CompanyPaymentInfo").
		Where("id_user = ? AND clock_in >= ? AND clock_in < ?", idUser, start, end).
		Order("clock_in ASC").
		Find(&entries).Error; err != nil {

		return nil, err
	}

	results := make([]model.TimeEntryDurationMonth, 0, len(entries))

	// Compute durations
	for _, entry := range entries {
		totalMinutes := uint32(entry.ClockOut.Sub(entry.ClockIn).Minutes())
		hours := totalMinutes / 60
		minutes := totalMinutes % 60

		hourlyRate := 0.0
		if entry.CompanyPaymentInfo != nil {
			hourlyRate = entry.CompanyPaymentInfo.HourlyRate
		}

		results = append(results, model.TimeEntryDurationMonth{
			ID:                   entry.ID,
			ClockIn:              entry.ClockIn,
			ClockOut:             entry.ClockOut,
			TotalDurationMinutes: totalMinutes,
			Hours:                hours,
			Minutes:              minutes,
			Paid:                 entry.Paid,
			HourlyRate:           hourlyRate,
		})
	}

	return results, nil
}

func (r *timeEntryRepo) Update(entry *model.TimeEntry) error {
	// Ensure the record exists and belongs to the same user
	var existing model.TimeEntry
	if err := r.db.First(&existing, entry.ID).Error; err != nil {
		return err
	}
	if existing.IdUser != entry.IdUser {
		return gorm.ErrRecordNotFound
	}

	// Update all fields except ID and CreatedAt
	return r.db.Model(&existing).Updates(map[string]interface{}{
		"clock_in":                entry.ClockIn,
		"clock_out":               entry.ClockOut,
		"paid":                    entry.Paid,
		"id_user":                 entry.IdUser,
		"company_payment_info_id": entry.CompanyPaymentInfoID,
	}).Error
}

func (r *timeEntryRepo) CountEntries(idUser uint, paid bool) (int, error) {
	var count int64
	err := r.db.Model(&model.TimeEntry{}).
		Where("id_user = ? AND paid = ?", idUser, paid).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
