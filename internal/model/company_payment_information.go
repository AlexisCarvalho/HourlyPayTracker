package model

type CompanyPaymentInformation struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	CompanyName string  `json:"company_name"`
	HourlyRate  float64 `json:"hourly_rate"`
}
