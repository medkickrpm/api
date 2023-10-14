package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type UserVerification struct {
	UUID      string    `json:"uuid" gorm:"primary_key;not null" example:"1"`
	UserID    *uint     `json:"user_id" gorm:"not null" example:"1"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (uv *UserVerification) CreateUserVerification() error {
	if err := database.DB.Create(&uv).Error; err != nil {
		return err
	}
	return nil
}

func GetUserVerifications() ([]UserVerification, error) {
	var userVerifications []UserVerification
	if err := database.DB.Find(&userVerifications).Error; err != nil {
		return nil, err
	}
	return userVerifications, nil
}

func (uv *UserVerification) GetUserVerification() error {
	if err := database.DB.Preload("User").Where("uuid = ?", uv.UUID).First(&uv).Error; err != nil {
		return err
	}
	return nil
}

func (uv *UserVerification) DeleteUserVerification() error {
	if err := database.DB.Delete(&uv).Error; err != nil {
		return err
	}
	return nil
}
