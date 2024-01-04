package user

import (
	gsheet "MedKick-backend/pkg/GSheet"
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/s3"
	"MedKick-backend/pkg/sendgrid"
	"MedKick-backend/pkg/validator"
	"MedKick-backend/utils"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
	"gorm.io/gorm"
)

type CreateRequest struct {
	FirstName         string `json:"first_name" validate:"required"`
	LastName          string `json:"last_name" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	Phone             string `json:"phone" validate:"required"`
	Password          string `json:"password" validate:"required"`
	Role              string `json:"role" validate:"required"` // Roles: admin, doctor, patient, doctornv, patientnv (nv = not verified email)
	DOB               string `json:"dob" validate:"required"`
	Location          string `json:"location" validate:"required"`
	City              string `json:"city"`
	ZipCode           string `json:"zipcode"`
	State             string `json:"state"`
	Country           string `json:"country"`
	InsuranceProvider string `json:"insurance_provider" validate:"required"`
	InsuranceID       string `json:"insurance_id" validate:"required"`
	OrganizationID    uint   `json:"organization_id" validate:"required"`
	Provider          string `json:"provider"`
}

func isValidRole(role string) bool {
	validRoles := map[string]bool{
		"admin":        true,
		"patient":      true,
		"care_manager": true,
		"org_admin":    true,
	}
	return validRoles[role]
}

// createUser godoc
// @Summary Create User
// @Description ADMIN ONLY - Create User
// @Tags User
// @Accept json
// @Produce json
// @Param create body CreateRequest true "Create Request"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user [post]
func createUser(c echo.Context) error {
	var request CreateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	if err := validator.Validate.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to validate request, make sure all fields are filled out correctly",
		})
	}

	if !isValidRole(request.Role) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid role, must be 'admin', 'doctor', 'nurse', 'patient', 'doctornv', 'nursenv', 'care_manager', 'org_admin', or 'patientnv'.",
		})
	}

	// Check if user already exists
	existingUser := models.User{
		Email: request.Email,
		Phone: request.Phone,
	}
	if err := existingUser.GetUser(); err == nil {
		if request.Email == existingUser.Email {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Email already exists",
			})

		}
	}

	if err := existingUser.GetUserByPhone(); err == nil {
		if request.Phone == existingUser.Phone {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Phone already exists",
			})
		}
	}

	// Check if organization exists
	org := models.Organization{
		ID: request.OrganizationID,
	}
	if err := org.GetOrganization(); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Organization does not exist",
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
		City:              request.City,
		ZipCode:           request.ZipCode,
		State:             request.State,
		Country:           request.Country,
		InsuranceProvider: request.InsuranceProvider,
		InsuranceID:       request.InsuranceID,
		OrganizationID:    &request.OrganizationID,
		Provider:          request.Provider,
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

	// Create User Verification
	uv := models.UserVerification{
		UUID:   uuid.NewString(),
		UserID: u.ID,
	}
	if err := uv.CreateUserVerification(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create user verification",
		})
	}

	body := fmt.Sprintf("<p>Please verify your email by clicking this link: <a href=\"https://api.medkick.air.business/v1/auth/validate/%s\">Click Me</a><p>", uv.UUID)
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

// getUser godoc
// @Summary Get User(s)
// @Description Gets users, if ID is specified, gets specific user, if ID is "all", gets all users
// @Tags User
// @Accept json
// @Produce json
// @Param id path string false "User ID"
// @Param filter query string false "Role Filter" Enums(admin, doctor, nurse, patient, doctornv, nursenv, patientnv)
// @Param filter query string false "Status Filter" Enums(critical, warning)
// @Success 200 {object} []models.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id} [get]
func getUser(c echo.Context) error {
	id := c.Param("id")

	filter := c.QueryParam("filter")
	status := c.QueryParam("status")

	if !isValidRole(filter) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid filter",
		})
	}

	self := middleware.GetSelf(c)
	if id == "all" {
		if self.Role != "admin" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: "Unauthorized",
			})
		}

		if filter == "" {
			users, err := models.GetUsers()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}
			return c.JSON(http.StatusOK, users)
		} else {
			users, err := models.GetUsersWithRole(filter)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}

			if filter == "patient" && self.OrganizationID != nil && status != "" {
				filteredPatients, err := filterCriticalPatient(users, status)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
						Error: err.Error(),
					})
				}
				return c.JSON(http.StatusOK, filteredPatients)
			}

			return c.JSON(http.StatusOK, users)
		}
	}
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusOK, self)
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

	if self.Role == "admin" || (self.Role == "care_manager" && *self.OrganizationID == *u.OrganizationID) || (self.Role == "org_admin" && *self.OrganizationID == *u.OrganizationID) {
		fmt.Println("I am admin or something")
		return c.JSON(http.StatusOK, u)
	}

	if *self.ID == *u.ID {
		return c.JSON(http.StatusOK, u)
	}

	return c.JSON(http.StatusForbidden, dto.ErrorResponse{
		Error: "Unauthorized",
	})
}

// getPatients godoc
// @Summary Get Patients(s)
// @Description Gets patients, if ID is specified, gets specific patient, if ID is "all", gets all patients
// @Tags User
// @Accept json
// @Produce json
// @Param id path string false "Patient ID"
// @Param org query string false "Org Filter"
// @Success 200 {object} []models.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /patient/{id} [get]
func getPatients(c echo.Context) error {

	id := c.Param("id")
	org := c.QueryParam("org")

	self := middleware.GetSelf(c)

	if id == "all" {
		if self.Role != "admin" && org == "" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: "Unauthorized",
			})
		}

		if self.Role == "admin" && org == "" {
			users, err := models.GetAllPatients()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get Patients",
				})
			}
			return c.JSON(http.StatusOK, users)
		} else {

			idInt, err := strconv.ParseUint(org, 10, 32)
			if err != nil {
				return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error: "Invalid ID",
				})
			}

			users, err := models.GetAllPatientsWithOrg(idInt)

			if (self.Role != "org_admin" && self.Role != "care_manager") || uint64(*self.OrganizationID) != idInt {
				return c.JSON(http.StatusForbidden, dto.ErrorResponse{
					Error: "Unauthorized",
				})
			}

			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}

			if len(users) > 0 {
				return c.JSON(http.StatusOK, users)
			} else {
				return c.JSON(http.StatusNotFound, dto.ErrorResponse{
					Error: "Patients not found against the organization",
				})
			}
		}
	}

	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get users",
		})
	}

	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}

	user, err := models.GetPatient(uint(idInt))
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Failed to get user",
		})
	}

	if self.Role == "admin" || (self.Role == "doctor" && self.OrganizationID == &user.Organization.ID) {
		return c.JSON(http.StatusOK, user)
	}

	return c.JSON(http.StatusForbidden, dto.ErrorResponse{
		Error: "Unauthorized",
	})
}

// getUsersInOrg godoc
// @Summary Get Users in Organization
// @Description ADMIN & DOCTOR ONLY - if ID is specified, gets users in that organization, if ID is not specified, gets users in self's organization
// @Tags User
// @Accept json
// @Produce json
// @Param id path int false "Organization ID"
// @Param filter query string false "Role Filter" Enums(admin, doctor, nurse, patient, doctornv, nursenv, patientnv)
// @Param filter query string false "Status Filter" Enums(critical, warning)
// @Success 200 {object} []models.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/org/{id} [get]
func getUsersInOrg(c echo.Context) error {
	orgId := c.Param("id")

	filter := c.QueryParam("filter")
	status := c.QueryParam("status")

	if !isValidRole(filter) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid filter",
		})
	}

	self := middleware.GetSelf(c)

	if orgId == "all" {
		if self.Role != "admin" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: "Unauthorized",
			})
		}

		if filter == "" {
			users, err := models.GetUsers()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}

			// Filter out patients
			var filteredUsers []models.User
			for _, user := range users {
				if user.Role != "patient" {
					filteredUsers = append(filteredUsers, user)
				}
			}
			return c.JSON(http.StatusOK, filteredUsers)
		} else {
			users, err := models.GetUsersWithRole(filter)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}

			return c.JSON(http.StatusOK, users)
		}

	}

	if orgId == "" {
		if self.Role == "admin" {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Must specify organization ID",
			})
		}
		if filter == "" {
			users, err := models.GetUsersInOrg(self.OrganizationID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}

			// Filter out admin
			var filteredUsers []models.User
			for _, user := range users {
				if user.Role != "admin" {
					filteredUsers = append(filteredUsers, user)
				}
			}

			return c.JSON(http.StatusOK, filteredUsers)
		} else {
			if filter == "admin" {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: "Unauthorized",
				})
			}

			users, err := models.GetUsersInOrgWithRole(self.OrganizationID, filter)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}
			return c.JSON(http.StatusOK, users)
		}
	}

	orgInt, err := strconv.ParseUint(orgId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}
	orgUint := uint(orgInt)
	/*
		Temperary disabing htis condition to verify the endpoint working
		if self.OrganizationID != &orgUint && self.Role != "admin" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: "Unauthorized",
			})
		}
	*/

	if filter == "" {
		users, err := models.GetUsersInOrg(&orgUint)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get users",
			})
		}

		return c.JSON(http.StatusOK, users)
	} else {
		if filter == "admin" {
			if self.Role != "admin" {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: "Unauthorized",
				})
			}
		}

		users, err := models.GetUsersInOrgWithRole(&orgUint, filter)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get users",
			})
		}

		if filter == "patient" && status != "" {
			filteredPatients, err := filterCriticalPatient(users, status)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: err.Error(),
				})
			}
			return c.JSON(http.StatusOK, filteredPatients)
		}

		return c.JSON(http.StatusOK, users)
	}
}

// getDevicesInUser godoc
// @Summary Get Devices in User
// @Description If ID is specified, gets devices in that user, if ID is not specified, gets devices in self
// @Tags User
// @Accept json
// @Produce json
// @Param id path int false "User ID"
// @Success 200 {object} []models.Device
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/devices [get]
func getDevicesInUser(c echo.Context) error {
	id := c.Param("id")
	self := middleware.GetSelf(c)

	if id == "" {
		devices, err := models.GetDevicesByUser(*self.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get devices",
			})
		}

		return c.JSON(http.StatusOK, devices)
	} else {
		idInt, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Failed to convert id to uint",
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

		if self.Role == "admin" || (self.Role == "org_admin" && *self.OrganizationID == *u.OrganizationID) || (self.Role == "care_manager" && *self.OrganizationID == *u.OrganizationID) || *self.ID == idUint {
			devices, err := models.GetDevicesByUser(idUint)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get devices",
				})
			}

			return c.JSON(http.StatusOK, devices)
		}
	}

	return c.JSON(http.StatusForbidden, dto.ErrorResponse{
		Error: "Unauthorized",
	})
}

// getInteractionsInUser godoc
// @Summary Get Interactions in User
// @Description If ID is specified, gets interactions in that user, if ID is not specified, gets interactions in self
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Success 200 {object} []models.Interaction
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/interactions [get]
func getInteractionsInUser(c echo.Context) error {
	self := middleware.GetSelf(c)

	startDateRaw := c.QueryParam("start_date")
	endDateRaw := c.QueryParam("end_date")

	//convert start_date and end_date to time.Time
	startDate, err := time.Parse("2006-01-02", startDateRaw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to parse start_date",
		})
	}

	var endDate time.Time

	if endDateRaw == "" {
		endDate = time.Now()
	} else {

		// Adjust endDate to the end of the day

		endDate, err = time.Parse("2006-01-02", endDateRaw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Failed to parse end_date",
			})
		}
		endDate = endDate.Add(24 * time.Hour).Add(-time.Second)

	}
	fmt.Println("startDate: ", startDate)
	fmt.Println("endDate: ", endDate)

	// Make sure startDate is before endDate
	if startDate.After(endDate) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Start date must be before end date",
		})
	}

	// Make sure startDate is before present day
	if startDate.After(time.Now()) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Start date must be before present day",
		})
	}

	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to convert id to uint",
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

	if self.Role == "admin" || ((self.Role == "care_manager" || self.Role == "org_admin") && *self.OrganizationID == *u.OrganizationID) || *self.ID == idUint {
		if startDateRaw != "" {
			interactions, err := models.GetInteractionsByUserBetweenDates(idUint, startDate, endDate)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get interactions",
				})
			}

			return c.JSON(http.StatusOK, interactions)
		} else {
			interactions, err := models.GetInteractionsByUser(idUint)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get interactions",
				})
			}

			return c.JSON(http.StatusOK, interactions)
		}
	}

	return c.JSON(http.StatusForbidden, dto.ErrorResponse{
		Error: "Unauthorized",
	})
}

// getTotalInteractionDuration godoc
// @Summary Total user interaction duration
// @Description Get the sum of the interaction durations for a user
// @Tags Interaction
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Success 200 {object} uint
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/interactions/duration [get]
func getTotalInteractionDuration(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}

	idUint := uint(idInt)
	self := middleware.GetSelf(c)
	if self.Role == "patient" && *self.ID != idUint {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error: "Forbidden",
		})
	}

	u := models.User{
		ID: &idUint,
	}

	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get user",
		})
	}

	startDateRaw := c.QueryParam("start_date")
	endDateRaw := c.QueryParam("end_date")

	//convert start_date and end_date to time.Time
	startDate, err := time.Parse("2006-01-02", startDateRaw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to parse start_date",
		})
	}

	var endDate time.Time

	if endDateRaw == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", endDateRaw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Failed to parse end_date",
			})
		}
	}

	// Make sure startDate is before endDate
	if startDate.After(endDate) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Start date must be before end date",
		})
	}

	// Make sure startDate is before present day
	if startDate.After(time.Now()) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Start date must be before present day",
		})
	}

	if self.Role == "admin" || ((self.Role == "doctor" || self.Role == "nurse") && *self.OrganizationID == *u.OrganizationID) || *self.ID == idUint {
		if startDateRaw != "" {
			interactions, err := models.GetInteractionsByUser(idUint)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get interactions",
				})
			}

			var totalDuration uint
			for _, interaction := range interactions {
				totalDuration += interaction.Duration
			}

			return c.JSON(http.StatusOK, totalDuration)
		} else {
			interactions, err := models.GetInteractionsByUserBetweenDates(idUint, startDate, endDate)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get interactions",
				})
			}

			var totalDuration uint
			for _, interaction := range interactions {
				totalDuration += interaction.Duration
			}

			return c.JSON(http.StatusOK, totalDuration)
		}
	}

	return c.JSON(http.StatusForbidden, dto.ErrorResponse{
		Error: "Forbidden",
	})
}

// getCarePlans godoc
// @Summary Get care plans in User
// @Description If ID is specified, gets care plans in that user, if ID is not specified, gets care plans in self
// @Tags User
// @Accept json
// @Produce json
// @Param id path int false "User ID"
// @Success 200 {object} dto.CareplanSheetResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id}/careplans [get]
func getCarePlansInUser(c echo.Context) error {
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}

	idUint := uint(idInt)

	user := models.User{
		ID: &idUint,
	}

	if err := user.GetUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get user with this id",
		})
	}

	spreadsheetID := "1ZMmc0Sv74GVg6PRKLhAmlqbVE02rPfw7s8IXDGOrtPI"
	spreadsheet, err := gsheet.Service.FetchSpreadsheet(spreadsheetID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to fetch spreadsheet",
		})
	}

	sheet, err := spreadsheet.SheetByTitle("Careplans")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get sheet",
		})
	}

	// get all the rows
	sheet.Rows = sheet.Rows[1:]

	formatedUserDOB, err := utils.ConvertDateFormat(user.DOB, "01-02-2006")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to convert date format for comparison",
		})
	}

	userFirstName := strings.ToLower(user.FirstName)
	userLastName := strings.ToLower(user.LastName)
	userFirstName = strings.Trim(userFirstName, " ")
	userLastName = strings.Trim(userLastName, " ")

	// check if user exists in sheet
	for _, row := range sheet.Rows {
		if strings.Trim(strings.ToLower(row[3].Value), " ") == userFirstName && strings.Trim(strings.ToLower(row[4].Value), " ") == userLastName && row[5].Value == formatedUserDOB {

			// map row to CareplanSheetResponse
			var response dto.CareplanSheetResponse
			val := reflect.ValueOf(&response).Elem()

			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)

				if field.CanSet() && i < len(row) {
					switch field.Kind() {
					case reflect.String:
						field.SetString(row[i].Value)
					}
				}
			}

			return c.JSON(http.StatusOK, response)
		}
	}

	return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Error: "Found no record for this user in GSheet",
	})

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
	City              string `json:"city"`
	ZipCode           string `json:"zipcode"`
	State             string `json:"state"`
	Country           string `json:"country"`
	InsuranceProvider string `json:"insurance_provider"`
	InsuranceID       string `json:"insurance_id"`
	OrganizationID    *uint  `json:"organization_id"`
	Provider          string `json:"provider,omitempty"`
}

// updateUser godoc
// @Summary Update User
// @Description Updates user, if ID is specified, updates specific user, if ID is not specified, updates self
// @Tags User
// @Accept json
// @Produce json
// @Param id path string false "User ID"
// @Param update body UpdateRequest true "Update Request"
// @Success 200 {object} models.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id} [patch]
func updateUser(c echo.Context) error {
	id := c.Param("id")

	var request UpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	if !isValidRole(request.Role) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid role, must be 'admin', 'doctor', 'nurse', 'patient', 'doctornv', 'nursenv', 'care_manager', 'org_admin', or 'patientnv'.",
		})
	}

	self := middleware.GetSelf(c)

	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}
	idUint := uint(idInt)

	if id == "" || idUint == *self.ID {
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
		if request.City != "" {
			self.City = request.City
		}
		if request.ZipCode != "" {
			self.ZipCode = request.ZipCode
		}
		if request.State != "" {
			self.State = request.State
		}
		if request.Country != "" {
			self.Country = request.Country
		}
		if request.InsuranceProvider != "" {
			self.InsuranceProvider = request.InsuranceProvider
		}
		if request.InsuranceID != "" {
			self.InsuranceID = request.InsuranceID
		}

		self.Provider = request.Provider

		if err := self.UpdateUser(); err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to update user",
			})
		}

		return c.JSON(http.StatusOK, self)
	}

	u := models.User{
		ID: &idUint,
	}
	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get user",
		})
	}

	if self.Role == "admin" || self.Role == "org_admin" || (self.Role == "care_manager" && *self.OrganizationID == *u.OrganizationID) {
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
		if request.City != "" {
			u.City = request.City
		}
		if request.ZipCode != "" {
			u.ZipCode = request.ZipCode
		}
		if request.State != "" {
			u.State = request.State
		}
		if request.Country != "" {
			u.Country = request.Country
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
				return c.JSON(http.StatusForbidden, dto.ErrorResponse{
					Error: "Forbidden",
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

	return c.JSON(http.StatusForbidden, dto.ErrorResponse{
		Error: "Forbidden",
	})
}

// validateUser godoc
// @Summary Validate User
// @Description Validates user
// @Tags User
// @Accept json
// @Produce json
// @Success 307
// @Failure 307
// @Router /auth/validate/{id} [get]
func validateUser(c echo.Context) error {
	uuid := c.Param("id")

	uv := models.UserVerification{
		UUID: uuid,
	}
	if err := uv.GetUserVerification(); err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "https://www.medkick.air.business/user/failed")
	}

	u := uv.User
	u.Role = u.Role[:len(u.Role)-2] // Remove the nv

	if err := u.UpdateUser(); err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "https://www.medkick.air.business/user/failed")
	}

	if err := uv.DeleteUserVerification(); err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "https://www.medkick.air.business/user/failed")
	}

	return c.Redirect(http.StatusTemporaryRedirect, "https://www.medkick.air.business/user/verified")
}

// deleteUser godoc
// @Summary Delete User
// @Description Admin & Doctor ONLY - Delete User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/{id} [delete]
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

// countUser godoc
// @Summary Count Users
// @Description ADMIN ONLY - Count Users
// @Tags User
// @Accept json
// @Produce json
// @Param filter query string false "Role Filter" Enums(admin, doctor, nurse, patient, doctornv, nursenv, patientnv)
// @Param status query string false "Status Filter" Enums(critical, warning)
// @Success 200 {object} map[string]int64
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/count [get]
func countUser(c echo.Context) error {
	filter := c.QueryParam("filter")
	status := c.QueryParam("status")

	if filter != "" && filter != "admin" && filter != "doctor" && filter != "patient" && filter != "nurse" && filter != "doctornv" && filter != "nursenv" && filter != "patientnv" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid filter",
		})
	}

	self := middleware.GetSelf(c)

	if self.Role != "admin" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	var userCount int64
	var err error

	if filter == "" {
		userCount, err = models.CountUsers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get users",
			})
		}
	} else {
		if filter != "patient" || status == "" {
			userCount, err = models.CountUsersWithRole(filter)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}
		} else {

			users, err := models.GetUsersWithRole(filter)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}

			filteredPatients, err := filterCriticalPatient(users, status)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: err.Error(),
				})
			}

			userCount = int64(len(filteredPatients))
		}
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"count": userCount,
	})
}

// countUserInOrg godoc
// @Summary Count Users in Organization
// @Description Count Users in Organization
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param filter query string false "Role Filter" Enums(admin, doctor, nurse, patient, doctornv, nursenv, patientnv)
// @Param status query string false "Status Filter" Enums(critical, warning)
// @Success 200 {object} map[string]int64
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/org/{id}/count [get]
func countUserInOrg(c echo.Context) error {
	orgId := c.Param("id")

	filter := c.QueryParam("filter")
	status := c.QueryParam("status")

	if filter != "" && filter != "admin" && filter != "doctor" && filter != "patient" && filter != "nurse" && filter != "doctornv" && filter != "nursenv" && filter != "patientnv" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid filter",
		})
	}

	var userCount int64
	var err error

	orgInt, err := strconv.ParseUint(orgId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid ID",
		})
	}
	orgUint := uint(orgInt)

	if filter == "" {
		userCount, err = models.CountUsersInOrg(orgUint)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: "Failed to get users",
			})
		}
	} else {
		if filter != "patient" || status == "" {
			userCount, err = models.CountUsersWithRoleInOrg(orgUint, filter)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}
		} else {
			users, err := models.GetUsersInOrgWithRole(&orgUint, filter)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Failed to get users",
				})
			}

			filteredPatients, err := filterCriticalPatient(users, status)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: err.Error(),
				})
			}

			userCount = int64(len(filteredPatients))
		}
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"count": userCount,
	})
}

func filterCriticalPatient(users []models.User, status string) (filteredPatients []models.User, err error) {
	var patientList []uint
	for _, user := range users {
		patientList = append(patientList, *user.ID)
	}

	filteredPatients = make([]models.User, 0)

	telemetryData, err := models.GetPatientTelemetryData(patientList)
	if err != nil {
		return nil, errors.New("failed to get telemetry data")
	}

	latestTelemetryData := models.GetLatestPatientTelemetryData(telemetryData)

	thresholdList, err := models.ListAlertThresholds(patientList)
	if err != nil {
		return nil, errors.New("failed to get alert thresholds")
	}

	patientStatusFunc := models.GetPatientStatusFunc(thresholdList)
	patientSelected := make(map[uint]struct{})
	for _, data := range latestTelemetryData {
		isCritical, isWarning := patientStatusFunc(data)
		if (status == "critical" && isCritical) || (status == "warning" && isWarning) {
			patientSelected[data.PatientID] = struct{}{}
		}
	}

	for _, user := range users {
		if _, ok := patientSelected[*user.ID]; ok {
			filteredPatients = append(filteredPatients, user)
		}
	}

	return filteredPatients, nil
}

// verifyField godoc
// @Summary Verify if field already exists Field
// @Description Verify if field already exists
// @Description when is_available is true that means the phone/email does not already exist in the DB and can be used and If the response is false that means the value already exists
// @Tags User
// @Accept json
// @Produce json
// @Param field query string true "Field" Enums(email, phone)
// @Param value query string true "Value"
// @Success 200 {object} dto.VerifyUserFieldResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /user/verifyUserField [get]
func verifyUserField(c echo.Context) error {
	field := c.QueryParam("field")
	value := c.QueryParam("value")

	if field != "email" && field != "phone" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid field",
		})
	}

	if value == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid value",
		})
	}

	exists, err := models.VerifyUserField(field, value)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to verify field",
		})
	}

	if !exists {
		return c.JSON(http.StatusOK, dto.VerifyUserFieldResponse{
			IsAvailable: true,
		})
	}

	return c.JSON(http.StatusOK, dto.VerifyUserFieldResponse{
		IsAvailable: false,
	})
}

// uploadUserAvatar godoc
// @Summary Upload User Avatar
// @Description Upload User Avatar
// @Tags User
// @Accept json
// @Produce json
// @Param avatar formData file true "Avatar"
// @Success 200 {object} string
// @Failure 400 {object} dto.ErrorResponse
// @Router /user/avatar [post]
func uploadUserAvatar(c echo.Context) error {
	file, err := c.FormFile("avatar")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to get file",
		})
	}

	// Check if file extension is valid
	extension := filepath.Ext(file.Filename)
	if extension != ".png" && extension != ".jpg" && extension != ".jpeg" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid file type",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to open file",
		})
	}
	defer src.Close()

	// Decode the image
	img, _, err := image.Decode(src)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to decode image",
		})
	}

	// Get current user
	self := middleware.GetSelf(c)

	// Generate file name with user first name, last name, and timestamp
	fileName := fmt.Sprintf("%s_%s_%d%s", strings.TrimSpace(self.FirstName), strings.TrimSpace(self.LastName), time.Now().Unix(), extension)
	uploadLocation := s3.AvatarFolder + fileName

	// Resize the image
	resizedImg := resize.Resize(200, 0, img, resize.Lanczos3)

	// Create a new file
	thumbnailPath := "/tmp/" + fileName
	thumbnail, err := os.Create(thumbnailPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create new file",
		})
	}
	defer os.Remove(thumbnailPath)
	defer thumbnail.Close()

	// Encode the image to the new file
	var encodeErr error
	switch extension {
	case ".png":
		encodeErr = png.Encode(thumbnail, resizedImg)
	case ".jpg", ".jpeg":
		encodeErr = jpeg.Encode(thumbnail, resizedImg, &jpeg.Options{Quality: 75})
	}
	if encodeErr != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to encode image",
		})
	}

	// Open the thumbnail file
	thumbnailFile, err := os.Open(thumbnailPath)
	if err != nil {
		return err
	}
	defer thumbnailFile.Close()

	// Upload to S3
	if err := s3.UploadFile(uploadLocation, thumbnailFile); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to upload file",
		})
	}

	// Update user avatar src
	if err := models.UpdateUserAvatar(*self.ID, fileName); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update user avatar src",
		})
	}

	return c.JSON(http.StatusOK, fileName)
}

// getUserAvatar godoc
// @Summary Get User Avatar
// @Description Get User Avatar
// @Tags User
// @Accept json
// @Produce json
// @Param avatarPath path string true "Avatar Path"
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Router /user/avatar/{avatarPath} [get]
func getUserAvatar(c echo.Context) error {
	avatarPath := c.Param("avatarPath")

	// get avatar from s3
	file, err := s3.DownloadFile(s3.AvatarFolder + avatarPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get avatar",
		})
	}

	defer file.Body.Close()

	// get file extension
	extension := filepath.Ext(avatarPath)

	// set content type
	switch extension {
	case ".png":
		c.Response().Header().Set("Content-Type", "image/png")
	case ".jpg":
		c.Response().Header().Set("Content-Type", "image/jpg")
	case ".jpeg":
		c.Response().Header().Set("Content-Type", "image/jpeg")
	}

	// copy file to response
	if _, err := io.Copy(c.Response(), file.Body); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get avatar",
		})
	}

	return nil
}
