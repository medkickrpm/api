package models

import (
	"MedKick-backend/pkg/database"
	"time"

	"gorm.io/datatypes"
)

type TelemetryAlert struct {
	ID             uint                 `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	OrganizationID uint                 `json:"organization_id" gorm:"not null" example:"1"`
	Organization   *Organization        `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	DeviceType     DeviceType           `json:"device_type" example:"BloodPressure"`
	DeviceID       uint                 `json:"device_id" example:"1"`
	Device         *Device              `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	TelemetryID    uint                 `json:"telemetry_id" example:"1"`
	Telemetry      *DeviceTelemetryData `json:"telemetry,omitempty" gorm:"foreignKey:TelemetryID"`
	PatientID      uint                 `json:"patient_id" example:"1"`
	Patient        *User                `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	AlertType      AlertType            `json:"alert_type" example:"WarningHigh"`
	Data           datatypes.JSONMap    `json:"data" example:"{\"value\": 120}"`
	IsActive       bool                 `json:"is_active" example:"true"`
	ResolvedByID   *uint                `json:"resolved_by_id,omitempty" example:"1"`
	ResolvedBy     *User                `json:"resolved_by,omitempty" gorm:"foreignKey:ResolvedByID"`
	MeasuredAt     time.Time            `json:"measured_at" example:"2021-01-01T00:00:00Z"`
	IsAutoResolved bool                 `json:"is_auto_resolved,omitempty" example:"true"`
	ResolvedAt     *time.Time           `json:"resolved_at,omitempty" example:"2021-01-01T00:00:00Z"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type AlertType string

const (
	AlertCritical AlertType = "Critical"
	AlertWarning  AlertType = "Warning"
	AlertOk       AlertType = "Ok"
)

func (t *TelemetryAlert) GetTelemetryAlert() error {
	db := database.DB.Model(&TelemetryAlert{})

	if t.ID != 0 {
		db = db.Where("id = ?", t.ID)
	}

	if t.DeviceID != 0 {
		db = db.Where("device_id = ?", t.DeviceID)
	}

	db = db.Where("is_active = ?", t.IsActive)
	db = db.Where("is_auto_resolved = ?", t.IsAutoResolved)

	if t.PatientID != 0 {
		db = db.Where("patient_id = ?", t.PatientID)
	}

	if err := db.First(&t).Error; err != nil {
		return err
	}

	return nil
}

func (t *TelemetryAlert) InsertTelemetryAlert() error {
	db := database.DB.Model(&TelemetryAlert{})
	if err := db.Create(&t).Error; err != nil {
		return err
	}
	return nil
}

func (t *TelemetryAlert) ResolveTelemetryAlert() error {
	db := database.DB.Model(&TelemetryAlert{})
	db = db.Where("id = ?", t.ID)
	db = db.Where("is_active = ?", true)
	db = db.Where("is_auto_resolved = ?", false)

	if err := db.UpdateColumns(map[string]interface{}{
		"is_active":      false,
		"resolved_by_id": t.ResolvedByID,
		"resolved_at":    time.Now().UTC(),
	}).Error; err != nil {
		return err
	}

	return nil
}

func ListTelemetryAlerts(org uint, isActive bool, pagination PageReq, sort SortReq) ([]TelemetryAlert, error) {
	var telemetryAlerts []TelemetryAlert
	db := database.DB.Model(&TelemetryAlert{})
	db = db.Where("organization_id = ?", org)
	db = db.Where("is_active = ?", isActive)
	db = db.Where("is_auto_resolved = ?", false)

	db.Scopes(pagination.Paginate())
	db.Scopes(sort.Sort())

	if err := db.Preload("Patient").Preload("ResolvedBy").Find(&telemetryAlerts).Error; err != nil {
		return nil, err
	}

	return telemetryAlerts, nil
}

func CountTelemetryAlerts(org uint, isActive bool) (int64, error) {
	var count int64
	db := database.DB.Model(&TelemetryAlert{})
	db = db.Where("organization_id = ?", org)
	db = db.Where("is_active = ?", isActive)
	db = db.Where("is_auto_resolved = ?", false)

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
