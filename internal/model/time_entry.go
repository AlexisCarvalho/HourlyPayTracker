package model

import (
	"time"

	"gorm.io/gorm"
)

type TimeEntry struct {
	gorm.Model
	IdUser   uint      `json:"id_user"`
	ClockIn  time.Time `json:"clock_in"`
	ClockOut time.Time `json:"clock_out"`
	Paid     bool      `json:"paid"`
}

type TimeEntryDuration struct {
	ID                   uint      `json:"id"`
	ClockIn              time.Time `json:"clock_in"`
	ClockOut             time.Time `json:"clock_out"`
	TotalDurationMinutes uint32    `json:"total_duration_minutes"`
	Hours                uint32    `json:"hours"`
	Minutes              uint32    `json:"minutes"`
}

type TimeEntryDurationMonth struct {
	ID                   uint      `json:"id"`
	ClockIn              time.Time `json:"clock_in"`
	ClockOut             time.Time `json:"clock_out"`
	TotalDurationMinutes uint32    `json:"total_duration_minutes"`
	Hours                uint32    `json:"hours"`
	Minutes              uint32    `json:"minutes"`
	Paid                 bool      `json:"paid"`
}
