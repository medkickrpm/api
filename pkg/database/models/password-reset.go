package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type PasswordReset struct {
	UUID      string    `json:"uuid" gorm:"primary_key;not null" example:"1"`
	UserID    *uint     `json:"user_id" gorm:"not null" example:"1"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (pr *PasswordReset) CreatePasswordReset() error {
	if err := database.DB.Create(&pr).Error; err != nil {
		return err
	}
	return nil
}

func GetPasswordResets() ([]PasswordReset, error) {
	var passwordResets []PasswordReset
	if err := database.DB.Find(&passwordResets).Error; err != nil {
		return nil, err
	}
	return passwordResets, nil
}

func (pr *PasswordReset) GetPasswordReset() error {
	if err := database.DB.Where("uuid = ?", pr.UUID).First(&pr).Error; err != nil {
		return err
	}
	return nil
}

func (pr *PasswordReset) DeletePasswordReset() error {
	if err := database.DB.Delete(&pr).Error; err != nil {
		return err
	}
	return nil
}
