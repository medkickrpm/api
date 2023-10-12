package device

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/validator"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"time"
)

type MioData struct {
	DataType           string `json:"data_type" validate:"required"`
	IMEI               string `json:"imei" validate:"required"`
	SerialNumber       string `json:"sn"`
	Iccid              string `json:"iccid" validate:"required"`
	User               uint   `json:"user"`
	SystolicBP         uint   `json:"sys"`
	DiastolicBP        uint   `json:"dia"`
	Pulse              uint   `json:"pul"`
	IrregularHeartBeat bool   `json:"ihb"`
	HandShaking        bool   `json:"hand"`
	TripleMeasure      bool   `json:"tri"`
	Battery            uint   `json:"bat" validate:"required"`
	Signal             uint   `json:"sig" validate:"required"`
	Timestamp          int64  `json:"ts"`
	Timezone           string `json:"tz" validate:"required"`
	UID                string `json:"uid"`
	Weight             uint   `json:"wt"`
	WeightStableTime   uint   `json:"wet"`
	WeightLockCount    uint   `json:"lts"`
	UploadTime         int64  `json:"upload_time"`
	BloodGlucose       uint   `json:"data"`
	Unit               uint   `json:"unit"`
	TestPaperType      uint   `json:"sample"`
	SampleType         uint   `json:"sample_type"`
	Meal               uint   `json:"meal"`
	SignalLevel        uint   `json:"sig_lvl"`
	Uptime             int64  `json:"uptime"`
}

type MioStatus struct {
	DataType         string `json:"data_type" validate:"required"`
	IMEI             string `json:"imei" validate:"required"`
	Battery          uint   `json:"bat"`
	Timezone         string `json:"tz"`
	NetworkOperators string `json:"ops"`
	NetworkFormat    string `json:"net"`
	Signal           uint   `json:"sig"`
	SOCTemperature   int    `json:"tp"`
	MeasureCount     uint   `json:"me_num"`
	AttachTime       int64  `json:"at_t"`
}

type RequestTelemetry struct {
	DeviceID    string  `json:"deviceId" validate:"required"`
	IsTest      bool    `json:"isTest"`
	ModelNumber string  `json:"modelNumber" validate:"required"`
	Data        MioData `json:"data"`
	CreatedAt   uint    `json:"createdAt" validate:"required"`
}

type RequestStatus struct {
	DeviceID    string    `json:"deviceId" validate:"required"`
	IsTest      bool      `json:"isTest"`
	ModelNumber string    `json:"modelNumber" validate:"required"`
	Status      MioStatus `json:"status"`
	CreatedAt   uint      `json:"createdAt" validate:"required"`
}

// ingestTelemetry godoc
// @Summary Ingest Data
// @Description Mio Connect Data Ingestion Endpoint (Webhook)
// @Tags Mio (DO NOT USE)
// @Accept json
// @Produce json
// @Param create body Request true "Request"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /mio/forwardtelemetry [post]
func ingestTelemetry(c echo.Context) error {
	//Verify API Key from header
	apiKey := c.Request().Header.Get("X-MIO-KEY")
	if apiKey != os.Getenv("MIO_API_KEY") {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	var req RequestTelemetry
	if err := c.Bind(&req); err != nil {
		fmt.Println("1")
		fmt.Println(err.Error())
		return err
	}

	if err := validator.Validate.Struct(req); err != nil {
		fmt.Println("2")
		println(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if req.IsTest {
		fmt.Printf("Test data received: %+v\n", req)
		return c.NoContent(http.StatusNoContent)
	}

	device := &models.Device{
		IMEI: req.Data.IMEI,
	}
	if err := device.GetDeviceByIMEI(); err != nil {
		log.Errorf("Failed to get device by IMEI: %s", err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get device by IMEI",
		})
	}

	if device.Name == "" {
		switch req.Data.DataType {
		case "bpm_gen2_measure":
			device.Name = "Sphygmomanometer"
		case "scale_gen2_measure":
			device.Name = "Weight Scale"
		case "bgm_gen1_measure":
			device.Name = "Blood Glucose Meter"
		}

		if err := device.UpdateDevice(); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to update device",
			})
		}
	}

	if err := device.UpdateBattery(req.Data.Battery); err != nil {
		log.Errorf("Failed to update device battery: %s", err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update device battery",
		})
	}

	if req.Data.DataType == "bpm_gen2_measure" {
		t := time.Unix(req.Data.Timestamp, 0)
		dtd := &models.DeviceTelemetryData{
			SystolicBP:         req.Data.SystolicBP,
			DiastolicBP:        req.Data.DiastolicBP,
			Pulse:              req.Data.Pulse,
			IrregularHeartBeat: req.Data.IrregularHeartBeat,
			HandShaking:        req.Data.HandShaking,
			TripleMeasurement:  req.Data.TripleMeasure,
			DeviceID:           device.ID,
			MeasuredAt:         t,
		}
		if err := dtd.CreateDeviceTelemetryData(); err != nil {
			log.Errorf("Failed to create device telemetry data: %s", err)
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to create device telemetry data",
			})
		}
		return c.NoContent(http.StatusNoContent)
	}
	if req.Data.DataType == "scale_gen2_measure" {
		t := time.Unix(req.Data.UploadTime, 0)
		dtd := &models.DeviceTelemetryData{
			Weight:           req.Data.Weight,
			WeightStableTime: req.Data.WeightStableTime,
			WeightLockCount:  req.Data.WeightLockCount,
			DeviceID:         device.ID,
			MeasuredAt:       t,
		}
		if err := dtd.CreateDeviceTelemetryData(); err != nil {
			log.Errorf("Failed to create device telemetry data: %s", err)
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to create device telemetry data",
			})
		}
		return c.NoContent(http.StatusNoContent)
	}
	if req.Data.DataType == "bgm_gen1_measure" {
		t := time.Unix(req.Data.Uptime, 0)
		unit := ""
		if req.Data.Unit == 1 {
			unit = "mmol/L"
		} else if req.Data.Unit == 2 {
			unit = "mg/dL"
		} else {
			unit = "Unknown"
		}

		testPaper := ""
		if req.Data.TestPaperType == 1 {
			testPaper = "GOD"
		} else if req.Data.TestPaperType == 2 {
			testPaper = "GDH"
		} else {
			testPaper = "Unknown"
		}

		sampleType := ""
		if req.Data.SampleType == 1 {
			sampleType = "blood or resistance"
		} else if req.Data.SampleType == 2 {
			sampleType = "quality control liquid"
		} else {
			sampleType = "sample is invalid"
		}

		meal := ""
		if req.Data.Meal == 1 {
			meal = "before meal"
		} else if req.Data.Meal == 2 {
			meal = "after meal"
		} else {
			meal = "Unknown"
		}

		dtd := &models.DeviceTelemetryData{
			BloodGlucose: req.Data.BloodGlucose,
			Unit:         unit,
			TestPaper:    testPaper,
			SampleType:   sampleType,
			Meal:         meal,
			DeviceID:     device.ID,
			MeasuredAt:   t,
		}

		if err := dtd.CreateDeviceTelemetryData(); err != nil {
			log.Errorf("Failed to create device telemetry data: %s", err)
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to create device telemetry data",
			})
		}
		return c.NoContent(http.StatusNoContent)
	}

	log.Warnf("Unknown data type: %s", req.Data.DataType)
	return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Error: fmt.Sprintf("Unknown data type, %s", req.Data.DataType),
	})
}

// ingestStatus godoc
// @Summary Ingest Status
// @Description Mio Connect Status Ingestion Endpoint (Webhook)
// @Tags Mio (DO NOT USE)
// @Accept json
// @Produce json
// @Param create body Request true "Request"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /mio/forwardstatus [post]
func ingestStatus(c echo.Context) error {
	//Verify API Key from header
	apiKey := c.Request().Header.Get("X-MIO-KEY")
	if apiKey != os.Getenv("MIO_API_KEY") {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	var req RequestStatus
	if err := c.Bind(&req); err != nil {
		fmt.Println("1")
		fmt.Println(err.Error())
		return err
	}

	if err := validator.Validate.Struct(req); err != nil {
		fmt.Println("2")
		println(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if req.IsTest {
		fmt.Printf("Test data received: %+v\n", req)
		return c.NoContent(http.StatusNoContent)
	}

	device := &models.Device{
		IMEI: req.Status.IMEI,
	}
	if err := device.GetDeviceByIMEI(); err != nil {
		log.Errorf("Failed to get device by IMEI: %s", err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get device by IMEI",
		})
	}

	if device.Name == "" {
		switch req.Status.DataType {
		case "bpm_gen2_status":
			device.Name = "Sphygmomanometer"
		case "scale_gen2_status":
			device.Name = "Weight Scale"
		case "bgm_gen1_status":
			device.Name = "Blood Glucose Meter"
		}

		if err := device.UpdateDevice(); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to update device",
			})
		}
	}

	if req.Status.DataType == "bpm_gen2_status" {
		t := time.Unix(req.Status.AttachTime, 0)
		dsd := &models.DeviceStatusData{
			Timezone:      req.Status.Timezone,
			NetworkOps:    req.Status.NetworkOperators,
			NetworkFormat: req.Status.NetworkFormat,
			Signal:        req.Status.Signal,
			Temperature:   req.Status.SOCTemperature,
			AttachTime:    t,
		}

		if err := dsd.CreateDeviceStatusData(); err != nil {
			log.Errorf("Failed to create device status data: %s", err)
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to create device status data",
			})
		}

		return c.NoContent(http.StatusNoContent)
	}
	if req.Status.DataType == "scale_gen2_status" {
		t := time.Unix(req.Status.AttachTime, 0)
		dsd := &models.DeviceStatusData{
			Timezone:      req.Status.Timezone,
			NetworkOps:    req.Status.NetworkOperators,
			NetworkFormat: req.Status.NetworkFormat,
			Signal:        req.Status.Signal,
			Temperature:   req.Status.SOCTemperature,
			AttachTime:    t,
		}

		if err := dsd.CreateDeviceStatusData(); err != nil {
			log.Errorf("Failed to create device status data: %s", err)
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to create device status data",
			})
		}

		return c.NoContent(http.StatusNoContent)
	}
	if req.Status.DataType == "bgm_gen1_status" {
		t := time.Unix(req.Status.AttachTime, 0)
		dsd := &models.DeviceStatusData{
			Timezone:      req.Status.Timezone,
			NetworkOps:    req.Status.NetworkOperators,
			NetworkFormat: req.Status.NetworkFormat,
			Signal:        req.Status.Signal,
			Temperature:   req.Status.SOCTemperature,
			AttachTime:    t,
		}

		if err := dsd.CreateDeviceStatusData(); err != nil {
			log.Errorf("Failed to create device status data: %s", err)
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to create device status data",
			})
		}

		return c.NoContent(http.StatusNoContent)
	}

	log.Warnf("Unknown data type: %s", req.Status.DataType)
	return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Error: fmt.Sprintf("Unknown data type, %s", req.Status.DataType),
	})
}
