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
	if err := database.DB.Preload("Device").Find(&deviceTelemetryData).Error; err != nil {
		return nil, err
	}

	return deviceTelemetryData, nil
}

func GetDeviceTelemetryDataByDevice(deviceId uint) ([]DeviceTelemetryData, error) {
	var deviceTelemetryData []DeviceTelemetryData
	if err := database.DB.Preload("Device").Where("device_id = ?", deviceId).Find(&deviceTelemetryData).Error; err != nil {
		return nil, err
	}

	return deviceTelemetryData, nil
}

func GetDeviceTelemetryDataByDeviceBetweenDates(deviceId uint, startDate, endDate time.Time) ([]DeviceTelemetryData, error) {
	var deviceTelemetryData []DeviceTelemetryData
	if err := database.DB.Preload("Device").Where("device_id = ? AND measured_at BETWEEN ? AND ?", deviceId, startDate, endDate).Find(&deviceTelemetryData).Error; err != nil {
		return nil, err
	}

	return deviceTelemetryData, nil
}

func GetLatestDeviceTelemetryDataByDevice(deviceId uint) (DeviceTelemetryData, error) {
	var deviceTelemetryData DeviceTelemetryData
	if err := database.DB.Preload("Device").Where("device_id = ?", deviceId).Last(&deviceTelemetryData).Error; err != nil {
		return deviceTelemetryData, err
	}

	return deviceTelemetryData, nil
}

func (d *DeviceTelemetryData) GetDeviceTelemetryData() error {
	if err := database.DB.Preload("Device").Where("id = ?", d.ID).First(&d).Error; err != nil {
		return err
	}
	return nil
}

func GetNumberOfTelemetryEntriesThisWeek(deviceId uint) (int64, error) {
	var count int64
	if err := database.DB.Model(&DeviceTelemetryData{}).Where("device_id = ? AND measured_at BETWEEN ? AND ?", deviceId, time.Now().AddDate(0, 0, -7), time.Now()).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
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

type DeviceTelemetryDataForPatient struct {
	PatientID uint
	DeviceTelemetryData
}

func GetPatientTelemetryData(patientIDs []uint) ([]DeviceTelemetryDataForPatient, error) {
	var result []struct {
		DeviceTelemetryData
		UserID uint
	}

	selects := []string{
		"device_telemetry_data.*",
		"devices.user_id as user_id",
	}

	err := database.DB.Model(&DeviceTelemetryData{}).
		Select(selects).
		Joins("JOIN devices ON device_telemetry_data.device_id = devices.id").
		Where("devices.user_id IN ?", patientIDs).
		Order("devices.user_id, device_telemetry_data.measured_at DESC").
		Find(&result).Error

	if err != nil {
		return nil, err
	}

	var telemetryData []DeviceTelemetryDataForPatient

	for _, r := range result {
		telemetryData = append(telemetryData, DeviceTelemetryDataForPatient{
			PatientID:           r.UserID,
			DeviceTelemetryData: r.DeviceTelemetryData,
		})
	}

	return telemetryData, nil
}

func GetPatientStatusFunc(thresholds []AlertThreshold) func(data DeviceTelemetryData) (isCritical, isWarning bool) {

	var systolicBPThreshold, diastolicBPThreshold AlertThreshold

	for _, t := range thresholds {
		if t.DeviceType == BloodPressure {
			if t.MeasurementType == Systolic {
				systolicBPThreshold = t
			} else if t.MeasurementType == Diastolic {
				diastolicBPThreshold = t
			}
		}
	}

	return func(data DeviceTelemetryData) (isCritical, isWarning bool) {
		if systolicBPThreshold.CriticalLow != nil {
			if data.SystolicBP < *systolicBPThreshold.CriticalLow || data.SystolicBP > *systolicBPThreshold.CriticalHigh {
				return true, false
			}
		}
		if diastolicBPThreshold.CriticalLow != nil {
			if data.DiastolicBP < *diastolicBPThreshold.CriticalLow || data.DiastolicBP > *diastolicBPThreshold.CriticalHigh {
				return true, false
			}
		}
		if systolicBPThreshold.WarningLow != nil {
			if data.SystolicBP < *systolicBPThreshold.WarningLow || data.SystolicBP > *systolicBPThreshold.WarningHigh {
				return false, true
			}
		}
		if diastolicBPThreshold.WarningLow != nil {
			if data.DiastolicBP < *diastolicBPThreshold.WarningLow || data.DiastolicBP > *diastolicBPThreshold.WarningHigh {
				return false, true
			}
		}
		return false, false
	}
}
