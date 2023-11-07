package organization

import (
	"MedKick-backend/pkg/database/models"
	"time"
)

type InteractionSettingData struct {
	SettingType models.InteractionSettingType `json:"setting_type" validate:"required,oneof=ColorThreshold"`
	Value       int64                         `json:"value"`
}

type TelemetryAlertResponse struct {
	AlertID     uint                   `json:"alert_id" example:"1"`
	PatientID   uint                   `json:"patient_id" example:"1"`
	PatientName string                 `json:"patient_name" example:"John Doe"`
	TelemetryID uint                   `json:"telemetry_id" example:"1"`
	PhoneNumber string                 `json:"phone_number" example:"08123456789"`
	Vitals      map[string]interface{} `json:"vitals"`
	Status      models.AlertType       `json:"status" example:"Critical"`
	IsActive    bool                   `json:"is_active" example:"true"`
	ResolvedBy  string                 `json:"resolved_by,omitempty" example:"John Doe"`
	Time        time.Time              `json:"time" example:"2021-01-01T00:00:00Z"`
}

func convertModelToResponse(data []models.TelemetryAlert) []TelemetryAlertResponse {
	res := make([]TelemetryAlertResponse, 0)
	for _, d := range data {
		rd := TelemetryAlertResponse{
			AlertID:     d.ID,
			PatientID:   d.PatientID,
			TelemetryID: d.TelemetryID,
			Vitals:      d.Data,
			Status:      d.AlertType,
			IsActive:    d.IsActive,
			Time:        d.MeasuredAt,
		}
		if d.Patient != nil {
			rd.PatientName = d.Patient.FirstName + " " + d.Patient.LastName
			rd.PhoneNumber = d.Patient.Phone
		}
		if d.ResolvedBy != nil {
			rd.ResolvedBy = d.ResolvedBy.FirstName + " " + d.ResolvedBy.LastName
		}

		res = append(res, rd)
	}

	return res
}
