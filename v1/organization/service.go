package organization

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

// listServices godoc
// @Summary List Services
// @Description List Services
// @Tags Organization
// @Accept json
// @Produce json
// @Success 200 {object} []models.Service
// @Failure 500 {object} dto.ErrorResponse
// @Router /services [get]
func listServices(c echo.Context) error {
	services, err := models.ListServices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list services",
		})
	}

	return c.JSON(http.StatusOK, services)
}
