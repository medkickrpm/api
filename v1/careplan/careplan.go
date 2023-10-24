package careplan

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/s3"
	"MedKick-backend/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CreateRequest struct {
	UserID   uint  `json:"user_id" validate:"required"`
	DoctorID *uint `json:"doctor_id"`
}

// createCareplan godoc
// @Summary Create a care plan
// @Description Creates a care plan for the provided user
// @Tags Careplan
// @Accept json
// @Produce json
// @Param create body CreateRequest true "Create Request"
// @Success 201 {object} models.CarePlan
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /careplan [post]
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
	if self.Role == "doctor" || self.Role == "nurse" {
		request.DoctorID = self.ID
	}
	if self.Role == "admin" {
		if request.DoctorID == nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Doctor ID is required",
			})
		}
	}

	careplan := &models.CarePlan{
		UserID:   request.UserID,
		DoctorID: *request.DoctorID,
		URL:      "",
	}
	if err := careplan.CreateCarePlan(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create care plan",
		})
	}

	return c.JSON(http.StatusCreated, careplan)
}

// uploadCareplan godoc
// @Summary Upload Careplan
// @Description Uploads a careplan for the provided user
// @Tags Careplan
// @Param id path int true "Careplan ID"
// @Param file formData file true "File"
// @Success 200 {object} []models.CarePlan
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /careplan/{id} [put]
func uploadCareplan(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "File is required",
		})
	}
	// Capped at 100MB
	if file.Size > 100*1024*1024 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "File is too large",
		})
	}
	if file.Size == 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "File is empty",
		})
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

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
	if self.Role == "doctor" || self.Role == "nurse" {
		if careplan.User.OrganizationID != self.OrganizationID {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Unauthorized: You are not in the same organization as the patient",
			})
		}
	}

	if careplan.URL != "" {
		if err := s3.DeleteFile(careplan.URL); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to delete old careplan",
			})
		}
	}

	careplan.URL = "careplan/" + strconv.FormatUint(uint64(careplan.ID), 10) + "/" + file.Filename
	if err := s3.UploadFile(careplan.URL, src); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to upload file",
		})
	}

	if err := careplan.UpdateCarePlan(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update care plan",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Care plan uploaded successfully",
	})
}

// getCareplan godoc
// @Summary Get Careplan(s)
// @Description Returns a list of careplans if no ID is provided, otherwise returns the careplan with the provided ID
// @Tags Careplan
// @Accept json
// @Produce json
// @Param id path int false "Careplan ID"
// @Success 200 {object} []models.CarePlan
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /careplan/{id} [get]
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
		if self.Role == "doctor" || self.Role == "nurse" {
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
	if self.Role == "doctor" || self.Role == "nurse" {
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

// downloadCareplan godoc
// @Summary Download Careplan
// @Description Download the careplan with the provided ID
// @Tags Careplan
// @Accept json
// @Produce json
// @Param id path int true "Careplan ID"
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /careplan/{id}/file [get]
func downloadCareplan(c echo.Context) error {
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
	if self.Role == "doctor" || self.Role == "nurse" {
		if careplan.User.OrganizationID != self.OrganizationID {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Unauthorized: You are not in the same organization as the patient",
			})
		}
	}

	file, err := s3.DownloadFile(careplan.URL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to download file",
		})
	}

	defer file.Body.Close()

	return c.Stream(http.StatusOK, *file.ContentType, file.Body)
}

// deleteCareplan godoc
// @Summary Delete a Careplan
// @Description Deletes the careplan with the provided ID
// @Tags Careplan
// @Accept json
// @Produce json
// @Param id path int true "Careplan ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /careplan/{id} [delete]
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
	if self.Role == "doctor" || self.Role == "nurse" {
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
