package models

import (
	"MedKick-backend/pkg/database"
	"time"

	"gorm.io/gorm/clause"
)

type InteractionSetting struct {
	ID             uint                   `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	OrganizationID uint                   `json:"organization_id" gorm:"index:,unique,composite:type; not null" example:"1"`
	Type           InteractionSettingType `json:"type" gorm:"index:,unique,composite:type; not null" example:"ColorThreshold"`
	Value          int64                  `json:"value" gorm:"not null" example:"1"`
	CreatedAt      time.Time              `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt      time.Time              `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type InteractionSettingType string

const (
	ColorThreshold InteractionSettingType = "ColorThreshold"
)

func (i *InteractionSetting) UpsertInteractionSetting() error {
	db := database.DB.Model(&InteractionSetting{})
	db = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "organization_id"}, {Name: "type"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"value",
		}),
	})

	return db.Create(&i).Error
}

func (i *InteractionSetting) GetInteractionSetting() error {
	if err := database.DB.Where("organization_id = ? AND type = ?", i.OrganizationID, i.Type).First(&i).Error; err != nil {
		return err
	}
	return nil
}
