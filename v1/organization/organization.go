package organization

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Address2 string `json:"address2"`
	City     string `json:"city" validate:"required"`
	State    string `json:"state" validate:"required"`
	Zip      string `json:"zip" validate:"required"`
	Country  string `json:"country" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

func createOrganization(c echo.Context) error {
	var req CreateRequest
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

	o := models.Organization{
		Name:     req.Name,
		Address:  req.Address,
		Address2: req.Address2,
		City:     req.City,
		State:    req.State,
		Zip:      req.Zip,
		Country:  req.Country,
		Phone:    req.Phone,
	}

	if err := o.CreateOrganization(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create organization",
		})
	}

	return c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "Organization created successfully",
	})
}

func getOrganization(c echo.Context) error {
	self := middleware.GetSelf(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	idUint := uint(id)

	if self.Role == "patient" || self.Role == "doctor" {
		idUint = *self.OrganizationID
	}

	o := models.Organization{
		ID: idUint,
	}

	if err := o.GetOrganization(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get organization",
		})
	}

	return c.JSON(http.StatusOK, o)
}

type UpdateRequest struct {
	Name     string `json:"name" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Address2 string `json:"address2"`
	City     string `json:"city" validate:"required"`
	State    string `json:"state" validate:"required"`
	Zip      string `json:"zip" validate:"required"`
	Country  string `json:"country" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

func updateOrganization(c echo.Context) error {
	var req UpdateRequest
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

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	idUint := uint(id)

	if self.Role == "doctor" {
		idUint = *self.OrganizationID
	}

	o := models.Organization{
		ID: idUint,
	}

	if err := o.GetOrganization(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get organization",
		})
	}

	if req.Name != "" {
		o.Name = req.Name
	}
	if req.Address != "" {
		o.Address = req.Address
	}
	if req.Address2 != "" {
		o.Address2 = req.Address2
	}
	if req.City != "" {
		o.City = req.City
	}
	if req.State != "" {
		o.State = req.State
	}
	if req.Zip != "" {
		o.Zip = req.Zip
	}
	if req.Country != "" {
		o.Country = req.Country
	}
	if req.Phone != "" {
		o.Phone = req.Phone
	}

	if err := o.UpdateOrganization(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update organization",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Organization updated successfully",
	})
}

func deleteOrganization(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	o := models.Organization{
		ID: uint(id),
	}

	if err := o.DeleteOrganization(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to delete organization",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Organization deleted successfully",
	})
}
