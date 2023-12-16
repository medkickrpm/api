package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type Device struct {
	ID                  uint                  `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	Name                string                `json:"name" gorm:"not null" example:"Sphygmomanometer/Weight Scale/Blood Glucose Meter"`
	ModelNumber         string                `json:"model_number" gorm:"not null" example:"123456"`
	IMEI                string                `json:"imei" gorm:"not null" example:"123456789"`
	SerialNumber        string                `json:"serial_number" gorm:"not null" example:"123456789"`
	BatteryLevel        uint                  `json:"battery_level" gorm:"not null" example:"100"`
	SignalStrength      string                `json:"signal_strength" gorm:"not null" example:"100"`
	FirmwareVersion     string                `json:"firmware_version" gorm:"not null" example:"1.0.0"`
	UserID              uint                  `json:"user_id" example:"1"`
	User                User                  `json:"-" gorm:"foreignKey:UserID"`
	DeviceTelemetryData []DeviceTelemetryData `json:"DeviceTelemetryData,omitempty"`
	CreatedAt           time.Time             `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt           time.Time             `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (d *Device) CreateDevice() error {
	if err := database.DB.Create(&d).Error; err != nil {
		return err
	}
	return nil
}

func GetDevices() ([]Device, error) {
	var devices []Device
	if err := database.DB.Preload("User").Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

func GetDevicesByOrganization(organizationId uint) ([]Device, error) {
	var devices []Device

	if err := database.DB.Preload("User").Joins("JOIN users on devices.user_id = users.id").Where("users.organization_id = ?", organizationId).Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

func GetDevicesByUser(userId uint) ([]Device, error) {
	var devices []Device
	if err := database.DB.Preload("User").Where("user_id = ?", userId).Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

func (d *Device) GetDeviceByIMEI() error {
	if err := database.DB.Preload("User").Where("imei = ?", d.IMEI).First(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *Device) GetDevice() error {
	if err := database.DB.Preload("User").Where("id = ?", d.ID).First(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *Device) UpdateDevice() error {
	if err := database.DB.Save(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *Device) DeleteDevice() error {
	if err := database.DB.Delete(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *Device) UpdateBattery(batteryLevel uint) error {
	if err := database.DB.Model(&Device{}).Where("id = ?", d.ID).Update("battery_level", batteryLevel).Error; err != nil {
		return err
	}
	return nil
}

func (d *Device) GetAvailableDevices() ([]Device, error) {
	var devices []Device

	if err := database.DB.Where("user_id IS NULL").Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}
