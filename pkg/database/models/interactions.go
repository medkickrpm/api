package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type Interaction struct {
	ID          uint      `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	UserID      uint      `json:"user_id" gorm:"not null" example:"1"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	DoctorID    uint      `json:"doctor_id" gorm:"not null" example:"1"`
	Doctor      User      `json:"doctor" gorm:"foreignKey:DoctorID"`
	Duration    uint      `json:"duration" gorm:"not null" example:"30"`
	Notes       string    `json:"notes" gorm:"not null" example:"Patient is doing well"`
	SessionDate time.Time `json:"session_date" gorm:"not null" example:"2021-01-01T00:00:00Z"`
	CreatedAt   time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (i *Interaction) CreateInteraction() error {
	if err := database.DB.Create(&i).Error; err != nil {
		return err
	}
	return nil
}

func GetInteractions() ([]Interaction, error) {
	var interactions []Interaction
	if err := database.DB.Find(&interactions).Error; err != nil {
		return nil, err
	}

	return interactions, nil
}

func GetInteractionsByUser(userId uint) ([]Interaction, error) {
	var interactions []Interaction
	if err := database.DB.Where("user_id = ?", userId).Find(&interactions).Error; err != nil {
		return nil, err
	}

	return interactions, nil
}

func GetInteractionsByDoctor(doctorId uint) ([]Interaction, error) {
	var interactions []Interaction
	if err := database.DB.Where("doctor_id = ?", doctorId).Find(&interactions).Error; err != nil {
		return nil, err
	}

	return interactions, nil
}

func (i *Interaction) GetInteraction() error {
	if err := database.DB.Where("id = ?", i.ID).First(&i).Error; err != nil {
		return err
	}
	return nil
}

func (i *Interaction) UpdateInteraction() error {
	if err := database.DB.Save(&i).Error; err != nil {
		return err
	}
	return nil
}

func (i *Interaction) DeleteInteraction() error {
	if err := database.DB.Delete(&i).Error; err != nil {
		return err
	}
	return nil
}
