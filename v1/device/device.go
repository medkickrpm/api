package device

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/validator"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DeviceAssignRequest struct {
	DeviceID uint `json:"device_id" validate:"required"`
	UserID   uint `json:"user_id" validate:"required"`
}

// getDevice godoc
// @Summary Get Devices
// @Description Get devices by id, set id to 'all' to get all devices
// @Tags Devices
// @Accept json
// @Produce json
// @Param id path string false "Device ID"
// @Success 200 {object} []models.Device
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /device/{id} [get]
func getDevice(c echo.Context) error {
	self := middleware.GetSelf(c)

	id := c.Param("id")

	if id == "all" || id == "" {
		if self.Role == "admin" {
			devices, err := models.GetDevices()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get devices",
				})
			}
			return c.JSON(http.StatusOK, devices)
		} else if self.Role == "doctor" || self.Role == "nurse" {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Doctors/Nurses cannot get all devices",
			})
		} else {
			devices, err := models.GetDevicesByUser(*self.ID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get devices by user",
				})
			}
			return c.JSON(http.StatusOK, devices)
		}
	}

	// Convert id to uint
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to convert id to uint",
		})
	}

	device := &models.Device{
		ID: uint(idInt),
	}
	if err := device.GetDevice(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get device by device id",
		})
	}

	if self.Role == "patient" && device.UserID != *self.ID {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error: "Forbidden",
		})
	}

	return c.JSON(http.StatusOK, device)
}

// GetAvailableDevices godoc
// @Summary Get Available Devices
// @Description Get Available Devices
// @Tags Devices
// @Produce json
// @Success 200 {object} []models.DeviceDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /device/available-devices [get]
func GetAvailableDevices(c echo.Context) error {
	device := &models.Device{}
	devices, err := device.GetAvailableDevices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get available devices",
		})
	}
	return c.JSON(http.StatusOK, devices)
}

// AssignDevice godoc
// @Summary Assign Device
// @Description Assign Device
// @Tags Devices
// @Accept json
// @Produce json
// @Param request body DeviceAssignRequest true "Assign Device"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /device/assign-device [patch]
func AssignDevice(c echo.Context) error {
	var request DeviceAssignRequest
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

	// Check if device id is valid
	device := &models.Device{
		ID: request.DeviceID,
	}
	if err := device.GetDevice(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get device by device id",
		})
	}

	// Check if user id is valid
	user := &models.User{
		ID: &request.UserID,
	}
	if err := user.GetUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get user by user id",
		})
	}

	// Assign device to user
	device = &models.Device{}

	if err := device.AssignDeviceToUser(request.DeviceID, request.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to assign device to user",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully assigned device to user",
	})
}

type UpdateRequest struct {
	Name            string `json:"name"`
	ModelNumber     string `json:"model_number"`
	IMEI            string `json:"imei"`
	SerialNumber    string `json:"serial_number"`
	BatteryLevel    *uint  `json:"battery_level"`
	SignalStrength  string `json:"signal_strength"`
	FirmwareVersion string `json:"firmware_version"`
	UserID          *uint  `json:"user_id"`
}

// updateDevice godoc
// @Summary Update Device
// @Description Update Device
// @Tags Devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param request body UpdateRequest true "Update Device"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /device/{id} [patch]
func updateDevice(c echo.Context) error {
	self := middleware.GetSelf(c)

	id := c.Param("id")

	// Convert id to uint
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to convert id to uint",
		})
	}

	device := &models.Device{
		ID: uint(idInt),
	}
	if err := device.GetDevice(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get device by device id",
		})
	}

	if self.Role == "patient" && device.UserID != *self.ID {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error: "Forbidden",
		})
	}

	var request UpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to bind request",
		})
	}

	if request.Name != "" {
		device.Name = request.Name
	}
	if request.ModelNumber != "" {
		device.ModelNumber = request.ModelNumber
	}
	if request.IMEI != "" {
		device.IMEI = request.IMEI
	}
	if request.SerialNumber != "" {
		device.SerialNumber = request.SerialNumber
	}
	if request.BatteryLevel != nil {
		device.BatteryLevel = *request.BatteryLevel
	}
	if request.SignalStrength != "" {
		device.SignalStrength = request.SignalStrength
	}
	if request.FirmwareVersion != "" {
		device.FirmwareVersion = request.FirmwareVersion
	}
	if request.UserID != nil {
		device.UserID = *request.UserID
	}

	if err := device.UpdateDevice(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update device",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully updated device",
	})
}

// deleteDevice godoc
// @Summary Delete Device
// @Description Delete Device
// @Tags Devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /device/{id} [delete]
func deleteDevice(c echo.Context) error {
	self := middleware.GetSelf(c)

	id := c.Param("id")

	// Convert id to uint
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to convert id to uint",
		})
	}

	device := &models.Device{
		ID: uint(idInt),
	}
	if err := device.GetDevice(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get device by device id",
		})
	}

	if self.Role == "patient" && device.UserID != *self.ID {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error: "Forbidden",
		})
	}

	if err := device.DeleteDevice(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to delete device",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully deleted device",
	})
}
