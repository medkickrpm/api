package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type Organization struct {
	ID        uint      `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	Name      string    `json:"name" gorm:"not null" example:"John Hopkins"`
	Address   string    `json:"address" gorm:"not null" example:"123 Main St"`
	Address2  string    `json:"address2" gorm:"not null" example:"Apt 1"`
	City      string    `json:"city" gorm:"not null" example:"Baltimore"`
	State     string    `json:"state" gorm:"not null" example:"MD"`
	Zip       string    `json:"zip" gorm:"not null" example:"12345"`
	Country   string    `json:"country" gorm:"not null" example:"USA"`
	Phone     string    `json:"phone" gorm:"not null" example:"08123456789"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (o *Organization) CreateOrganization() error {
	if err := database.DB.Create(&o).Error; err != nil {
		return err
	}
	return nil
}

func GetOrganizations() ([]Organization, error) {
	var organizations []Organization
	if err := database.DB.Find(&organizations).Error; err != nil {
		return nil, err
	}

	return organizations, nil
}

func (o *Organization) GetOrganization() error {
	if err := database.DB.Where("id = ?", o.ID).First(&o).Error; err != nil {
		return err
	}
	return nil
}

func (o *Organization) UpdateOrganization() error {
	if err := database.DB.Save(&o).Error; err != nil {
		return err
	}
	return nil
}

func (o *Organization) DeleteOrganization() error {
	if err := database.DB.Delete(&o).Error; err != nil {
		return err
	}
	return nil
}
