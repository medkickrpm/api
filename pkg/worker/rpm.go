package worker

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/database/models"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

func processCPTCode99453(...uint) error {
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
	db = database.DB.Model(&models.Device{}).
		Joins("JOIN device_telemetry_data ON device_telemetry_data.device_id = devices.id").
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

	var bills []models.Bill
	monthNumber := getMonthNumberFrom2023()
	for _, patientID := range filteredPatientList {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99453": monthNumber}); err != nil {
			log.Error(err)
		}

		bills = append(bills, models.Bill{
			PatientID:   patientID,
			ServiceCode: "RPM",
			CPTCode:     99453,
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

func processCPTCode99454(...uint) error {
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
	db = database.DB.Model(&models.Device{}).
		Joins("JOIN device_telemetry_data ON device_telemetry_data.device_id = devices.id").
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

	db = database.DB.Model(&models.Device{}).
		Joins("JOIN device_telemetry_data ON device_telemetry_data.device_id = devices.id").
		Where("devices.user_id IN (?)", filteredPatientList).
		Where("device_telemetry_data.measured_at >= ?", startDate).
		Group("devices.user_id")

	if err := db.Pluck("devices.user_id", &filteredPatientList).Error; err != nil {
		return err
	}

	var bills []models.Bill
	for _, patientID := range filteredPatientList {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99454": monthNumber}); err != nil {
			log.Error(err)
		}

		bills = append(bills, models.Bill{
			PatientID:   patientID,
			ServiceCode: "RPM",
			CPTCode:     99454,
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

	fmt.Println("Total Patients: (99457)", len(patientList))

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
	db = database.DB.Model(&models.Device{}).
		Joins("JOIN device_telemetry_data ON device_telemetry_data.device_id = devices.id").
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

	var bills []models.Bill
	for _, patientID := range filteredPatientList2 {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99457": monthNumber}); err != nil {
			log.Error(err)
		}

		bills = append(bills, models.Bill{
			PatientID:   patientID,
			ServiceCode: "RPM",
			CPTCode:     99457,
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

func processCPTCode99458(patientIDs ...uint) error {
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
		Where("last_bill_entries.c99458 < ? OR last_bill_entries.c99458 IS NULL", monthNumber)

	if len(patientIDs) > 0 {
		db = db.Where("patient_services.patient_id IN (?)", patientIDs)
	}

	if err := db.Pluck("patient_services.patient_id", &patientList).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients: (99458)", len(patientList))

	if len(patientList) == 0 {
		return nil
	}

	year, month, _ := time.Now().UTC().Date()
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	var filteredPatientWithInteraction []struct {
		UserID   uint `gorm:"column:user_id"`
		Duration uint `gorm:"column:duration"`
	}

	db = database.DB.Model(&models.Interaction{}).
		Select("user_id, SUM(duration) as duration").
		Where("user_id IN (?)", patientList).
		Where("session_date >= ?", startDate).
		Group("user_id").
		Having("SUM(duration) >= ?", 40*60)

	if err := db.Find(&filteredPatientWithInteraction).Error; err != nil {
		return err
	}

	if len(filteredPatientWithInteraction) == 0 {
		return nil
	}

	var filteredPatientList []uint
	for _, patient := range filteredPatientWithInteraction {
		filteredPatientList = append(filteredPatientList, patient.UserID)
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

	if len(filteredPatientList2) == 0 {
		fmt.Println("Total Patients for Billing: (99458)", len(filteredPatientList2))
		return nil
	}

	filteredPatientList2Map := make(map[uint]struct{})
	for _, patientID := range filteredPatientList2 {
		filteredPatientList2Map[patientID] = struct{}{}
	}

	// number of bills for cpt code 99458 for each patient from filteredPatientList2 from date startDate
	var patientBill []struct {
		PatientID uint `gorm:"column:patient_id"`
		Count     int  `gorm:"column:count"`
	}

	db = database.DB.Model(&models.Bill{}).
		Select("patient_id, COUNT(*) as count").
		Where("patient_id IN (?)", filteredPatientList2).
		Where("entry_at >= ?", startDate).
		Where("cpt_code = ?", 99458).
		Group("patient_id")

	if err := db.Find(&patientBill).Error; err != nil {
		return err
	}

	var patientBillMap = make(map[uint]int)
	for _, bill := range patientBill {
		patientBillMap[bill.PatientID] = bill.Count
	}

	var bills []models.Bill
	var billCompletedMap = make(map[uint]struct{})
	var billedPatientCount int

	for _, patient := range filteredPatientWithInteraction {
		if _, ok := filteredPatientList2Map[patient.UserID]; ok {
			requiredBillCount := 1
			if patient.Duration >= 60*60 {
				requiredBillCount = 2
				billCompletedMap[patient.UserID] = struct{}{}
			}

			actualBillCount := patientBillMap[patient.UserID]

			for i := 0; i < requiredBillCount-actualBillCount; i++ {
				bills = append(bills, models.Bill{
					PatientID:   patient.UserID,
					ServiceCode: "RPM",
					CPTCode:     99458,
					EntryAt:     time.Now().UTC(),
				})
			}

			if requiredBillCount > actualBillCount {
				billedPatientCount++
			}
		}
	}

	if len(bills) > 0 {
		fmt.Println("Total Patients for Billing: (99458)", billedPatientCount)
		if err := database.DB.Create(&bills).Error; err != nil {
			return err
		}
	}

	for patientID := range billCompletedMap {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99458": monthNumber}); err != nil {
			log.Error(err)
		}
	}

	return nil
}
