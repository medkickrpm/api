package user

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

// upsertDiagnoses godoc
// @Summary Upsert Diagnoses
// @Description Upsert Diagnoses
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param upsert body DiagnosisData true "Upsert Request"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/diagnoses [put]
func upsertDiagnoses(c echo.Context) error {
	var req struct {
		PatientID uint `json:"-" param:"id" validate:"required"`
		DiagnosisData
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

	diagnosisCodes, err := models.ListDiagnosisCodes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list ICT10 codes",
		})
	}

	allDiagnosesMap := make(map[string]models.Diagnosis)
	for _, d := range diagnosisCodes {
		allDiagnosesMap[d.Code] = d
	}

	patientDiagnoses, err := models.GetPatientDiagnoses(req.PatientID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get patient diagnoses",
		})
	}
	existingDiagnosisMap := make(map[string]struct{})
	for _, d := range patientDiagnoses {
		existingDiagnosisMap[d.Diagnosis.Code] = struct{}{}
	}

	toInsertDiagnoses := make([]models.PatientDiagnosis, 0)
	toDeleteDiagnoses := make([]uint, 0)

	newDiagnosesMap := make(map[string]struct{})

	for _, diagnosis := range req.Diagnoses {
		newDiagnosesMap[diagnosis] = struct{}{}

		if _, ok := allDiagnosesMap[diagnosis]; !ok {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Invalid Diagnosis code",
			})
		}

		if _, ok := existingDiagnosisMap[diagnosis]; !ok {
			toInsertDiagnoses = append(toInsertDiagnoses, models.PatientDiagnosis{
				UserID:      req.PatientID,
				DiagnosisID: allDiagnosesMap[diagnosis].ID,
			})
		}
	}

	for _, diagnosis := range patientDiagnoses {
		if _, ok := newDiagnosesMap[diagnosis.Diagnosis.Code]; !ok {
			toDeleteDiagnoses = append(toDeleteDiagnoses, diagnosis.DiagnosisID)
		}
	}

	if len(toDeleteDiagnoses) > 0 {
		if err = models.DeletePatientDiagnoses(req.PatientID, toDeleteDiagnoses); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to delete patient diagnoses",
			})
		}
	}

	if len(toInsertDiagnoses) > 0 {
		if err = models.CreatePatientDiagnoses(toInsertDiagnoses); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to create patient diagnoses",
			})
		}
	}

	return c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "Diagnoses upsert successful",
	})
}

// getDiagnoses godoc
// @Summary Get Diagnoses
// @Description Get Diagnoses
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} []string "Diagnoses"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/diagnoses [get]
func getDiagnoses(c echo.Context) error {
	req := struct {
		PatientID uint `json:"-" param:"id" validate:"required"`
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

	patientDiagnoses, err := models.GetPatientDiagnoses(req.PatientID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get patient diagnoses",
		})
	}

	diagnoses := make([]string, 0)
	for _, d := range patientDiagnoses {
		diagnoses = append(diagnoses, d.Diagnosis.Code)
	}

	return c.JSON(http.StatusOK, diagnoses)
}
