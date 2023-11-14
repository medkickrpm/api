package migration

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/database/models"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase(database.Config())

	err = database.DB.AutoMigrate(
		&models.User{},
		&models.PasswordReset{},
		&models.Organization{},
		&models.CarePlan{},
		&models.Interaction{},
		&models.Device{},
		&models.DeviceStatusData{},
		&models.DeviceTelemetryData{},
		&models.DeviceLogData{},
		&models.UserVerification{},
		&models.AlertThreshold{},
		&models.InteractionSetting{},
		&models.TelemetryAlert{},
		&models.Service{},
		&models.PatientService{},
		&models.LastBillEntry{},
		&models.Bill{},
	)
	if err != nil {
		panic("Could not migrate database")
	}

	fmt.Println("Database migrated successfully.")
}
