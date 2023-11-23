package worker

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/database/models"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

func processCPTCode99484(patientIDs ...uint) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic: ", err)
		}
	}()

	monthNumber := getMonthNumberFrom2023()

	var patientList []uint

	db := database.DB.Model(&models.PatientService{}).
		Joins("JOIN services ON services.id = patient_services.service_id").
		Where("services.is_enabled = ?", true).
		Where("services.code = ?", "BHI").
		Where("patient_services.ended_at IS NULL").
		Joins("LEFT JOIN last_bill_entries ON last_bill_entries.patient_id = patient_services.patient_id").
		Where("last_bill_entries.c99484 < ? OR last_bill_entries.c99484 IS NULL", monthNumber)

	if len(patientIDs) > 0 {
		db = db.Where("patient_services.patient_id IN (?)", patientIDs)
	}

	if err := db.Pluck("patient_services.patient_id", &patientList).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients: (99484)", len(patientList))

	if len(patientList) == 0 {
		return nil
	}

	startDate := getStartDateOfMonth()

	var filteredPatientList []uint

	db = database.DB.Model(&models.Interaction{}).
		Where("user_id IN (?)", patientList).
		Where("session_date >= ?", startDate).
		Where("cost_category = ?", "BHI").
		Group("user_id").
		Having("SUM(duration) >= ?", 20*60)

	if err := db.Pluck("user_id", &filteredPatientList).Error; err != nil {
		return err
	}

	if len(filteredPatientList) == 0 {
		return nil
	}

	var filteredPatientList2 []uint
	db = database.DB.Model(&models.Device{}).
		Joins("JOIN device_telemetry_data ON device_telemetry_data.device_id = devices.id").
		Where("devices.user_id IN (?)", filteredPatientList).
		Where("device_telemetry_data.measured_at >= ?", startDate).
		Group("devices.user_id")

	if err := db.Pluck("devices.user_id", &filteredPatientList2).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients for Billing: (99484)", len(filteredPatientList2))

	if len(filteredPatientList2) == 0 {
		return nil
	}

	var bills []models.Bill
	for _, patientID := range filteredPatientList2 {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99484": monthNumber}); err != nil {
			log.Error(err)
		}

		bills = append(bills, models.Bill{
			PatientID:   patientID,
			ServiceCode: "BHI",
			CPTCode:     99484,
			EntryAt:     time.Now().UTC(),
		})
	}

	if len(bills) > 0 {
		if err := models.CreateBill(bills); err != nil {
			return err
		}
	}

	return nil
}
