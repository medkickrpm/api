package worker

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/database/models"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

func processCPTCode99453() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic: ", err)
		}
	}()

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

	fmt.Println("Total Patients (99453): ", len(patientList))

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

	fmt.Println("Total Patients for Billing (99453): ", len(filteredPatientList))

	if len(filteredPatientList) == 0 {
		return nil
	}

	monthNumber := getMonthNumberFrom2023()
	for _, patientID := range filteredPatientList {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99453": monthNumber}); err != nil {
			log.Error(err)
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

func processCPTCode99454() error {
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
		Where("services.code = ?", "RPM").
		Where("patient_services.ended_at IS NULL").
		Joins("LEFT JOIN last_bill_entries ON last_bill_entries.patient_id = patient_services.patient_id").
		Where("last_bill_entries.c99454 < ? OR last_bill_entries.c99454 IS NULL", monthNumber)

	if err := db.Pluck("patient_services.patient_id", &patientList).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients: (99454)", len(patientList))

	if len(patientList) == 0 {
		return nil
	}

	year, month, _ := time.Now().UTC().Date()
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	var filteredPatientList []uint
	db = database.DB.Model(&models.Device{}).Joins("JOIN device_telemetry_data ON device_telemetry_data.device_id = devices.id").
		Where("devices.user_id IN (?)", patientList).
		Where("device_telemetry_data.measured_at >= ?", startDate).
		Group("devices.user_id").
		Having("COUNT(DISTINCT DATE(device_telemetry_data.measured_at)) >= 16")

	if err := db.Pluck("devices.user_id", &filteredPatientList).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients for Billing: (99454)", len(filteredPatientList))

	if len(filteredPatientList) == 0 {
		return nil
	}

	db = database.DB.Model(&models.Device{}).Joins("JOIN device_telemetry_data ON device_telemetry_data.device_id = devices.id").
		Where("devices.user_id IN (?)", filteredPatientList).
		Where("device_telemetry_data.measured_at >= ?", startDate).
		Group("devices.user_id")

	if err := db.Pluck("devices.user_id", &filteredPatientList).Error; err != nil {
		return err
	}

	for _, patientID := range filteredPatientList {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99454": monthNumber}); err != nil {
			log.Error(err)
		}

		bill := models.Bill{
			PatientID:   patientID,
			ServiceCode: "RPM",
			CPTCode:     99454,
			EntryAt:     time.Now().UTC(),
		}

		if err := bill.CreateBill(); err != nil {
			log.Error(err)
		}
	}

	return nil
}

func processCPTCode99457(patientIDs ...uint) error {
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
		Where("services.code = ?", "RPM").
		Where("patient_services.ended_at IS NULL").
		Joins("LEFT JOIN last_bill_entries ON last_bill_entries.patient_id = patient_services.patient_id").
		Where("last_bill_entries.c99457 < ? OR last_bill_entries.c99457 IS NULL", monthNumber)

	if len(patientIDs) > 0 {
		db = db.Where("patient_services.patient_id IN (?)", patientIDs)
	}

	if err := db.Pluck("patient_services.patient_id", &patientList).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients: (99454)", len(patientList))

	if len(patientList) == 0 {
		return nil
	}

	year, month, _ := time.Now().UTC().Date()
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	var filteredPatientList []uint

	db = database.DB.Model(&models.Interaction{}).
		Where("user_id IN (?)", patientList).
		Where("session_date >= ?", startDate).
		Group("user_id").
		Having("SUM(duration) >= ?", 20*60)

	if err := db.Pluck("user_id", &filteredPatientList).Error; err != nil {
		return err
	}

	if len(filteredPatientList) == 0 {
		return nil
	}

	var filteredPatientList2 []uint
	db = database.DB.Model(&models.Device{}).Joins("JOIN device_telemetry_data ON device_telemetry_data.device_id = devices.id").
		Where("devices.user_id IN (?)", filteredPatientList).
		Where("device_telemetry_data.measured_at >= ?", startDate).
		Group("devices.user_id")

	if err := db.Pluck("devices.user_id", &filteredPatientList2).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients for Billing: (99457)", len(filteredPatientList2))

	if len(filteredPatientList2) == 0 {
		return nil
	}

	for _, patientID := range filteredPatientList2 {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99457": monthNumber}); err != nil {
			log.Error(err)
		}

		bill := models.Bill{
			PatientID:   patientID,
			ServiceCode: "RPM",
			CPTCode:     99457,
			EntryAt:     time.Now().UTC(),
		}

		if err := bill.CreateBill(); err != nil {
			log.Error(err)
		}
	}

	return nil
}
