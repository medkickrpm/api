package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type DeviceStatusData struct {
	ID uint `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	//Device
	Timezone      string    `json:"timezone" gorm:"null" example:"UTC+6"`
	NetworkOps    string    `json:"network_ops" gorm:"null" example:"T-Mobile;Verizon"`
	NetworkFormat string    `json:"network_format" gorm:"null" example:"GSM;eMTC;NB-IoT"`
	Signal        uint      `json:"signal" gorm:"null" example:"100"`
	Temperature   int       `json:"temperature" gorm:"null" example:"100"`
	MeasureCount  uint      `json:"measure_count" gorm:"null" example:"100"`
	AttachTime    time.Time `json:"attach_time" gorm:"null" example:"100"`

	DeviceID  uint      `json:"device_id" gorm:"not null" example:"1"`
	Device    Device    `json:"device" gorm:"foreignKey:DeviceID"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (d *DeviceStatusData) CreateDeviceStatusData() error {
	if err := database.DB.Create(&d).Error; err != nil {
		return err
	}
	return nil
}

func GetDeviceStatusData() ([]DeviceStatusData, error) {
	var deviceStatusData []DeviceStatusData
	if err := database.DB.Find(&deviceStatusData).Error; err != nil {
		return nil, err
	}

	return deviceStatusData, nil
}

func GetDeviceStatusDataByDevice(deviceId uint) ([]DeviceStatusData, error) {
	var deviceStatusData []DeviceStatusData
	if err := database.DB.Where("device_id = ?", deviceId).Find(&deviceStatusData).Error; err != nil {
		return nil, err
	}

	return deviceStatusData, nil
}

func (d *DeviceStatusData) GetDeviceStatusData() error {
	if err := database.DB.Where("id = ?", d.ID).First(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *DeviceStatusData) UpdateDeviceStatusData() error {
	if err := database.DB.Save(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *DeviceStatusData) DeleteDeviceStatusData() error {
	if err := database.DB.Delete(&d).Error; err != nil {
		return err
	}
	return nil
}
