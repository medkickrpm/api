package organization

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

// listTelemetryAlert godoc
// @Summary List Telemetry Alert
// @Description List Telemetry Alert
// @Tags Organization
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param status query string false "Status"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param sort_by query string false "Sort By"
// @Param sort_direction query string false "Sort Direction"
// @Success 200 {object} TelemetryAlertResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organization/{id}/telemetry-alert [get]
func listTelemetryAlert(c echo.Context) error {
	param := struct {
		OrganizationID uint   `param:"id"`
		Status         string `query:"status"`
		models.PageReq
		models.SortReq
	}{
		PageReq: models.NewPageReq(),
		SortReq: models.NewSortReq(),
	}

	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	self := middleware.GetSelf(c)
	if self.Role == "doctor" || self.Role == "nurse" {
		param.OrganizationID = *self.OrganizationID
	}

	o := models.Organization{
		ID: param.OrganizationID,
	}

	if err := o.GetOrganization(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get organization",
		})
	}

	isActive := true

	if param.Status == "inactive" {
		isActive = false
	}

	data, err := models.ListTelemetryAlerts(param.OrganizationID, isActive, param.PageReq, param.SortReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get telemetry alert",
		})
	}

	total, err := models.CountTelemetryAlerts(param.OrganizationID, isActive)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get telemetry alert",
		})
	}

	response := map[string]interface{}{
		"data":  convertModelToResponse(data),
		"page":  param.Page,
		"size":  len(data),
		"total": total,
	}

	return c.JSON(http.StatusOK, response)
}

// resolveTelemetryAlert godoc
// @Summary Resolve Telemetry Alert
// @Description Resolve Telemetry Alert
// @Tags Organization
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param alert path int true "Alert ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organization/{id}/telemetry-alert/{alert}/resolve [patch]
func resolvedTelemetryAlert(c echo.Context) error {
	var param struct {
		OrganizationID uint `param:"id"`
		AlertID        uint `param:"alert"`
	}

	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	self := middleware.GetSelf(c)
	if self.Role == "doctor" || self.Role == "nurse" {
		param.OrganizationID = *self.OrganizationID
	}

	o := models.Organization{
		ID: param.OrganizationID,
	}

	if err := o.GetOrganization(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get organization",
		})
	}

	t := models.TelemetryAlert{
		ID:             param.AlertID,
		IsActive:       true,
		IsAutoResolved: false,
	}

	if err := t.GetTelemetryAlert(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get telemetry alert",
		})
	}

	t.IsActive = false
	t.ResolvedByID = self.ID

	if err := t.ResolveTelemetryAlert(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update telemetry alert",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully resolved telemetry alert",
	})
}
