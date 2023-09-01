package careplan

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
	UserID   uint  `json:"user_id" validate:"required"`
	DoctorID *uint `json:"doctor_id"`
}

// Returns a presigned URL for the user to upload their care plan
func createCareplan(c echo.Context) error {
	var request CreateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	if err := validator.Validate.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	self := middleware.GetSelf(c)
	if self.Role == "doctor" {
		request.DoctorID = self.ID
	}
	if self.Role == "admin" {
		if request.DoctorID == nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Doctor ID is required",
			})
		}
	}

	// Create presigned URL here...
	url := ""

	careplan := &models.CarePlan{
		UserID:   request.UserID,
		DoctorID: *request.DoctorID,
		URL:      url,
	}
	if err := careplan.CreateCarePlan(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create care plan",
		})
	}

	return c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "Care plan created successfully",
	})
}

func getCareplan(c echo.Context) error {
	id := c.Param("id")

	if id != "" {
		idUint, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Invalid ID",
			})
		}

		careplan := &models.CarePlan{
			ID: uint(idUint),
		}
		if err := careplan.GetCarePlan(); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get care plan",
			})
		}

		self := middleware.GetSelf(c)
		if self.Role == "patient" {
			if &careplan.UserID != self.ID {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: "Unauthorized: You are not the patient",
				})
			}
		}
		if self.Role == "doctor" {
			if careplan.User.OrganizationID != self.OrganizationID {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: "Unauthorized: You are not in the same organization as the patient",
				})
			}
		}

		return c.JSON(http.StatusOK, careplan)
	}

	self := middleware.GetSelf(c)
	if self.Role == "admin" {
		careplans, err := models.GetCarePlans()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get care plans",
			})
		}

		return c.JSON(http.StatusOK, careplans)
	}
	if self.Role == "doctor" {
		careplans, err := models.GetCarePlansByDoctorID(*self.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get care plans",
			})
		}

		return c.JSON(http.StatusOK, careplans)
	}
	if self.Role == "patient" {
		careplans, err := models.GetCarePlansByUserID(*self.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get care plans",
			})
		}

		return c.JSON(http.StatusOK, careplans)
	}

	return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
		Error: "Err. grabbing careplan data, please report this to an admin",
	})
}

func deleteCareplan(c echo.Context) error {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}

	careplan := &models.CarePlan{
		ID: uint(idUint),
	}
	if err := careplan.GetCarePlan(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get care plan",
		})
	}

	self := middleware.GetSelf(c)
	if self.Role == "doctor" {
		if careplan.User.OrganizationID != self.OrganizationID {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Unauthorized: You are not in the same organization as the patient",
			})
		}
	}

	if err := careplan.DeleteCarePlan(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to delete care plan",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Care plan deleted successfully",
	})
}
