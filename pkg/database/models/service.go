package models

import "time"

type Service struct {
	ID          uint   `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	Code        string `json:"service_code" gorm:"type:varchar(10); not null; uniqueIndex" example:"RPM"`
	Name        string `json:"service_name" gorm:"not null" example:"Remote Patient Monitoring"`
	IsEnabled   bool   `json:"is_enabled" gorm:"not null" example:"true"`
	Description string `json:"description" example:"Remote Patient Monitoring"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type PatientService struct {
	ID        uint    `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	PatientID uint    `json:"patient_id" gorm:"index:,unique,composite:patient_service; not null" example:"1"`
	Patient   User    `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	ServiceID uint    `json:"service_id" gorm:"index:,unique,composite:patient_service; not null" example:"1"`
	Service   Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`

	StartedAt time.Time  `json:"started_at" gorm:"not null" example:"2021-01-01T00:00:00Z"`
	EndedAt   *time.Time `json:"ended_at" gorm:"default:null" example:"2021-01-01T00:00:00Z"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}
