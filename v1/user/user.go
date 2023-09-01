package user

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/sendgrid"
	"MedKick-backend/pkg/validator"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CreateRequest struct {
	FirstName         string `json:"first_name" validate:"required"`
	LastName          string `json:"last_name" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	Phone             string `json:"phone" validate:"required"`
	Password          string `json:"password" validate:"required"`
	Role              string `json:"role" validate:"required"`
	DOB               string `json:"dob" validate:"required"`
	Location          string `json:"location" validate:"required"`
	InsuranceProvider string `json:"insurance_provider" validate:"required"`
	InsuranceID       string `json:"insurance_id" validate:"required"`
	OrganizationID    uint   `json:"organization_id" validate:"required"`
}

// Roles: admin, doctor, patient, doctornv, patientnv (nv = not verified email)
func createUser(c echo.Context) error {
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

	if request.Role != "admin" && request.Role != "doctor" && request.Role != "patient" && request.Role != "doctornv" && request.Role != "patientnv" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid role, must be 'admin', 'doctor', 'patient', 'doctornv', or 'patientnv'.",
		})
	}

	// Check if user already exists
	existingUser := models.User{
		Email: request.Email,
	}
	if err := existingUser.GetUser(); err == nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "User already exists",
		})
	}

	u := models.User{
		FirstName:         request.FirstName,
		LastName:          request.LastName,
		Email:             request.Email,
		Phone:             request.Phone,
		Password:          request.Password,
		Role:              request.Role,
		DOB:               request.DOB,
		Location:          request.Location,
		InsuranceProvider: request.InsuranceProvider,
		InsuranceID:       request.InsuranceID,
		OrganizationID:    &request.OrganizationID,
	}

	// Hash password & Create User
	if err := u.HashPassword(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to hash password",
		})
	}

	if err := u.CreateUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create user",
		})
	}

	// TODO - Send verification email
	body := "Please verify your email by clicking this link: https://med-kick.com/verify-email"
	subject := "MedKick Email Verification"
	if err := sendgrid.SendEmail(fmt.Sprintf("%s %s", u.FirstName, u.LastName), u.Email, subject, body); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to send verification email",
		})
	}

	return c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "Successfully created user",
	})
}

func getUser(c echo.Context) error {
	id := c.Param("id")

	self := middleware.GetSelf(c)
	if id == "all" {
		if self.Role != "admin" {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Unauthorized",
			})
		}

		users, err := models.GetUsers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get users",
			})
		}
		return c.JSON(http.StatusCreated, users)
	}
	if id == "" {
		return c.JSON(http.StatusCreated, self)
	}

	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}
	idUint := uint(idInt)
	u := models.User{
		ID: &idUint,
	}
	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get user",
		})
	}

	if self.Role == "admin" || (self.Role == "doctor" && self.OrganizationID == u.OrganizationID) {
		return c.JSON(http.StatusCreated, u)
	}

	return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
		Error: "Unauthorized",
	})
}

func getUsersInOrg(c echo.Context) error {
	orgId := c.Param("id")

	self := middleware.GetSelf(c)

	if orgId == "" {
		if self.Role == "admin" {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Must specify organization ID",
			})
		}
		users, err := models.GetUsersInOrg(self.OrganizationID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get users",
			})
		}
		return c.JSON(http.StatusOK, users)
	}

	orgInt, err := strconv.ParseUint(orgId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}
	orgUint := uint(orgInt)
	if self.OrganizationID != &orgUint && self.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	users, err := models.GetUsersInOrg(&orgUint)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get users",
		})
	}
	return c.JSON(http.StatusOK, users)
}

type UpdateRequest struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Password          string `json:"password"`
	Role              string `json:"role"`
	DOB               string `json:"dob"`
	Location          string `json:"location"`
	InsuranceProvider string `json:"insurance_provider"`
	InsuranceID       string `json:"insurance_id"`
	OrganizationID    *uint  `json:"organization_id"`
}

func updateUser(c echo.Context) error {
	id := c.Param("id")

	var request UpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	self := middleware.GetSelf(c)
	if id == "" {
		if request.FirstName != "" {
			self.FirstName = request.FirstName
		}
		if request.LastName != "" {
			self.LastName = request.LastName
		}
		if request.Email != "" {
			self.Email = request.Email
		}
		if request.Phone != "" {
			self.Phone = request.Phone
		}
		if request.Password != "" {
			self.Password = request.Password
			err := self.HashPassword()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to hash password",
				})
			}
		}
		if request.Role != "" {
			self.Role = request.Role
		}
		if request.DOB != "" {
			self.DOB = request.DOB
		}
		if request.Location != "" {
			self.Location = request.Location
		}
		if request.InsuranceProvider != "" {
			self.InsuranceProvider = request.InsuranceProvider
		}
		if request.InsuranceID != "" {
			self.InsuranceID = request.InsuranceID
		}

		if err := self.UpdateUser(); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to update user",
			})
		}
	}

	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}
	idUint := uint(idInt)
	u := models.User{
		ID: &idUint,
	}
	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get user",
		})
	}

	if self.Role == "admin" || (self.Role == "doctor" && self.OrganizationID == u.OrganizationID) {
		if request.FirstName != "" {
			u.FirstName = request.FirstName
		}
		if request.LastName != "" {
			u.LastName = request.LastName
		}
		if request.Email != "" {
			u.Email = request.Email
		}
		if request.Phone != "" {
			u.Phone = request.Phone
		}
		if request.Password != "" {
			u.Password = request.Password
			err := u.HashPassword()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to hash password",
				})
			}
		}
		if request.Role != "" {
			u.Role = request.Role
		}
		if request.DOB != "" {
			u.DOB = request.DOB
		}
		if request.Location != "" {
			u.Location = request.Location
		}
		if request.InsuranceProvider != "" {
			u.InsuranceProvider = request.InsuranceProvider
		}
		if request.InsuranceID != "" {
			u.InsuranceID = request.InsuranceID
		}

		if request.OrganizationID != nil {
			if self.Role == "admin" {
				u.OrganizationID = request.OrganizationID
			} else {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: "Unauthorized",
				})
			}
		}

		if err := u.UpdateUser(); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to update user",
			})
		}
		return c.JSON(http.StatusOK, u)
	}

	return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
		Error: "Unauthorized",
	})
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}
	idUint := uint(idInt)

	u := models.User{
		ID: &idUint,
	}
	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get user",
		})
	}

	self := middleware.GetSelf(c)
	if self.Role != "admin" && self.OrganizationID != u.OrganizationID {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	if err := u.DeleteUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to delete user",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully deleted user",
	})
}
