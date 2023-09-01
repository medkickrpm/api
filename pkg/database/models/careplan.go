package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type CarePlan struct {
	ID        uint      `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	UserID    uint      `json:"user_id" gorm:"not null" example:"1"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	DoctorID  uint      `json:"doctor_id" gorm:"not null" example:"1"`
	Doctor    User      `json:"doctor" gorm:"foreignKey:DoctorID"`
	URL       string    `json:"url" gorm:"not null" example:"https://cdn.med-kick.com/xxx.pdf"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (c *CarePlan) CreateCarePlan() error {
	if err := database.DB.Create(&c).Error; err != nil {
		return err
	}
	return nil
}

func GetCarePlans() ([]CarePlan, error) {
	var careplans []CarePlan
	if err := database.DB.Find(&careplans).Error; err != nil {
		return nil, err
	}

	return careplans, nil
}

func GetCarePlansByUserID(id uint) ([]CarePlan, error) {
	var careplans []CarePlan
	if err := database.DB.Where("user_id = ?", id).Find(&careplans).Error; err != nil {
		return nil, err
	}

	return careplans, nil
}

func GetCarePlansByDoctorID(id uint) ([]CarePlan, error) {
	var careplans []CarePlan
	if err := database.DB.Where("doctor_id = ?", id).Find(&careplans).Error; err != nil {
		return nil, err
	}

	return careplans, nil
}

func (c *CarePlan) GetCarePlan() error {
	if err := database.DB.Where("id = ?", c.ID).First(&c).Error; err != nil {
		return err
	}
	return nil
}

func (c *CarePlan) UpdateCarePlan() error {
	if err := database.DB.Save(&c).Error; err != nil {
		return err
	}
	return nil
}

func (c *CarePlan) DeleteCarePlan() error {
	if err := database.DB.Delete(&c).Error; err != nil {
		return err
	}
	return nil
}
