package models

import (
	"MedKick-backend/pkg/database"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID                *uint        `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	FirstName         string       `json:"first_name" gorm:"not null" example:"John"`
	LastName          string       `json:"last_name" gorm:"not null" example:"Doe"`
	Email             string       `json:"email" gorm:"not null;unique"`
	Phone             string       `json:"phone" gorm:"not null;unique" example:"08123456789"`
	Password          string       `json:"password" gorm:"not null" example:"123456"`
	Role              string       `json:"role" gorm:"not null" example:"admin"`
	DOB               string       `json:"dob" gorm:"not null" example:"2000-01-01"`
	Location          string       `json:"Location" gorm:"not null" example:"Dallas, TX"`
	AvatarSRC         string       `json:"avatar_src" gorm:"not null" example:"https://cdn.med-kick.com/xxx.jpg"`
	InsuranceProvider string       `json:"insurance_provider" gorm:"not null" example:"Aetna"`
	InsuranceID       string       `json:"insurance_id" gorm:"not null" example:"123456789"`
	OrganizationID    *uint        `json:"organization_id" gorm:"null" example:"1"`
	Organization      Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	CreatedAt         time.Time    `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt         time.Time    `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (u *User) CreateUser() error {
	if err := database.DB.Create(&u).Error; err != nil {
		return err
	}
	return nil
}

func GetUsers() ([]User, error) {
	var users []User
	if err := database.DB.Preload("Organization").Find(&users).Error; err != nil {
		return nil, err
	}

	for i := range users {
		users[i].SanitizeUser()
	}

	return users, nil
}

func GetUsersInOrg(orgId *uint) ([]User, error) {
	var users []User
	if err := database.DB.Where("organization_id = ?", orgId).Preload("Organization").Find(&users).Error; err != nil {
		return nil, err
	}

	for i := range users {
		users[i].SanitizeUser()
	}

	return users, nil
}

func (u *User) GetUser() error {
	if u.Email != "" {
		if err := database.DB.Where("email = ?", u.Email).Preload("Organization").First(&u).Error; err != nil {
			return err
		}
		u.SanitizeUser()
		return nil
	}
	if err := database.DB.Where("id = ?", u.ID).Preload("Organization").First(&u).Error; err != nil {
		return err
	}
	u.SanitizeUser()
	return nil
}

func (u *User) GetUserRaw() error {
	if u.Email != "" {
		if err := database.DB.Where("email = ?", u.Email).Preload("Organization").First(&u).Error; err != nil {
			return err
		}
		return nil
	}
	if err := database.DB.Where("id = ?", u.ID).Preload("Organization").First(&u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateUser() error {
	if err := database.DB.Save(&u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) DeleteUser() error {
	if err := database.DB.Delete(&u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) SanitizeUser() {
	u.Password = ""
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return errors.New("[User-Model] failed to hash password")
	}

	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) UpdatePassword(newPassword string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return errors.New("[User-Model] failed to hash password")
	}

	if err := database.DB.Model(&u).Update("password", string(bytes)).Error; err != nil {
		return err
	}

	return nil
}
