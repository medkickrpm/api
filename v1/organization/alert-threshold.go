package organization

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/validator"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MeasurementRequest struct {
	MeasurementType models.MeasurementType `json:"measurement_type" validate:"required,oneof=Systolic Diastolic Pulse Weight"`
	CriticalLow     *uint                  `json:"critical_low"`
	WarningLow      *uint                  `json:"warning_low"`
	WarningHigh     *uint                  `json:"warning_high"`
	CriticalHigh    *uint                  `json:"critical_high"`
}

func (m MeasurementRequest) validate() error {
	if m.CriticalLow == nil && m.WarningLow == nil && m.WarningHigh == nil && m.CriticalHigh == nil {
		return errors.New("at least one threshold must be set")
	}

	if m.CriticalLow != nil && m.CriticalHigh != nil && *m.CriticalLow > *m.CriticalHigh {
		return errors.New("critical low must be less than critical high")
	}

	if m.WarningLow != nil && m.WarningHigh != nil && *m.WarningLow > *m.WarningHigh {
		return errors.New("warning low must be less than warning high")
	}

	if m.CriticalLow != nil && m.WarningLow != nil && *m.CriticalLow > *m.WarningLow {
		return errors.New("critical low must be less than warning low")
	}

	if m.CriticalHigh != nil && m.WarningHigh != nil && *m.CriticalHigh < *m.WarningHigh {
		return errors.New("critical high must be greater than warning high")
	}

	return nil
}

type CreateAlertThresholdRequest struct {
	DeviceType   models.DeviceType    `json:"device_type" validate:"required,oneof=BloodPressure BloodGlucose WeightScale"`
	Measurements []MeasurementRequest `json:"measurements" validate:"required,min=1,dive,required"`
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
		if err := measurement.validate(); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: err.Error(),
			})
		}
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
