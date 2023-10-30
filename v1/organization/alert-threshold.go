package organization

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

func upsertAlertThreshold(c echo.Context) error {
	var req struct {
		OrganizationID uint `json:"-" param:"id" validate:"required"`
		AlertThresholdData
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

	if err := models.UpsertAlertThresholds(alertThresholds); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to upsert alert thresholds",
		})
	}

	return c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "Alert thresholds upsert successful",
	})
}

func listAlertThresholds(c echo.Context) error {
	var req struct {
		OrganizationID uint `json:"-" param:"id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	self := middleware.GetSelf(c)
	if self.Role == "doctor" || self.Role == "nurse" {
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

	alertThresholds, err := models.ListAlertThresholds(req.OrganizationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get alert thresholds",
		})
	}

	return c.JSON(http.StatusOK, convertAlertThresholdModelToResponse(alertThresholds))
}
