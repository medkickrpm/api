package models

import (
	"MedKick-backend/pkg/database"
	"gorm.io/gorm/clause"
	"time"
)

type AlertThreshold struct {
	ID              uint            `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	OrganizationID  uint            `json:"organization_id" gorm:"index:,unique,composite:measurement; not null" example:"1"`
	Organization    Organization    `json:"organization" gorm:"foreignKey:OrganizationID"`
	DeviceType      DeviceType      `json:"device_type" gorm:"index:,unique,composite:measurement; not null" example:"BloodPressure"`
	MeasurementType MeasurementType `json:"measurement_type" gorm:"index:,unique,composite:measurement; not null" example:"Systolic"`
	CriticalLow     *uint           `json:"critical_low" gorm:"default:null" example:"60"`
	WarningLow      *uint           `json:"warning_low" gorm:"default:null" example:"80"`
	WarningHigh     *uint           `json:"warning_high" gorm:"default:null" example:"120"`
	CriticalHigh    *uint           `json:"critical_high" gorm:"default:null" example:"140"`
	Note            string          `json:"note" gorm:"default:null" example:"This is a note"`
	CreatedAt       time.Time       `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt       time.Time       `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type DeviceType string

const (
	BloodPressure DeviceType = "BloodPressure"
	BloodGlucose  DeviceType = "BloodGlucose"
	WeightScale   DeviceType = "WeightScale"
)

type MeasurementType string

const (
	Systolic  MeasurementType = "Systolic"
	Diastolic MeasurementType = "Diastolic"
	Pulse     MeasurementType = "Pulse"
	Weight    MeasurementType = "Weight"
)

func ListAlertThresholds(orgId uint) ([]AlertThreshold, error) {
	var alertThresholds []AlertThreshold
	if err := database.DB.Where("organization_id = ?", orgId).Find(&alertThresholds).Error; err != nil {
		return nil, err
	}

	return alertThresholds, nil
}

func UpsertAlertThresholds(alertThresholds []AlertThreshold) error {
	db := database.DB.Model(&AlertThreshold{})
	// Conflict with  OrganizationID, DeviceType, MeasurementType then update all
	db = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "organization_id"}, {Name: "device_type"}, {Name: "measurement_type"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"critical_low",
			"warning_low",
			"warning_high",
			"critical_high",
		}),
	})
	if err := db.Create(&alertThresholds).Error; err != nil {
		return err
	}
	return nil
}
