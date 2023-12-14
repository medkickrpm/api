package models

import (
	"MedKick-backend/pkg/database"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	Provider          string       `json:"provider,omitempty" example:"Test Provider"`
	Device            []Device
	PatientDiagnosis  []PatientDiagnosis
	Interaction       []Interaction
	CreatedAt         time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt         time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type DeviceTelemetryDataResponse struct {
	ID uint `json:"id"`
	//Sphygmomanometer
	SystolicBP         uint `json:"systolic_bp"`
	DiastolicBP        uint `json:"diastolic_bp"`
	Pulse              uint `json:"pulse"`
	IrregularHeartBeat bool `json:"irregular_heartbeat"`
	HandShaking        bool `json:"hand_shaking"`
	TripleMeasurement  bool `json:"triple_measurement"`
	//Weight Scale
	Weight           uint `json:"weight"`
	WeightStableTime uint `json:"weight_stable_time"`
	WeightLockCount  uint `json:"weight_lock_count"`
	//Blood Glucose Meter
	BloodGlucose uint   `json:"blood_glucose"`
	Unit         string `json:"unit"`
	TestPaper    string `json:"test_paper"`
	SampleType   string `json:"sample_type"`
	Meal         string `json:"meal"`

	DeviceID   uint      `json:"device_id"`
	MeasuredAt time.Time `json:"measured_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DeviceResponse struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	ModelNumber         string `json:"model_number"`
	IMEI                string `json:"imei"`
	SerialNumber        string `json:"serial_number"`
	BatteryLevel        uint   `json:"battery_level"`
	SignalStrength      string `json:"signal_strength"`
	FirmwareVersion     string `json:"firmware_version"`
	UserID              uint   `json:"user_id"`
	DeviceTelemetryData *DeviceTelemetryDataResponse
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type DignosesResponse struct {
	UserID      uint      `json:"user_id" gorm:"primaryKey"`
	DiagnosisID uint      `json:"diagnosis_id" gorm:"primaryKey"`
	Diagnosis   Diagnosis `json:"diagnosis" gorm:"foreignKey:DiagnosisID"`
	CreatedAt   time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
}

type InteractionResponse struct {
	ID           uint      `json:"id" gorm:"primary_key;auto_increment" example:"1"`
	UserID       uint      `json:"user_id" gorm:"not null" example:"1"`
	DoctorID     uint      `json:"doctor_id" gorm:"not null" example:"1"`
	Doctor       User      `json:"doctor" gorm:"foreignKey:DoctorID"`
	Duration     uint      `json:"duration" gorm:"not null" example:"30"`
	Notes        string    `json:"notes" gorm:"not null" example:"Patient is doing well"`
	CostCategory string    `json:"cost_category" gorm:"not null" example:""`
	SessionDate  time.Time `json:"session_date" gorm:"not null" example:"2021-01-01T00:00:00Z"`
	CreatedAt    time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type MainInterActionsResponse struct {
	TotalDuration uint
	Readings      int
	ReadingDate   time.Time
}

type UserResponse struct {
	ID                uint               `json:"id"`
	FirstName         string             `json:"first_name"`
	LastName          string             `json:"last_name"`
	Email             string             `json:"email"`
	Role              string             `json:"role"`
	DOB               string             `json:"dob"`
	Location          string             `json:"location"`
	AvatarSrc         string             `json:"avatar_src"`
	InsuranceProvider string             `json:"insurance_provider"`
	InsuranceID       string             `json:"insurance_id"`
	Organization      Organization       `json:"organization"`
	PatientDiagnosis  []DignosesResponse `json:"patient_diagnosis"`
	Devices           []DeviceResponse   `json:"devices"`
	// Interactions      MainInterActionsResponse `json:"interactions,omitempty"`
	TotalDuration uint      `json:"total_duration" example:"30"`
	Readings      int       `json:"readings" example:"1"`
	ReadingDate   time.Time `json:"reading_date" example:"2021-01-01T00:00:00Z"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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

func GetAllPatients() ([]UserResponse, error) {
	var userResponses []UserResponse
	var users []User
	// Set the date range for the current month
	startDate := time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC) // First day of the current month
	endDate := time.Now()

	if err := database.DB.
		Where("role = 'patient'").
		Select("id", "first_name", "last_name", "email", "phone", "password", "role", "dob", "Location", "avatar_src", "insurance_provider", "insurance_id", "organization_id", "created_at", "updated_at").
		Preload("Organization").
		Preload("Device", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "model_number", "imei", "serial_number", "battery_level", "signal_strength", "firmware_version", "user_id").
				Preload("DeviceTelemetryData", func(db *gorm.DB) *gorm.DB {
					return db.Order("created_at desc").Limit(2)
				}) // Specify the fields you want from the Devices table
		}).
		Preload("PatientDiagnosis", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Diagnosis")
		}).
		Preload("Interaction", func(db *gorm.DB) *gorm.DB {
			return db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}).
		Find(&users).Error; err != nil {
		return nil, err
	}

	for _, user := range users {
		userResponses = append(userResponses, user.SanitizedUserResponse())
	}

	return userResponses, nil
}

func GetPatient(id uint) (*UserResponse, error) {
	var userResponses UserResponse
	var users *User
	// Set the date range for the current month
	startDate := time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC) // First day of the current month
	endDate := time.Now()

	if err := database.DB.
		Where("id = ?", id).
		Where("role = 'patient'").
		Select("id", "first_name", "last_name", "email", "phone", "password", "role", "dob", "Location", "avatar_src", "insurance_provider", "insurance_id", "organization_id", "created_at", "updated_at").
		Preload("Organization").
		Preload("Device", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "model_number", "imei", "serial_number", "battery_level", "signal_strength", "firmware_version", "user_id").
				Preload("DeviceTelemetryData", func(db *gorm.DB) *gorm.DB {
					return db.Order("created_at desc").Limit(2)
				}) // Specify the fields you want from the Devices table
		}).
		Preload("PatientDiagnosis", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Diagnosis")
		}).
		Preload("Interaction", func(db *gorm.DB) *gorm.DB {
			return db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}).
		Find(&users).Error; err != nil {
		return nil, err
	}

	if users.ID == nil {
		return nil, errors.New("User Not Found")
	} else {
		userResponses = users.SanitizedUserResponse()

		return &userResponses, nil
	}

}

func GetAllPatientsWithOrg(orgId uint64) ([]UserResponse, error) {
	var userResponses []UserResponse
	var users []User
	// Set the date range for the current month
	startDate := time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC) // First day of the current month
	endDate := time.Now()

	if err := database.DB.
		Where("role = 'patient'").
		Where("organization_id = ?", orgId).
		Select("id", "first_name", "last_name", "email", "phone", "password", "role", "dob", "Location", "avatar_src", "insurance_provider", "insurance_id", "organization_id", "created_at", "updated_at").
		Preload("Organization").
		Preload("Device", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "model_number", "imei", "serial_number", "battery_level", "signal_strength", "firmware_version", "user_id").
				Preload("DeviceTelemetryData", func(db *gorm.DB) *gorm.DB {
					return db.Order("created_at desc").Limit(2)
				}) // Specify the fields you want from the Devices table
		}).
		Preload("PatientDiagnosis", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Diagnosis")
		}).
		Preload("Interaction", func(db *gorm.DB) *gorm.DB {
			return db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}).
		Find(&users).Error; err != nil {
		return nil, err
	}

	for _, user := range users {
		userResponses = append(userResponses, user.SanitizedUserResponse())
	}

	return userResponses, nil
}

func CountUsers() (int64, error) {
	var count int64
	if err := database.DB.Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func CountUsersInOrg(orgId uint) (int64, error) {
	var count int64
	if err := database.DB.Model(&User{}).Where("organization_id = ?", orgId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetUsersWithRole(role string) ([]User, error) {
	var users []User
	if err := database.DB.Where("role = ?", role).Preload("Organization").Find(&users).Error; err != nil {
		return nil, err
	}

	for i := range users {
		users[i].SanitizeUser()
	}

	return users, nil
}

func CountUsersWithRole(role string) (int64, error) {
	var count int64
	if err := database.DB.Model(&User{}).Where("role = ?", role).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func CountUsersWithRoleInOrg(orgId uint, role string) (int64, error) {
	var count int64
	if err := database.DB.Model(&User{}).Where("organization_id = ? AND role = ?", orgId, role).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func GetUsersInOrg(orgId *uint) ([]User, error) {
	var users []User
	if err := database.DB.Where("organization_id = ? AND role !='patient'", orgId).Preload("Organization").Find(&users).Error; err != nil {
		return nil, err
	}

	for i := range users {
		users[i].SanitizeUser()
	}

	return users, nil
}

func GetUsersInOrgWithRole(orgId *uint, role string) ([]User, error) {
	var users []User
	if err := database.DB.Where("organization_id = ? AND role = ?", orgId, role).Preload("Organization").Find(&users).Error; err != nil {
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

// check if a user already exist with the same phone number
func (u *User) GetUserByPhone() error {
	var user User
	if err := database.DB.Where("phone = ?", u.Phone).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserV2() (*UserResponse, error) {
	var user User
	if u.Email != "" {
		if err := database.DB.Where("email = ?", u.Email).Preload("Organization").First(&user).Error; err != nil {
			return nil, err
		}
		userResponse := user.SanitizedUserResponse()
		return &userResponse, nil
	}
	if err := database.DB.Where("id = ?", u.ID).Preload("Organization").Preload("Device", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "model_number", "imei", "serial_number", "battery_level", "signal_strength", "firmware_version", "user_id") // Specify the fields you want from the Devices table
	}).First(&user).Error; err != nil {
		return nil, err
	}

	userResponse := user.SanitizedUserResponse()
	return &userResponse, nil
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

	u.SanitizeUser()
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

func (user *User) SanitizedUserResponse() UserResponse {
	var devices []DeviceResponse
	var interactions []InteractionResponse
	var reading time.Time
	for index, device := range user.Device {
		dataExists := false // Set your condition here
		var DeviceTelemetries DeviceTelemetryDataResponse
		if len(device.DeviceTelemetryData) > 0 {
			dataExists = true
			telemetry := device.DeviceTelemetryData[0]
			DeviceTelemetries = DeviceTelemetryDataResponse{
				ID:                 telemetry.ID,
				SystolicBP:         telemetry.SystolicBP,
				DiastolicBP:        telemetry.DiastolicBP,
				Pulse:              telemetry.Pulse,
				IrregularHeartBeat: telemetry.IrregularHeartBeat,
				HandShaking:        telemetry.HandShaking,
				TripleMeasurement:  telemetry.TripleMeasurement,
				Weight:             telemetry.Weight,
				WeightStableTime:   telemetry.WeightStableTime,
				WeightLockCount:    telemetry.WeightLockCount,
				BloodGlucose:       telemetry.BloodGlucose,
				Unit:               telemetry.Unit,
				TestPaper:          telemetry.TestPaper,
				SampleType:         telemetry.SampleType,
				Meal:               telemetry.Meal,
				DeviceID:           telemetry.DeviceID,
				MeasuredAt:         telemetry.MeasuredAt,
				CreatedAt:          telemetry.CreatedAt,
				UpdatedAt:          telemetry.UpdatedAt,
			}
			reading = telemetry.MeasuredAt
		}

		devices = append(devices, DeviceResponse{
			ID:              device.ID,
			Name:            device.Name,
			ModelNumber:     device.ModelNumber,
			IMEI:            device.IMEI,
			SerialNumber:    device.SerialNumber,
			BatteryLevel:    device.BatteryLevel,
			SignalStrength:  device.SignalStrength,
			FirmwareVersion: device.FirmwareVersion,
			UserID:          device.UserID,
			CreatedAt:       device.CreatedAt,
			UpdatedAt:       device.UpdatedAt,
		})

		if dataExists {
			devices[index].DeviceTelemetryData = &DeviceTelemetries

		} else {
			devices[index].DeviceTelemetryData = nil
		}
	}

	var Dignoses []DignosesResponse
	for _, dignoses := range user.PatientDiagnosis {
		Dignoses = append(Dignoses, DignosesResponse{
			UserID:      dignoses.UserID,
			DiagnosisID: dignoses.DiagnosisID,
			Diagnosis:   dignoses.Diagnosis,
			CreatedAt:   dignoses.CreatedAt,
		})
	}

	var duration uint = 0
	for _, interaction := range user.Interaction {
		interactions = append(interactions, InteractionResponse{
			ID:           interaction.ID,
			UserID:       interaction.UserID,
			DoctorID:     interaction.DoctorID,
			Doctor:       interaction.Doctor,
			Duration:     interaction.Duration,
			Notes:        interaction.Notes,
			CostCategory: interaction.CostCategory,
			SessionDate:  interaction.SessionDate,
			CreatedAt:    interaction.CreatedAt,
			UpdatedAt:    interaction.UpdatedAt,
		})
		duration += interaction.Duration
	}

	response := UserResponse{
		ID:                *user.ID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		Role:              user.Role,
		DOB:               user.DOB,
		Location:          user.Location,
		AvatarSrc:         user.AvatarSRC,
		InsuranceProvider: user.InsuranceProvider,
		InsuranceID:       user.InsuranceID,
		Organization: Organization{
			ID:        user.Organization.ID,
			Name:      user.Organization.Name,
			Address:   user.Organization.Address,
			Address2:  user.Organization.Address2,
			City:      user.Organization.City,
			State:     user.Organization.State,
			Zip:       user.Organization.Zip,
			Country:   user.Organization.Country,
			Phone:     user.Organization.Phone,
			CreatedAt: user.Organization.CreatedAt,
			UpdatedAt: user.Organization.UpdatedAt,
		},
		PatientDiagnosis: Dignoses,
		TotalDuration:    duration,
		Readings:         len(interactions),
		ReadingDate:      reading,
		Devices:          devices,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}

	return response
}
