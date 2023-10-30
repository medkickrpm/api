package organization

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MeasurementRequest struct {
	MeasurementType models.MeasurementType `json:"measurement_type" validate:"required"`
	CriticalHigh    uint                   `json:"critical_high" validate:"required"`
	WarningHigh     *uint                  `json:"warning_high,omitempty"`
	WarningLow      *uint                  `json:"warning_low,omitempty"`
	CriticalLow     uint                   `json:"critical_low" validate:"required"`
}

type CreateAlertThresholdRequest struct {
	DeviceType   models.DeviceType    `json:"device_type" validate:"required"`
	Measurements []MeasurementRequest `json:"measurements" validate:"required,dive"`
}

func createAlertThreshold(c echo.Context) error {
	var req struct {
		OrganizationID uint `json:"-" param:"id" validate:"required"`
		CreateAlertThresholdRequest
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	if err := validator.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	self := middleware.GetSelf(c)
	if self.Role == "doctor" {
		req.OrganizationID = *self.OrganizationID
	}

	o := models.Organization{
		ID: req.OrganizationID,
	}

	if err := o.GetOrganization(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get organization",
		})
	}

	var alertThresholds []models.AlertThreshold

	for _, measurement := range req.Measurements {
		alertThresholds = append(alertThresholds, models.AlertThreshold{
			OrganizationID:  req.OrganizationID,
			DeviceType:      req.DeviceType,
			MeasurementType: measurement.MeasurementType,
			CriticalHigh:    measurement.CriticalHigh,
			WarningHigh:     measurement.WarningHigh,
			WarningLow:      measurement.WarningLow,
			CriticalLow:     measurement.CriticalLow,
		})
	}

	if err := models.CreateAlertThresholds(alertThresholds); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create alert thresholds",
		})
	}

	return c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "Alert thresholds created",
	})
}
