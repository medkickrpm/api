package user

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/validator"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// upsertPatientServices godoc
// @Summary Upsert Patient Services
// @Description Upsert Patient Services
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param upsert body PatientServiceData true "Upsert Request"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/patient-service [put]
func upsertPatientServices(c echo.Context) error {
	var req struct {
		PatientID uint `json:"-" param:"id"`
		PatientServiceData
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

	// return active services by the patient
	patientServices, err := models.ListPatientServices(req.PatientID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list patient services",
		})
	}

	// return all available active services
	allServices, err := models.ListServices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list services",
		})
	}

	// all services in a map
	serviceMap := make(map[string]models.Service)
	for _, service := range allServices {
		serviceMap[service.Code] = service
	}

	// new requested services in a map
	newServiceMap := make(map[string]bool)
	for _, service := range req.Services {
		newServiceMap[service] = true
	}

	if newServiceMap["CCM"] == newServiceMap["PCM"] {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Cannot have both CCM and PCM",
		})
	}

	// old services in a map
	oldServiceMap := make(map[string]struct{})

	currentTime := time.Now().UTC()

	var toUpsert []models.PatientService

	// if a service is not in the new requested services, set the ended_at to current time meaning ended
	for _, patientService := range patientServices {
		oldServiceMap[patientService.Service.Code] = struct{}{}
		if _, ok := newServiceMap[patientService.Service.Code]; !ok {
			patientService.EndedAt = &currentTime
			patientService.UpdatedAt = currentTime
			toUpsert = append(toUpsert, patientService)
		}
	}

	// if a service is not in the old services, set the started_at to current time meaning started
	for _, service := range req.Services {
		if _, ok := oldServiceMap[service]; !ok {
			if svc, found := serviceMap[service]; found {
				toUpsert = append(toUpsert, models.PatientService{
					PatientID: req.PatientID,
					ServiceID: svc.ID,
					StartedAt: currentTime,
				})
			}
		}
	}

	if len(toUpsert) > 0 {
		if err = models.UpsertPatientServices(toUpsert); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to upsert patient services",
			})
		}
	}

	return c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "Successfully upserted patient services",
	})
}

// listPatientServices godoc
// @Summary List Patient Services
// @Description List Patient Services
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} PatientServiceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/patient-service [get]
func listPatientServices(c echo.Context) error {
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

	patientServices, err := models.ListPatientServices(req.PatientID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list patient services",
		})
	}

	return c.JSON(http.StatusOK, convertPatientServiceModelToResponse(patientServices))
}
