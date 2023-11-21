package models

import (
	"MedKick-backend/pkg/database"
	"time"
)

type Diagnosis struct {
	ID   uint   `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	Code string `json:"code" gorm:"type:varchar(10); not null; uniqueIndex" example:"I10"`
}

func (Diagnosis) TableName() string {
	return "diagnoses"
}

func ListDiagnosisCodes() ([]Diagnosis, error) {
	var codes []Diagnosis
	if err := database.DB.Find(&codes).Error; err != nil {
		return nil, err
	}

	return codes, nil
}

type PatientDiagnosis struct {
	UserID      uint      `json:"user_id" gorm:"primaryKey"`
	Patient     User      `json:"patient" gorm:"foreignKey:UserID"`
	DiagnosisID uint      `json:"diagnosis_id" gorm:"primaryKey"`
	Diagnosis   Diagnosis `json:"diagnosis" gorm:"foreignKey:DiagnosisID"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
}

func (PatientDiagnosis) TableName() string {
	return "patient_diagnoses"
}

func DeletePatientDiagnoses(userId uint, diagnoses []uint) error {
	db := database.DB.Model(&PatientDiagnosis{})
	db = db.Where("user_id = ?", userId)
	if len(diagnoses) > 0 {
		db = db.Where("diagnosis_id IN ?", diagnoses)
	}
	err := db.Delete(&PatientDiagnosis{}).Error
	if err != nil {
		return err
	}
	return nil
}

func CreatePatientDiagnoses(diagnoses []PatientDiagnosis) error {
	err := database.DB.Model(&PatientDiagnosis{}).Create(&diagnoses).Error
	if err != nil {
		return err
	}
	return nil
}

func GetPatientDiagnoses(userId uint) ([]PatientDiagnosis, error) {
	var diagnoses []PatientDiagnosis
	db := database.DB.Model(&PatientDiagnosis{})
	db = db.Preload("Diagnosis")
	db = db.Where("user_id = ?", userId)

	err := db.Find(&diagnoses).Error
	if err != nil {
		return nil, err
	}
	return diagnoses, nil
}

func ListPatientDiagnosesCodeByPatientIDs(patientIDs []uint) (map[uint]string, error) {
	var data []struct {
		UserID uint
		Codes  string
	}

	if err := database.DB.Model(&PatientDiagnosis{}).
		Select("user_id, GROUP_CONCAT(diagnoses.code) AS codes").
		Joins("JOIN diagnoses ON patient_diagnoses.diagnosis_id = diagnoses.id").
		Where("user_id IN (?)", patientIDs).
		Group("user_id").
		Scan(&data).Error; err != nil {
		return nil, err
	}

	var result = make(map[uint]string)
	for _, d := range data {
		result[d.UserID] = d.Codes
	}

	return result, nil
}
