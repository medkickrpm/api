package device

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// getStatus godoc
// @Summary Get Status Data
// @Description Get Status Data for given ID
// @Tags Mio
// @Accept json
// @Produce json
// @Param id path string false "Device ID"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Success 200 {object} []models.DeviceStatusData
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /mio/status/{id} [get]
func getStatus(c echo.Context) error {
	self := middleware.GetSelf(c)

	id := c.Param("id")

	// Convert id to uint
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to convert id to uint",
		})
	}

	device := &models.Device{
		ID: uint(idInt),
	}
	if err := device.GetDevice(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get device by device id",
		})
	}

	if self.Role == "patient" && device.UserID != *self.ID {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error: "Forbidden",
		})
	}

	startDateRaw := c.QueryParam("start_date")
	endDateRaw := c.QueryParam("end_date")

	//convert start_date and end_date to time.Time
	startDate, err := time.Parse("2006-01-02", startDateRaw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to parse start_date",
		})
	}

	var endDate time.Time

	if endDateRaw == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", endDateRaw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Failed to parse end_date",
			})
		}
	}

	// Make sure startDate is before endDate
	if startDate.After(endDate) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Start date must be before end date",
		})
	}

	// Make sure startDate is before present day
	if startDate.After(time.Now()) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Start date must be before present day",
		})
	}

	statuses, err := models.GetDeviceStatusDataByDeviceBetweenDates(device.ID, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get device telemetry data by device id",
		})
	}

	return c.JSON(http.StatusOK, statuses)
}
