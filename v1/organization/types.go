package organization

import (
	"MedKick-backend/pkg/database/models"
	"errors"
)

type MeasurementData struct {
	MeasurementType models.MeasurementType `json:"measurement_type" validate:"required,oneof=Systolic Diastolic Pulse Weight"`
	CriticalLow     *uint                  `json:"critical_low"`
	WarningLow      *uint                  `json:"warning_low"`
	WarningHigh     *uint                  `json:"warning_high"`
	CriticalHigh    *uint                  `json:"critical_high"`
}

func (m MeasurementData) validate() error {
	if m.CriticalLow != nil && m.CriticalHigh != nil && *m.CriticalLow > *m.CriticalHigh {
		return errors.New("critical low must be less than critical high")
	}

	if m.WarningLow != nil && m.WarningHigh != nil && *m.WarningLow > *m.WarningHigh {
		return errors.New("warning low must be less than warning high")
	}

	if m.CriticalLow != nil && m.WarningLow != nil && *m.CriticalLow > *m.WarningLow {
		return errors.New("critical low must be less than warning low")
	}

	if m.CriticalHigh != nil && m.WarningHigh != nil && *m.CriticalHigh < *m.WarningHigh {
		return errors.New("critical high must be greater than warning high")
	}

	return nil
}

type AlertThresholdData struct {
	DeviceType   models.DeviceType `json:"device_type" validate:"required,oneof=BloodPressure BloodGlucose WeightScale"`
	Measurements []MeasurementData `json:"measurements" validate:"required,min=1,dive,required"`
	Note         string            `json:"note"`
}

func convertAlertThresholdModelToResponse(data []models.AlertThreshold) []AlertThresholdData {
	deviceMap := make(map[models.DeviceType][]MeasurementData)

	for _, d := range data {
		deviceMap[d.DeviceType] = append(deviceMap[d.DeviceType], MeasurementData{
			MeasurementType: d.MeasurementType,
			CriticalLow:     d.CriticalLow,
			WarningLow:      d.WarningLow,
			WarningHigh:     d.WarningHigh,
			CriticalHigh:    d.CriticalHigh,
		})
	}

	response := make([]AlertThresholdData, 0)

	for deviceType, measurements := range deviceMap {
		threshold := AlertThresholdData{
			DeviceType:   deviceType,
			Measurements: measurements,
		}
		if len(measurements) > 0 {
			threshold.Note = data[0].Note
		}
		response = append(response, threshold)
	}

	return response
}
