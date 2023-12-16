package device

import (
	"MedKick-backend/pkg/echo/middleware"

	"github.com/labstack/echo/v4"
)

func Routes(r *echo.Group) {
	r.POST("/mio/forwardtelemetry", ingestTelemetry)
	r.POST("/mio/forwardstatus", ingestStatus)

	r.GET("/mio/telemetry/:id", getTelemetry, middleware.NotGuest)
	r.GET("/mio/telemetry/:id/latest", getLatestTelemetry, middleware.NotGuest)
	r.GET("/mio/telemetry/:id/count", getNumberOfTelemetryEntriesThisWeek, middleware.NotGuest)

	r.GET("/mio/status/:id", getStatus, middleware.NotGuest)

	r.GET("/device/available-devices", GetAvailableDevices, middleware.NotGuest)

	r.GET("/device/:id", getDevice, middleware.NotGuest)
	r.PATCH("/device/:id", updateDevice, middleware.NotGuest)
	r.DELETE("/device/:id", deleteDevice, middleware.NotGuest, middleware.HasRole("admin"))
}
