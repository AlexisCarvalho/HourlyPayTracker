package model

type User struct {
	ID                 uint                       `gorm:"primaryKey" json:"id"`
	Name               string                     `json:"name"`
	Code               string                     `json:"code" gorm:"unique"`
	Password           string                     `json:"password"`
	PreferredCompanyID *uint                      `json:"preferred_company_id"`
	PreferredCompany   *CompanyPaymentInformation `json:"preferred_company" gorm:"foreignKey:PreferredCompanyID"`
}
