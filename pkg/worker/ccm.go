package worker

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/database/models"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

func processCPTCode99490(patientIDs ...uint) error {
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
		Where("services.code = ?", "CCM").
		Where("patient_services.ended_at IS NULL").
		Joins("LEFT JOIN last_bill_entries ON last_bill_entries.patient_id = patient_services.patient_id").
		Where("last_bill_entries.c99490 < ? OR last_bill_entries.c99490 IS NULL", monthNumber)

	if len(patientIDs) > 0 {
		db = db.Where("patient_services.patient_id IN (?)", patientIDs)
	}

	if err := db.Pluck("patient_services.patient_id", &patientList).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients: (99490)", len(patientList))

	if len(patientList) == 0 {
		return nil
	}

	startDate := getStartDateOfMonth()

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

	fmt.Println("Total Patients for Billing: (99490)", len(filteredPatientList2))

	if len(filteredPatientList2) == 0 {
		return nil
	}

	var bills []models.Bill
	for _, patientID := range filteredPatientList2 {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99490": monthNumber}); err != nil {
			log.Error(err)
		}

		bills = append(bills, models.Bill{
			PatientID:   patientID,
			ServiceCode: "CCM",
			CPTCode:     99490,
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

func processCPTCode99439(patientIDs ...uint) error {
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
		Where("services.code = ?", "CCM").
		Where("patient_services.ended_at IS NULL").
		Joins("LEFT JOIN last_bill_entries ON last_bill_entries.patient_id = patient_services.patient_id").
		Where("last_bill_entries.c99439 < ? OR last_bill_entries.c99439 IS NULL", monthNumber)

	if len(patientIDs) > 0 {
		db = db.Where("patient_services.patient_id IN (?)", patientIDs)
	}

	if err := db.Pluck("patient_services.patient_id", &patientList).Error; err != nil {
		return err
	}

	fmt.Println("Total Patients: (99439)", len(patientList))

	if len(patientList) == 0 {
		return nil
	}

	startDate := getStartDateOfMonth()

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
		fmt.Println("Total Patients for Billing: (99439)", len(filteredPatientList2))
		return nil
	}

	filteredPatientList2Map := make(map[uint]struct{})
	for _, patientID := range filteredPatientList2 {
		filteredPatientList2Map[patientID] = struct{}{}
	}

	var patientBill []struct {
		PatientID uint `gorm:"column:patient_id"`
		Count     int  `gorm:"column:count"`
	}

	db = database.DB.Model(&models.Bill{}).
		Select("patient_id, COUNT(*) as count").
		Where("patient_id IN (?)", filteredPatientList2).
		Where("entry_at >= ?", startDate).
		Where("cpt_code = ?", 99439).
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
					ServiceCode: "CCM",
					CPTCode:     99439,
					EntryAt:     time.Now().UTC(),
				})
			}

			if requiredBillCount > actualBillCount {
				billedPatientCount++
			}
		}
	}

	if len(bills) > 0 {
		fmt.Println("Total Patients for Billing: (99439)", billedPatientCount)
		if err := database.DB.Create(&bills).Error; err != nil {
			return err
		}
	}

	for patientID := range billCompletedMap {
		if err := models.UpdateLastBillEntry(patientID, map[string]interface{}{"c99439": monthNumber}); err != nil {
			log.Error(err)
		}
	}

	return nil
}
