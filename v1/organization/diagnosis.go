package organization

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

// listDiagnosisCodes godoc
// @Summary List Diagnosis Codes
// @Description List Diagnosis Codes
// @Tags Organization
// @Accept json
// @Produce json
// @Success 200 {object} []models.Diagnosis
// @Failure 500 {object} dto.ErrorResponse
// @Router /diagnoses [get]
func listDiagnosisCodes(c echo.Context) error {
	codes, err := models.ListDiagnosisCodes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list diagnosis codes",
		})
	}

	return c.JSON(http.StatusOK, codes)
}
