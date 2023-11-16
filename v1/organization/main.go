package organization

import (
	"MedKick-backend/pkg/echo/middleware"

	"github.com/labstack/echo/v4"
)

func Routes(r *echo.Group) {
	r.POST("/organization", createOrganization, middleware.NotGuest, middleware.HasRole("admin"))
	r.GET("/organization/:id", getOrganization, middleware.NotGuest)
	r.PATCH("/organization/:id", updateOrganization, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))
	r.DELETE("/organization/:id", deleteOrganization, middleware.NotGuest, middleware.HasRole("admin"))

	r.GET("/organization/:id/devices", getDevicesInOrganization, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))

	r.PUT("/organization/:id/interaction-setting", upsertInteractionSetting, middleware.NotGuest, middleware.HasRole("doctor", "admin"))
	r.GET("/organization/:id/interaction-setting", getInteractionSetting, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))

	r.GET("/organization/:id/telemetry-alert", listTelemetryAlert, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))
	r.PATCH("/organization/:id/telemetry-alert/:alert/resolve", resolvedTelemetryAlert, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))

	r.GET("/organization/:id/billing-report", getBillingReport, middleware.NotGuest, middleware.HasRole("admin"))
}
