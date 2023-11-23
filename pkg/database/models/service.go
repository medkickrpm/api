package models

import (
	"MedKick-backend/pkg/database"
	"time"

	"gorm.io/gorm/clause"
)

type Service struct {
	ID          uint   `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	Code        string `json:"service_code" gorm:"type:varchar(10); not null; uniqueIndex" example:"RPM"`
	Name        string `json:"service_name" gorm:"not null" example:"Remote Patient Monitoring"`
	IsEnabled   bool   `json:"is_enabled" gorm:"not null" example:"true"`
	Description string `json:"description" example:"Remote Patient Monitoring"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func ListServices() ([]Service, error) {
	var services []Service

	db := database.DB.Model(&Service{})
	if err := db.Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}

type PatientService struct {
	ID        uint    `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	PatientID uint    `json:"patient_id" gorm:"not null" example:"1"`
	Patient   User    `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	ServiceID uint    `json:"service_id" gorm:"not null" example:"1"`
	Service   Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	Status    bool    `json:"status" gorm:"not null; DEFAULT:1" example:"true"`

	StartedAt time.Time  `json:"started_at" gorm:"not null" example:"2021-01-01T00:00:00Z"`
	EndedAt   *time.Time `json:"ended_at" gorm:"default:null" example:"2021-01-01T00:00:00Z"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func ListPatientServices(patientID uint, status string, pagination PageReq, sort SortReq) ([]PatientService, error) {
	var patientServices []PatientService
	db := database.DB.Model(&PatientService{})
	db = db.Preload("Service")
	db = db.Where("patient_id = ?", patientID)

	if status == "active" {
		db = db.Where("status = ?", true)
	} else if status == "inactive" {
		db = db.Where("status = ?", false)
	}

	db = db.Scopes(pagination.Paginate())
	db = db.Scopes(sort.Sort())

	if err := db.Find(&patientServices).Error; err != nil {
		return nil, err
	}

	return patientServices, nil
}

func CountPatientServices(patientID uint, status string) (int, error) {
	var count int64

	db := database.DB.Model(&PatientService{})
	db = db.Preload("Service")
	db = db.Where("patient_id = ?", patientID)

	if status == "active" {
		db = db.Where("status = ?", true)
	} else if status == "inactive" {
		db = db.Where("status = ?", false)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func UpsertPatientServices(services []PatientService) error {
	// If patient_id & service_id already exists, update ended_at and started_at
	// otherwise, insert new row

	db := database.DB.Model(&PatientService{})
	// Conflict with  PatientID, ServiceID then update all
	db = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"status",
			"ended_at",
			"updated_at",
		}),
	})

	if err := db.Create(&services).Error; err != nil {
		return err
	}

	return nil
}
