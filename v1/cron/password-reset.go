package cron

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Token string `json:"token" validate:"required"`
}

// clearPasswordResetTokens clears all password reset tokens that are older than 24 hours
func clearPasswordResetTokens(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to bind request",
		})
	}

	if err := validator.Validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request",
		})
	}

	if req.Token != os.Getenv("CRON_SECRET") {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Invalid token",
		})
	}

	pwdResets, err := models.GetPasswordResets()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get password resets",
		})
	}

	for _, pwdReset := range pwdResets {
		if pwdReset.CreatedAt.Add(24 * time.Hour).Before(time.Now()) {
			err := pwdReset.DeletePasswordReset()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to delete password reset",
				})
			}
		}
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully cleared password resets",
	})
}
