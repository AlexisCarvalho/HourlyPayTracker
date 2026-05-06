package model

import (
	"time"
)

type TimeEntry struct {
	ID                   uint                       `gorm:"primaryKey" json:"id"`
	IdUser               uint                       `json:"id_user"`
	CompanyPaymentInfoID uint                       `json:"company_payment_info_id"`
	CompanyPaymentInfo   *CompanyPaymentInformation `json:"company_payment_info" gorm:"foreignKey:CompanyPaymentInfoID"`
	ClockIn              time.Time                  `json:"clock_in"`
	ClockOut             time.Time                  `json:"clock_out"`
	Paid                 bool                       `json:"paid"`
}

type TimeEntryDuration struct {
	ID                   uint      `json:"id"`
	ClockIn              time.Time `json:"clock_in"`
	ClockOut             time.Time `json:"clock_out"`
	TotalDurationMinutes uint32    `json:"total_duration_minutes"`
	Hours                uint32    `json:"hours"`
	Minutes              uint32    `json:"minutes"`
	HourlyRate           float64   `json:"hourly_rate"`
}

type TimeEntryDurationMonth struct {
	ID                   uint      `json:"id"`
	ClockIn              time.Time `json:"clock_in"`
	ClockOut             time.Time `json:"clock_out"`
	TotalDurationMinutes uint32    `json:"total_duration_minutes"`
	Hours                uint32    `json:"hours"`
	Minutes              uint32    `json:"minutes"`
	Paid                 bool      `json:"paid"`
	HourlyRate           float64   `json:"hourly_rate"`
}
