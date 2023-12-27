package user

import (
	"MedKick-backend/pkg/echo/middleware"

	"github.com/labstack/echo/v4"
)

func Routes(r *echo.Group) {
	r.POST("/auth/login", login)
	r.GET("/auth/logout", logout, middleware.NotGuest)
	r.POST("/auth/register", register)
	r.POST("/auth/reset-password", resetPassword)
	r.POST("/auth/verify-reset-password", verifyResetPassword)

	r.GET("/auth/validate/:id", validateUser)

	r.POST("/user", createUser, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager"))
	r.GET("/user/:id", getUser, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))
	r.GET("/patient/:id", getPatients, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager"))
	r.GET("/user", getUser, middleware.NotGuest)
	r.GET("/user/count", countUser, middleware.NotGuest)
	r.GET("/user/org/:id", getUsersInOrg, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager"))
	r.GET("/user/org/:id/count", countUserInOrg, middleware.NotGuest, middleware.HasRole("admin", "doctor", "nurse"))
	r.PATCH("/user/:id", updateUser, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager"))
	r.DELETE("/user/:id", deleteUser, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager"))

	r.GET("/user/:id/devices", getDevicesInUser, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))
	r.GET("/user/:id/interactions", getInteractionsInUser, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))
	r.GET("/user/:id/interactions/duration", getTotalInteractionDuration, middleware.NotGuest)
	r.GET("/user/:id/careplans", getCarePlansInUser, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))

	r.PUT("/user/:id/alert-threshold", upsertAlertThreshold, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))
	r.GET("/user/:id/alert-threshold", listAlertThresholds, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))

	r.PUT("/user/:id/diagnoses", upsertDiagnoses, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))
	r.GET("/user/:id/diagnoses", getDiagnoses, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))
	r.PUT("/user/:id/patient-service", upsertPatientServices, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))
	r.GET("/user/:id/patient-service", listPatientServices, middleware.NotGuest, middleware.HasRole("admin", "org_admin", "care_manager", "patient"))
}
