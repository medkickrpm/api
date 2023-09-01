package user

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/sendgrid"
	"MedKick-backend/pkg/validator"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func login(c echo.Context) error {
	var request LoginRequest
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

	u := models.User{
		Email: request.Email,
	}
	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "User not found",
		})
	}

	if valid := u.CheckPassword(request.Password); !valid {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid password",
		})
	}

	session, err := middleware.Store.Get(c.Request(), "medkick-session")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get session",
		})
	}

	session.Values["user-id"] = u.ID
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to save session",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully logged in",
	})
}

func logout(c echo.Context) error {
	session, err := middleware.Store.Get(c.Request(), "medkick-session")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get session",
		})
	}

	session.Values["user-id"] = nil
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to save session",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully logged out",
	})
}

type RegisterRequest struct {
	FirstName         string `json:"first_name" validate:"required"`
	LastName          string `json:"last_name" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	Phone             string `json:"phone" validate:"required"`
	Password          string `json:"password" validate:"required"`
	DOB               string `json:"dob" validate:"required"`
	Location          string `json:"location" validate:"required"`
	InsuranceProvider string `json:"insurance_provider" validate:"required"`
	InsuranceID       string `json:"insurance_id" validate:"required"`
	OrganizationID    uint   `json:"organization_id" validate:"required"`
}

func register(c echo.Context) error {
	var request RegisterRequest
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
		Role:              "patientnv",
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

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func resetPassword(c echo.Context) error {
	var request ResetPasswordRequest
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

	u := models.User{
		Email: request.Email,
	}
	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "User not found",
		})
	}

	// Create password reset token
	pwdReset := models.PasswordReset{
		UUID:   uuid.NewString(),
		UserID: u.ID,
	}
	if err := pwdReset.CreatePasswordReset(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create password reset token",
		})
	}

	// TODO - Send password reset email
	body := "Please reset your password by clicking this link: https://med-kick.com/reset-password"
	subject := "MedKick Password Reset"
	if err := sendgrid.SendEmail(fmt.Sprintf("%s %s", u.FirstName, u.LastName), u.Email, subject, body); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to send password reset email",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully sent password reset email",
	})
}

type VerifyResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	UUID     string `json:"uuid" validate:"required"`
}

func verifyResetPassword(c echo.Context) error {
	var request VerifyResetPasswordRequest
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

	u := models.User{
		Email: request.Email,
	}
	if err := u.GetUser(); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "User not found",
		})
	}

	pwdReset := models.PasswordReset{
		UUID: request.UUID,
	}
	if err := pwdReset.GetPasswordReset(); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid token",
		})
	}

	if *pwdReset.UserID != *u.ID {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid token",
		})
	}

	u.Password = request.Password
	if err := u.HashPassword(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to hash password",
		})
	}

	if err := u.UpdateUser(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update user",
		})
	}

	if err := pwdReset.DeletePasswordReset(); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to delete password reset token",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully reset password",
	})
}
