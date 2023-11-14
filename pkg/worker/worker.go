package worker

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/database/models"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

func ProcessCPTCode99453() error {
	var patientList []uint

	db := database.DB.Model(&models.PatientService{}).
		Joins("JOIN services ON services.id = patient_services.service_id").
		Where("services.is_enabled = ?", true).
		Where("services.code = ?", "RPM").
		Where("patient_services.ended_at IS NULL").
		Joins("LEFT JOIN last_bill_entries ON last_bill_entries.patient_id = patient_services.patient_id").
		Where("last_bill_entries.c99453 = 0 OR last_bill_entries.c99453 IS NULL")

	if err := db.Pluck("patient_services.patient_id", &patientList).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients: ", len(patientList))

	if len(patientList) == 0 {
		return nil
	}

	var filteredPatientList []uint
	db = database.DB.Model(&models.Device{}).Joins("JOIN device_telemetry_data ON device_telemetry_data.device_id = devices.id").
		Where("devices.user_id IN (?)", patientList).
		Where("device_telemetry_data.measured_at < ?", time.Now().UTC().AddDate(0, 0, -16)).
		Group("devices.user_id")

	if err := db.Pluck("devices.user_id", &filteredPatientList).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients for Billing: ", len(filteredPatientList))

	if len(filteredPatientList) == 0 {
		return nil
	}

	monthNumber := getMonthNumberFrom2023()
	for _, patientID := range filteredPatientList {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99453": monthNumber}); err != nil {
			return err
		}

		bill := models.Bill{
			PatientID:   patientID,
			ServiceCode: "RPM",
			CPTCode:     99453,
			EntryAt:     time.Now().UTC(),
		}

		if err := bill.CreateBill(); err != nil {
			log.Error(err)
		}
	}

	return nil
}
