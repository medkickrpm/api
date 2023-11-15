package models

import (
	"MedKick-backend/pkg/database"
	"errors"
	"time"

	"gorm.io/gorm"
)

type LastBillEntry struct {
	ID        uint `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	PatientID uint `json:"patient_id" gorm:"uniqueIndex; not null" example:"1"`
	Patient   User `json:"patient,omitempty" gorm:"foreignKey:PatientID"`

	C99453 int `json:"c99453" gorm:"type:smallint; not null; default: 0" example:"1"`
	C99454 int `json:"c99454" gorm:"type:smallint; not null; default: 0" example:"1"`
	C99457 int `json:"c99457" gorm:"type:smallint; not null; default: 0" example:"1"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func UpdateLastBillEntry(patientID uint, update map[string]interface{}) error {
	update["patient_id"] = patientID
	db := database.DB.Model(&LastBillEntry{}).Where("patient_id = ?", patientID)
	resp := db.Updates(update)
	if (resp.Error != nil && errors.Is(resp.Error, gorm.ErrRecordNotFound)) || resp.RowsAffected == 0 {
		update["created_at"] = time.Now().UTC()
		update["updated_at"] = time.Now().UTC()
		if err := database.DB.Model(&LastBillEntry{}).Create(update).Error; err != nil {
			return err
		}
	}

	return nil
}

type Bill struct {
	ID          uint      `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	PatientID   uint      `json:"patient_id" gorm:"not null" example:"1"`
	Patient     User      `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	ServiceCode string    `json:"service_code" gorm:"not null" example:"RPM"`
	CPTCode     int64     `json:"cpt_code" gorm:"not null" example:"1"`
	EntryAt     time.Time `json:"entry_at" gorm:"not null" example:"2021-01-01T00:00:00Z"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (b *Bill) CreateBill() error {
	if err := database.DB.Create(&b).Error; err != nil {
		return err
	}
	return nil
}
