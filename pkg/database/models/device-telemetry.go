package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type DeviceTelemetryData struct {
	ID uint `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	//Sphygmomanometer
	SystolicBP         uint `json:"systolic_bp" gorm:"null" example:"120"`
	DiastolicBP        uint `json:"diastolic_bp" gorm:"null" example:"80"`
	Pulse              uint `json:"pulse" gorm:"null" example:"80"`
	IrregularHeartBeat bool `json:"irregular_heartbeat" gorm:"null" example:"false"`
	HandShaking        bool `json:"hand_shaking" gorm:"null" example:"false"`
	TripleMeasurement  bool `json:"triple_measurement" gorm:"null" example:"false"`
	//Weight Scale
	Weight           uint `json:"weight" gorm:"null" example:"80"`
	WeightStableTime uint `json:"weight_stable_time" gorm:"null" example:"5"`
	WeightLockCount  uint `json:"weight_lock_count" gorm:"null" example:"3"`
	//Blood Glucose Meter
	BloodGlucose uint   `json:"blood_glucose" gorm:"null" example:"80"`
	Unit         string `json:"unit" gorm:"null" example:"mg/dL"`
	TestPaper    string `json:"test_paper" gorm:"null" example:"1. GOD; 2. GDH"`
	SampleType   string `json:"sample_type" gorm:"null" example:"1. Blood or Resistance; 2. Quality Control Liquid; 3. Sample is invalid"`
	Meal         string `json:"meal" gorm:"null" example:"1. Before Meal; 2. After Meal"`

	DeviceID   uint      `json:"device_id" gorm:"not null" example:"1"`
	Device     Device    `json:"device" gorm:"foreignKey:DeviceID"`
	MeasuredAt time.Time `json:"measured_at" example:"2021-01-01T00:00:00Z"`
	CreatedAt  time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (d *DeviceTelemetryData) CreateDeviceTelemetryData() error {
	if err := database.DB.Create(&d).Error; err != nil {
		return err
	}
	return nil
}

func GetDeviceTelemetryData() ([]DeviceTelemetryData, error) {
	var deviceTelemetryData []DeviceTelemetryData
	if err := database.DB.Find(&deviceTelemetryData).Error; err != nil {
		return nil, err
	}

	return deviceTelemetryData, nil
}

func GetDeviceTelemetryDataByDevice(deviceId uint) ([]DeviceTelemetryData, error) {
	var deviceTelemetryData []DeviceTelemetryData
	if err := database.DB.Where("device_id = ?", deviceId).Find(&deviceTelemetryData).Error; err != nil {
		return nil, err
	}

	return deviceTelemetryData, nil
}

func (d *DeviceTelemetryData) GetDeviceTelemetryData() error {
	if err := database.DB.Where("id = ?", d.ID).First(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *DeviceTelemetryData) UpdateDeviceTelemetryData() error {
	if err := database.DB.Save(&d).Error; err != nil {
		return err
	}
	return nil
}

func (d *DeviceTelemetryData) DeleteDeviceTelemetryData() error {
	if err := database.DB.Delete(&d).Error; err != nil {
		return err
	}
	return nil
}
