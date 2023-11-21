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

	r.POST("/user", createUser, middleware.NotGuest, middleware.HasRole("admin"))
	r.GET("/user/:id", getUser, middleware.NotGuest)
	r.GET("/user", getUser, middleware.NotGuest)
	r.GET("/user/count", countUser, middleware.NotGuest)
	r.GET("/user/org/:id", getUsersInOrg, middleware.NotGuest, middleware.HasRole("admin", "doctor", "nurse"))
	r.GET("/user/org/:id/count", countUserInOrg, middleware.NotGuest, middleware.HasRole("admin", "doctor", "nurse"))
	r.PATCH("/user/:id", updateUser, middleware.NotGuest)
	r.DELETE("/user/:id", deleteUser, middleware.NotGuest, middleware.HasRole("admin", "doctor"))

	r.GET("/user/:id/devices", getDevicesInUser, middleware.NotGuest)
	r.GET("/user/:id/interactions", getInteractionsInUser, middleware.NotGuest)
	r.GET("/user/:id/interactions/duration", getTotalInteractionDuration, middleware.NotGuest)
	r.GET("/user/:id/careplans", getCarePlansInUser, middleware.NotGuest)

	r.PUT("/user/:id/alert-threshold", upsertAlertThreshold, middleware.NotGuest, middleware.HasRole("doctor", "admin"))
	r.GET("/user/:id/alert-threshold", listAlertThresholds, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))

	r.PUT("/user/:id/diagnoses", upsertDiagnoses, middleware.NotGuest, middleware.HasRole("admin", "doctor", "nurse"))
	r.GET("/user/:id/diagnoses", getDiagnoses, middleware.NotGuest, middleware.HasRole("admin", "doctor", "nurse"))
}
