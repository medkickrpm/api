package device

import (
	"MedKick-backend/pkg/echo/middleware"
	"github.com/labstack/echo/v4"
)

func Routes(r *echo.Group) {
	r.POST("/mio/forwardtelemetry", ingestTelemetry)

	r.GET("/mio/telemetry/:id", getTelemetry, middleware.NotGuest)

	r.GET("/device/:id", getDevice, middleware.NotGuest)
	r.PATCH("/device/:id", updateDevice, middleware.NotGuest)
	r.DELETE("/device/:id", deleteDevice, middleware.NotGuest, middleware.HasRole("admin"))
}
