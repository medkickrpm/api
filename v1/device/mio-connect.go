package device

import (
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/validator"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
)

type MioData struct {
	DataType           string `json:"dataType" validate:"required"`
	IMEI               string `json:"imei" validate:"required"`
	SerialNumber       string `json:"sn""`
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
	Timestamp          uint   `json:"ts"`
	Timezone           string `json:"tz" validate:"required"`
	UID                string `json:"uid"`
	Weight             uint   `json:"wt"`
	WeightStableTime   uint   `json:"wet"`
	WeightLockCount    uint   `json:"lts"`
	UploadTime         uint   `json:"upload_time"`
	BloodGlucose       uint   `json:"data"`
	Unit               uint   `json:"unit"`
	TestPaperType      uint   `json:"sample"`
	SampleType         uint   `json:"sample_type"`
	Meal               uint   `json:"meal"`
	SignalLevel        uint   `json:"sig_lvl"`
	Uptime             uint   `json:"uptime"`
}

type MioStatus struct {
	DataType         string `json:"dataType" validate:"required"`
	IMEI             string `json:"imei" validate:"required"`
	Battery          uint   `json:"bat" validate:"required"`
	Timezone         string `json:"tz" validate:"required"`
	NetworkOperators string `json:"ops"`
	NetworkFormat    string `json:"net"`
	Signal           uint   `json:"sig"`
	SOCTemperature   uint   `json:"tp"`
	MeasureCount     uint   `json:"me_num"`
	AttachTime       uint   `json:"at_t"`
}

type Request struct {
	DeviceID    string    `json:"deviceId" validate:"required"`
	isTest      bool      `json:"isTest" validate:"required"`
	ModelNumber string    `json:"modelNumber" validate:"required"`
	Data        MioData   `json:"data"`
	Status      MioStatus `json:"status"`
	CreatedAt   string    `json:"createdAt" validate:"required"`
}

func IngestData(c echo.Context) error {
	//Verify API Key from header
	apiKey := c.Request().Header.Get("X-MIO-KEY")
	if apiKey != os.Getenv("MIO_API_KEY") {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := validator.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if req.Data.DataType == "bpm_gen2_measure" {

	} else if req.Data.DataType == "scale_gen2_measure" {

	} else if req.Data.DataType == "bgm_gen1_measure" {

	} else if req.Data.DataType == "bpm_gen2_status" {

	} else if req.Data.DataType == "scale_gen2_status" {

	} else if req.Data.DataType == "bgm_gen1_status" {

	} else if req.Data.DataType == "bpm_gen2_log" {

	} else if req.Data.DataType == "bgm_gen1_heartbeat" {

	} else {
		log.Warnf("Unknown data type: %s", req.Data.DataType)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: fmt.Sprintf("Unknown data type, %s", req.Data.DataType),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
