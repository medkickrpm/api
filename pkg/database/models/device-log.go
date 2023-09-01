package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type DeviceLogData struct {
	ID uint `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	//Device
	Battery         string    `json:"battery" gorm:"null" example:"100"`
	HardwareVersion string    `json:"hardware_version" gorm:"null" example:"1.0.0"`
	MCUVersion      string    `json:"mcu_version" gorm:"null" example:"1.0.0"`
	AppVersion      string    `json:"app_version" gorm:"null" example:"1.0.0"`
	ModemVersion    string    `json:"modem_version" gorm:"null" example:"1.0.0"`
	BPMAlgo         string    `json:"bpm_algo" gorm:"null" example:"1.0.0"`
	DeviceID        uint      `json:"device_id" gorm:"not null" example:"1"`
	Device          Device    `json:"device" gorm:"foreignKey:DeviceID"`
	CreatedAt       time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt       time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (d *DeviceLogData) CreateDeviceLogData() error {
	if err := database.DB.Create(&d).Error; err != nil {
		return err
	}
	return nil
}

func GetDeviceLogData() ([]DeviceLogData, error) {
	var deviceLogData []DeviceLogData
	if err := database.DB.Find(&deviceLogData).Error; err != nil {
		return nil, err
	}

	return deviceLogData, nil
}

func GetDeviceLogDataByDevice(deviceId uint) ([]DeviceLogData, error) {
	var deviceLogData []DeviceLogData
	if err := database.DB.Where("device_id = ?", deviceId).Find(&deviceLogData).Error; err != nil {
		return nil, err
	}

	return deviceLogData, nil
}

func (d *DeviceLogData) GetDeviceLogData() error {
	if err := database.DB.Where("id = ?", d.ID).First(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *DeviceLogData) UpdateDeviceLogData() error {
	if err := database.DB.Save(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *DeviceLogData) DeleteDeviceLogData() error {
	if err := database.DB.Delete(&d).Error; err != nil {
		return err
	}
	return nil
}
