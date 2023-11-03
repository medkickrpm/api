package organization

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

// upsertInteractionSetting godoc
// @Summary Upsert Interaction Setting
// @Description Upsert Interaction Setting
// @Tags Organization
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param upsert body InteractionSettingData true "Upsert Request"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organization/{id}/interaction-setting [put]
func upsertInteractionSetting(c echo.Context) error {
	req := struct {
		OrganizationID uint `json:"-" param:"id" validate:"required"`
		InteractionSettingData
	}{}

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

	interactionSetting := models.InteractionSetting{
		OrganizationID: req.OrganizationID,
		Type:           req.SettingType,
		Value:          req.Value,
	}

	if err := interactionSetting.UpsertInteractionSetting(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to upsert interaction setting",
		})
	}

	return c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "Interaction setting upsert successful",
	})
}

// getInteractionSetting godoc
// @Summary Get Interaction Setting
// @Description Get Interaction Setting
// @Tags Organization
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param filter query string true "Filter"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organization/{id}/interaction-setting [get]
func getInteractionSetting(c echo.Context) error {
	req := struct {
		OrganizationID uint   `json:"-" param:"id" validate:"required"`
		Filter         string `json:"-" query:"filter" validate:"required,oneof=ColorThreshold"`
	}{}

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

	interactionSetting := models.InteractionSetting{
		OrganizationID: req.OrganizationID,
		Type:           models.InteractionSettingType(req.Filter),
	}

	if err := interactionSetting.GetInteractionSetting(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get interaction setting",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"value": interactionSetting.Value,
	})
}
