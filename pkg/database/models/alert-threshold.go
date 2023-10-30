package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type AlertThreshold struct {
	ID              uint            `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	OrganizationID  uint            `json:"organization_id" gorm:"index:,unique,composite:measurement; not null" example:"1"`
	Organization    Organization    `json:"organization" gorm:"foreignKey:OrganizationID"`
	DeviceType      DeviceType      `json:"device_type" gorm:"index:,unique,composite:measurement; not null" example:"BloodPressure"`
	MeasurementType MeasurementType `json:"measurement_type" gorm:"index:,unique,composite:measurement; not null" example:"Systolic"`
	CriticalHigh    uint            `json:"critical_high" example:"140"`
	WarningHigh     *uint           `json:"warning_high" gorm:"default:null" example:"120"`
	WarningLow      *uint           `json:"warning_low" gorm:"default:null" example:"80"`
	CriticalLow     uint            `json:"critical_low" example:"60"`
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

func CreateAlertThresholds(alertThresholds []AlertThreshold) error {
	if err := database.DB.Create(&alertThresholds).Error; err != nil {
		return err
	}
	return nil
}
