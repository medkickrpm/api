package cron

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/validator"
	"MedKick-backend/pkg/worker"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// triggerCptWorker godoc
// @Summary Trigger CPT Worker
// @Description CRON ONLY - Triggers CPT Worker
// @Tags CRON
// @Accept json
// @Produce json
// @Param CronToken body Request true "Token Request"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cron/trigger-cpt-worker [post]
func triggerCptWorker(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to bind request",
		})
	}

	if err := validator.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request",
		})
	}

	if req.Token != os.Getenv("CRON_SECRET") {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Invalid token",
		})
	}

	worker.TriggerCPTWorker()

	return c.NoContent(http.StatusNoContent)
}

// clearTestBillings godoc
// @Summary Clear Test Billings
// @Description CRON ONLY - Clears all test billings
// @Tags CRON
// @Accept json
// @Produce json
// @Param CronToken body Request true "Token Request"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cron/clear-test-billings [post]
func clearTestBillings(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to bind request",
		})
	}

	if err := validator.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request",
		})
	}

	if req.Token != os.Getenv("CRON_SECRET") {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Invalid token",
		})
	}

	testPatient := uint(11)
	if err := models.DeleteLastBillEntry(testPatient); err != nil {
		log.Error(err)
	}

	year, month, _ := time.Now().UTC().Date()
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	if err := models.DeleteBillByPatientIDInRange(testPatient, startDate, endDate); err != nil {
		log.Error(err)
	}

	return c.NoContent(http.StatusNoContent)
}
