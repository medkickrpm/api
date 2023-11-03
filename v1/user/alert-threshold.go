package user

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

// upsertAlertThreshold godoc
// @Summary Upsert Alert Threshold
// @Description Upsert Alert Threshold
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param upsert body AlertThresholdData true "Upsert Request"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/alert-threshold [put]
func upsertAlertThreshold(c echo.Context) error {
	var req struct {
		PatientID uint `json:"-" param:"id" validate:"required"`
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

	u := models.User{
		ID: &req.PatientID,
	}

	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get user",
		})
	}

	if u.Role != "patient" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "User is not a patient",
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
			PatientID:       req.PatientID,
			DeviceType:      req.DeviceType,
			MeasurementType: measurement.MeasurementType,
			CriticalHigh:    measurement.CriticalHigh,
			WarningHigh:     measurement.WarningHigh,
			WarningLow:      measurement.WarningLow,
			CriticalLow:     measurement.CriticalLow,
			Note:            req.Note,
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

// listAlertThresholds godoc
// @Summary List Alert Thresholds
// @Description List Alert Thresholds
// @Tags User
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} []AlertThresholdData
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/alert-threshold [get]
func listAlertThresholds(c echo.Context) error {
	var req struct {
		PatientID uint `json:"-" param:"id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	u := models.User{
		ID: &req.PatientID,
	}

	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get user",
		})
	}

	if u.Role != "patient" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "User is not a patient",
		})
	}

	alertThresholds, err := models.ListAlertThresholds([]uint{req.PatientID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get alert thresholds",
		})
	}

	return c.JSON(http.StatusOK, convertAlertThresholdModelToResponse(alertThresholds))
}
