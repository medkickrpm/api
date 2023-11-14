package worker

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/database/models"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/go-co-op/gocron"
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

	fmt.Println("Total Patients for Billing: (9945)", len(filteredPatientList))

	if len(filteredPatientList) == 0 {
		return nil
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

func RunCPTWorker() {
	s := gocron.NewScheduler(time.UTC)

	job99453, _ := s.Tag("99453").Every(1).Day().At("07:00").Do(func() {
		if err := processCPTCode99453(); err != nil {
			fmt.Println(err)
		}
	})

	// run every day after 15th of each month
	job99454, _ := s.Tag("99454").Cron("0 7 16-31 * *").Do(func() {
		if err := processCPTCode99454(); err != nil {
			fmt.Println(err)
		}
	})

	s.RunAllWithDelay(time.Second * 5)

	_, _ = s.Tag("Worker").Every(6).Hour().Do(func() {
		fmt.Println("Run CPT Worker At:")
		fmt.Println("	99453: ", job99453.NextRun())
		fmt.Println("	99454: ", job99454.NextRun())
	})

	s.StartBlocking()
}
