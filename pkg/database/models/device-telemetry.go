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
	PatientID  uint
	DeviceName string
	DeviceType DeviceType
	DeviceTelemetryData
}

func GetPatientTelemetryData(patientIDs []uint) ([]DeviceTelemetryDataForPatient, error) {
	var result []struct {
		DeviceTelemetryData
		DeviceName string
		UserID     uint
	}

	selects := []string{
		"device_telemetry_data.*",
		"devices.name as device_name",
		"devices.user_id as user_id",
	}

	err := database.DB.Model(&DeviceTelemetryData{}).
		Select(selects).
		Joins("JOIN devices ON device_telemetry_data.device_id = devices.id").
		Where("devices.user_id IN ?", patientIDs).
		Where("device_telemetry_data.measured_at > ?", time.Now().AddDate(0, 0, -30)).
		Order("devices.user_id, device_telemetry_data.measured_at DESC").
		Find(&result).Error

	if err != nil {
		return nil, err
	}

	var telemetryData []DeviceTelemetryDataForPatient

	for _, r := range result {
		telemetryData = append(telemetryData, DeviceTelemetryDataForPatient{
			PatientID:           r.UserID,
			DeviceName:          r.DeviceName,
			DeviceTelemetryData: r.DeviceTelemetryData,
		})
	}

	return telemetryData, nil
}

func GetLatestPatientTelemetryData(data []DeviceTelemetryDataForPatient) []DeviceTelemetryDataForPatient {
	sphygmomanometerData := make(map[uint]DeviceTelemetryDataForPatient)
	weightScaleData := make(map[uint]DeviceTelemetryDataForPatient)
	bloodGlucoseMeterData := make(map[uint]DeviceTelemetryDataForPatient)

	for _, d := range data {
		if d.DeviceName == "Sphygmomanometer" {
			if _, ok := sphygmomanometerData[d.PatientID]; !ok {
				d.DeviceType = BloodPressure
				sphygmomanometerData[d.PatientID] = d
			}
		} else if d.DeviceName == "Weight Scale" {
			if _, ok := weightScaleData[d.PatientID]; !ok {
				d.DeviceType = WeightScale
				weightScaleData[d.PatientID] = d
			}
		} else if d.DeviceName == "Blood Glucose Meter" {
			if _, ok := bloodGlucoseMeterData[d.PatientID]; !ok {
				d.DeviceType = BloodGlucose
				bloodGlucoseMeterData[d.PatientID] = d
			}
		}
	}

	var result []DeviceTelemetryDataForPatient

	for _, d := range sphygmomanometerData {
		result = append(result, d)
	}

	for _, d := range weightScaleData {
		result = append(result, d)
	}

	for _, d := range bloodGlucoseMeterData {
		result = append(result, d)
	}

	return result
}

func GetPatientStatusFunc(thresholds []AlertThreshold) func(data DeviceTelemetryDataForPatient) (isCritical bool, isWarning bool) {

	systolicBPThresholdMap := make(map[uint]AlertThreshold)
	diastolicBPThresholdMap := make(map[uint]AlertThreshold)
	weightThresholdMap := make(map[uint]AlertThreshold)

	for _, t := range thresholds {
		if t.DeviceType == BloodPressure {
			if t.MeasurementType == Systolic {
				systolicBPThresholdMap[t.PatientID] = t
			} else if t.MeasurementType == Diastolic {
				diastolicBPThresholdMap[t.PatientID] = t
			}
		}
		if t.DeviceType == WeightScale {
			if t.MeasurementType == Weight {
				weightThresholdMap[t.PatientID] = t
			}
		}
	}

	return func(data DeviceTelemetryDataForPatient) (isCritical, isWarning bool) {
		systolicBPThreshold := systolicBPThresholdMap[data.PatientID]
		diastolicBPThreshold := diastolicBPThresholdMap[data.PatientID]
		weightThreshold := weightThresholdMap[data.PatientID]

		if data.DeviceType == BloodPressure {
			if (systolicBPThreshold.CriticalLow != nil && data.SystolicBP < *systolicBPThreshold.CriticalLow) ||
				(systolicBPThreshold.CriticalHigh != nil && data.SystolicBP > *systolicBPThreshold.CriticalHigh) ||
				(diastolicBPThreshold.CriticalLow != nil && data.DiastolicBP < *diastolicBPThreshold.CriticalLow) ||
				(diastolicBPThreshold.CriticalHigh != nil && data.DiastolicBP > *diastolicBPThreshold.CriticalHigh) {
				return true, false
			}

			if (systolicBPThreshold.WarningLow != nil && data.SystolicBP < *systolicBPThreshold.WarningLow) ||
				(systolicBPThreshold.WarningHigh != nil && data.SystolicBP > *systolicBPThreshold.WarningHigh) ||
				(diastolicBPThreshold.WarningLow != nil && data.DiastolicBP < *diastolicBPThreshold.WarningLow) ||
				(diastolicBPThreshold.WarningHigh != nil && data.DiastolicBP > *diastolicBPThreshold.WarningHigh) {
				return false, true
			}
		}

		if data.DeviceType == WeightScale {
			if (weightThreshold.CriticalLow != nil && data.Weight < *weightThreshold.CriticalLow) ||
				(weightThreshold.CriticalHigh != nil && data.Weight > *weightThreshold.CriticalHigh) {
				return true, false
			}

			if (weightThreshold.WarningLow != nil && data.Weight < *weightThreshold.WarningLow) ||
				(weightThreshold.WarningHigh != nil && data.Weight > *weightThreshold.WarningHigh) {
				return false, true
			}
		}

		return false, false
	}
}

func (d DeviceTelemetryData) GetStatusByPatientThreshold(deviceType DeviceType, thresholds []AlertThreshold) AlertType {

	alertType := AlertOk

	if deviceType == BloodPressure {
		for _, t := range thresholds {
			if t.DeviceType == BloodPressure {
				if t.MeasurementType == Systolic {
					if (t.CriticalLow != nil && d.SystolicBP < *t.CriticalLow) ||
						(t.CriticalHigh != nil && d.SystolicBP > *t.CriticalHigh) {
						alertType = AlertCritical
					}

					if (t.WarningLow != nil && d.SystolicBP < *t.WarningLow) ||
						(t.WarningHigh != nil && d.SystolicBP > *t.WarningHigh) {
						if alertType != AlertCritical {
							alertType = AlertWarning
						}
					}
				}

				if t.MeasurementType == Diastolic {
					if (t.CriticalLow != nil && d.DiastolicBP < *t.CriticalLow) ||
						(t.CriticalHigh != nil && d.DiastolicBP > *t.CriticalHigh) {
						alertType = AlertCritical
					}

					if (t.WarningLow != nil && d.DiastolicBP < *t.WarningLow) ||
						(t.WarningHigh != nil && d.DiastolicBP > *t.WarningHigh) {
						if alertType != AlertCritical {
							alertType = AlertWarning
						}
					}
				}
			}
		}
	}

	return alertType
}
