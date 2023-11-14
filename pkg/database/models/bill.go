package models

import "time"

type LastBillEntry struct {
	ID        uint `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	PatientID uint `json:"patient_id" gorm:"uniqueIndex; not null" example:"1"`
	Patient   User `json:"patient,omitempty" gorm:"foreignKey:PatientID"`

	C99453 int `json:"c99453" gorm:"type:smallint; not null" example:"1"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type Bill struct {
	ID        uint      `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	PatientID uint      `json:"patient_id" gorm:"not null" example:"1"`
	Patient   User      `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	ServiceID uint      `json:"service_id" gorm:"not null" example:"1"`
	Service   Service   `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	CPTCode   int64     `json:"cpt_code" gorm:"not null" example:"1"`
	EntryAt   time.Time `json:"entry_at" gorm:"not null" example:"2021-01-01T00:00:00Z"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}
