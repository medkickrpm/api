package organization

import (
	"MedKick-backend/pkg/database/models"
)

type InteractionSettingData struct {
	SettingType models.InteractionSettingType `json:"setting_type" validate:"required,oneof=ColorThreshold"`
	Value       int64                         `json:"value"`
}
